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
//
//  K210 UART features (CPR reg)
//  ----------------------------
//	FIFO_MODE                16
//	DMA_EXTRA                yes
//	UART_ADD_ENCODED_PARAMS  yes
//	SHADOW                   yes
//	FIFO_STAT                no
//	FIFO_ACCESS              no
//	NEW_FEAT                 yes
//	SIR_LP_MODE              no
//	SIR_MODE                 yes
//	THRE_MODE                yes
//	AFCE_MODE                no
//	APB_DATA_WIDTH           32

// Periph represents UART peripheral.
type Periph struct {
	rbr_dll_thr mmio.U32
	dlh_ier     mmio.U32
	fcr_iir     mmio.U32
	lcr         mmio.U32
	mcr         mmio.U32 // AFCE bit not implemented (AFCE_MODE=no)
	lsr         mmio.U32
	msr         mmio.U32
	scr         mmio.U32
	lpdll       mmio.U32 // not implemented (SIR_LP_MODE=no)
	lpdlh       mmio.U32 // not implemented (SIR_LP_MODE=no)
	_           [2]uint32
	srbr_sthr   [16]mmio.U32
	far         mmio.U32 // not implemented (FIFO_ACCESS=no)
	tfr         mmio.U32 // not implemented (FIFO_ACCESS=no)
	rfw         mmio.U32 // not implemented (FIFO_ACCESS=no)
	usr         mmio.U32 // only BUSY bit implemented (FIFO_STAT=no)
	tfl         mmio.U32 // not implemented (FIFO_STAT=no)
	rfl         mmio.U32 // not implemented (FIFO_STAT=no)
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
	PE  LineConf = 1 << 3 // parity enable
	EPS LineConf = 3 << 3 // even parity select
	BC  LineConf = 1 << 6 // break control bit
	dla          = 1 << 7 // divisor latch access bit
)

func (p *Periph) LineConf() LineConf {
	return LineConf(p.lcr.LoadBits(uint32(W8 + S2 + PE + EPS + BC)))
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
	//AFCE ModeConf = 1 << 5 // auto flow controll enable bit
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

func (p *Periph) SoftReset(txfifo, rxfifo, uart bool) {
	var reset uint32
	if txfifo {
		reset = 1 << 2
	}
	if rxfifo {
		reset |= 1 << 1
	}
	if uart {
		reset |= 1 << 0
	}
	p.srr.Store(reset)
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
	p.stet.Store(uint32(tft&TFT8) >> 4)
}

func (p *Periph) RFT() FIFOConf {
	return FIFOConf(p.stet.Load() << 6)
}

func (p *Periph) SetRFT(rft FIFOConf) {
	p.srt.Store(uint32(rft&RFT14) >> 6)
}

type Int uint8

// Interrupts listed in decreasing priority. The clearing method is given at the
// end of each event description (in parentheses).
const (
	RxStatus    Int = 6 // overrun/parity/framing err, break (call LineStatus)
	RxReady     Int = 4 // Rx data available (read below the trigger level)
	TxReady     Int = 2 // Tx ready for next data (cleared by Event method)
	ModemStatus Int = 0 // modem/flow control signal (call ModemStatus)
	BusyFault   Int = 7 // SetLineConf while the UART is busy (call Status1)
	None        Int = 1 // no interrupt
)

// Event returns the highest priority enabled event. If rxto is true the RxReady
// event was generated because there was no in or out activity on non-empty FIFO
// (below Rx trigger level) for 4 word period.
func (p *Periph) Int() (i Int, rxto bool) {
	iir := p.fcr_iir.Load()
	return Int(iir & 7), iir&15 == 12
}

// IntConf
type IntConf uint8

const (
	RxStatusEn    IntConf = 1 << 2 // enable RxStatus event
	RxReadyEn     IntConf = 1 << 0 // enbale RxReady event
	TxReadyEn     IntConf = 1 << 1 // enable TxReady event
	ModemStatusEn IntConf = 1 << 3 // enable ModemStatus event
	PTIME         IntConf = 1 << 7 // programmable Tx intrerrupt mode enable
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

type Status uint8

const (
	RxNotEmpty Status = 1 << 0 // at least one received word can be read
	LINBreak   Status = 1 << 4 // break sequence detected
	TxEmpty    Status = 1 << 5 // Tx hold register empty (use if !FE || !PTIME)
	TxFull     Status = 1 << 5 // Tx FIFO full (use if FE && PTIME)
	TxDone     Status = 1 << 6 // Tx done (shift register is empty)
)

type Error uint8

const (
	ErrOverrun Error = 1 << 1 // data lost because of no free space in hardware
	ErrParity  Error = 1 << 2 // parity error detected
	ErrFraming Error = 1 << 3 // no valid stop bit detected in the receive data
	ErrRxFIFO  Error = 1 << 7 // FIFOmode && (ErrParity||ErrFraming||LINBreak)

	ErrAll = ErrOverrun | ErrParity | ErrFraming | ErrRxFIFO
)

// Status returns the line status bits. It clears LINBreak event and all errors.
func (p *Periph) Status() (Status, Error) {
	lsr := p.lsr.Load()
	return Status(lsr) &^ Status(ErrAll), Error(lsr) & ErrAll
}

type Status1 uint8

const (
	Busy           Status1 = 1 << 0
	//TxFIFONotFull  Status1 = 1 << 1
	//TxFIFOEmpty    Status1 = 1 << 2
	//RxFIFONotEmpty Status1 = 1 << 3
	//RxFIFOFull     Status1 = 1 << 4
)

// Status1 returns the UART status bits.
func (p *Periph) Status1() Status1 {
	return Status1(p.usr.Load())
}

func (p *Periph) Load() int {
	return int(p.rbr_dll_thr.Load())
}

func (p *Periph) Store(d int) {
	p.rbr_dll_thr.Store(uint32(d))
}

func (p *Periph) CPR() uint32 {
	return p.cpr.Load()
}