package main

import (
	"github.com/embeddedgo/kendryte/p/fpioa"
	"github.com/embeddedgo/kendryte/p/gpio"
	"github.com/embeddedgo/kendryte/p/sysctl"
	"github.com/embeddedgo/kendryte/p/uart"
)

var d1, d2 uint64

func delay1() {
	end := d1 + 1e7
	for d1 != end {
		d1++
	}
}

func delay2() {
	end := d2 + 3e7
	for d2 != end {
		d2++
	}
}

func main() {
	sctl := sysctl.SYSCTL()
	sctl.APB0_CLK_EN().Set()
	sctl.CLK_EN_PERI.Store(sysctl.FPIOA_CLK_EN | sysctl.GPIO_CLK_EN)

	const (
		FUNC_GPIO1 = 57
		FUNC_GPIO2 = 58
		FUNC_GPIO3 = 59
	)

	green := &fpioa.FPIOA().IO[12]
	green.Store(FUNC_GPIO1<<fpioa.CH_SELn | 15<<fpioa.DSn | fpioa.OE_EN |
		fpioa.PD)
	red := &fpioa.FPIOA().IO[13]
	red.Store(FUNC_GPIO2<<fpioa.CH_SELn | 15<<fpioa.DSn | fpioa.OE_EN |
		fpioa.PD)
	blue := &fpioa.FPIOA().IO[14]
	blue.Store(FUNC_GPIO3<<fpioa.CH_SELn | 15<<fpioa.DSn | fpioa.OE_EN |
		fpioa.PD)

	GPIO := gpio.GPIO()
	GPIO.DATA_OUTPUT.SetBits(1<<1 + 1<<2 + 1<<3)
	GPIO.DIRECTION.SetBits(1<<1 + 1<<2 + 1<<3)
	for {
		GPIO.DATA_OUTPUT.SetBits(1 << 3)
		GPIO.DATA_OUTPUT.ClearBits(1 << 1)
		delay2()
		GPIO.DATA_OUTPUT.SetBits(1 << 1)
		GPIO.DATA_OUTPUT.ClearBits(1 << 2)
		delay2()
		GPIO.DATA_OUTPUT.SetBits(1 << 2)
		GPIO.DATA_OUTPUT.ClearBits(1 << 3)
		delay2()
		uartPutc('A')
		uartPutc('\r')
		uartPutc('\n')
	}
}

func uartPutc(c byte) {
	u := uart.UART3()
	for u.LSR.LoadBits(1<<5) != 0 {
	}
	u.RBR_DLL_THR.Store(uart.RBR_DLL_THR(c))
}

func uartInit() {
	SYSCTL := sysctl.SYSCTL()
	SYSCTL.APB0_CLK_EN().Set()
	SYSCTL.UART3_CLK_EN().Set()

	u := uart.UART3()
	u.LCR.SetBits(1 << 7)
	//...
}
