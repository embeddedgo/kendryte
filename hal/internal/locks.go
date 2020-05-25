// Copyright 2020 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package internal

import "sync"

var MX struct {
	SYSCTL struct {
		CLK_EN_CENT sync.Mutex
		APB0_CLK_EN int

		CLK_EN_PERI sync.Mutex
		PERI_RESET  sync.Mutex
	}
}
