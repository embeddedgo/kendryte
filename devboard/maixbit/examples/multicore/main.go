// Copyright 2020 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	_ "github.com/embeddedgo/kendryte/devboard/maixbit/board/init"
)

func main() {
	for {
		print("Hello!\r\n")
	}
}

func hartid() int
