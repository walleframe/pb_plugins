// Copyright 2024 aggronmagi <czy463@163.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package tpl

import (
	"bytes"
	"go/format"
	"os/exec"

	"golang.org/x/tools/imports"
)

type FormatFunc func(in []byte) ([]byte, error)

func GoFormat(in []byte) (out []byte, err error) {
	return format.Source(in)
}

func GoimportsCmdFormat(in []byte) (out []byte, err error) {
	ib := bytes.NewBuffer(in)
	ob := bytes.NewBuffer(nil)
	cmd := exec.Command("goimports")
	cmd.Stdin = ib
	cmd.Stdout = ob

	err = cmd.Run()
	if err != nil {
		return nil, err
	}

	return ob.Bytes(), err
}

func GoFormat2(in []byte) ([]byte, error) {
	opts := &imports.Options{
		TabIndent: false,
		TabWidth:  4,
		Fragment:  true,
		Comments:  true,
	}
	data, err := imports.Process("", in, opts)
	if err != nil {
		return nil, err
	}
	return format.Source(data)
}

// OptionGoimportsFormtat use the command line `goimports` command to
// format first. If an error is reported, use the default formatting method.
func OptionGoimportsFormtat(in []byte) ([]byte, error) {
	out := bytes.NewBuffer(nil)
	cmd := exec.Command("goimports")
	cmd.Stdin = bytes.NewBuffer(in)
	cmd.Stdout = out
	err := cmd.Run()
	if err == nil {
		return out.Bytes(), nil
	}
	opts := &imports.Options{
		TabIndent: false,
		TabWidth:  4,
		Fragment:  true,
		Comments:  true,
	}
	data, err := imports.Process("", in, opts)
	if err != nil {
		return nil, err
	}
	return format.Source(data)
}
