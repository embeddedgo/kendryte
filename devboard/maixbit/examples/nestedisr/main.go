// Copyright 2020 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"embedded/rtos"
	"time"

	"github.com/embeddedgo/kendryte/hal/fpioa"
	"github.com/embeddedgo/kendryte/hal/gpiohs"
	"github.com/embeddedgo/kendryte/hal/irq"

	"github.com/embeddedgo/kendryte/devboard/maixbit/board/buttons"
	"github.com/embeddedgo/kendryte/devboard/maixbit/board/leds"
)

func main() {
	btn := buttons.User.Pin()
	btn.Setup(fpioa.GPIOHS0 | fpioa.EnIE) // set button pin as gpio.Pin0
	fpioa.Pin(9).Setup(fpioa.GPIOHS1 | fpioa.EnIE | fpioa.EnOE | fpioa.DriveH8L5)

	p := gpiohs.P(0)

	p.InpEn.Store(gpiohs.Pin0 | gpiohs.Pin1)  // enable input on both pins
	p.OutEn.Store(gpiohs.Pin1)                // enable output on Pin1
	p.FallIP.Store(gpiohs.Pin0 | gpiohs.Pin1) // clear falling edge pending bits
	p.RiseIP.Store(gpiohs.Pin1)               // clear rising edge pending bit
	p.FallIE.Store(gpiohs.Pin0 | gpiohs.Pin1) // enable IRQ on falling edge
	p.RiseIE.Store(gpiohs.Pin1)               // enable IRQ on rising edge

	irq.GPIOHS1.Enable(rtos.IntPrioMid, irq.M0)
	irq.GPIOHS0.Enable(rtos.IntPrioLow, irq.M0)

	for {
		leds.Green.Set(leds.Green.Get() + 1)
		println(p.InpVal.Load())
		time.Sleep(time.Second)
	}
}

//go:interrupthandler
func GPIOHS0_Handler() {
	p := gpiohs.P(0)
	p.FallIP.Store(gpiohs.Pin0)
	before := leds.Blue.Get()
	out := p.OutVal.Load()
	p.OutVal.Store(out ^ gpiohs.Pin1)
	for leds.Blue.Get() == before {
	}
}

//go:interrupthandler
func GPIOHS1_Handler() {
	p := gpiohs.P(0)
	p.FallIP.Store(gpiohs.Pin1)
	p.RiseIP.Store(gpiohs.Pin1)
	leds.Blue.Set(leds.Blue.Get() + 1)
}
