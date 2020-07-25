// Copyright 2020 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package timer

type PWM struct {
	*Channel
}

func NewPWM(ch *Channel) *PWM {
	return &PWM{Channel: ch}
}

func (d *PWM) Enable() {
	d.Periph().EnableClock()

	d.control.SetBits(pwmEnable | enable | userMode | interruptMask)
}

func (d *PWM) SetFrequency(frequency float64, duty float64) {
	clk := float64(d.Periph().Bus().Clock() * 2)

	if frequency < 0 || frequency > 2147483647 {
		panic("pwm: frequency outside of 32bit range")
	}
	if duty < 0 || duty > 1 {
		panic("pwm: duty cycle must be 0.0-1.0")
	}
	period := uint32(clk / frequency)
	percent := uint32(duty * float64(period))

	d.load_count.Store(period - percent)
	d.Periph().load_count2[d.n()].Store(percent)
}
