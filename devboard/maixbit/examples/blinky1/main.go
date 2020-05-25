package main

import (
	"embedded/arch/riscv/systim"
	"time"

	"github.com/embeddedgo/kendryte/p/fpioa"
	"github.com/embeddedgo/kendryte/p/gpio"
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

	sctl.APB0_CLK_EN().Set()
	sctl.CLK_EN_PERI.Store(sysctl.FPIOA_CLK_EN | sysctl.GPIO_CLK_EN)

	const (
		FUNC_GPIO1 = 57
		FUNC_GPIO2 = 58
		FUNC_GPIO3 = 59
		FUNC_GPIO4 = 60
	)

	green := &fpioa.FPIOA().IO[12]
	green.Store(FUNC_GPIO1<<fpioa.CH_SELn | 15<<fpioa.DSn | fpioa.OE_EN)
	red := &fpioa.FPIOA().IO[13]
	red.Store(FUNC_GPIO2<<fpioa.CH_SELn | 15<<fpioa.DSn | fpioa.OE_EN)
	blue := &fpioa.FPIOA().IO[14]
	blue.Store(FUNC_GPIO3<<fpioa.CH_SELn | 15<<fpioa.DSn | fpioa.OE_EN)

	pin6 := &fpioa.FPIOA().IO[6]
	pin6.Store(FUNC_GPIO4<<fpioa.CH_SELn | 15<<fpioa.DSn | fpioa.OE_EN)

	GPIO := gpio.GPIO()
	GPIO.DATA_OUTPUT.SetBits(1<<1 + 1<<2 + 1<<3 + 1<<4)
	GPIO.DIRECTION.SetBits(1<<1 + 1<<2 + 1<<3 + 1<<4)

	for {
		GPIO.DATA_OUTPUT.SetBits(1 << 3)
		GPIO.DATA_OUTPUT.ClearBits(1 << 1)
		time.Sleep(time.Second)

		GPIO.DATA_OUTPUT.SetBits(1 << 1)
		GPIO.DATA_OUTPUT.ClearBits(1 << 2)
		time.Sleep(time.Second)

		GPIO.DATA_OUTPUT.SetBits(1 << 2)
		GPIO.DATA_OUTPUT.ClearBits(1 << 3)
		time.Sleep(time.Second)

		GPIO.DATA_OUTPUT.SetBits(1 << 4)
		time.Sleep(time.Second)
		GPIO.DATA_OUTPUT.ClearBits(1 << 4)
	}
}
