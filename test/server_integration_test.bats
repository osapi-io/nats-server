# Copyright (c) 2025 John Dewey

# Permission is hereby granted, free of charge, to any person obtaining a copy
# of this software and associated documentation files (the "Software"), to
# deal in the Software without restriction, including without limitation the
# rights to use, copy, modify, merge, publish, distribute, sublicense, and/or
# sell copies of the Software, and to permit persons to whom the Software is
# furnished to do so, subject to the following conditions:

# The above copyright notice and this permission notice shall be included in
# all copies or substantial portions of the Software.

# THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
# IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
# FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
# AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
# LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
# FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER
# DEALINGS IN THE SOFTWARE.

load "libs/common.bash"

setup() {
	cd ../examples/simple-server
	start_server
}

teardown() {
	cd -
	stop_server
}

@test "publish and subscribe to a test message" {
	PIN=$(date +"%Y%m%d%H%M%S")

	nats sub test-subject >output.txt &
	SUB_PID=$!
	sleep 1

	nats pub test-subject "PIN: $PIN"

	sleep 1

	grep "PIN: $PIN" output.txt
	[ "$?" -eq 0 ]

	kill $SUB_PID
}
