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
	enable        uint32 = 0x01 << 0
	freeMode      uint32 = 0x00 << 1
	userMode      uint32 = 0x01 << 1
	interruptMask uint32 = 0x01 << 2
	pwmEnable     uint32 = 0x01 << 3
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

func (p *Periph) ClearIRQs() {
	p.eoi.Load()
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

// Periph returns the parent timer of this channel
func (c *Channel) Periph() *Periph {
	return TIMER(int(c.p()))
}

// SetLowTicks assigns the load_count register the number of clock cycles to
// count down. i.e. the interval when operating as a timer, or the high period
// when operating as a PWM.
// Conversion from time units to ticks can be calculated by retrieaving the
// peripherals bus clock with: myChannel.Periph().Bus().Clock()
func (c *Channel) SetLowTicks(ticks int) {
	if ticks < 0 || ticks > 2147483647 {
		panic("timer: period outside of 32bit range")
	}
	c.load_count.Store(uint32(ticks))
}

// SetHighTicks assigns the load_count2 register the number of clock cycles to
// count down while driving the PWM output low. Total timer interval is the sum
// of load_count and load_count2
func (c *Channel) SetHighTicks(ticks int) {
	if ticks < 0 || ticks > 2147483647 {
		panic("timer: period outside of 32bit range")
	}
	c.Periph().load_count2[c.n()].Store(uint32(ticks))
}

// DutyRegs returns pointers to the timers load_count and load_count2 registers
// for manipulation inside of interrupts or DMA operations.
func (c *Channel) DutyRegs() (lowTicks, highTicks *mmio.U32) {
	lowTicks = &c.load_count
	highTicks = &c.Periph().load_count2[c.n()]
	return
}

// EnableIRQ unsets the interrupt mask and enables the channel, if a interval
// has not been set the timer will be defaulted to 10ms.
func (c *Channel) EnableIRQ() {
	// Avoid interrupt storm if frequency has not been set
	if c.load_count.Load() == 0 {
		c.SetLowTicks(2000000) // 10ms
	}

	// Clear any existing ISRs
	c.ClearIRQ()

	// Enable timer in user mode, unset interrupt mask if it was set
	c.control.SetBits(enable | userMode)
	c.control.ClearBits(interruptMask)
}

func (c *Channel) DisableIRQ() {
	c.control.SetBits(interruptMask)
}

// ClearIRQ marks the interrupt as complete
func (c *Channel) ClearIRQ() {
	c.eoi.Load()
}
