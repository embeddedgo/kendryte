// Copyright 2020 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"embedded/rtos"
	"math"
	"time"

	_ "github.com/embeddedgo/kendryte/devboard/maixbit/board/init"
	"github.com/embeddedgo/kendryte/hal/fpioa"
	"github.com/embeddedgo/kendryte/hal/irq"
	"github.com/embeddedgo/kendryte/hal/timer"
)

var tickCount uint64

//go:interrupthandler
func TIMER0A_Handler() {
	timer.TIMER(0).ClearIRQ(0)
	timer.TIMER(0).ClearIRQ(1)

	tickCount++
}

func main() {
	freq := 100.0 // Hz

	// Pin assignment
	ch1 := fpioa.Pin(12)
	ch2 := fpioa.Pin(13)
	ch3 := fpioa.Pin(14)
	ch1.Setup(fpioa.TIMER0_TOGGLE2 | fpioa.EnOE | fpioa.DriveH34L23)
	ch2.Setup(fpioa.TIMER0_TOGGLE1 | fpioa.EnOE | fpioa.DriveH34L23)
	ch3.Setup(fpioa.TIMER0_TOGGLE3 | fpioa.EnOE | fpioa.DriveH34L23)

	// Peripheral
	p := timer.TIMER(0)

	// Driver instance
	r := timer.PWM(p.Channel(0))
	g := timer.PWM(p.Channel(1))
	b := timer.PWM(p.Channel(2))

	// Set frequency now so we don't have to wait around for the first period to expire
	r.SetFrequency(freq, .5)
	g.SetFrequency(freq, .5)
	b.SetFrequency(freq, .5)

	// Enable the timer and PWM function
	r.Enable()
	g.Enable()
	b.Enable()

	// Enable an ISR for the green channel
	g.EnableIRQ()
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
