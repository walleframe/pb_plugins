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
package keyarg

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/walleframe/pb_plugins/pkg/utils"
)

func TestMatchKey(t *testing.T) {
	var datas = []struct {
		name    string
		input   string
		outputs []KeyArg
		err     string
	}{
		// source arg match
		{
			name:    "normal1",
			input:   "normal",
			outputs: []KeyArg{&SourceArg{source: "normal"}},
		},
		{
			name:    "normal2",
			input:   "normal:k2",
			outputs: []KeyArg{&SourceArg{source: "normal:k2"}},
		},
		{
			name:    "normal3",
			input:   "nv:v2:v3",
			outputs: []KeyArg{&SourceArg{source: "nv:v2:v3"}},
		},
		// time arg match
		{
			name:    "time-day",
			input:   "@daystamp",
			outputs: []KeyArg{&TimeArg{fun: "DayStamp"}},
		},
		{
			name:    "time-month",
			input:   "@monTHstamp+8",
			outputs: []KeyArg{&TimeArg{fun: "MonthStamp", op: " + 8"}},
		},
		{
			name:    "time-weak",
			input:   "@weekstamp - 10",
			outputs: []KeyArg{&TimeArg{fun: "WeekStamp", op: " - 10"}},
		},
		{
			name:    "time-year",
			input:   "@yearstamp",
			outputs: []KeyArg{&TimeArg{fun: "YearStamp"}},
		},
		// go type arg match
		{
			name:  "gotype-normal",
			input: "$int8:$int16:$int32:$int64:$uint8:$uint16:$uint32:$uint64:$bool:$string:$float32:$float64",
			outputs: []KeyArg{
				&NumberArg{typ: "int8", arg: "arg1"},
				&NumberArg{typ: "int16", arg: "arg2"},
				&NumberArg{typ: "int32", arg: "arg3"},
				&NumberArg{typ: "int64", arg: "arg4"},
				&NumberArg{typ: "uint8", arg: "arg5"},
				&NumberArg{typ: "uint16", arg: "arg6"},
				&NumberArg{typ: "uint32", arg: "arg7"},
				&NumberArg{typ: "uint64", arg: "arg8"},
				&NumberArg{typ: "bool", arg: "arg9"},
				&StringArg{arg: "arg10"},
				&NumberArg{typ: "float32", arg: "arg11"},
				&NumberArg{typ: "float64", arg: "arg12"},
			},
		},
		{
			name:  "gotype-op/name",
			input: "$uid=int64:$k2=uint32%10:$uint8-16",
			outputs: []KeyArg{
				&NumberArg{typ: "int64", arg: "uid"},
				&NumberArg{typ: "uint32", arg: "k2", op: " % 10"},
				&NumberArg{typ: "uint8", arg: "arg1", op: " - 16"},
			},
		},
		// mix args
		{
			name:  "mix-1",
			input: "ud:xy:$uid=int64:$k2=uint32%10:$uint8-16:@daystamp+3600",
			outputs: []KeyArg{
				&SourceArg{source: "ud:xy"},
				&NumberArg{typ: "int64", arg: "uid"},
				&NumberArg{typ: "uint32", arg: "k2", op: " % 10"},
				&NumberArg{typ: "uint8", arg: "arg1", op: " - 16"},
				&TimeArg{fun: "DayStamp", op: " + 3600"},
			},
		},
		// error check
		{
			name:  "not match any",
			input: "#xx:x",
			err:   "contain invalid char",
		},
		{
			name:  "conflict arg name",
			input: "$uid=int32:$uid=int32",
			err:   "conflict arg name",
		},
		{
			name:  "invalid op",
			input: "$uid=int32%:$uid=int32",
			err:   "invalid syntax",
		},
	}

	for _, data := range datas {
		t.Run(data.name, func(t *testing.T) {
			args, err := MatchKey(data.input, nil)
			if len(data.err) > 0 {
				assert.NotNil(t, err, "except an error")
				assert.Contains(t, err.Error(), data.err, "error match")
				t.Log(utils.Sdump(args, "outputs"))
			} else {
				assert.Nil(t, err, "need no error")
				assert.EqualValues(t, data.outputs, args, "compare output values")
			}
		})
	}
}
