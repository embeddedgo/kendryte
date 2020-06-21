// Copyright 2020 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"time"

	_ "github.com/embeddedgo/kendryte/devboard/maixbit/board/init"
)

const n = 64

var a, b [n]uint64

func main() {
	t1 := time.Now()
	for i := 0; i < 1000*128/n; i++ {
		a = b
		b = a
		a = b
		b = a
		a = b
		b = a
		a = b
		b = a
	}
	t2 := time.Now()
	println(t2.Sub(t1).String())
}

// Results (K210 416 MHz):
//
// n=4:   duff=14.564423ms, loop=18.489062ms
// n=64:  duff=10.153846ms, loop=15.215384ms
// n=128: duff=10.011779ms, loop=15.026923ms
//
// loop means ssaConfig.noDuffDevice=true
