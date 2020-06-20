// Copyright 2020 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"embedded/rtos"
	"math/rand"
	"runtime"

	_ "github.com/embeddedgo/kendryte/devboard/maixbit/board/init"
)

type report struct {
	tid, hartid int
}

var ch = make(chan report, 3)

func thread(tid int) {
	runtime.LockOSThread()
	rtos.SetPrivLevel(0)
	for {
		loop(1e6 + rand.Intn(1e6))
		ch <- report{tid, hartid()}
	}
}

func main() {
	var lasthart [30]int
	for i := range lasthart {
		go thread(i)
	}
	runtime.LockOSThread()
	rtos.SetPrivLevel(0)
	for r := range ch {
		lasthart[r.tid] = r.hartid
		hid := hartid()
		print(hid>>1, hid&1)
		for _, hid := range lasthart {
			print(" ", hid>>1, hid&1)
		}
		println()
	}
}

func hartid() int
func loop(n int)
