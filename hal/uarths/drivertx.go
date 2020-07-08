// Copyright 2020 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package uarths

import (
	"runtime"
	"sync/atomic"
	"unsafe"
)

const txMin = 1 // to minmize number of IRQs, can be >1 to ensure Tx continuity

func (d *Driver) EnableTx() {
	cfg, _ := d.p.TxConf()
	d.p.SetTxConf(cfg|TxEn, 0)
}

func (d *Driver) DisableTx() {
	cfg, _ := d.p.TxConf()
	d.p.SetTxConf(cfg&^TxEn, 0)
}

// WriteString works like Write but accepts string instead of byte slice.
func (d *Driver) WriteString(s string) (n int, err error) {
	if len(s) == 0 {
		return 0, nil
	}
	// write as many bytes as possible in thread context
	for {
		if d.p.TxFull() {
			break
		}
		d.p.Store(int(s[n]))
		if n++; n == len(s) {
			return
		}
	}
	// the remaining part will be written by ISR
	d.txdata = s
	d.txn = n
	m := 9 - n
	if m < txMin {
		m = txMin
	} else if m > 7 {
		m = 7
	}
	d.txdone.Clear()
	d.p.SetTxMinCnt(m)
	d.p.EnableIRQ(TxMin)
	if !d.txdone.Sleep(d.timeoutTx) {
		d.p.SetTxMinCnt(0)
		d.txdata = d.txdata[:0]
		for atomic.LoadUint32(&d.isr) == isrTx {
			runtime.Gosched()
		}
		err = ErrTimeout
	}
	d.p.DisableIRQ(TxMin)
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

func (d *Driver) Flush() error {
	d.p.SetTxMinCnt(1)
	for d.p.Events()&TxMin == 0 {
		runtime.Gosched()
	}
	d.p.SetTxMinCnt(0)
	return nil
}
