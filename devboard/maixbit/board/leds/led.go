// Copyright 2020 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package leds

import (
	"github.com/embeddedgo/kendryte/hal/fpioa"

	_ "github.com/embeddedgo/kendryte/devboard/maixbit/board/system"
)

// Onboard LEDs
const (
	Green LED = 12
	Red   LED = 13
	Blue  LED = 14

	User = Green
)

type LED uint8

func (d LED) SetOn() {
	fpioa.Pin(d).Setup(fpioa.CONSTANT | fpioa.DriveH34L23 | fpioa.EnOE |
		fpioa.InvOE | fpioa.EnIE | fpioa.InvIE)
}

func (d LED) SetOff() {
	fpioa.Pin(d).Setup(fpioa.CONSTANT | fpioa.DriveH34L23 | fpioa.EnOE |
		fpioa.InvOE | fpioa.InvDO | fpioa.EnIE | fpioa.InvIE)
}

func (d LED) Set(on int) {
	if on&1 != 0 {
		d.SetOn()
	} else {
		d.SetOff()
	}
}

func (d LED) Get() int       { return fpioa.Pin(d).Load() ^ 1 }
func (d LED) Pin() fpioa.Pin { return fpioa.Pin(d) }

func init() {
	Red.SetOff()
	Green.SetOff()
	Blue.SetOff()
}
