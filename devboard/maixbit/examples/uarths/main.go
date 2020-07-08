// Copyright 2020 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"embedded/rtos"

	"github.com/embeddedgo/kendryte/hal/fpioa"
	"github.com/embeddedgo/kendryte/hal/irq"
	"github.com/embeddedgo/kendryte/hal/uarths"

	_ "github.com/embeddedgo/kendryte/devboard/maixbit/board/init"
)

var u *uarths.Driver

func main() {
	rx := fpioa.Pin(4)
	tx := fpioa.Pin(5)
	rx.Setup(fpioa.UARTHS_RX | fpioa.EnIE | fpioa.Schmitt)
	tx.Setup(fpioa.UARTHS_TX | fpioa.DriveH34L23 | fpioa.EnOE)

	u = uarths.NewDriver(uarths.UARTHS(1))
	u.SetBaudrate(9600)
	u.EnableTx()

	p := u.Periph()
	for {
		_, ok := p.Load()
		if !ok {
			break
		}
	}

	irq.UARTHS.Enable(rtos.IntPrioLow, irq.M0)

	for {
		u.WriteString("*0123456789abcdef0123456789abcdef0123456789abcdef0*\r\n")
	}
}

//go:interrupthandler
func UARTHS_Handler() {
	u.ISR()
}
