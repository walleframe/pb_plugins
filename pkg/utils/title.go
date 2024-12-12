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
package utils

import (
	"strings"
	"unsafe"
)

func Title(in string) (out string) {
	if strings.Contains(in, ".") {
		list := strings.Split(in, ".")
		list[len(list)-1] = Title(list[len(list)-1])
		return strings.Join(list, ".")
	}
	list := strings.Split(in, "_")
	for _, v := range list {
		out += strings.Title(v)
	}
	return
}

// string转bytes
func UnsafeStr2Bytes(s string) []byte {
	return unsafe.Slice(unsafe.StringData(s), len(s))
}

// []byte转string
func UnsafeBytes2Str(b []byte) string {
	return unsafe.String(unsafe.SliceData(b), len(b))
}

func MergeNotSameSlice(l1, l2 []string) (res []string) {
	for _, v := range l1 {
		if v == "" {
			continue
		}
		res = append(res, strings.ToLower(v))
	}
	for _, v := range l2 {
		if v == "" {
			continue
		}
		find := false
		for _, v2 := range res {
			if v == v2 {
				find = true
				break
			}
		}
		if find {
			continue
		}
		res = append(res, strings.ToLower(v))
	}
	return
}
