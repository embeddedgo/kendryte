// Copyright 2020 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package gpiohs provides interface to configure and control GPIOHS peripheral.
// GPIOHS can controll up to 32 FPIOA pins.
package gpiohs

import (
	"embedded/mmio"
	"unsafe"

	"github.com/embeddedgo/kendryte/hal/internal"
	"github.com/embeddedgo/kendryte/p/bus"
	"github.com/embeddedgo/kendryte/p/mmap"
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
	Pin8
	Pin9
	Pin10
	Pin11
	Pin12
	Pin13
	Pin14
	Pin15
	Pin16
	Pin17
	Pin18
	Pin19
	Pin20
	Pin21
	Pin22
	Pin23
	Pin24
	Pin25
	Pin26
	Pin27
	Pin28
	Pin29
	Pin30
	Pin31
)

type PinReg struct{ U32 mmio.U32 }

func (r *PinReg) Load() Pins      { return Pins(r.U32.Load()) }
func (r *PinReg) Store(pins Pins) { r.U32.Store(uint32(pins)) }
func (r *PinReg) Set(pins Pins) {
	internal.AtomicSetBits(&r.U32, uint32(pins))
}
func (r *PinReg) Clear(pins Pins) {
	internal.AtomicClearBits(&r.U32, uint32(pins))
}
func (r *PinReg) Toggle(pins Pins) {
	internal.AtomicToggleBits(&r.U32, uint32(pins))
}

// SiFive GPIO

type Port struct {
	InpVal PinReg // data input
	InpEn  PinReg // input enable
	OutEn  PinReg // output enable
	OutVal PinReg // data output (can be read to get the last written value)
	pullup PinReg // internal pull-up enable
	drive  PinReg // drive strength
	RiseIE PinReg // rise interrupt enable
	RiseIP PinReg // rise interrupt pending
	FallIE PinReg // fall interrupt enable
	FallIP PinReg // fall interrupt pending
	HighIE PinReg // high interrupt enable
	HighIP PinReg // high interrupt pending
	LowIE  PinReg // low interrupt enable
	LowIP  PinReg // low interrupt pending
	iofEn  PinReg // enable hardware driven functions
	iofSel PinReg // select hardware driven function
	OutXor PinReg // invert
}

func P(n int) *Port {
	if n != 0 {
		panic("gpiohs: bad port number")
	}
	return (*Port)(unsafe.Pointer(mmap.GPIOHS_BASE))
}

func (p *Port) Bus() bus.Bus {
	return bus.TileLink
}
