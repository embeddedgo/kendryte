// Copyright 2020 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package gpio provides interface to configure and control GPIO peripheral.
// GPIO can controll up to 8 FPIOA pins.
package gpio

import (
	"embedded/mmio"
	"time"
	"unsafe"

	"github.com/embeddedgo/kendryte/hal/internal"
	"github.com/embeddedgo/kendryte/p/bus"
	"github.com/embeddedgo/kendryte/p/mmap"
	"github.com/embeddedgo/kendryte/p/sysctl"
)

type Pins uint32

const (
	Pin0 Pins = 1 << iota
	Pin1
	Pin2
	Pin3
	Pin4
	Pin5
	Pin6
	Pin7
)

type PinReg struct{ U32 mmio.U32 }

func (r *PinReg) Load() Pins      { return Pins(r.U32.Load()) }
func (r *PinReg) Store(pins Pins) { r.U32.Store(uint32(pins)) }

// Synopsys DW_apb_gpio

type Port struct {
	OutVal       PinReg // values stored are output on the output pins
	Dir          PinReg // sets pins as outputs
	source       uint32
	_            [9]uint32
	IntEn        PinReg // enable pins to detect interrupt conditions
	IntMask      PinReg // disable generateing interrupts by pins
	IntEdge      PinReg // configure pins as edge sensitive
	IntPol       PinReg // configure pins as active high level / rising edge
	IntStatus    PinReg // lists active interrupt requests (IntRaw &^ IntMask)
	IntRaw       PinReg // lists raw interrupt requests (before masking)
	IntDebounce  PinReg // require active singal for 2 cycles of gpio_db_clk
	IntClear     PinReg // clear interrupt requests
	InpVal       PinReg // values of external signal on input pins
	_            [3]uint32
	IntLevelSync PinReg   // sync level-sensitive interrupts to l4_mp_clk
	IdCode       mmio.U32 // chip identification
	IntBothEdge  PinReg   // detect both edges on edge sensitive pins
}

func P(n int) *Port {
	if n != 0 {
		panic("gpio: bad port number")
	}
	return (*Port)(unsafe.Pointer(mmap.GPIO_BASE))
}

func (p *Port) Bus() bus.Bus {
	return bus.APB0
}

func (p *Port) EnableClock() {
	sc := sysctl.SYSCTL()
	mx := &internal.MX.SYSCTL

	mx.CLK_EN_CENT.Lock()
	if mx.APB0_CLK_EN == 0 {
		sc.APB0_CLK_EN().Set()
	}
	mx.APB0_CLK_EN++
	mx.CLK_EN_CENT.Unlock()

	mx.CLK_EN_PERI.Lock()
	sc.CLK_EN_PERI.SetBits(sysctl.GPIO_CLK_EN)
	mx.CLK_EN_PERI.Unlock()
}

func (p *Port) DisableClock() {
	sc := sysctl.SYSCTL()
	mx := &internal.MX.SYSCTL

	mx.CLK_EN_PERI.Lock()
	sc.CLK_EN_PERI.ClearBits(sysctl.GPIO_CLK_EN)
	mx.CLK_EN_PERI.Unlock()

	mx.CLK_EN_CENT.Lock()
	mx.APB0_CLK_EN--
	if mx.APB0_CLK_EN == 0 {
		sc.APB0_CLK_EN().Clear()
	}
	mx.CLK_EN_CENT.Unlock()
}

func (p *Port) Reset() {
	sc := sysctl.SYSCTL()
	mx := &internal.MX.SYSCTL

	mx.PERI_RESET.Lock()
	sc.PERI_RESET.SetBits(sysctl.GPIO_RESET)
	mx.PERI_RESET.Unlock()

	time.Sleep(10 * time.Microsecond)

	mx.PERI_RESET.Lock()
	sc.PERI_RESET.ClearBits(sysctl.GPIO_RESET)
	mx.PERI_RESET.Unlock()
}
