// Copyright 2020 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"time"

	_ "github.com/embeddedgo/kendryte/devboard/maixbit/board/init"
)

const n = 4

var (
	a = new([n]uint64)
	b = new([n]uint64)
)

func main() {
	for {
		t1 := time.Now()
		for i := 0; i < 1000*128/n; i++ {
			*a = *b
			*b = *a
			*a = *b
			*b = *a
			*a = *b
			*b = *a
			*a = *b
			*b = *a
		}
		t2 := time.Now()
		println(t2.Sub(t1).String())
		time.Sleep(time.Second)
	}
}

// Results (K210 416 MHz):
//
// n=2:   duff=19.6ms, loop=27.6ms, inline=12.2ms
// n=3:   duff=16.6ms, loop=23.9ms, inline=10.0ms
// n=4:   duff=15.1ms, loop=22.0ms, inline=9.46ms
// n=64:  duff=10.7ms, loop=17.0ms
// n=128: duff=10.0ms, loop=16.1ms
