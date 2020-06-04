// Copyright 2020 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"github.com/embeddedgo/kendryte/hal/fpioa"
	"github.com/embeddedgo/kendryte/p/bus"
	"github.com/embeddedgo/kendryte/p/sysctl"
	"github.com/embeddedgo/kendryte/p/uart"
)

func putc(c byte) {

}

func main() {
	fpioa.EnableClock()
	rx := fpioa.Pin(4)
	tx := fpioa.Pin(5)
	rx.Setup(fpioa.UART3_RX | fpioa.EnIE | fpioa.Schmitt)
	tx.Setup(fpioa.UART3_TX | fpioa.DriveH34L23 | fpioa.EnOE)

	sysctl.SYSCTL().UART3_CLK_EN().Set()

	baudRate := 115200
	div := uint32(bus.APB0.Clock() / int64(baudRate))
	dlh := div >> 12
	dll := (div - dlh<<12) >> 4
	dlf := div - dlh<<12 - dll<<4

	u := uart.UART3()
	u.LCR.SetBits(1 << 7)
	u.DLH_IER.U32.Store(dlh)
	u.RBR_DLL_THR.U32.Store(dll)
	u.DLF.U32.Store(dlf)
	u.LCR.U32.Store(0)

	//...

}
