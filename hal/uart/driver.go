// Copyright 2020 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package uart

import (
	"embedded/rtos"
	"sync"
	"sync/atomic"
	"time"
)

type DriverError uint8

const (
	// ErrBufOverflow is returned if one or more received bytes has been dropped
	// because of the lack of free space in the driver's receive buffer.
	ErrBufOverflow DriverError = iota + 1

	// ErrTimeout is returned if timeout occured. It means that the read/write
	// operation has been interrupted. In case of write you can not determine
	// the exact number of bytes sent to the remote party.
	ErrTimeout
)

// Error implements error interface.
func (e DriverError) Error() string {
	switch e {
	case ErrBufOverflow:
		return "uart: buffer overflow"
	case ErrTimeout:
		return "uart: timeout"
	}
	return ""
}

type Driver struct {
	p *Periph

	// rx state
	rxbuf   []byte
	nextr   uint32
	nextw   uint32
	rxcmd   uint32
	rxerr   uint32
	rxready rtos.Note

	// tx state
	txdata string
	txn    int
	txdone rtos.Note

	isr       uint32
	timeoutRx time.Duration
	timeoutTx time.Duration
	mx        sync.Mutex
}

// NewDriver returns a new driver for p.
func NewDriver(p *Periph) *Driver {
	return &Driver{p: p, timeoutRx: -1, timeoutTx: -1}
}

// Config is an unified configuration bitfield intended to be used by
//  Driver's Setup method. It combines bits from Conf1 and Conf2 bitfields.
type Config uint16

const (
	Word5b = Config(W5) // 5-bit data word
	Word6b = Config(W6) // 6-bit data word
	Word7b = Config(W7) // 7-bit data word
	Word8b = Config(W8) // 8-bit data word
	Stop2b = Config(S2) // 2 stop bits for 6 to 8-bit word, 1.5 for 5-bit word

	ParOdd  = Config(PE)       // parity control enabled: odd
	ParEven = Config(PE | EPS) // parity control enabled: even

	//HWFC = Config(RTS|AFCE) << 8 // hardware flow controll using RTS/CTS
	Loop = Config(LB) << 8       // loop-back diagnostic mode
	SIR  = Config(SIRE) << 8     // IrDA SIR (serial infrared) mode
)

func (d *Driver) Periph() *Periph {
	return d.p
}

func (d *Driver) Setup(cfg Config, baudrate int) {
	d.p.EnableClock()
	d.p.Reset()
	d.p.SetLineConf(LineConf(cfg))
	d.p.SetModeConf(ModeConf(cfg >> 8))
	d.p.SetFIFOConf(FE | TFT2 | RFT1)
	d.p.SetIntConf(PTIME)
	d.p.SetBaudrate(baudrate)
}

// SetBaudrate configures UART speed.
func (d *Driver) SetBaudrate(baudrate int) {
	d.p.SetBaudrate(baudrate)
}

// SetReadTimeout sets the read timeout used by Read* functions.
func (d *Driver) SetReadTimeout(timeout time.Duration) {
	d.timeoutRx = timeout
}

// SetWriteTimeout sets the write timeout used by Write* functions.
func (d *Driver) SetWriteTimeout(timeout time.Duration) {
	d.timeoutTx = timeout
}

const (
	isrNone = iota
	isrRx
	isrTx
)

// ISR handles UART interrupts.
func (d *Driver) ISR() {
	for {
		ir, _ := d.p.Int()
		switch ir {
		case RxReady:
			// todo
		case TxReady:
			atomic.StoreUint32(&d.isr, isrTx)
			for {
				if d.txn >= len(d.txdata) {
					d.txdone.Wakeup()
					break
				}
				if d.p.Status1()&TxFIFONotFull == 0 {
					m := len(d.txdata) - d.txn
					if m <= 12 {
						tft := TFT4
						if m <= 8 {
							tft = TFT8
						}
						d.p.SetTFT(tft)
					}
					break
				}
				d.p.Store(int(d.txdata[d.txn]))
				d.txn++
			}
			atomic.StoreUint32(&d.isr, isrNone)
		default:
			return
		}
	}
}
