// Copyright 2020 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"embedded/rtos"
	"runtime"
	"time"

	"github.com/embeddedgo/kendryte/devboard/maixbit/board/leds"
)

func blink(led leds.LED, n int) {
	runtime.LockOSThread()
	rtos.SetPrivLevel(0)
	for {
		if hartid() != 0 {
			led.SetOn()
			leds.Green.SetOn()
			time.Sleep(11 * time.Millisecond)
			//delay(1e6)
			leds.Green.SetOff()
			led.SetOff()
		} else {
			led.SetOn()
			time.Sleep(23 * time.Millisecond)
			//delay(1e6)
			led.SetOff()
		}
		delay(n)
	}
}

func main() {
	go blink(leds.Blue, 5e7)
	blink(leds.Red, 13e7)
}

func hartid() int
func delay(n int)
