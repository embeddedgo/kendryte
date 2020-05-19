package main

import (
	"github.com/embeddedgo/kendryte/p/fpioa"
	"github.com/embeddedgo/kendryte/p/gpio"
	"github.com/embeddedgo/kendryte/p/sysctl"
)

var d uint64

func delay() {
	stop := d + 1e7
	for d != stop {
		d++
	}
}

func main() {
	SYSCTL := sysctl.SYSCTL()
	SYSCTL.APB0_CLK_EN().Set()
	SYSCTL.CLK_EN_PERI.Store(sysctl.FPIOA_CLK_EN | sysctl.GPIO_CLK_EN)

	const FUNC_GPIO3 = 59

	redPin := &fpioa.FPIOA().IO[13]
	redPin.Store(FUNC_GPIO3<<fpioa.CH_SELn | 15<<fpioa.DSn | fpioa.OE_EN |
		fpioa.PD)

	GPIO := gpio.GPIO()
	GPIO.DIRECTION.SetBits(1 << 3)
	GPIO.DATA_OUTPUT.SetBits(1 << 3)
	for {
		delay()
		GPIO.DATA_OUTPUT.ClearBits(1 << 3)
		delay()
		GPIO.DATA_OUTPUT.SetBits(1 << 3)
	}
}
