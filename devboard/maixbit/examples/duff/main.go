// Copyright 2020 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"time"

	"github.com/embeddedgo/kendryte/hal/fpioa"
	"github.com/embeddedgo/kendryte/hal/uart"

	_ "github.com/embeddedgo/kendryte/devboard/maixbit/board/init"
)

const n = 64

var a, b [n]uint64

func main() {
	t1 := time.Now()
	for i := 0; i < 1000*128/n; i++ {
		a = b
		b = a
		a = b
		b = a
		a = b
		b = a
		a = b
		b = a
	}
	t2 := time.Now()
	fmt.Fprintln(dbg, t2.Sub(t1), "\r")
}

// Results (K210 416 MHz):
//
// n=4:   duff=14.564423ms, loop=18.489062ms
// n=64:  duff=10.153846ms, loop=15.215384ms
// n=128: duff=10.011779ms, loop=15.026923ms
//
// loop means ssaConfig.noDuffDevice=true

var dbg *Serial

func init() {
	rx := fpioa.Pin(4)
	tx := fpioa.Pin(5)
	rx.Setup(fpioa.UART3_RX | fpioa.EnIE | fpioa.Schmitt)
	tx.Setup(fpioa.UART3_TX | fpioa.DriveH34L23 | fpioa.EnOE)

	u := uart.UART(3)
	u.EnableClock()
	u.Reset()
	u.SetConf1(uart.Word8b)
	u.SetConf2(0)
	u.SetConf3(uart.FE | uart.CRF | uart.CTF | uart.TFT8 | uart.RFT1)
	u.SetConf4(uart.PTIME)
	u.SetBaudrate(115200)

	dbg = &Serial{u}
}

type Serial struct {
	p *uart.Periph
}

func (s *Serial) WriteByte(b byte) {
	for {
		if ev, _ := s.p.Status(); ev&uart.TxFull == 0 {
			break
		}
	}
	s.p.Store(int(b))
}

func (s *Serial) Write(p []byte) (int, error) {
	for _, b := range p {
		s.WriteByte(b)
	}
	return len(p), nil
}

func (s *Serial) WriteString(p string) (int, error) {
	for i := 0; i < len(p); i++ {
		s.WriteByte(p[i])
	}
	return len(p), nil
}
