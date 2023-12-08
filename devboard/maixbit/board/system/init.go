// Copyright 2020 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package system

import (
	"embedded/arch/riscv/systim"
	"embedded/rtos"
	"os"
	"runtime"
	"syscall"
	_ "unsafe"

	"github.com/embeddedgo/fs/termfs"
	"github.com/embeddedgo/kendryte/hal/fpioa"
	"github.com/embeddedgo/kendryte/hal/uart"
	"github.com/embeddedgo/kendryte/hal/uart/uart3"
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

	setupConsole()
	runtime.GOMAXPROCS(2)
}

func setupConsole() {
	u := uart3.Driver()
	u.UsePin(fpioa.Pin(4), uart.RXD)
	u.UsePin(fpioa.Pin(5), uart.TXD)
	u.Setup(uart.Word8b, 115200)
	rtos.SetSystemWriter(write)
	u.EnableRx(64)

	// Setup a serial console (standard input and output).
	con := termfs.New("UART3", u, u)
	con.SetCharMap(termfs.InCRLF | termfs.OutLFCRLF)
	con.SetEcho(true)
	con.SetLineMode(true, 256)
	rtos.Mount(con, "/dev/console")
	var err error
	os.Stdin, err = os.OpenFile("/dev/console", syscall.O_RDONLY, 0)
	checkErr(err)
	os.Stdout, err = os.OpenFile("/dev/console", syscall.O_WRONLY, 0)
	checkErr(err)
	os.Stderr = os.Stdout
}

func checkErr(err error) {
	if err != nil {
		println(err.Error())
		os.Exit(1)
	}
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
	u := uart3.Driver()
	for {
		if ev, _ := u.Periph().Status(); ev&uart.TxFull == 0 {
			u.Periph().Store(int(b))
			return
		}
	}
}
