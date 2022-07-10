#!/bin/sh

INTERFACE=ftdi/cjmcu-2232hl
TARGET=k210
OOCD=openocd-kendryte
RESET=none

. $(emgo env GOROOT)/../scripts/debug-oocd.sh $@
