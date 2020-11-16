#!/bin/sh

INTERFACE=ftdi/cjmcu-2232hl
TARGET=k210
OOCD=openocd-kendryte
RESET=none

. ../../../../../scripts/debug-oocd.sh $@
