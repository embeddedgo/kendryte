// Copyright 2020 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"embedded/rtos"
	"fmt"
	"time"

	"github.com/embeddedgo/kendryte/hal/fpioa"
	"github.com/embeddedgo/kendryte/hal/irq"
	"github.com/embeddedgo/kendryte/hal/uarths"

	_ "github.com/embeddedgo/kendryte/devboard/maixbit/board/init"
)

var u *uarths.Driver

func main() {
	rx := fpioa.Pin(4)
	tx := fpioa.Pin(5)
	rx.Setup(fpioa.UARTHS_RX | fpioa.EnIE | fpioa.Schmitt)
	tx.Setup(fpioa.UARTHS_TX | fpioa.DriveH34L23 | fpioa.EnOE)

	u = uarths.NewDriver(uarths.UARTHS(1))
	u.SetBaudrate(750e3)
	u.EnableTx()
	u.EnableRx(nil)

	irq.UARTHS.Enable(rtos.IntPrioLow, irq.M0)

	n := 40
	s := "00000000001111111111222222222233333333334444444444" +
		"555555555566666666667777777777\r\n"
	br := u.Periph().Baudrate()
	for k := 0; k < 0; k++ {
		t := time.Now()
		for i := 0; i < n; i++ {
			u.WriteString(s)
		}
		dt := int(time.Now().Sub(t))
		lps := (n*1e9 + dt/2) / dt
		bps := (n*len(s)*1e9 + dt/2) / dt
		fmt.Fprintf(u, "br: %d b/s (%d B/s),  speed: %d line/s (%d B/s)\r\n",
			br, br/8, lps, bps)
		time.Sleep(2 * time.Second)
	}

	buf := make([]byte, 128)
	for {
		n, err := u.Read(buf)
		if n != 0 {
			u.WriteString("inp: ")
			u.Write(buf[:n])
			u.WriteString("\r\n")
		}
		if err != nil {
			u.WriteString("err: " + err.Error() + "\r\n")
		}
	}

}

//go:interrupthandler
func UARTHS_Handler() {
	u.ISR()
}
