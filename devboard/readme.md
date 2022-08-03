## Support for Kendryte K210 development boards

### Supported boards

#### [maixbit](maixbit)

Sipeed [Maix Bit v2](maixbit/doc/board.jpg) development board based on Kendryte [K210](https://s3.cn-north-1.amazonaws.com.cn/dl.kendryte.com/documents/kendryte_datasheet_20181011163248_en.pdf) SOC

![Maix Bit v2](maixbit/doc/board.jpg)

### Directory structure

Every supported board has its directory here. Each of them contains three subdirectories:

#### board

The `board` directory contains Go packages that provide the interface to the peripherals available on the board such as buttons, LEDs, etc.
There is also an `init` package that imported configures the whole system for typical usage.
Usually you don't need to import it explicitly because any other package in the `board` directory imports `init` itself to ensure proper operation.

#### doc

The `doc` directory may contain documentation, diagrams, photos and other useful resources.

#### examples

The `examples` directory contains simple programs that demonstrate how to use the available Go packages to work with the board and various external peripherals connected to it.
