#!/bin/sh

GOTARGET=k210
GOMEM=0x80000000:6M

. $(emgo env GOROOT)/../scripts/build.sh $@
