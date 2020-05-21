package main

import (
	"embedded/mmio"
	"unsafe"

	"github.com/embeddedgo/kendryte/p/fpioa"
	"github.com/embeddedgo/kendryte/p/gpio"
	"github.com/embeddedgo/kendryte/p/sysctl"
	"github.com/embeddedgo/kendryte/p/uart"
)

func MulDiv(x, m, d uint64) uint64 {
	divx := x / d
	modx := x - divx*d
	divm := m / d
	modm := m - divm*d
	return divx*m + modx*divm + modx*modm/d
}

func nanotime() int64 {
	sc := sysctl.SYSCTL()
	cpuHz := uint64(26e6)
	if sc.ACLK_SEL().Load() != 0 {
		pll := sc.PLL[0].Load()
		r := uint64(pll&sysctl.CLKR)>>sysctl.CLKRn + 1
		f := uint64(pll&sysctl.CLKF)>>sysctl.CLKFn + 1
		od := uint64(pll&sysctl.CLKOD)>>sysctl.CLKODn + 1
		pllHz := 26e6 * f / (r * od)
		aclkDivSel := sc.ACLK_DIVIDER_SEL().Load() >> sysctl.ACLK_DIVIDER_SELn
		cpuHz = pllHz / 2 << aclkDivSel
	}
	clintHz := cpuHz / 50
	mtime := (*mmio.U64)(unsafe.Pointer(uintptr(0x2000000 + 0xBFF8))).Load()
	return int64(MulDiv(mtime, 1e9, clintHz))
}

func waitUntil(t int64) {
	for nanotime() < t {
	}
}

var lastt uint64

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

	t := nanotime()
	for {
		GPIO.DATA_OUTPUT.SetBits(1 << 3)
		GPIO.DATA_OUTPUT.ClearBits(1 << 1)
		t += 1e9
		waitUntil(t)
		GPIO.DATA_OUTPUT.SetBits(1 << 1)
		GPIO.DATA_OUTPUT.ClearBits(1 << 2)
		t += 1e9
		waitUntil(t)
		GPIO.DATA_OUTPUT.SetBits(1 << 2)
		GPIO.DATA_OUTPUT.ClearBits(1 << 3)
		t += 1e9
		waitUntil(t)
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
