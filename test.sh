#!/bin/bash

# The MIT License (MIT)

# Copyright (c) 2014-2016 George Lester

# Permission is hereby granted, free of charge, to any person obtaining a copy
# of this software and associated documentation files (the "Software"), to deal
# in the Software without restriction, including without limitation the rights
# to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
# copies of the Software, and to permit persons to whom the Software is
# furnished to do so, subject to the following conditions:

# The above copyright notice and this permission notice shall be included in all
# copies or substantial portions of the Software.

# THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
# IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
# FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
# AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
# LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
# OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
# SOFTWARE.

# Script that runs tests, code coverage, and benchmarks all at once.
# Builds a symlink in /tmp, mostly to avoid messing with GOPATH at the user's shell level.

TEMPORARY_PATH="/tmp/govaluate_test"
SRC_PATH="${TEMPORARY_PATH}/src"
FULL_PATH="${TEMPORARY_PATH}/src/govaluate"

# set up temporary directory
rm -rf "${FULL_PATH}"
mkdir -p "${SRC_PATH}"

ln -s $(pwd) "${FULL_PATH}"
export GOPATH="${TEMPORARY_PATH}"

pushd "${TEMPORARY_PATH}/src/govaluate"

# run the actual tests.
export GOVALUATE_TORTURE_TEST="true"
go test -bench=. -benchmem #-coverprofile coverage.out
status=$?

if [ "${status}" != 0 ];
then
	exit $status
fi

# coverage
# disabled because travis go1.4 seems not to support it suddenly?
#go tool cover -func=coverage.out

popd
