// Copyright 2021 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package uart3

import (
	_ "unsafe"

	"github.com/embeddedgo/kendryte/hal/uart"
	"github.com/embeddedgo/kendryte/hal/uart/internal"
)

var driver *uart.Driver

// Driver returns a ready to use driver for UART3 peripheral.
func Driver() *uart.Driver {
	if driver == nil {
		driver = internal.UART(3)
	}
	return driver
}

//go:interrupthandler
func _UART3_Handler() { driver.ISR() }

//go:linkname _UART3_Handler IRQ13_Handler
