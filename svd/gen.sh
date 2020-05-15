#!/bin/sh

set -e

cd ../../../embeddedgo/kendryte/hal
hal=$(pwd)
cd ../p
rm -rf *

svdxgen github.com/embeddedgo/kendryte/p ../svd/*.svd
