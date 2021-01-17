// Copyright 2021 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package internal

import (
	"embedded/rtos"

	"github.com/embeddedgo/kendryte/hal/irq"
	"github.com/embeddedgo/kendryte/hal/uart"
)

// Driver returns a ready to use driver for UART1 peripheral.
func UART(n int) *uart.Driver {
	driver := uart.NewDriver(uart.UART(n)) // must before ir.Enable
	ctx := irq.M0
	ir := irq.UART1 + rtos.IRQ(n - 1)
	if ir&1 != 0 {
		ctx = irq.M1
	}
	ir.Enable(rtos.IntPrioLow, ctx)
	return driver
}
