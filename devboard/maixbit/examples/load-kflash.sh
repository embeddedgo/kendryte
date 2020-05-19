#!/bin/sh

set -e

name=$(basename $(pwd))
objcopy -O binary $name.elf $name.bin
kflash -p /dev/ttyUSB0 -b 1500000 -B bit_mic $name.bin
rm $name.bin