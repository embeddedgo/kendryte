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

	"github.com/embeddedgo/kendryte/devboard/maixbit/board/leds"
)

func main() {
	// configure IRQ pins
	cfg := fpioa.EnIE | fpioa.EnOE | fpioa.DriveH8L5
	fpioa.Pin(10).Setup(fpioa.GPIOHS0 | cfg)
	fpioa.Pin(11).Setup(fpioa.GPIOHS1 | cfg)

	p := gpiohs.P(0)
	irqPins := gpiohs.Pin0 | gpiohs.Pin1

	// set IRQ pins low
	p.OutVal.Clear(irqPins)

	// enable both directions
	p.InpEn.Set(irqPins)
	p.OutEn.Set(irqPins)

	// clear edge detectors pending bits
	p.FallIP.Store(irqPins)
	p.RiseIP.Store(irqPins)

	// generate IRQ on both edges
	p.FallIE.Set(irqPins)
	p.RiseIE.Set(irqPins)

	// enable interrupts in PLIC
	irq.GPIOHS0.Enable(rtos.IntPrioLow, irq.M0)
	irq.GPIOHS1.Enable(rtos.IntPrioLow, irq.M0)

	for {
		leds.Red.SetOn()
		time.Sleep(time.Second)
		p.OutVal.Toggle(gpiohs.Pin0)
		time.Sleep(2 * time.Second)
		leds.Red.SetOff()
		time.Sleep(time.Second)
	}
}

//go:interrupthandler
func GPIOHS0_Handler() {
	p := gpiohs.P(0)
	p.FallIP.Store(gpiohs.Pin0)
	p.RiseIP.Store(gpiohs.Pin0)
	for i := 0; i < 3e6; i++ {
		leds.Green.SetOn()
	}
	p.OutVal.Toggle(gpiohs.Pin1)
	for i := 0; i < 3e6; i++ {
		leds.Green.SetOn()
	}
	leds.Green.SetOff()
}

//go:interrupthandler
func GPIOHS1_Handler() {
	p := gpiohs.P(0)
	p.FallIP.Store(gpiohs.Pin1)
	p.RiseIP.Store(gpiohs.Pin1)
	for i := 0; i < 3e6; i++ {
		leds.Blue.SetOn()
	}
	leds.Blue.SetOff()
}
