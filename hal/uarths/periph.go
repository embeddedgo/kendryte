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

type RxConf uint16

const (
	RxEn RxConf = 1 << 0
)

func (p *Periph) RxConf() (c RxConf, maxCnt int) {
	v := p.rxctrl.Load()
	return RxConf(v & 0xFFFF), int(v>>16) & 7
}

func (p *Periph) SetRxConf(c RxConf, maxCnt int) {
	checkCnt(maxCnt)
	p.rxctrl.Store(uint32(c) | uint32(maxCnt<<16))
}

func (p *Periph) SetRxMaxCnt(maxCnt int) {
	checkCnt(maxCnt)
	p.rxctrl.StoreBits(7<<16, uint32(maxCnt<<16))
}

func (p *Periph) EnableRx()  { p.rxctrl.SetBits(uint32(RxEn)) }
func (p *Periph) DisableRx() { p.rxctrl.ClearBits(uint32(RxEn)) }

type TxConf uint16

const (
	TxEn     TxConf = 1 << 0
	TxStop2b TxConf = 1 << 1
)

func (p *Periph) TxConf() (c TxConf, minCnt int) {
	v := p.txctrl.Load()
	return TxConf(v & 0xFFFF), int(v>>16) & 7
}

func (p *Periph) SetTxConf(c TxConf, minCnt int) {
	checkCnt(minCnt)
	p.txctrl.Store(uint32(c) | uint32(minCnt<<16))
}

func (p *Periph) SetTxMinCnt(minCnt int) {
	checkCnt(minCnt)
	p.txctrl.StoreBits(7<<16, uint32(minCnt<<16))
}

func (p *Periph) EnableTx()  { p.txctrl.SetBits(uint32(TxEn)) }
func (p *Periph) DisableTx() { p.txctrl.ClearBits(uint32(TxEn)) }

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
	TxMin Event = 1 << 0 // raised if len(Tx FIFO) < minCnt
	RxMax Event = 1 << 1 // raised if len(Rx FIFO) > maxCnt
)

func (p *Periph) Events() Event {
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
