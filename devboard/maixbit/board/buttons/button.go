// Copyright 2020 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package buttons

import (
	"github.com/embeddedgo/kendryte/hal/fpioa"

	_ "github.com/embeddedgo/kendryte/devboard/maixbit/board/init"
)

// Onboard buttons
const (
	BOOT Button = 16 // PA0

	User = BOOT
)

type Button uint8

func (b Button) Read() int      { return fpioa.Pin(b).Load() ^ 1 }
func (b Button) Pin() fpioa.Pin { return fpioa.Pin(b) }
