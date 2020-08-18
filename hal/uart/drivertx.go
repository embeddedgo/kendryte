// Copyright 2020 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package uart

import (
	"runtime"
	"sync/atomic"
	"time"
	"unsafe"
)

// WriteString works like Write but accepts string instead of byte slice.
func (d *Driver) WriteString(s string) (n int, err error) {
	if len(s) == 0 {
		return 0, nil
	}
	d.txdata = s
	d.txn = 1
	d.txdone.Clear()
	d.p.Store(int(s[0])) // there are at least 8 free bytes in Tx FIFO
	if !d.txdone.Sleep(d.timeoutTx) {
		d.txdata = d.txdata[:0]
		for atomic.LoadUint32(&d.isr) == isrTx {
			runtime.Gosched()
		}
		err = ErrTimeout
	}
	d.txdata = ""
	return d.txn, err
}

// Write sends bytes from p to the remote party. It return the number of bytes
// sent and error if detected. It does not provide any guarantee that the bytes
// sent were received by the remote party.
func (d *Driver) Write(p []byte) (int, error) {
	return d.WriteString(*(*string)(unsafe.Pointer(&p)))
}

// WriteByte sends one byte to the remote party. See Write for more information.
func (d *Driver) WriteByte(b byte) (err error) {
	s := struct {
		p *byte
		n int
	}{&b, 1}
	_, err = d.WriteString(*(*string)(unsafe.Pointer(&s)))
	return
}

// Flush waits until the hardware Tx buffer is empty. You cannot rely on Flush
// to for example set the external transceiver to Rx mode. It gives you no
// guarantee that all bits have been physically transmitted but only ensures the
// internal Tx FIFO is empty. You can check the TxDone bit to ensure the Tx
// shift register is empty but it also does not guarantee the last (stop) bit
// have left the external pin and reached the remote party or has been
// processed by external transceiver.
func (d *Driver) Flush() error {
	d.txdone.Clear()
	d.p.SetTFT(TFT0)
	d.txdone.Sleep(8 * 10 * time.Second / time.Duration(d.p.Baudrate()))
	d.p.SetTFT(TFT8)
	return nil
}
