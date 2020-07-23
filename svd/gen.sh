#!/bin/sh

set -e

cd ../../../embeddedgo/kendryte/hal
hal=$(pwd)
cd ../p
rm -rf *

svdxgen github.com/embeddedgo/kendryte/p ../svd/*.svd

for p in sysctl fpioa gpio gpiohs uart uarths timer spi; do
	cd $p
	xgen *.go
	GOOS=noos GOARCH=riscv64 go build -tags k210
	cd ..
done

perlscript='
s/ = \d/ rtos.IRQ$&/g;
s/package irq/$&

import "embedded\/rtos"

const (
	M0 rtos.IntCtx = 0 \/\/ machine mode on hart 0
	S0 rtos.IntCtx = 1 \/\/ supervisor mode on hart 0
	M1 rtos.IntCtx = 2 \/\/ machine mode on hart 1
	S1 rtos.IntCtx = 3 \/\/ supervisor mode on hart 1
)/;
'

cd $hal/irq
rm -f *
cp ../../p/irq/* .
perl -pi -e "$perlscript" *.go
