// Copyright 2020 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package uart

import (
	"runtime"
	"sync/atomic"
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
	d.p.SetTFT(TFT2)
	d.p.Store(int(s[0]))
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

/*
func (d *Driver) Flush() error {
	d.p.SetTFT(TFT0)
	return nil
}
*/
