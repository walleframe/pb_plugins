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
