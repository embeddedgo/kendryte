// Copyright 2020 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"embedded/rtos"

	"github.com/embeddedgo/kendryte/hal/irq"

	_ "github.com/embeddedgo/kendryte/devboard/maixbit/board/init"
)

func main() {
	irq.GPIO.Enable(rtos.IntPrioLow, 0)
	// WIP...
}


//go:interrupthandler
func GPIO_Handler() {
}