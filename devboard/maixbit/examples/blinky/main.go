// Copyright 2020 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"time"

	"github.com/embeddedgo/kendryte/devboard/maixbit/board/buttons"
	"github.com/embeddedgo/kendryte/devboard/maixbit/board/leds"
)

func delay() {
	if buttons.User.Read() != 0 {
		time.Sleep(time.Second / 8)
	} else {
		time.Sleep(time.Second / 2)
	}
}

func main() {
	for {
		leds.Blue.SetOff()
		leds.Red.SetOn()
		delay()

		leds.Red.SetOff()
		leds.Green.SetOn()
		delay()

		leds.Green.SetOff()
		leds.Blue.SetOn()
		delay()
	}
}
