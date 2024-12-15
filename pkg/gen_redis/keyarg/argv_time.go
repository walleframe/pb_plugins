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
	"fmt"
	"strings"
)

var (
	WTimePkg = "github.com/walleframe/walle/util/wtime"
)

type TimeArg struct {
	EmptyArg
	fun string
	op  string
}

// 需要导入的包
func (x *TimeArg) Imports() []string {
	return []string{WTimePkg}
}

// 格式化代码
func (x *TimeArg) FormatCode(obj string) (code string) {
	if strings.HasSuffix(x.fun, "String") {
		return fmt.Sprintf(`%s.WriteString(wtime.%s())`, obj, x.fun)
	}
	return fmt.Sprintf(`%s.WriteInt64(wtime.%s()%s)`, obj, x.fun, x.op)
}

// TimeMatch 时间参数匹配
type TimeMatch struct{}

func (x TimeMatch) Match(st KeyStructer, k ArgIndex, arg string) (_ KeyArg, err error) {
	// 非@开头，不匹配
	if !strings.HasPrefix(arg, "@") {
		return nil, nil
	}

	fun := trimParentheses(strings.TrimPrefix(arg, "@"))
	op := ""
	fun, op, err = splitKeyOp(fun)
	if err != nil {
		return nil, err
	}
	switch strings.ToLower(fun) {
	case "daystamp":
		fun = "DayStamp"
	case "weekstamp":
		fun = "WeekStamp"
	case "monthstamp":
		fun = "MonthStamp"
	case "yearstamp":
		fun = "YearStamp"
	case "curday":
		fun = "CurDayString"
	case "nextday":
		fun = "NextDayString"
	case "curweek":
		fun = "CurWeekString"
	case "nextweek":
		fun = "NextWeekString"
	case "curmonth":
		fun = "CurMonthString"
	case "nextmonth":
		fun = "NextMonthString"
	case "curyear":
		fun = "CurYearString"
	case "nextyear":
		fun = "NextYearString"
	default:
		err = fmt.Errorf("[%s] is invalid wtime functions.", fun)
		return
	}
	if strings.HasSuffix(fun, "String") && op != "" {
		err = fmt.Errorf("arg:[%s] string format not support number operation", arg)
		return
	}
	return &TimeArg{fun: fun, op: op}, nil
}
