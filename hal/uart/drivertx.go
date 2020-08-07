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
