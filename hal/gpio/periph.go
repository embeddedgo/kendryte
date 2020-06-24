// Copyright 2020 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

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

// Synopsys DW_apb_gpio

type Port struct {
	dataOutput         mmio.U32
	direction          mmio.U32
	source             mmio.U32
	_                  [9]uint32
	interruptEnable    mmio.U32
	interruptMask      mmio.U32
	interruptLevel     mmio.U32
	interruptPolarity  mmio.U32
	interruptStatus    mmio.U32
	interruptStatusRaw mmio.U32
	interruptDebounce  mmio.U32
	interruptClear     mmio.U32
	dataInput          mmio.U32
	_                  [3]uint32
	syncLevel          mmio.U32
	idCode             mmio.U32
	interruptBothedge  mmio.U32
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

// Load returns input value of all pins.
func (p *Port) Load() Pins {
	return Pins(p.dataInput.Load())
}

// LoadOut returns output value of all pins.
func (p *Port) LoadOut() Pins {
	return Pins(p.dataOutput.Load())
}

// Store sets output value of all pins to value specified by val.
func (p *Port) Store(val Pins) {
	p.dataOutput.Store(uint32(val))
}
