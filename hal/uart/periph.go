// Copyright 2020 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package uart

import (
	"embedded/mmio"
	"unsafe"

	"github.com/embeddedgo/kendryte/hal/internal"
	"github.com/embeddedgo/kendryte/p/bus"
	"github.com/embeddedgo/kendryte/p/mmap"
	"github.com/embeddedgo/kendryte/p/sysctl"
)

// Periph represents UART peripheral.
type Periph struct {
	rbr_dll_thr mmio.U32
	dlh_ier     mmio.U32
	fcr_iir     mmio.U32
	lcr         mmio.U32
	mcr         mmio.U32
	lsr         mmio.U32
	msr         mmio.U32
	scr         mmio.U32
	lpdll       mmio.U32
	lpdlh       mmio.U32
	_           [2]uint32
	srbr_sthr   [16]mmio.U32
	far         mmio.U32
	tfr         mmio.U32
	rfw         mmio.U32
	usr         mmio.U32
	tfl         mmio.U32
	rfl         mmio.U32
	srr         mmio.U32
	srts        mmio.U32
	sbcr        mmio.U32
	sdmam       mmio.U32
	sfe         mmio.U32
	srt         mmio.U32
	stet        mmio.U32
	htx         mmio.U32
	dmasa       mmio.U32
	tcr         mmio.U32
	deen        mmio.U32
	reen        mmio.U32
	det         mmio.U32
	tat         mmio.U32
	dlf         mmio.U32
	rar         mmio.U32
	tar         mmio.U32
	lcr_ext     mmio.U32
	_           [9]uint32
	cpr         mmio.U32
	ucv         mmio.U32
	ctr         mmio.U32
}

func UART(n int) *Periph {
	if n < 1 || n > 3 {
		panic("uart: bad number")
	}
	return (*Periph)(unsafe.Pointer(mmap.UART1_BASE + uintptr(n-1)*0x10000))
}

func (p *Periph) Bus() bus.Bus {
	return bus.APB0
}

func (p *Periph) n() uintptr {
	return (uintptr(unsafe.Pointer(p)) - mmap.UART1_BASE) / 0x10000
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
	sc.CLK_EN_PERI.SetBits(sysctl.UART1_CLK_EN << p.n())
	mx.CLK_EN_PERI.Unlock()
}

func (p *Periph) DisableClock() {
	sc := sysctl.SYSCTL()
	mx := &internal.MX.SYSCTL

	mx.CLK_EN_PERI.Lock()
	sc.CLK_EN_PERI.ClearBits(sysctl.UART1_CLK_EN << p.n())
	mx.CLK_EN_PERI.Unlock()

	mx.CLK_EN_CENT.Lock()
	mx.APB0_CLK_EN--
	if mx.APB0_CLK_EN == 0 {
		sc.APB0_CLK_EN().Clear()
	}
	mx.CLK_EN_CENT.Unlock()
}

type Conf1 uint8

const (
	Word5b Conf1 = 0 << 0 // 5-bit data word
	Word6b Conf1 = 1 << 0 // 6-bit data word
	Word7b Conf1 = 2 << 0 // 7-bit data word
	Word8b Conf1 = 3 << 0 // 8-bit data word

	Stop2b Conf1 = 1 << 2 // 2 stop bits for 6 to 8-bit word, 1.5 for 5-bit word

	ParOdd  Conf1 = 1 << 3 // parity control enabled: odd.
	ParEven Conf1 = 3 << 3 // parity control enabled: even

	Break Conf1 = 1 << 6 // break control bit

	dla = 1 << 7 // divisor latch access bit
)

func (p *Periph) Conf1() Conf1 {
	return Conf1(p.lcr.LoadBits(uint32(Word8b + Stop2b + ParEven + Break)))
}

func (p *Periph) SetConf1(c Conf1) {
	p.lcr.Store(uint32(c))
}

type Conf2 uint8

const (
	DTR  Conf2 = 1 << 0 // directly control of DTR output
	RTS  Conf2 = 1 << 1 // directly control of RTS output
	LB   Conf2 = 1 << 4 // put the UART into loop-back diagnostic mode
	AFCE Conf2 = 1 << 5 // auto flow controll enable bit
	SIRE Conf2 = 1 << 6 // SIR mode enable bit
)

func (p *Periph) Conf2() Conf2 {
	return Conf2(p.mcr.Load())
}

func (p *Periph) SetConf2(c Conf2) {
	p.lcr.Store(uint32(c))
}

type Conf3 uint8

const (
	FE     Conf3 = 1 << 0 // enable FIFO mode
	CRF    Conf3 = 1 << 1 // reset and clear Rx FIFO, self clearing bit
	CTF    Conf3 = 1 << 2 // reset and clear Tx FIFO, self clearing bit
	DMAM1  Conf3 = 1 << 3 // dma mode 1
	TFT0   Conf3 = 0 << 4 // empty Tx FIFO interrupt threshold
	TFT2   Conf3 = 1 << 4 // 2 words in Tx FIFO interrupt threshold
	TFT1_4 Conf3 = 2 << 4 // 1/4 Tx FIFO interrupt threshold (4 words)
	TFT1_2 Conf3 = 3 << 4 // 1/2 Tx FIFO interrupt threshold (8 words)
	RFT1   Conf3 = 0 << 6 // 1 word Rx FIFO interrupt threshold
	RFT1_4 Conf3 = 1 << 6 // 1/4 Rx FIFO interrupt threshold (4 words)
	RFT1_2 Conf3 = 2 << 6 // 1/2 Rx FIFO interrupt threshold (8 words)
	RFT_2  Conf3 = 3 << 6 // 2 less than full Rx FIFO interrupt threshold
)

type Conf4 uint8

const (
	PTIME Conf4 = 1 << 7 // programmable Tx intrerrupt mode enable
)

func (p *Periph) SetBaudrate(br int) {
	div := bus.APB0.Clock() / int64(br)
	if uint64(div) >= 1<<20 {
		panic("uart: bad baudrate")
	}
	p.lcr.SetBits(dla)
	p.dlh_ier.Store(uint32(div >> 12))
	p.rbr_dll_thr.Store(uint32(div >> 4 & 0xFF))
	p.far.Store(uint32(div & 0xF))
	p.lcr.ClearBits(dla)
}

type Event uint8

const (
	RxNotEmpty Event = 1 << 0 // at least one received word can be read
	Break      Event = 1 << 4 // break sequence detected
	TxEmpty    Event = 1 << 5 // Tx hold register empty !(FIFO mode && PTIME)
	TxFull     Event = 1 << 5 // Tx FIFO full (FIFO mode && PTIME)
	TxDone     Event = 1 << 6 // transmssion complete (shift register is empty)
)

type Error uint8

const (
	ErrOverrun Error = 1 << 1
	ErrParity  Error = 1 << 2
	ErrFraming Error = 1 << 3
	ErrRxFIFO  Error = 1 << 7
)

func (p *Periph) Status() (Status, Error) {
	lsr := p.lsr.Load()
}

type Int uint8

const (
	RxNot
)

func (p *Periph) Load() int {
	return int(p.rbr_dll_thr.Load())
}

func (p *Periph) Store(d int) {
	return p.rbr_dll_thr.Store(uint32(d))
}
