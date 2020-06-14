// Copyright 2020 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package init

import (
	"embedded/arch/riscv/systim"
	"runtime"

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
	runtime.GOMAXPROCS(2)
}
