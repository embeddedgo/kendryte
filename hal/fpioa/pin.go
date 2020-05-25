// Copyright 2020 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fpioa

import (
	"embedded/mmio"
	"time"
	"unsafe"

	"github.com/embeddedgo/kendryte/hal/internal"
	"github.com/embeddedgo/kendryte/p/mmap"
	"github.com/embeddedgo/kendryte/p/sysctl"
)

func EnableClock() {
	sc := sysctl.SYSCTL()
	mx := &internal.MX.SYSCTL

	mx.CLK_EN_CENT.Lock()
	if mx.APB0_CLK_EN == 0 {
		sc.APB0_CLK_EN().Set()
	}
	mx.APB0_CLK_EN++
	mx.CLK_EN_CENT.Unlock()

	mx.CLK_EN_PERI.Lock()
	sc.FPIOA_CLK_EN().Set()
	mx.CLK_EN_PERI.Unlock()
}

func DisableClock() {
	sc := sysctl.SYSCTL()
	mx := &internal.MX.SYSCTL

	mx.CLK_EN_PERI.Lock()
	sc.FPIOA_CLK_EN().Clear()
	mx.CLK_EN_PERI.Unlock()

	mx.CLK_EN_CENT.Lock()
	mx.APB0_CLK_EN--
	if mx.APB0_CLK_EN == 0 {
		sc.APB0_CLK_EN().Clear()
	}
	mx.CLK_EN_CENT.Unlock()
}

func Reset() {
	sc := sysctl.SYSCTL()
	mx := &internal.MX.SYSCTL

	mx.PERI_RESET.Lock()
	sc.FPIOA_RESET().Set()
	mx.PERI_RESET.Unlock()

	time.Sleep(10 * time.Microsecond)

	mx.PERI_RESET.Lock()
	sc.FPIOA_RESET().Clear()
	mx.PERI_RESET.Unlock()
}

type periph struct {
	io     [48]mmio.U32
	tieEn  [8]mmio.U32
	tieVal [8]mmio.U32
}

func p() *periph {
	return (*periph)(unsafe.Pointer(mmap.FPIOA_BASE))
}

// Pin represents a FPIOA pin.
type Pin uint32

type Config uint32

const (
	Func Config = 0xFF << 0 // Function number (see function list below).
	None        = RESV0     // No function

	Drive       Config = 0xF << 8 // Driving strength [mA]:
	DriveH8L5   Config = 0 << 8   // src:  5.0-11.2, sink:  3.2- 8.3
	DriveH11L8  Config = 1 << 8   // src:  7.5-16.8, sink:  4.7-12.3
	DriveH15L11 Config = 2 << 8   // src: 10.0-22.3, sink:  6.3-16.4
	DriveH19L13 Config = 3 << 8   // src: 12.4-27.8, sink:  7.8-20.2
	DriveH23L16 Config = 4 << 8   // src: 14.9-33.3, sink:  9.4-24.2
	DriveH26L18 Config = 5 << 8   // src: 17.4-38.7, sink: 10.9-28.1
	DriveH30L21 Config = 6 << 8   // src: 19.8-44.1, sink: 12.4-31.8
	DriveH34L23 Config = 7 << 8   // src: 22.3-49.5, sink: 13.9-35.5

	EnOE  Config = 1 << 12 // Enable peripheral OE (output enable) signal
	InvOE Config = 1 << 13 // Invert peripheral OE signal (before AND with EnOE)
	OutOE Config = 1 << 14 // Output OE signal instead of data signal
	InvDO Config = 1 << 15 // Invert the output signal

	PullUp   Config = 1 << 16 // Enable internal pull-up
	PullDown Config = 1 << 17 // Enable internal pull-down
	SlewRate Config = 1 << 19 // Reduce the output signal slew rate

	EnIE  Config = 1 << 20 // Enable peripheral IE (input enable) signal
	InvIE Config = 1 << 21 // Invert peripheral IE signal (before AND with EnIE)
	InvDI Config = 1 << 22 // Invert the input signal

	Schmitt Config = 1 << 23 // Schmitt trigger

	DI Config = 1 << 31 // Current value of input sgnal
)

// Setup configures pin. It can also be used for simple data output (see Set and
// Clear function for more information).
func (pin Pin) Setup(cfg Config) {
	p().io[pin].Store(uint32(cfg))
}

// Config returns the pin configuration and state. It can be used for simple
// data input (see Load for more information).
func (pin Pin) Config() Config {
	return Config(p().io[pin].Load())
}

// Set atomically sets the FPIOA pin to high value. It is the equivalent of
//
//	pin.Setup(CONSTANT | DriveH34L23 | EnOE | InvOE | InvDO)
//
// Use the GPIO or GPIOHS peripherals to control more pins at the same time.
func (pin Pin) Set() {
	p().io[pin].Store(uint32(CONSTANT | DriveH34L23 | EnOE | InvOE | InvDO))
}

// Clear atomically sets the FPIOA pin to low value. It is the equivalent of
//
//	pin.Setup(CONSTANT | DriveH34L23 | EnOE | InvOE)
//
// Use the GPIO or GPIOHS peripherals to control more pins at the same time.
func (pin Pin) Clear() {
	p().io[pin].Store(uint32(CONSTANT | DriveH34L23 | EnOE | InvOE))
}

// Load returns the input value of the FPIOA pin. It is the equivalent of
//
//	int(pin.Config() >> 31)
//
// Use the GPIO or GPIOHS peripherals to read more pins at the same time.
func (pin Pin) Load() int {
	return int(p().io[pin].Load() >> 31)
}

// Function list
const (
	JTAG_TCLK      Config = 0   // JTAG Test Clock
	JTAG_TDI       Config = 1   // JTAG Test Data In
	JTAG_TMS       Config = 2   // JTAG Test Mode Select
	JTAG_TDO       Config = 3   // JTAG Test Data Out
	SPI0_D0        Config = 4   // SPI0 Data 0
	SPI0_D1        Config = 5   // SPI0 Data 1
	SPI0_D2        Config = 6   // SPI0 Data 2
	SPI0_D3        Config = 7   // SPI0 Data 3
	SPI0_D4        Config = 8   // SPI0 Data 4
	SPI0_D5        Config = 9   // SPI0 Data 5
	SPI0_D6        Config = 10  // SPI0 Data 6
	SPI0_D7        Config = 11  // SPI0 Data 7
	SPI0_SS0       Config = 12  // SPI0 Chip Select 0
	SPI0_SS1       Config = 13  // SPI0 Chip Select 1
	SPI0_SS2       Config = 14  // SPI0 Chip Select 2
	SPI0_SS3       Config = 15  // SPI0 Chip Select 3
	SPI0_ARB       Config = 16  // SPI0 Arbitration
	SPI0_SCLK      Config = 17  // SPI0 Serial Clock
	UARTHS_RX      Config = 18  // UART High speed Receiver
	UARTHS_TX      Config = 19  // UART High speed Transmitter
	RESV6          Config = 20  // Reserved function
	RESV7          Config = 21  // Reserved function
	CLK_SPI1       Config = 22  // Clock SPI1
	CLK_I2C1       Config = 23  // Clock I2C1
	GPIOHS0        Config = 24  // GPIO High speed 0
	GPIOHS1        Config = 25  // GPIO High speed 1
	GPIOHS2        Config = 26  // GPIO High speed 2
	GPIOHS3        Config = 27  // GPIO High speed 3
	GPIOHS4        Config = 28  // GPIO High speed 4
	GPIOHS5        Config = 29  // GPIO High speed 5
	GPIOHS6        Config = 30  // GPIO High speed 6
	GPIOHS7        Config = 31  // GPIO High speed 7
	GPIOHS8        Config = 32  // GPIO High speed 8
	GPIOHS9        Config = 33  // GPIO High speed 9
	GPIOHS10       Config = 34  // GPIO High speed 10
	GPIOHS11       Config = 35  // GPIO High speed 11
	GPIOHS12       Config = 36  // GPIO High speed 12
	GPIOHS13       Config = 37  // GPIO High speed 13
	GPIOHS14       Config = 38  // GPIO High speed 14
	GPIOHS15       Config = 39  // GPIO High speed 15
	GPIOHS16       Config = 40  // GPIO High speed 16
	GPIOHS17       Config = 41  // GPIO High speed 17
	GPIOHS18       Config = 42  // GPIO High speed 18
	GPIOHS19       Config = 43  // GPIO High speed 19
	GPIOHS20       Config = 44  // GPIO High speed 20
	GPIOHS21       Config = 45  // GPIO High speed 21
	GPIOHS22       Config = 46  // GPIO High speed 22
	GPIOHS23       Config = 47  // GPIO High speed 23
	GPIOHS24       Config = 48  // GPIO High speed 24
	GPIOHS25       Config = 49  // GPIO High speed 25
	GPIOHS26       Config = 50  // GPIO High speed 26
	GPIOHS27       Config = 51  // GPIO High speed 27
	GPIOHS28       Config = 52  // GPIO High speed 28
	GPIOHS29       Config = 53  // GPIO High speed 29
	GPIOHS30       Config = 54  // GPIO High speed 30
	GPIOHS31       Config = 55  // GPIO High speed 31
	GPIO0          Config = 56  // GPIO pin 0
	GPIO1          Config = 57  // GPIO pin 1
	GPIO2          Config = 58  // GPIO pin 2
	GPIO3          Config = 59  // GPIO pin 3
	GPIO4          Config = 60  // GPIO pin 4
	GPIO5          Config = 61  // GPIO pin 5
	GPIO6          Config = 62  // GPIO pin 6
	GPIO7          Config = 63  // GPIO pin 7
	UART1_RX       Config = 64  // UART1 Receiver
	UART1_TX       Config = 65  // UART1 Transmitter
	UART2_RX       Config = 66  // UART2 Receiver
	UART2_TX       Config = 67  // UART2 Transmitter
	UART3_RX       Config = 68  // UART3 Receiver
	UART3_TX       Config = 69  // UART3 Transmitter
	SPI1_D0        Config = 70  // SPI1 Data 0
	SPI1_D1        Config = 71  // SPI1 Data 1
	SPI1_D2        Config = 72  // SPI1 Data 2
	SPI1_D3        Config = 73  // SPI1 Data 3
	SPI1_D4        Config = 74  // SPI1 Data 4
	SPI1_D5        Config = 75  // SPI1 Data 5
	SPI1_D6        Config = 76  // SPI1 Data 6
	SPI1_D7        Config = 77  // SPI1 Data 7
	SPI1_SS0       Config = 78  // SPI1 Chip Select 0
	SPI1_SS1       Config = 79  // SPI1 Chip Select 1
	SPI1_SS2       Config = 80  // SPI1 Chip Select 2
	SPI1_SS3       Config = 81  // SPI1 Chip Select 3
	SPI1_ARB       Config = 82  // SPI1 Arbitration
	SPI1_SCLK      Config = 83  // SPI1 Serial Clock
	SPI_SLAVE_D0   Config = 84  // SPI Slave Data 0
	SPI_SLAVE_SS   Config = 85  // SPI Slave Select
	SPI_SLAVE_SCLK Config = 86  // SPI Slave Serial Clock
	I2S0_MCLK      Config = 87  // I2S0 Master Clock
	I2S0_SCLK      Config = 88  // I2S0 Serial Clock(BCLK)
	I2S0_WS        Config = 89  // I2S0 Word Select(LRCLK)
	I2S0_IN_D0     Config = 90  // I2S0 Serial Data Input 0
	I2S0_IN_D1     Config = 91  // I2S0 Serial Data Input 1
	I2S0_IN_D2     Config = 92  // I2S0 Serial Data Input 2
	I2S0_IN_D3     Config = 93  // I2S0 Serial Data Input 3
	I2S0_OUT_D0    Config = 94  // I2S0 Serial Data Output 0
	I2S0_OUT_D1    Config = 95  // I2S0 Serial Data Output 1
	I2S0_OUT_D2    Config = 96  // I2S0 Serial Data Output 2
	I2S0_OUT_D3    Config = 97  // I2S0 Serial Data Output 3
	I2S1_MCLK      Config = 98  // I2S1 Master Clock
	I2S1_SCLK      Config = 99  // I2S1 Serial Clock(BCLK)
	I2S1_WS        Config = 100 // I2S1 Word Select(LRCLK)
	I2S1_IN_D0     Config = 101 // I2S1 Serial Data Input 0
	I2S1_IN_D1     Config = 102 // I2S1 Serial Data Input 1
	I2S1_IN_D2     Config = 103 // I2S1 Serial Data Input 2
	I2S1_IN_D3     Config = 104 // I2S1 Serial Data Input 3
	I2S1_OUT_D0    Config = 105 // I2S1 Serial Data Output 0
	I2S1_OUT_D1    Config = 106 // I2S1 Serial Data Output 1
	I2S1_OUT_D2    Config = 107 // I2S1 Serial Data Output 2
	I2S1_OUT_D3    Config = 108 // I2S1 Serial Data Output 3
	I2S2_MCLK      Config = 109 // I2S2 Master Clock
	I2S2_SCLK      Config = 110 // I2S2 Serial Clock(BCLK)
	I2S2_WS        Config = 111 // I2S2 Word Select(LRCLK)
	I2S2_IN_D0     Config = 112 // I2S2 Serial Data Input 0
	I2S2_IN_D1     Config = 113 // I2S2 Serial Data Input 1
	I2S2_IN_D2     Config = 114 // I2S2 Serial Data Input 2
	I2S2_IN_D3     Config = 115 // I2S2 Serial Data Input 3
	I2S2_OUT_D0    Config = 116 // I2S2 Serial Data Output 0
	I2S2_OUT_D1    Config = 117 // I2S2 Serial Data Output 1
	I2S2_OUT_D2    Config = 118 // I2S2 Serial Data Output 2
	I2S2_OUT_D3    Config = 119 // I2S2 Serial Data Output 3
	RESV0          Config = 120 // Reserved function
	RESV1          Config = 121 // Reserved function
	RESV2          Config = 122 // Reserved function
	RESV3          Config = 123 // Reserved function
	RESV4          Config = 124 // Reserved function
	RESV5          Config = 125 // Reserved function
	I2C0_SCLK      Config = 126 // I2C0 Serial Clock
	I2C0_SDA       Config = 127 // I2C0 Serial Data
	I2C1_SCLK      Config = 128 // I2C1 Serial Clock
	I2C1_SDA       Config = 129 // I2C1 Serial Data
	I2C2_SCLK      Config = 130 // I2C2 Serial Clock
	I2C2_SDA       Config = 131 // I2C2 Serial Data
	CMOS_XCLK      Config = 132 // DVP System Clock
	CMOS_RST       Config = 133 // DVP System Reset
	CMOS_PWDN      Config = 134 // DVP Power Down Mode
	CMOS_VSYNC     Config = 135 // DVP Vertical Sync
	CMOS_HREF      Config = 136 // DVP Horizontal Reference output
	CMOS_PCLK      Config = 137 // Pixel Clock
	CMOS_D0        Config = 138 // Data Bit 0
	CMOS_D1        Config = 139 // Data Bit 1
	CMOS_D2        Config = 140 // Data Bit 2
	CMOS_D3        Config = 141 // Data Bit 3
	CMOS_D4        Config = 142 // Data Bit 4
	CMOS_D5        Config = 143 // Data Bit 5
	CMOS_D6        Config = 144 // Data Bit 6
	CMOS_D7        Config = 145 // Data Bit 7
	SCCB_SCLK      Config = 146 // SCCB Serial Clock
	SCCB_SDA       Config = 147 // SCCB Serial Data
	UART1_CTS      Config = 148 // UART1 Clear To Send
	UART1_DSR      Config = 149 // UART1 Data Set Ready
	UART1_DCD      Config = 150 // UART1 Data Carrier Detect
	UART1_RI       Config = 151 // UART1 Ring Indicator
	UART1_SIR_IN   Config = 152 // UART1 Serial Infrared Input
	UART1_DTR      Config = 153 // UART1 Data Terminal Ready
	UART1_RTS      Config = 154 // UART1 Request To Send
	UART1_OUT2     Config = 155 // UART1 User-designated Output 2
	UART1_OUT1     Config = 156 // UART1 User-designated Output 1
	UART1_SIR_OUT  Config = 157 // UART1 Serial Infrared Output
	UART1_BAUD     Config = 158 // UART1 Transmit Clock Output
	UART1_RE       Config = 159 // UART1 Receiver Output Enable
	UART1_DE       Config = 160 // UART1 Driver Output Enable
	UART1_RS485_EN Config = 161 // UART1 RS485 Enable
	UART2_CTS      Config = 162 // UART2 Clear To Send
	UART2_DSR      Config = 163 // UART2 Data Set Ready
	UART2_DCD      Config = 164 // UART2 Data Carrier Detect
	UART2_RI       Config = 165 // UART2 Ring Indicator
	UART2_SIR_IN   Config = 166 // UART2 Serial Infrared Input
	UART2_DTR      Config = 167 // UART2 Data Terminal Ready
	UART2_RTS      Config = 168 // UART2 Request To Send
	UART2_OUT2     Config = 169 // UART2 User-designated Output 2
	UART2_OUT1     Config = 170 // UART2 User-designated Output 1
	UART2_SIR_OUT  Config = 171 // UART2 Serial Infrared Output
	UART2_BAUD     Config = 172 // UART2 Transmit Clock Output
	UART2_RE       Config = 173 // UART2 Receiver Output Enable
	UART2_DE       Config = 174 // UART2 Driver Output Enable
	UART2_RS485_EN Config = 175 // UART2 RS485 Enable
	UART3_CTS      Config = 176 // UART3 Clear To Send
	UART3_DSR      Config = 177 // UART3 Data Set Ready
	UART3_DCD      Config = 178 // UART3 Data Carrier Detect
	UART3_RI       Config = 179 // UART3 Ring Indicator
	UART3_SIR_IN   Config = 180 // UART3 Serial Infrared Input
	UART3_DTR      Config = 181 // UART3 Data Terminal Ready
	UART3_RTS      Config = 182 // UART3 Request To Send
	UART3_OUT2     Config = 183 // UART3 User-designated Output 2
	UART3_OUT1     Config = 184 // UART3 User-designated Output 1
	UART3_SIR_OUT  Config = 185 // UART3 Serial Infrared Output
	UART3_BAUD     Config = 186 // UART3 Transmit Clock Output
	UART3_RE       Config = 187 // UART3 Receiver Output Enable
	UART3_DE       Config = 188 // UART3 Driver Output Enable
	UART3_RS485_EN Config = 189 // UART3 RS485 Enable
	TIMER0_TOGGLE1 Config = 190 // TIMER0 Toggle Output 1
	TIMER0_TOGGLE2 Config = 191 // TIMER0 Toggle Output 2
	TIMER0_TOGGLE3 Config = 192 // TIMER0 Toggle Output 3
	TIMER0_TOGGLE4 Config = 193 // TIMER0 Toggle Output 4
	TIMER1_TOGGLE1 Config = 194 // TIMER1 Toggle Output 1
	TIMER1_TOGGLE2 Config = 195 // TIMER1 Toggle Output 2
	TIMER1_TOGGLE3 Config = 196 // TIMER1 Toggle Output 3
	TIMER1_TOGGLE4 Config = 197 // TIMER1 Toggle Output 4
	TIMER2_TOGGLE1 Config = 198 // TIMER2 Toggle Output 1
	TIMER2_TOGGLE2 Config = 199 // TIMER2 Toggle Output 2
	TIMER2_TOGGLE3 Config = 200 // TIMER2 Toggle Output 3
	TIMER2_TOGGLE4 Config = 201 // TIMER2 Toggle Output 4
	CLK_SPI2       Config = 202 // Clock SPI2
	CLK_I2C2       Config = 203 // Clock I2C2
	INTERNAL0      Config = 204 // Internal function signal 0
	INTERNAL1      Config = 205 // Internal function signal 1
	INTERNAL2      Config = 206 // Internal function signal 2
	INTERNAL3      Config = 207 // Internal function signal 3
	INTERNAL4      Config = 208 // Internal function signal 4
	INTERNAL5      Config = 209 // Internal function signal 5
	INTERNAL6      Config = 210 // Internal function signal 6
	INTERNAL7      Config = 211 // Internal function signal 7
	INTERNAL8      Config = 212 // Internal function signal 8
	INTERNAL9      Config = 213 // Internal function signal 9
	INTERNAL10     Config = 214 // Internal function signal 10
	INTERNAL11     Config = 215 // Internal function signal 11
	INTERNAL12     Config = 216 // Internal function signal 12
	INTERNAL13     Config = 217 // Internal function signal 13
	INTERNAL14     Config = 218 // Internal function signal 14
	INTERNAL15     Config = 219 // Internal function signal 15
	INTERNAL16     Config = 220 // Internal function signal 16
	INTERNAL17     Config = 221 // Internal function signal 17
	CONSTANT       Config = 222 // Constant function
	INTERNAL18     Config = 223 // Internal function signal 18
	DEBUG0         Config = 224 // Debug function 0
	DEBUG1         Config = 225 // Debug function 1
	DEBUG2         Config = 226 // Debug function 2
	DEBUG3         Config = 227 // Debug function 3
	DEBUG4         Config = 228 // Debug function 4
	DEBUG5         Config = 229 // Debug function 5
	DEBUG6         Config = 230 // Debug function 6
	DEBUG7         Config = 231 // Debug function 7
	DEBUG8         Config = 232 // Debug function 8
	DEBUG9         Config = 233 // Debug function 9
	DEBUG10        Config = 234 // Debug function 10
	DEBUG11        Config = 235 // Debug function 11
	DEBUG12        Config = 236 // Debug function 12
	DEBUG13        Config = 237 // Debug function 13
	DEBUG14        Config = 238 // Debug function 14
	DEBUG15        Config = 239 // Debug function 15
	DEBUG16        Config = 240 // Debug function 16
	DEBUG17        Config = 241 // Debug function 17
	DEBUG18        Config = 242 // Debug function 18
	DEBUG19        Config = 243 // Debug function 19
	DEBUG20        Config = 244 // Debug function 20
	DEBUG21        Config = 245 // Debug function 21
	DEBUG22        Config = 246 // Debug function 22
	DEBUG23        Config = 247 // Debug function 23
	DEBUG24        Config = 248 // Debug function 24
	DEBUG25        Config = 249 // Debug function 25
	DEBUG26        Config = 250 // Debug function 26
	DEBUG27        Config = 251 // Debug function 27
	DEBUG28        Config = 252 // Debug function 28
	DEBUG29        Config = 253 // Debug function 29
	DEBUG30        Config = 254 // Debug function 30
	DEBUG31        Config = 255 // Debug function 31
)
