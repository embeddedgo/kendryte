// Copyright 2020 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package uarths

import (
	"embedded/mmio"
	"unsafe"

	"github.com/embeddedgo/kendryte/hal/internal"
	"github.com/embeddedgo/kendryte/p/bus"
	"github.com/embeddedgo/kendryte/p/mmap"
)

// SiFive UART

// Periph represents UARTHS peripheral.
type Periph struct {
	txd    mmio.U32
	rxd    mmio.U32
	txctrl mmio.U32
	rxctrl mmio.U32
	ie     mmio.U32
	ip     mmio.U32
	brdiv  mmio.U32
}

func UARTHS(n int) *Periph {
	if n != 1 {
		panic("uarths: bad number")
	}
	return (*Periph)(unsafe.Pointer(mmap.UARTHS_BASE))
}

func (p *Periph) Bus() bus.Bus {
	return bus.TileLink
}

type ConfTx uint8

const (
	TxEn   ConfTx = 1 << 0
	Stop2b ConfTx = 1 << 1
)

func (p *Periph) ConfTx() (c ConfTx, mincnt int) {
	ctrl := p.txctrl.Load()
	return ConfTx(ctrl & 3), int(ctrl >> 16 & 7)
}

func (p *Periph) SetConfTx(c ConfTx, mincnt int) {
	checkCnt(mincnt)
	p.txctrl.Store(uint32(c) | uint32(mincnt<<16))
}

// SetTxMinLevel sets the minium level of Tx queue below which the TxMin event
// is generated.
func (p *Periph) SetTxMinCnt(mincnt int) {
	checkCnt(mincnt)
	p.txctrl.StoreBits(7<<16, uint32(mincnt<<16))
}

type ConfRx uint8

const (
	RxEn ConfRx = 1 << 0
)

func (p *Periph) ConfRx() (c ConfRx, wm int) {
	ctrl := p.rxctrl.Load()
	return ConfRx(ctrl & 1), int(ctrl >> 16 & 7)
}

func (p *Periph) SetConfRx(c ConfTx, maxcnt int) {
	checkCnt(maxcnt)
	p.rxctrl.Store(uint32(c) | uint32(maxcnt<<16))
}

// SetRxMaxCnt sets the maximum level of Rx queue above which the RxMax event
// is generated.
func (p *Periph) SetRxMaxCnt(maxcnt int) {
	checkCnt(maxcnt)
	p.rxctrl.StoreBits(7<<16, uint32(maxcnt<<16))
}

func (p *Periph) Load() (rxd int, ok bool) {
	d := p.rxd.Load()
	return int(d & 0xFF), d>>31 == 0
}

func (p *Periph) Store(d int) {
	p.txd.Store(uint32(d))
}

func (p *Periph) TxFull() bool {
	return p.txd.Load()>>31 != 0
}

func (p *Periph) Baudrate() int {
	return int(bus.TileLink.Clock() / int64(p.brdiv.Load()+1))
}

func (p *Periph) SetBaudrate(br int) {
	div := bus.TileLink.Clock() / int64(br)
	if div < 16 || div > 1<<16 {
		panic("uarths: bad baudrate")
	}
	p.brdiv.Store(uint32(div - 1))
}

type Event uint8

const (
	TxMin Event = 1 << 0
	RxMax Event = 1 << 1
)

func (p *Periph) Event() Event {
	return Event(p.ip.Load())
}

func (p *Periph) EnableIRQ(ev Event) {
	internal.AtomicSetBits(&p.ie, uint32(ev))
}

func (p *Periph) DisableIRQ(ev Event) {
	internal.AtomicClearBits(&p.ie, uint32(ev))
}

func checkCnt(cnt int) {
	if uint(cnt) > 7 {
		panic("uarths: bad min/max queue level")
	}
}
