package main

import (
	"embedded/arch/riscv/systim"
	"time"

	"github.com/embeddedgo/kendryte/hal/fpioa"
	"github.com/embeddedgo/kendryte/p/sysctl"
)

func main() {
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
	systim.Setup(cpuHz / 50)

	fpioa.EnableClock()

	green := fpioa.Pin(12)
	red := fpioa.Pin(13)
	blue := fpioa.Pin(14)
	pin6 := fpioa.Pin(6)

	high := fpioa.FuncConst | fpioa.Drive | fpioa.OutEn
	low := fpioa.FuncConst | fpioa.Drive | fpioa.OutEn | fpioa.OutEnInv

	for {
		green.Setup(low)
		time.Sleep(100 * time.Millisecond)
		green.Setup(high)
		time.Sleep(time.Second)

		red.Setup(low)
		time.Sleep(100 * time.Millisecond)
		red.Setup(high)
		time.Sleep(time.Second)

		blue.Setup(low)
		time.Sleep(100 * time.Millisecond)
		blue.Setup(high)
		time.Sleep(time.Second)

		pin6.Setup(low)
		time.Sleep(1000 * time.Millisecond)
		pin6.Setup(high)
		time.Sleep(time.Second)
	}
}
