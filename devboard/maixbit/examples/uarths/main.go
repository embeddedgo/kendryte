// Copyright 2020 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"github.com/embeddedgo/kendryte/hal/fpioa"
	"github.com/embeddedgo/kendryte/hal/uarths"

	_ "github.com/embeddedgo/kendryte/devboard/maixbit/board/init"
)

func main() {
	rx := fpioa.Pin(4)
	tx := fpioa.Pin(5)
	rx.Setup(fpioa.UARTHS_RX | fpioa.EnIE | fpioa.Schmitt)
	tx.Setup(fpioa.UARTHS_TX | fpioa.DriveH34L23 | fpioa.EnOE)

	u := uarths.UARTHS(1)
	u.SetTxConf(uarths.TxEn, 2)
	u.SetBaudrate(1500e3)

	for {
		puts(u, "Hello, World!\r\n")
	}
}

func putc(u *uarths.Periph, c byte) {
	for u.TxFull() {
	}
	u.Store(int(c))
}

func puts(u *uarths.Periph, s string) {
	for i := 0; i < len(s); i++ {
		putc(u, s[i])
	}
}
