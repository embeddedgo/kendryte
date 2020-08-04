// Copyright 2020 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package uart

import (
	"embedded/mmio"
	"time"
	"unsafe"

	"github.com/embeddedgo/kendryte/hal/internal"
	"github.com/embeddedgo/kendryte/p/bus"
	"github.com/embeddedgo/kendryte/p/mmap"
	"github.com/embeddedgo/kendryte/p/sysctl"
)

// Synopsys DW_apb_uart

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

func (p *Periph) Reset() {
	sc := sysctl.SYSCTL()
	mx := &internal.MX.SYSCTL

	mx.PERI_RESET.Lock()
	sc.PERI_RESET.SetBits(sysctl.UART1_RESET << p.n())
	mx.PERI_RESET.Unlock()

	time.Sleep(10 * time.Microsecond)

	mx.PERI_RESET.Lock()
	sc.PERI_RESET.ClearBits(sysctl.UART1_RESET << p.n())
	mx.PERI_RESET.Unlock()
}

// LineConf (line configuration) represents a serial data word chracteristics
type LineConf uint8

const (
	W5 LineConf = 0 << 0 // 5-bit data word
	W6 LineConf = 1 << 0 // 6-bit data word
	W7 LineConf = 2 << 0 // 7-bit data word
	W8 LineConf = 3 << 0 // 8-bit data word

	S2  LineConf = 1 << 2 // 2 stop bits for 6 to 8-bit word, 1.5 for 5-bit word
	PEN LineConf = 1 << 3 // parity enable
	EPS LineConf = 3 << 3 // even parity select
	BC  LineConf = 1 << 6 // break control bit
	dla          = 1 << 7 // divisor latch access bit
)

func (p *Periph) LineConf() LineConf {
	return Conf1(p.lcr.LoadBits(uint32(W8b + S2b + Even + Break)))
}

func (p *Periph) SetLineConf(c LineConf) {
	p.lcr.Store(uint32(c))
}

// ModeConf (mode configuration) represents UART mode of operation
type ModeConf uint8

const (
	DTR  ModeConf = 1 << 0 // directly control of DTR output
	RTS  ModeConf = 1 << 1 // directly control of RTS output
	LB   ModeConf = 1 << 4 // put the UART into loop-back diagnostic mode
	AFCE ModeConf = 1 << 5 // auto flow controll enable bit
	SIRE ModeConf = 1 << 6 // IrDA SIR (serial infrared) mode enable bit
)

func (p *Periph) ModeConf() ModeConf {
	return ModeConf(p.mcr.Load())
}

func (p *Periph) SetModeConf(c ModeConf) {
	p.mcr.Store(uint32(c))
}

// FIFOConf (FIFO configuration).
type FIFOConf uint8

const (
	FE    FIFOConf = 1 << 0 // enable FIFO mode
	CRF   FIFOConf = 1 << 1 // reset and clear Rx FIFO, self clearing bit
	CTF   FIFOConf = 1 << 2 // reset and clear Tx FIFO, self clearing bit
	DMAM1 FIFOConf = 1 << 3 // dma mode 1
	TFT0  FIFOConf = 0 << 4 // empty Tx FIFO interrupt threshold
	TFT2  FIFOConf = 1 << 4 // 2 words in Tx FIFO interrupt threshold
	TFT4  FIFOConf = 2 << 4 // 1/4 Tx FIFO interrupt threshold (4 words)
	TFT8  FIFOConf = 3 << 4 // 1/2 Tx FIFO interrupt threshold (8 words)
	RFT1  FIFOConf = 0 << 6 // 1 word Rx FIFO interrupt threshold
	RFT4  FIFOConf = 1 << 6 // 1/4 Rx FIFO interrupt threshold (4 words)
	RFT8  FIFOConf = 2 << 6 // 1/2 Rx FIFO interrupt threshold (8 words)
	RFT14 FIFOConf = 3 << 6 // 2 less than full Rx FIFO interrupt threshold
)

func (p *Periph) SetFIFOConf(c FIFOConf) {
	p.fcr_iir.Store(uint32(c))
}

func (p *Periph) FE() FIFOConf {
	return FIFOConf(p.sfe.Load()) & FE
}

func (p *Periph) SetFE(fe FIFOConf) {
	p.sfe.Store(uint32(fe & FE))
}

func (p *Periph) TFT() FIFOConf {
	return FIFOConf(p.srt.Load()<<4) & TFT8
}

func (p *Periph) SetTFT(tft FIFOConf) {
	p.srt.Store(uint32(tft&TFT8) >> 4)
}

func (p *Periph) RxFIFOTrigger() FIFOConf {
	return FIFOConf(p.srt.Load() << 6)
}

func (p *Periph) SetRxFIFOTrigger(rft FIFOConf) {
	p.srt.Store(uint32(srt&RFT14) >> 6)
}

type IRQ uint8

const (
	DSSI IRQ = 0
	NOI  IRQ = 1
	TBEI IRQ = 2
	RBFI IRQ = 4
	LSI  IRQ = 6
	BI   IRQ = 7  // busy interrupt
	CTOI IRQ = 12 // character timout interrupt
)

func (p *Periph) IRQ() IRQ {

}

// IntConf (interrupt mode configuration)
type IntConf uint8

const (
	ERBFI IntConf = 1 << 0 // enbale Rx data available interrupt
	ETBEI IntConf = 1 << 1 // enable Tx holding register empty interrupt
	ELSI  IntConf = 1 << 2 // enable receiver line status interrupt
	EDSSI IntConf = 1 << 3 // enable modem status interrupt
	PTIME IntConf = 1 << 7 // programmable Tx intrerrupt mode enable
)

func (p *Periph) IntConf() IntConf {
	return IntConf(p.dlh_ier.Load())
}

func (p *Periph) SetIntConf(c IntConf) {
	p.dlh_ier.Store(uint32(c))
}

func (p *Periph) SetBaudrate(br int) {
	div := (p.Bus().Clock() + int64(br)/2) / int64(br)
	if uint64(div) >= 1<<20 {
		panic("uart: bad baudrate")
	}
	p.lcr.SetBits(dla)
	p.dlh_ier.Store(uint32(div >> 12))
	p.rbr_dll_thr.Store(uint32(div >> 4 & 0xFF))
	p.dlf.Store(uint32(div & 0xF))
	p.lcr.ClearBits(dla)
}

type Event uint8

const (
	RxNotEmpty Event = 1 << 0 // at least one received word can be read
	LINBreak   Event = 1 << 4 // break sequence detected
	TxEmpty    Event = 1 << 5 // Tx hold register empty (!FIFOmode || !PTIME)
	TxFull     Event = 1 << 5 // Tx FIFO full (FIFOmode && PTIME)
	TxDone     Event = 1 << 6 // Tx done (shift register is empty)
)

type Error uint8

const (
	ErrOverrun Error = 1 << 1 // data lost because of no free space in hardware
	ErrParity  Error = 1 << 2 // parity error detected
	ErrFraming Error = 1 << 3 // no valid stop bit detected in the receive data
	ErrRxFIFO  Error = 1 << 7 // FIFOmode && (ErrParity||ErrFraming||LINBreak)
	ErrAll           = ErrOverrun | ErrParity | ErrFraming | ErrRxFIFO
)

// Status returns active status bits. It clears LINBreak event and all errors.
func (p *Periph) Status() (Event, Error) {
	lsr := p.lsr.Load()
	return Event(lsr) &^ Event(ErrAll), Error(lsr) & ErrAll
}

func (p *Periph) Load() int {
	return int(p.rbr_dll_thr.Load())
}

func (p *Periph) Store(d int) {
	p.rbr_dll_thr.Store(uint32(d))
}
