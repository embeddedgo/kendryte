// Copyright 2020 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package uart

import (
	"runtime"
	"sync/atomic"
)

func (d *Driver) WriteString(s string) (n int, err error) {
	if len(s) == 0 {
		return 0, nil
	}
	// write as many bytes as possible in thread context
	for d.p.Status1()&TxFIFONotFull != 0 {
		d.p.Store(int(s[n]))
		if n++; n == len(s) {
			return
		}
	}
	// the remaining part will be written by ISR
	d.txdata = s
	d.txn = n
	m := len(s) - n
	tft := TFT2
	switch {
	case m <= 8:
		tft = TFT8
	case m <= 12:
		tft = TFT4
	}
	d.txdone.Clear()
	d.p.SetTFT(tft)
	d.mx.Lock()
	d.p.SetIntConf(d.p.IntConf() | TxReadyEn)
	d.mx.Unlock()
	ok := d.txdone.Sleep(d.timeoutTx)
	d.mx.Lock()
	d.p.SetIntConf(d.p.IntConf() &^ TxReadyEn)
	d.mx.Unlock()
	if !ok {
		d.txdata = d.txdata[:0]
		for atomic.LoadUint32(&d.isr) == isrTx {
			runtime.Gosched()
		}
		err = ErrTimeout
	}
	d.txdata = ""
	return d.txn, err
}
