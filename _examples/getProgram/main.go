// Copyright (c) 2016 Uber Technologies, Inc.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package main

import (
	"github.com/uber-go/hackeroni/h1"

	"golang.org/x/crypto/ssh/terminal"

	"bufio"
	"fmt"
	"os"
	"strings"
	"syscall"
)

func main() {

	fmt.Print("HackerOne API Identifier: ")
	r := bufio.NewReader(os.Stdin)
	identifier, _ := r.ReadString('\n')

	fmt.Print("HackerOne API Token: ")
	token, _ := terminal.ReadPassword(int(syscall.Stdin))
	fmt.Print("\n")

	tp := h1.APIAuthTransport{
		APIIdentifier: strings.TrimSpace(identifier),
		APIToken:      strings.TrimSpace(string(token)),
	}

	client := h1.NewClient(tp.Client())

	fmt.Print("Program ID: ")
	programID, _ := r.ReadString('\n')
	fmt.Print("\n")

	program, _, err := client.Program.Get(strings.TrimSpace(programID))
	if err != nil {
		panic(err)
	}

	fmt.Printf("Program: %#v\n", program)

}
