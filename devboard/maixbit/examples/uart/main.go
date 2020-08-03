// Copyright 2020 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"github.com/embeddedgo/kendryte/hal/fpioa"
	"github.com/embeddedgo/kendryte/hal/uart"

	_ "github.com/embeddedgo/kendryte/devboard/maixbit/board/init"
)

func main() {
	rx := fpioa.Pin(4)
	tx := fpioa.Pin(5)
	rx.Setup(fpioa.UART1_RX | fpioa.EnIE | fpioa.Schmitt)
	tx.Setup(fpioa.UART1_TX | fpioa.DriveH34L23 | fpioa.EnOE)

	u := uart.UART(1)
	u.EnableClock()
	u.Reset()
	u.SetConf1(uart.W8b)
	u.SetConf2(0)
	u.SetConf3(uart.FE | uart.CRF | uart.CTF | uart.TFT8 | uart.RFT1)
	u.SetConf4(uart.PTIME)
	u.SetBaudrate(115200)

	for {
		puts(u, "Hello, World!\r\n")
	}
}

func putc(u *uart.Periph, c byte) {
	for {
		if ev, _ := u.Status(); ev&uart.TxFull == 0 {
			break
		}
	}
	u.Store(int(c))
}

func puts(u *uart.Periph, s string) {
	for i := 0; i < len(s); i++ {
		putc(u, s[i])
	}
}
