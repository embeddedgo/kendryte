// Copyright 2021 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package uart

import "github.com/embeddedgo/kendryte/hal/fpioa"

type Signal uint8

const (
	RTSn Signal = iota
	TXD
	CTSn
	RXD
)

// UsePin is a helper function that can be used to configure FPIOA pins as
// required by UART peripheral.
func (d *Driver) UsePin(pin fpioa.Pin, sig Signal) {
	var cfg fpioa.Config
	if sig <= TXD {
		cfg = fpioa.DriveH34L23 | fpioa.EnOE
		if sig == TXD {
			cfg |= fpioa.UART1_TX + fpioa.Config(d.p.n()*2)
		} else {
			cfg |= fpioa.UART1_RTS + fpioa.Config(d.p.n()*14)
		}
	} else {
		cfg = fpioa.EnIE | fpioa.Schmitt
		if sig == RXD {
			cfg |= fpioa.UART1_RX + fpioa.Config(d.p.n()*2)
		} else {
			cfg |= fpioa.UART1_CTS + fpioa.Config(d.p.n()*14)
		}
	}
	pin.Setup(cfg)
}
