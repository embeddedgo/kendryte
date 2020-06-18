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
	tid, cpuid int
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
	var lastcpu [39]int
	for i := range lastcpu {
		go thread(i)
	}
	for r := range ch {
		lastcpu[r.tid] = r.cpuid
		for _, cpuid := range lastcpu {
			print(" ", cpuid)
		}
		print("\r\n")
	}
}

func hartid() int
func loop(n int)
