// Copyright 2020 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package timer

import (
	"embedded/mmio"
	"unsafe"

	"github.com/embeddedgo/kendryte/hal/internal"
	"github.com/embeddedgo/kendryte/p/bus"
	"github.com/embeddedgo/kendryte/p/mmap"
	"github.com/embeddedgo/kendryte/p/sysctl"
)

const (
	ENABLE     uint32 = 0x01 << 0 //+ ENABLE
	MODE       uint32 = 0x01 << 1 //+ MODE
	FREE       uint32 = 0x00 << 1 //  FREE_MODE
	USER       uint32 = 0x01 << 1 //  USER_MODE
	INTERRUPT  uint32 = 0x01 << 2 //+ INTERRUPT_MASK
	PWM_ENABLE uint32 = 0x01 << 3 //+ PWM_ENABLE
)

type Periph struct {
	ch           [4]Channel
	_            [20]mmio.U32
	intstat      mmio.U32
	eoi          mmio.U32
	raw_intstat  mmio.U32
	comp_version mmio.U32
	load_count2  [4]mmio.U32
}

func TIMER(n int) *Periph {
	if n < 0 || n > 2 {
		panic("timer: bad number")
	}
	return (*Periph)(unsafe.Pointer(mmap.TIMER0_BASE + uintptr(n)*0x10000))
}

func (p *Periph) Bus() bus.Bus {
	return bus.APB0
}

func (p *Periph) n() uintptr {
	return (uintptr(unsafe.Pointer(p)) - mmap.TIMER0_BASE) / 0x10000
}

func (p *Periph) EnableClock() {
	sc := sysctl.SYSCTL()
	mx := &internal.MX.SYSCTL

	mx.CLK_EN_CENT.Lock()
	if mx.APB0_CLK_EN == 0 {
		sc.APB0_CLK_EN().Set()
	}
	mx.APB0_CLK_EN++
	mx.CLK_EN_CENT.Unlock()

	mx.CLK_EN_PERI.Lock()
	sc.CLK_EN_PERI.SetBits(sysctl.TIMER0_CLK_EN << p.n())
	mx.CLK_EN_PERI.Unlock()
}

func (p *Periph) DisableClock() {
	sc := sysctl.SYSCTL()
	mx := &internal.MX.SYSCTL

	mx.CLK_EN_PERI.Lock()
	sc.CLK_EN_PERI.ClearBits(sysctl.TIMER0_CLK_EN << p.n())
	mx.CLK_EN_PERI.Unlock()

	mx.CLK_EN_CENT.Lock()
	mx.APB0_CLK_EN--
	if mx.APB0_CLK_EN == 0 {
		sc.APB0_CLK_EN().Clear()
	}
	mx.CLK_EN_CENT.Unlock()
}

func (p *Periph) SetInterval(ch int, nanoseconds int64) {
	clk := p.Bus().Clock()

	step := 1e9 / clk
	period := uint32(nanoseconds / step)

	if period < 0 || period > 2147483647 {
		panic("timer: period outside of 32bit range")
	}

	p.ch[ch].load_count.Store(period)
}

func (p *Periph) ClearIRQ(ch int) {
	p.ch[ch].eoi.Load()
}

func (p *Periph) Channel(ch int) *Channel {
	return &p.ch[ch]
}

type Channel struct {
	load_count mmio.U32
	current    mmio.U32
	control    mmio.U32
	eoi        mmio.U32
	intstat    mmio.U32
}

func (c *Channel) p() uintptr {
	return (uintptr(unsafe.Pointer(c)) - mmap.TIMER0_BASE) / 0x10000
}

func (c *Channel) n() uintptr {
	return (uintptr(unsafe.Pointer(c)) - (mmap.TIMER0_BASE + c.p()*0x10000)) / 0x14
}

func (c *Channel) Periph() *Periph {
	return TIMER(int(c.p()))
}

func (c *Channel) SetInterval(nanoseconds int64) {
	clk := c.Periph().Bus().Clock()

	step := 1e9 / clk
	period := uint32(nanoseconds / step)

	if period < 0 || period > 2147483647 {
		panic("timer: period outside of 32bit range")
	}

	c.load_count.Store(period)
}

func (c *Channel) EnableIRQ() {
	// Avoid interrupt storm if frequency has not been set
	if c.load_count.Load() == 0 {
		c.SetInterval(1e7)
	}

	// Clear any existing ISRs
	c.ClearIRQ()

	// Enable timer in user mode, unset interrupt mask if it was set
	c.control.SetBits(ENABLE | USER)
	c.control.ClearBits(INTERRUPT)
}

func (c *Channel) DisableIRQ() {
	c.control.SetBits(INTERRUPT)
}

func (c *Channel) ClearIRQ() {
	c.eoi.Load()
}
