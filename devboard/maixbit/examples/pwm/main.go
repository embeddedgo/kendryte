// Copyright 2020 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"embedded/rtos"
	"math"
	"time"

	"github.com/embeddedgo/kendryte/devboard/maixbit/board/leds"
	"github.com/embeddedgo/kendryte/hal/fpioa"
	"github.com/embeddedgo/kendryte/hal/irq"
	"github.com/embeddedgo/kendryte/hal/timer"
)

var p *timer.Periph
var tickCount uint64

//go:interrupthandler
func TIMER0A_Handler() {
	// Clear the interrupt once we're done with it
	p.Channel(0).ClearIRQ()

	tickCount++
}

func main() {
	freq := 100.0 // Hz

	// Pin assignment
	ch1 := fpioa.Pin(leds.Red)
	ch2 := fpioa.Pin(leds.Green)
	ch3 := fpioa.Pin(leds.Blue)
	ch1.Setup(fpioa.TIMER0_TOGGLE1 | fpioa.EnOE | fpioa.DriveH34L23)
	ch2.Setup(fpioa.TIMER0_TOGGLE2 | fpioa.EnOE | fpioa.DriveH34L23)
	ch3.Setup(fpioa.TIMER0_TOGGLE3 | fpioa.EnOE | fpioa.DriveH34L23)

	// Peripheral is timer0
	p = timer.TIMER(0)

	// Driver instance
	r := timer.NewPWM(p.Channel(0))
	g := timer.NewPWM(p.Channel(1))
	b := timer.NewPWM(p.Channel(2))

	// Set frequency now so we don't have to wait around for the first period to expire
	r.SetFrequency(freq, .5)
	g.SetFrequency(freq, .5)
	b.SetFrequency(freq, .5)

	// Enable the timer and PWM function
	r.Enable()
	g.Enable()
	b.Enable()

	// Enable an interrupt for the red channel
	// This is only for demonstration purposes where you might want to change
	// the duty cycle on each clock and can be ommited
	r.EnableIRQ()
	irq.TIMER0A.Enable(rtos.IntPrioLow, irq.M0)

	// Animate duty cycle
	dc := 0.0
	tick := time.NewTicker(10 * time.Millisecond)
	for range tick.C {
		r.SetFrequency(freq, math.Abs(math.Sin(dc)))
		g.SetFrequency(freq, math.Abs(math.Sin(dc+math.Pi)))
		b.SetFrequency(freq, math.Abs(math.Cos(dc)))
		dc += 0.02
	}
}
