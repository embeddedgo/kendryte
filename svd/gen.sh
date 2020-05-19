#!/bin/sh

set -e

cd ../../../embeddedgo/kendryte/hal
hal=$(pwd)
cd ../p
rm -rf *

svdxgen github.com/embeddedgo/kendryte/p ../svd/*.svd

for p in sysctl fpioa gpio gpiohs; do
	cd $p
	xgen *.go
	GOOS=noos GOARCH=riscv64 go build -tags k210
	cd ..
done

perlscript='
s/package irq/$&\n\nimport "embedded\/rtos"/;
s/ = \d/ rtos.IRQ$&/g;
'

cd $hal/irq
rm -f *
cp ../../p/irq/* .
perl -pi -e "$perlscript" *.go