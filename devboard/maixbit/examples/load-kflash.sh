#!/bin/sh

set -e

name=$(basename $(pwd))
objcopy -O binary $name.elf $name.bin
kflash -p /dev/ttyUSB0 -B bit_mic -b 750000 $@ $name.bin
rm $name.bin
