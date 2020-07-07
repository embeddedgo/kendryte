// Copyright 2020 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package uarths

import (
	"embedded/rtos"
	"runtime"
	"sync/atomic"
	"time"
	"unsafe"
)

type Driver struct {
	p *Periph

	// rx state
	rxbuf    []byte
	nextr    uint32
	nextw    uint32
	rxcmd    uint32
	rxready  rtos.Note
	overflow bool

	// tx state
	txdata string
	txn    int
	txdone rtos.Note

	isr       uint32
	timeoutRx time.Duration
	timeoutTx time.Duration
}

const (
	cmdNone = iota
	cmdWakeup
)

const (
	isrNone = iota
	isrRx
	isrTx
)

// NewDriver returns a new driver for p.
func NewDriver(p *Periph) *Driver {
	return &Driver{p: p, timeoutRx: -1, timeoutTx: -1}
}

func (d *Driver) Periph() *Periph {
	return d.p
}

// EnableRx enables the UARTHS receiver using the provided slice to buffer
// received data. At least 2-byte buffer is required, which is effectively one
// byte buffer because the other byte always remains unused for efficient
// checking of an empty state. You can not rely on 7-byte hardware buffer as
// extension of software buffer because for the performance reasons the ISR do
// not return until it reads all bytes from hardware. If the software buffer is
// full the ISR simply drops read bytes until there is no more data to read.
// EnableRx panics if the receiving is already enabled or rxbuf is too short.
func (d *Driver) EnableRx(rxbuf []byte) {
	if d.rxbuf != nil {
		panic("uarths: enabled before")
	}
	if len(rxbuf) < 2 {
		panic("uarths: rxbuf too short")
	}
	d.rxbuf = rxbuf
	d.nextr = 0
	d.nextw = 0
	cfg, _ := d.p.RxConf()
	cfg |= RxEn
	d.p.SetRxConf(cfg, 0)
}

// DisableRx disables the UART receiver. The receive buffer is returned and no
// longer referenced by driver. You can use the Periph.DisableRx if you want to
// temporary disable the receiver leaving the driver intact.
func (d *Driver) DisableRx() (rxbuf []byte) {
	d.p.SetRxConf(0, 0)
	d.p.DisableIRQ(RxMax)
	for {
		if _, ok := d.p.Load(); !ok {
			break
		}
	}
	for atomic.LoadUint32(&d.isr) == isrRx {
		runtime.Gosched()
	}
	rxbuf = d.rxbuf
	d.rxbuf = nil
	return
}

// ISR handles UARTHS interrupts. It supports the reading thread to run in
// parallel on another hart.
func (d *Driver) ISR() {
	if b, ok := d.p.Load(); ok {
		atomic.StoreUint32(&d.isr, isrRx)
		for {
			nextw := d.nextw
			d.rxbuf[nextw] = byte(b)
			if nextw++; int(nextw) == len(d.rxbuf) {
				nextw = 0
			}
			if nextw != atomic.LoadUint32(&d.nextr) {
				atomic.StoreUint32(&d.nextw, nextw)
				if atomic.CompareAndSwapUint32(&d.rxcmd, cmdWakeup, cmdNone) {
					d.rxready.Wakeup()
				}
			} else {
				d.overflow = true
			}
			if b, ok = d.p.Load(); !ok {
				break
			}
		}
		atomic.StoreUint32(&d.isr, isrNone)
	}

	if d.p.Events()&TxMin == 0 {
		atomic.StoreUint32(&d.isr, isrTx)
		if d.txn >= len(d.txdata) {
			cfg, _ := d.p.TxConf()
			d.p.SetTxConf(cfg&^TxEn, 0)
			d.txdone.Wakeup()
		}
		for d.txn < len(d.txdata) && !d.p.TxFull() {
			d.p.Store(int(d.txdata[d.txn]))
			d.txn++
		}
		atomic.StoreUint32(&d.isr, isrNone)
	}
}

// Len returns the number of bytes that are ready to read from Rx buffer.
func (d *Driver) Len() int {
	n := int(atomic.LoadUint32(&d.nextw)) - int(d.nextr)
	if n < 0 {
		n += len(d.rxbuf)
	}
	return n
}

func (d *Driver) waitRxData() int {
	nextw := atomic.LoadUint32(&d.nextw)
	if nextw != d.nextr {
		return int(nextw)
	}
	d.rxready.Clear()
	atomic.StoreUint32(&d.rxcmd, cmdWakeup)
	nextw = atomic.LoadUint32(&d.nextw)
	if nextw != d.nextr {
		if atomic.SwapUint32(&d.rxcmd, cmdNone) == cmdNone {
			d.rxready.Sleep(-1) // wait for the upcoming wake up
		}
		return int(nextw)
	}
	if !d.rxready.Sleep(d.timeoutRx) {
		if atomic.SwapUint32(&d.rxcmd, cmdNone) != cmdNone {
			return int(nextw)
		}
		d.rxready.Sleep(-1) // wait for the upcoming wake up
	}
	nextw = atomic.LoadUint32(&d.nextw)
	if nextw != d.nextr {
		return int(nextw)
	}
	panic("uarths: wakeup on empty buffer")
}

func (d *Driver) markDataRead(nextr int) error {
	if nextr >= len(d.rxbuf) {
		nextr -= len(d.rxbuf)
	}
	atomic.StoreUint32(&d.nextr, uint32(nextr))
	if d.overflow {
		d.overflow = false
		return ErrBufOverflow
	}
	return nil
}

// ReadByte reads one byte and returns error if detected. ReadByte blocks only
// if the internal buffer is empty (d.Len() > 0 ensure non-blocking read).
func (d *Driver) ReadByte() (b byte, err error) {
	nextw := d.waitRxData()
	nextr := int(d.nextr)
	if nextw == nextr {
		return 0, ErrTimeout
	}
	return d.rxbuf[nextr], d.markDataRead(nextr + 1)
}

// Read reads up to len(p) bytes into p. It returns number of bytes read and an
// error if detected. Read blocks only if the internal buffer is empty (d.Len()
// > 0 ensure non-blocking read).
func (d *Driver) Read(p []byte) (n int, err error) {
	if len(p) == 0 {
		return 0, nil
	}
	nextw := d.waitRxData()
	nextr := int(d.nextr)
	if nextw == nextr {
		return 0, ErrTimeout
	}
	if nextr <= nextw {
		n = copy(p, d.rxbuf[nextr:nextw])
	} else {
		n = copy(p, d.rxbuf[nextr:])
		if n < len(p) {
			n += copy(p[n:], d.rxbuf[:nextw])
		}
	}
	return n, d.markDataRead(nextr + n)
}

// WriteByte sends one byte to the remote party and returns an error if detected// WriteByte does not provide any guarantee that the byte sent was received by
// the remote party.
func (d *Driver) WriteByte(b byte) (err error) {
	d.txdone.Clear()
	cfg, _ := d.p.TxConf()
	d.p.SetTxConf(cfg|TxEn, 1)
	d.p.Store(int(b))
	d.p.EnableIRQ(TxMin)
	if !d.txdone.Sleep(d.timeoutTx) {
		cfg, _ := d.p.TxConf()
		d.p.SetTxConf(cfg&^TxEn, 0)
		err = ErrTimeout
	}
	d.p.DisableIRQ(TxMin)
	for atomic.LoadUint32(&d.isr) == isrTx {
		runtime.Gosched()
	}
	return
}

// WriteString works like Write.
func (d *Driver) WriteString(s string) (n int, err error) {
	if len(s) == 0 {
		return 0, nil
	}
	d.txdata = s
	d.txn = 1
	d.txdone.Clear()
	cfg, _ := d.p.TxConf()
	d.p.SetTxConf(cfg|TxEn, 1)
	d.p.Store(int(s[0]))
	d.p.EnableIRQ(TxMin)
	if !d.txdone.Sleep(d.timeoutTx) {
		cfg, _ := d.p.TxConf()
		d.p.SetTxConf(cfg&^TxEn, 0)
		err = ErrTimeout
	}
	d.p.DisableIRQ(TxMin)
	d.txdata = d.txdata[:0]
	for atomic.LoadUint32(&d.isr) == isrTx {
		runtime.Gosched()
	}
	d.txdata = ""
	return d.txn, err // BUG: in case of timeout the ISR can still run in multicore system
}

// Write sends bytes from p to the remote party. It return the number of bytes
// sent and error if detected. It does not provide any guarantee that the bytes
// sent were received by the remote party.
func (d *Driver) Write(p []byte) (int, error) {
	return d.WriteString(*(*string)(unsafe.Pointer(&p)))
}

// SetReadTimeout sets the read timeout used by Read* functions.
func (d *Driver) SetReadTimeout(timeout time.Duration) {
	d.timeoutRx = timeout
}

// SetWriteTimeout sets the write timeout used by Write* functions.
func (d *Driver) SetWriteTimeout(timeout time.Duration) {
	d.timeoutTx = timeout
}

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
