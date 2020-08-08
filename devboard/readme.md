## Support for Kendryte K210 development boards

### Directory structure

Every board directory contains a set of packages (in *board* subdirectory) that provides the interface to the peripherals available on the board (for now the support is modest: only LEDs and buttons).

The board/init package, when imported, configures the whole system for typical usage. If you use any other package from *board* directory the board/init package is imported implicitly to ensure the board is properly configured.

The *examples* subdirectory as the name suggests includes sample code, but also scripts and configuration that help to build, load and debug.

There is also *doc* subdirectory that contain useful information and other resources about this development board.

### Supported boards

[maixbit](maixbit): Sipeed Maix Bit v2 development board based on Kendryte [K210](https://s3.cn-north-1.amazonaws.com.cn/dl.kendryte.com/documents/kendryte_datasheet_20181011163248_en.pdf) SOC, [website](https://maixduino.sipeed.com/en/hardware/board.html)

![Maix Bit v2](maixbit/doc/board.jpg)
