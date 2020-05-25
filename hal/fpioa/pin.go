// Copyright 2020 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fpioa

import (
	"embedded/mmio"
	"time"
	"unsafe"

	"github.com/embeddedgo/kendryte/hal/internal"
	"github.com/embeddedgo/kendryte/p/mmap"
	"github.com/embeddedgo/kendryte/p/sysctl"
)

func EnableClock() {
	sc := sysctl.SYSCTL()
	mx := &internal.MX.SYSCTL

	mx.CLK_EN_CENT.Lock()
	if mx.APB0_CLK_EN == 0 {
		sc.APB0_CLK_EN().Set()
	}
	mx.APB0_CLK_EN++
	mx.CLK_EN_CENT.Unlock()

	mx.CLK_EN_PERI.Lock()
	sc.FPIOA_CLK_EN().Set()
	mx.CLK_EN_PERI.Unlock()
}

func DisableClock() {
	sc := sysctl.SYSCTL()
	mx := &internal.MX.SYSCTL

	mx.CLK_EN_PERI.Lock()
	sc.FPIOA_CLK_EN().Clear()
	mx.CLK_EN_PERI.Unlock()

	mx.CLK_EN_CENT.Lock()
	mx.APB0_CLK_EN--
	if mx.APB0_CLK_EN == 0 {
		sc.APB0_CLK_EN().Clear()
	}
	mx.CLK_EN_CENT.Unlock()
}

func Reset() {
	sc := sysctl.SYSCTL()
	mx := &internal.MX.SYSCTL

	mx.PERI_RESET.Lock()
	sc.FPIOA_RESET().Set()
	mx.PERI_RESET.Unlock()

	time.Sleep(10 * time.Microsecond)

	mx.PERI_RESET.Lock()
	sc.FPIOA_RESET().Clear()
	mx.PERI_RESET.Unlock()
}

type periph struct {
	io     [48]mmio.U32
	tieEn  [8]mmio.U32
	tieVal [8]mmio.U32
}

func p() *periph {
	return (*periph)(unsafe.Pointer(mmap.FPIOA_BASE))
}

// Pin represents a FPIOA pin.
type Pin uint32

type Config uint32

const (
	Func      Config = 0xFF << 0 // Function number
	FuncNone  Config = 120 << 0  // No function
	FuncConst Config = 222 << 0  // Constatnt signal

	Drive       Config = 0xF << 8 // Driving strength [mA]:
	DriveH8L5   Config = 0 << 8   // src:  5.0-11.2, sink:  3.2- 8.3
	DriveH11L8  Config = 1 << 8   // src:  7.5-16.8, sink:  4.7-12.3
	DriveH15L11 Config = 2 << 8   // src: 10.0-22.3, sink:  6.3-16.4
	DriveH19L13 Config = 3 << 8   // src: 12.4-27.8, sink:  7.8-20.2
	DriveH23L16 Config = 4 << 8   // src: 14.9-33.3, sink:  9.4-24.2
	DriveH26L18 Config = 5 << 8   // src: 17.4-38.7, sink: 10.9-28.1
	DriveH30L21 Config = 6 << 8   // src: 19.8-44.1, sink: 12.4-31.8
	DriveH34L23 Config = 7 << 8   // src: 22.3-49.5, sink: 13.9-35.5

	OutEn     Config = 1 << 12 // Output enable (manual mode)
	OutEnInv  Config = 1 << 13 // Invert output enable
	OutSel    Config = 1 << 14 // Output selection: 0-data, 1-OE signal
	OutSelInv Config = 1 << 15 // Invert the output selection

	PullUp   Config = 1 << 16 // Enable internal pull-up
	PullDown Config = 1 << 17 // Enable internal pull-down
	SlevRate Config = 1 << 19 // Slew rate control enable

	InpEn      Config = 1 << 20 // Input enable (manual mode)
	InpEnInv   Config = 1 << 21 // Invert input enable
	InpDataInv Config = 1 << 22 // Invert data input

	Schmitt Config = 1 << 23 // Schmitt trigger

	InpVal Config = 1 << 31 // Current pin input value
)

func (pin Pin) Setup(cfg Config) {
	p().io[pin].Store(uint32(cfg))
}

func (pin Pin) Config() Config {
	return Config(p().io[pin].Load())
}
