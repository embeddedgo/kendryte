// Copyright 2020 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package uart

func (d *Driver) EnableTx() {
	mx.Lock()
	d.p.SetIntConf(PTHRE | ETBEI | d.p.IntConf())
	mx.Unlock()
}
