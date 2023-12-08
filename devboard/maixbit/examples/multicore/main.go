// Copyright 2020 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"embedded/rtos"
	"math/rand"
	"runtime"

	_ "github.com/embeddedgo/kendryte/devboard/maixbit/board/system"
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
	//runtime.GOMAXPROCS(4)
	var lasthart [30]int
	for i := range lasthart {
		go thread(i)
	}
	runtime.LockOSThread()
	rtos.SetPrivLevel(0)
	for r := range ch {
		lasthart[r.tid] = r.hartid
		hid := hartid()
		print(hid>>8, hid&0xFF)
		for _, hid := range lasthart {
			print(" ", hid>>8, hid&0xFF)
		}
		println()
	}
}

func hartid() int
func loop(n int)
