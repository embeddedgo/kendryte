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

type periph struct {
	io     [48]mmio.U32
	tieEn  [8]mmio.U32
	tieVal [8]mmio.U32
}

func p() *periph {
	return (*periph)(unsafe.Pointer(mmap.FPIOA_BASE))
}

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

func init() {
	EnableClock()
	//Reset()
	Pin(0).Setup(JTAG_TCLK | EnIE | Schmitt)
	Pin(1).Setup(JTAG_TDI | EnIE | Schmitt)
	Pin(2).Setup(JTAG_TMS | EnIE | Schmitt)
	Pin(3).Setup(JTAG_TDO | DriveH34L23 | EnOE)
	for pin := Pin(4); pin < 48; pin++ {
		pin.Setup(None)
	}
}
