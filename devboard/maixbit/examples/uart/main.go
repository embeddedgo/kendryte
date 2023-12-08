// Copyright 2020 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"time"

	"github.com/embeddedgo/kendryte/hal/fpioa"
	"github.com/embeddedgo/kendryte/hal/uart"
	"github.com/embeddedgo/kendryte/hal/uart/uart1"

	_ "github.com/embeddedgo/kendryte/devboard/maixbit/board/system"
)

var u *uart.Driver

func main() {
	u := uart1.Driver()
	u.UsePin(fpioa.Pin(4), uart.RXD)
	u.UsePin(fpioa.Pin(5), uart.TXD)
	u.Setup(uart.Word8b, 750e3)
	u.EnableRx(64)

	// Speed test

	u.WriteString("\r\nWrite speed test...\r\n\n")
	time.Sleep(time.Second)

	n := 40
	s := "00000000001111111111222222222233333333334444444444" +
		"555555555566666666667777777777\r\n"
	br := u.Periph().Baudrate()
	for k := 0; k < 2; k++ {
		t := time.Now()
		for i := 0; i < n; i++ {
			u.WriteString(s)
		}
		dt := int64(time.Now().Sub(t))
		bps := (int64(n*len(s))*1e9 + dt/2) / dt
		fmt.Fprintf(u, "baud rate: %d b/s (%d B/s),  speed: %d b/s (%d B/s)\r\n\n",
			br, br/8, bps*8, bps)
		time.Sleep(2 * time.Second)
	}

	// Read test

	s = "<=[+](*)->0123456789abcdefghijklmnoprstuvwxyzABCDEFGHIJKLMNOPRSTUVWXYZ"
	u.WriteString(s)
	u.WriteString(s)
	u.WriteString("\r\n\nRead test. ")
	u.WriteString("Use keyboard or paste some text (eg. the above line).\r\n\n")

	buf := make([]byte, 128)
	for {
		n, err := u.Read(buf)
		if n != 0 {
			fmt.Fprintf(u, "read%3d: %s\r\n", n, buf[:n])
		}
		if err != nil {
			u.WriteString("err: " + err.Error() + "\r\n")
		}
	}
}
