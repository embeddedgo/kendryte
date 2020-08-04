// Copyright 2020 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package init

import (
	"embedded/arch/riscv/systim"
	"embedded/rtos"
	"runtime"

	"github.com/embeddedgo/kendryte/hal/fpioa"
	"github.com/embeddedgo/kendryte/hal/uart"
	"github.com/embeddedgo/kendryte/p/bus"
	"github.com/embeddedgo/kendryte/p/sysctl"
)

func init() {
	sctl := sysctl.SYSCTL()
	cpuHz := int64(26e6)
	if sctl.ACLK_SEL().Load() != 0 {
		pll := sctl.PLL[0].Load()
		r := uint64(pll&sysctl.CLKR)>>sysctl.CLKRn + 1
		f := uint64(pll&sysctl.CLKF)>>sysctl.CLKFn + 1
		od := uint64(pll&sysctl.CLKOD)>>sysctl.CLKODn + 1
		pllHz := 26e6 * f / (r * od)
		aclkDivSel := sctl.ACLK_DIVIDER_SEL().Load() >> sysctl.ACLK_DIVIDER_SELn
		cpuHz = int64(pllHz / 2 << aclkDivSel)
	}
	bus.Core.SetClock(cpuHz)
	bus.TileLink.SetClock(cpuHz)
	bus.AXI.SetClock(cpuHz)
	bus.AHB.SetClock(cpuHz)
	clksel0 := sctl.CLK_SEL0.Load()
	div := clksel0&(sysctl.APB0_CLK_SEL>>sysctl.APB0_CLK_SELn) + 1
	bus.APB0.SetClock(cpuHz / int64(div))
	div = clksel0&(sysctl.APB1_CLK_SEL>>sysctl.APB1_CLK_SELn) + 1
	bus.APB1.SetClock(cpuHz / int64(div))
	div = clksel0&(sysctl.APB2_CLK_SEL>>sysctl.APB2_CLK_SELn) + 1
	bus.APB2.SetClock(cpuHz / int64(div))
	systim.Setup(cpuHz / 50)

	setupSystemWriter()
	runtime.GOMAXPROCS(2)
}

const (
	dbgPin  = 5
	dbgUART = 3
	dbgFunc = fpioa.UART3_TX
)

func setupSystemWriter() {
	tx := fpioa.Pin(dbgPin)
	tx.Setup(dbgFunc | fpioa.DriveH34L23 | fpioa.EnOE)

	u := uart.UART(dbgUART)
	u.EnableClock()
	u.Reset()
	u.SetLineConf(uart.W8)
	u.SetFIFOConf(uart.FE | uart.CRF | uart.CTF | uart.TFT8 | uart.RFT1)
	u.SetIntConf(uart.PTIME)
	u.SetBaudrate(115200)

	rtos.SetSystemWriter(write)
}

func write(_ int, p []byte) int {
	for _, b := range p {
		if b == '\n' {
			writeByte('\r')
		}
		writeByte(b)
	}
	return len(p)
}

func writeByte(b byte) {
	u := uart.UART(dbgUART)
	for {
		if ev, _ := u.Status(); ev&uart.TxFull == 0 {
			u.Store(int(b))
			return
		}
	}
}
