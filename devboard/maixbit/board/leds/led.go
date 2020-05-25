// Copyright 2020 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package leds

import (
	"github.com/embeddedgo/kendryte/hal/fpioa"

	_ "github.com/embeddedgo/kendryte/devboard/maixbit/board/init"
)

// Onboard LEDs
const (
	Green LED = 12
	Red   LED = 13
	Blue  LED = 14

	User = Green
)

type LED uint8

func (d LED) SetOn()         { fpioa.Pin(d).Clear() }
func (d LED) SetOff()        { fpioa.Pin(d).Set() }
func (d LED) Pin() fpioa.Pin { return fpioa.Pin(d) }

func (d LED) Set(on int) {
	pin := fpioa.Pin(d)
	if on&1 == 0 {
		pin.Set()
	} else {
		pin.Clear()
	}
}

func init() {
	fpioa.EnableClock()
	Red.SetOff()
	Green.SetOff()
	Blue.SetOff()
}
