// Copyright 2020 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package setup

// According to the information in the U-Boot source code the K210 PLL seems to
// be True Circuits, Inc. General-Purpose PLL.
//
// The K210 has three PLLs. Every one has three factors: r, f, od. The PLL
// output frequency is the VCO (voltage controlled oscilator) frequency divided
// by od:
//
//	fout = vco / od
//
// The VCO frequency is controlled by the phase decector that compares the
// reference clock divided by r and the VCO clock divided f. The phase
// detector controls the VCO to achieve:
//
//	fref / r = vco / f
//
// This results in the following relationship between the reference and the
// output frequency:
//
//	fout = fref * f / (r * od)
