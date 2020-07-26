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

// SetFrequency assigns the PWM channel with a clock rate in Hz and duty cycle
// between 0.0 and 1.0
func (d *PWM) SetFrequency(frequency float64, duty float64) {
	clk := float64(d.Periph().Bus().Clock() * 2)

	if frequency < 0 || frequency > 2147483647 {
		panic("pwm: frequency outside of 32bit range")
	}
	if duty < 0 || duty > 1 {
		panic("pwm: duty cycle must be 0.0-1.0")
	}
	period := int(clk / frequency)
	percent := int(duty * float64(period))

	d.SetLowTicks(period - percent)
	d.SetHighTicks(percent)
}
