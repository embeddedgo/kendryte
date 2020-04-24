#!/bin/sh

export GOOS=noos
export GOARCH=riscv64

GOTARGET=k210
GOMEM=0x80000000:6M

go build -tags $GOTARGET -ldflags "-M $GOMEM" -o blinky.elf
