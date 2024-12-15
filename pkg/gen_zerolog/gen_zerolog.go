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
package gen_zerolog

import (
	_ "embed"

	"github.com/walleframe/pb_plugins/pkg/tpl"
)

//go:embed zerolog.tpl
var mainTpl string

type GenerateEnv struct {
	// SvcPkg  string
	// UtilPkg string
	// Charset string
	// Collate string
}

var Config = &GenerateEnv{
	// SvcPkg:  "github.com/walleframe/svc_db",
	// UtilPkg: "github.com/walleframe/walle/util",
	// Charset: "utf8mb4",
	// Collate: "utf8mb4_general_ci",
}

func ParseConfigFromEnv() {
	// // 如果环境变量设置了值, 读取作为为默认值. 优先使用传递的参数
	// utils.GetEnvString("ZAP_COLLATE", &Config.Collate)
	// utils.GetEnvString("ZAP_CHARSET", &Config.Charset)
	// utils.GetEnvString("ZAP_UTIL_PKG", &Config.UtilPkg)
	// utils.GetEnvString("ZAP_SVC_PKG", &Config.SvcPkg)
}

var gen *tpl.GoTemplate

func InitTemplate() (err error) {
	gen = tpl.NewTemplate("zerolog")
	gen.Funcs(map[string]interface{}{})
	err = gen.Parse(mainTpl)
	if err != nil {
		return err
	}

	return nil
}

type MapField struct {
	Type    string                // go里面,map的类型
	KeyFunc func(k string) string // key的函数
	ValFunc string                // value的函数
}

type Field struct {
	Key   string // 写入日志的字段名
	Name  string // go结构体字段名
	Func  string // 使用的函数方法
	Value func(obj string) string
	Map   *MapField // 是否是map
}

type Message struct {
	Name   string
	Type   string
	Fields []*Field
}

type ZlogMessages struct {
	tpl.GoObject

	Package string
	File    string

	Msgs []*Message
}

func Generate(zlog *ZlogMessages) (out []*tpl.BuildOutput, err error) {
	zlog.Version = "0.0.1"

	// genTpl, err := gen.Clone()
	// if err != nil {
	// 	return nil, err
	// }
	genTpl := gen

	genTpl.AddImportFunc(zlog)

	for _, pkg := range []string{
		"github.com/rs/zerolog",
	} {
		zlog.Import(pkg, "pkg")
	}
	// tbl.Import(Config.SvcPkg, "svc_db")
	// tbl.Import(Config.UtilPkg, "util")

	data, err := genTpl.Exec(zlog)
	if err != nil {
		return nil, err
	}

	out = append(out, &tpl.BuildOutput{
		File: zlog.File,
		Data: data,
	})

	return
}
