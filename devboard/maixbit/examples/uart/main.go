// Copyright 2020 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"embedded/rtos"

	"github.com/embeddedgo/kendryte/hal/fpioa"
	"github.com/embeddedgo/kendryte/hal/irq"
	"github.com/embeddedgo/kendryte/hal/uart"

	_ "github.com/embeddedgo/kendryte/devboard/maixbit/board/init"
)

var u *uart.Driver

func main() {
	println(uart.UART(1).CPR())

	rx := fpioa.Pin(4)
	tx := fpioa.Pin(5)
	rx.Setup(fpioa.UART1_RX | fpioa.EnIE | fpioa.Schmitt)
	tx.Setup(fpioa.UART1_TX | fpioa.DriveH34L23 | fpioa.EnOE)

	u = uart.NewDriver(uart.UART(1))
	u.Setup(uart.Word8b, 115200)

	irq.UART1.Enable(rtos.IntPrioLow, irq.M0)

	for {
		u.WriteString("0123456789abcdefghijklmnoprstuvwxyx\r\n")
	}
}

//go:interrupthandler
func UART1_Handler() { u.ISR() }
