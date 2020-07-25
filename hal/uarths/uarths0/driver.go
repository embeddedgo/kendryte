// Copyright 2020 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package uarths0

import (
	"embedded/rtos"
	_ "unsafe"

	"github.com/embeddedgo/kendryte/hal/irq"
	"github.com/embeddedgo/kendryte/hal/uarths"
)

var driver *uarths.Driver

// Driver returns a ready to use driver for UARTHS peripheral.
func Driver() *uarths.Driver {
	if driver == nil {
		driver = uarths.NewDriver(uarths.UARTHS(0))
		irq.UARTHS.Enable(rtos.IntPrioLow, irq.M0)
	}
	return driver
}

//go:interrupthandler
func _UARTHS_Handler() { driver.ISR() }

//go:linkname _UARTHS_Handler IRQ33_Handler
