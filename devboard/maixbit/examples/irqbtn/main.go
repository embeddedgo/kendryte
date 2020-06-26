// Copyright 2020 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"embedded/rtos"
	"time"

	"github.com/embeddedgo/kendryte/hal/fpioa"
	"github.com/embeddedgo/kendryte/hal/gpio"
	"github.com/embeddedgo/kendryte/hal/irq"

	"github.com/embeddedgo/kendryte/devboard/maixbit/board/buttons"
	"github.com/embeddedgo/kendryte/devboard/maixbit/board/leds"
)

func main() {
	btn := buttons.User.Pin()
	btn.Setup(fpioa.GPIO0 | fpioa.EnIE) // set button pin as gpio.Pin0.

	p := gpio.P(0)
	p.EnableClock()
	p.Reset()
	p.IntEn.Store(gpio.Pin0)       // enable interrupt detecton on Pin0
	p.IntEdge.Store(gpio.Pin0)     // configure Pin0 as edge sensitive
	p.IntDebounce.Store(gpio.Pin0) // enable debouncing (doesn't help much)

	irq.GPIO.Enable(rtos.IntPrioLow, irq.M0)

	for {
		leds.Green.SetOn()
		println(leds.Green.Get())
		time.Sleep(time.Second)
		leds.Green.SetOff()
		println(leds.Green.Get())
		time.Sleep(time.Second)
	}
}

//go:interrupthandler
func GPIO_Handler() {
	gpio.P(0).IntClear.Store(gpio.Pin0)
	leds.Blue.Set(leds.Blue.Get() + 1)
}
