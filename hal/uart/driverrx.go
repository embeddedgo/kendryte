// Copyright 2020 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package uart

import (
	"runtime"
	"sync/atomic"
)

// Len returns the number of bytes that are ready to read from Rx buffer.
func (d *Driver) Len() int {
	n := int(atomic.LoadUint32(&d.nextw)) - int(d.nextr)
	if n < 0 {
		n += len(d.rxbuf)
	}
	return n
}

// EnableRx enables the UART receiver. It allocates an internal ring buffer of
// bufLen size. In most cases bufLen = 64 is good choise. At least 2-byte buffer
// is required, which is effectively one byte buffer because the other byte
// always remains unused for efficient checking of an empty state. You cannot
// rely on 8-byte hardware buffer as an extension of the software buffer
// because for the performance reasons the ISR will no return until it has read
// all bytes from hardware. If the software buffer is full the ISR simply drops
// read bytes until there is no more data to read. EnableRx panics if the
// receiving is already enabled.
func (d *Driver) EnableRx(bufLen int) {
	if d.rxbuf != nil {
		panic("uarths: enabled before")
	}
	if bufLen < 2 {
		panic("uarths: rxbuf too short")
	}
	d.rxbuf = make([]byte, bufLen)
	d.nextr = 0
	d.nextw = 0
	d.p.SetRFT(RFT8)
	d.p.SetIntConf(PTIME | TxReadyEn | RxReadyEn)
}

// DisableRx disables the UART receiver and frees memory allocated for the
//internal ring buffer.
func (d *Driver) DisableRx() {
	d.p.SetIntConf(PTIME | TxReadyEn)
	for atomic.LoadUint32(&d.isr) == isrRx {
		runtime.Gosched()
	}
	d.rxbuf = nil
	return
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
	d.p.SetRFT(RFT1)       // RxReady interrupt after first received byte
	defer d.p.SetRFT(RFT8) // RxReady interrupt on half buffer (8 bytes)
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
	panic("uart: wakeup on empty buffer")
}

func (d *Driver) markDataRead(nextr int) error {
	if nextr >= len(d.rxbuf) {
		nextr -= len(d.rxbuf)
	}
	atomic.StoreUint32(&d.nextr, uint32(nextr))
	if rxerr := atomic.SwapUint32(&d.rxerr, 0); rxerr != 0 {
		if rxerr == 1 {
			return ErrBufOverflow
		}
		return Error(rxerr) // hardware errors mask ErrBufOverflow
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
