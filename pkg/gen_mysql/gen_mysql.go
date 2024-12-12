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
package gen_mysql

import (
	_ "embed"
	"strings"

	"log"
	"path/filepath"

	"github.com/walleframe/pb_plugins/pkg/tpl"
	"github.com/walleframe/pb_plugins/pkg/utils"
)

//go:embed mysql.tpl
var wdbTpl string

type GenerateEnv struct {
	SvcPkg    string
	UtilPkg   string
	CodePath  string
	Charset   string
	Collate   string
	DisableEx bool
}

var Config = &GenerateEnv{
	SvcPkg:   "github.com/walleframe/svc_db",
	UtilPkg:  "github.com/walleframe/walle/util",
	CodePath: "pkg/gen/mysqlop/",
	Charset:  "utf8mb4",
	Collate:  "utf8mb4_general_ci",
}

func ParseConfigFromEnv() {
	// 如果环境变量设置了值, 读取作为为默认值. 优先使用传递的参数
	utils.GetEnvString("MYSQL_COLLATE", &Config.Collate)
	utils.GetEnvString("MYSQL_CHARSET", &Config.Charset)
	utils.GetEnvString("MYSQL_OPCODE_PATH", &Config.CodePath)
	utils.GetEnvString("MYSQL_PKG_UTIL", &Config.UtilPkg)
	utils.GetEnvString("MYSQL_PKG_SVC", &Config.SvcPkg)
	utils.GetEnvBool("MYSQL_DISABLE_EX", &Config.DisableEx)
}

var gen *tpl.GoTemplate

func InitTemplate() (err error) {
	gen = tpl.NewTemplate("mysql")
	gen.Funcs(map[string]interface{}{
		"And": func(i int) string {
			if i == 0 {
				return ""
			}
			return " and "
		},
	})
	err = gen.Parse(wdbTpl)
	if err != nil {
		return err
	}

	return nil
}

func Generate(tbl *SqlTable) (out []*tpl.BuildOutput, err error) {
	tbl.SvcDB = filepath.Base(Config.SvcPkg)
	tbl.Version = "0.0.1"

	// genTpl, err := gen.Clone()
	// if err != nil {
	// 	return nil, err
	// }
	genTpl := gen

	genTpl.AddImportFunc(tbl)

	for _, pkg := range []string{
		"context",
		"database/sql",
		"errors",
		"strings",
		"fmt",
		"sync/atomic",
		"github.com/jmoiron/sqlx",
	} {
		tbl.Import(pkg, "pkg")
	}
	tbl.Import(Config.SvcPkg, "svc_db")
	tbl.Import(Config.UtilPkg, "util")

	data, err := genTpl.Exec(tbl)
	if err != nil {
		return nil, err
	}

	out = append(out, &tpl.BuildOutput{
		File: filepath.Join(Config.CodePath, strings.ToLower(tbl.DB), strings.ToLower(tbl.SqlTable)+".dbop.go"),
		Data: data,
	})

	log.Println(tbl.DB, tbl.Name, "generate mysql code success")

	return
}
