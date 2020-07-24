// Copyright 2020 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package timer

import (
	"github.com/embeddedgo/kendryte/p/bus"
	"github.com/embeddedgo/kendryte/p/timer"
)

type Driver struct {
	p *Periph

	channel int
}

func NewDriver(p *Periph, ch int) *Driver {
	d := Driver{
		p:       p,
		channel: ch,
	}

	return &d
}

func (d *Driver) EnablePWM() {
	d.p.EnableClock()

	bits := timer.PWM_ENABLE | timer.ENABLE | timer.USER | timer.INTERRUPT
	d.p.CH[d.channel].CONTROL.SetBits(bits)
}

func (d *Driver) EnableISR() {
	d.p.EnableClock()

	// Avoid interrupt storm if frequency has not been set
	if d.p.CH[d.channel].LOAD.Load() == 0 || d.p.LOAD_COUNT2[d.channel].Load() == 0 {
		d.SetFrequency(1, .5)
	}

	// Clear any existing ISRs
	d.p.ResetISR(d.channel)

	// Enable timer in user mode, unset interrupt mask if it was set
	d.p.CH[d.channel].CONTROL.SetBits(timer.ENABLE | timer.USER)
	d.p.CH[d.channel].CONTROL.ClearBits(timer.INTERRUPT)
}

func (d *Driver) SetFrequency(frequency int, duty float64) {
	clk := bus.APB0.Clock()

	if frequency < 0 || frequency > 2147483647 {
		panic("pwm: frequency outside of 32bit range")
	}
	if duty < 0 || duty > 1 {
		panic("pwm: duty cycle must be 0.0-1.0")
	}
	period := uint32(clk) / uint32(frequency)
	percent := uint32(duty * float64(period))

	d.p.CH[d.channel].LOAD.Store(timer.LOAD(period - percent))
	d.p.LOAD_COUNT2[d.channel].Store(timer.LOAD_COUNT2(percent))
}
