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
	"sync"

	"log"
	"path/filepath"

	"github.com/walleframe/pb_plugins/pkg/tpl"
	"github.com/walleframe/pb_plugins/pkg/utils"
)

//go:embed mysql_op.tpl
var strMysqlOp string

//go:embed mysql_def.tpl
var strMysqlDef string

//go:embed mysql_tx.tpl
var strMysqlTx string

type GenerateEnv struct {
	SvcPkg    string
	UtilPkg   string
	CodePath  string
	CodePkg   string
	Charset   string
	Collate   string
	DisableEx bool
}

var Config = &GenerateEnv{
	SvcPkg:   "github.com/walleframe/svc_db",
	UtilPkg:  "github.com/walleframe/walle/util",
	CodePkg:  "github.com/walleframe/svc_db/pkg/gen/mysqlop",
	CodePath: "pkg/gen/mysqlop/",
	Charset:  "utf8mb4",
	Collate:  "utf8mb4_general_ci",
}

func ParseConfigFromEnv() {
	// 如果环境变量设置了值, 读取作为为默认值. 优先使用传递的参数
	utils.GetEnvString("MYSQL_COLLATE", &Config.Collate)
	utils.GetEnvString("MYSQL_CHARSET", &Config.Charset)
	utils.GetEnvString("MYSQL_OPCODE_PATH", &Config.CodePath)
	utils.GetEnvString("MYSQL_OPCODE_PKG", &Config.CodePkg)
	utils.GetEnvString("MYSQL_PKG_UTIL", &Config.UtilPkg)
	utils.GetEnvString("MYSQL_PKG_SVC", &Config.SvcPkg)
	utils.GetEnvBool("MYSQL_DISABLE_EX", &Config.DisableEx)
}

// mysql操作代码
var genOperation *tpl.GoTemplate

// mysql定义代码
var genDefinition *tpl.GoTemplate

// mysql tx
var (
	genTx *tpl.GoTemplate
	txMap sync.Map
)

// 初始化模板
func InitTemplate() (err error) {
	genOperation = tpl.NewTemplate("mysql_op")
	genOperation.Funcs(map[string]interface{}{
		"And": func(i int) string {
			if i == 0 {
				return ""
			}
			return " and "
		},
	})
	err = genOperation.Parse(strMysqlOp)
	if err != nil {
		return err
	}

	genDefinition = tpl.NewTemplate("mysql_def")
	genDefinition.Funcs(map[string]interface{}{
		"And": func(i int) string {
			if i == 0 {
				return ""
			}
			return " and "
		},
	})
	err = genDefinition.Parse(strMysqlDef)
	if err != nil {
		return err
	}

	genTx = tpl.NewTemplate("mysql_tx")
	err = genTx.Parse(strMysqlTx)
	if err != nil {
		return err
	}

	return nil
}

func Generate(tbl *SqlTable) (out []*tpl.BuildOutput, err error) {
	tbl.SvcDB = filepath.Base(Config.SvcPkg)
	tbl.Version = "0.0.2"

	if _, load := txMap.LoadOrStore(tbl.DB, struct{}{}); !load {
		o3, err := generateTx(tbl)
		if err != nil {
			return nil, err
		}
		out = append(out, o3...)
	}

	o1, err := generateDefine(tbl)
	if err != nil {
		return nil, err
	}
	out = append(out, o1...)

	o2, err := generateOpreation(tbl)
	if err != nil {
		return nil, err
	}
	out = append(out, o2...)

	return
}

func generateDefine(tbl *SqlTable) (out []*tpl.BuildOutput, err error) {

	// genTpl, err := gen.Clone()
	// if err != nil {
	// 	return nil, err
	// }
	genTpl := genDefinition

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

	defPkg := utils.PascalToSnake(tbl.Name) + "_def"
	tbl.DefPkg = defPkg + "."

	data, err := genTpl.Exec(tbl)
	if err != nil {
		return nil, err
	}

	out = append(out, &tpl.BuildOutput{
		File: filepath.Join(Config.CodePath, strings.ToLower(tbl.DB), defPkg, strings.ToLower(tbl.SqlTable)+".go"),
		Data: data,
	})

	log.Printf("generate mysql define [%s.%s] success", tbl.DB, tbl.Name)
	backupName := tbl.Name
	backupTable := tbl.SqlTable

	for _, dup := range tbl.Duplication {
		tbl.SqlTable = dup.TableName
		tbl.Name = dup.OpName
		defPkg := utils.PascalToSnake(tbl.Name) + "_def"
		tbl.DefPkg = defPkg + "."
		data, err := genTpl.Exec(tbl)
		if err != nil {
			return nil, err
		}
		out = append(out, &tpl.BuildOutput{
			File: filepath.Join(Config.CodePath, strings.ToLower(tbl.DB), defPkg, strings.ToLower(tbl.SqlTable)+".go"),
			Data: data,
		})
		log.Printf("generate mysql duplication define [%s.%s] success", tbl.DB, tbl.Name)
	}
	tbl.Name = backupName
	tbl.SqlTable = backupTable

	return
}

func generateOpreation(tbl *SqlTable) (out []*tpl.BuildOutput, err error) {

	// genTpl, err := gen.Clone()
	// if err != nil {
	// 	return nil, err
	// }
	genTpl := genOperation

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

	defPkg := utils.PascalToSnake(tbl.Name) + "_def"
	tbl.DefPkg = defPkg + "."
	pkgPath := filepath.ToSlash(filepath.Join(Config.CodePkg, strings.ToLower(tbl.DB), defPkg))
	tbl.Import(pkgPath, "def")

	data, err := genTpl.Exec(tbl)
	if err != nil {
		tbl.Remove(pkgPath)
		return nil, err
	}
	tbl.Remove(pkgPath)

	out = append(out, &tpl.BuildOutput{
		File: filepath.Join(Config.CodePath, strings.ToLower(tbl.DB), strings.ToLower(tbl.SqlTable)+".dbop.go"),
		Data: data,
	})

	log.Printf("generate mysql code [%s.%s] success", tbl.DB, tbl.Name)

	backName := tbl.Name
	backTable := tbl.SqlTable

	for _, dup := range tbl.Duplication {
		tbl.SqlTable = dup.TableName
		tbl.Name = dup.OpName

		defPkg := utils.PascalToSnake(tbl.Name) + "_def"
		tbl.DefPkg = defPkg + "."
		pkgPath := filepath.Join(Config.CodePkg, strings.ToLower(tbl.DB), defPkg)
		tbl.Import(pkgPath, "def")

		data, err := genTpl.Exec(tbl)
		if err != nil {
			tbl.Remove(pkgPath)
			return nil, err
		}
		tbl.Remove(pkgPath)
		out = append(out, &tpl.BuildOutput{
			File: filepath.Join(Config.CodePath, strings.ToLower(tbl.DB), strings.ToLower(tbl.SqlTable)+".dbop.go"),
			Data: data,
		})
		log.Printf("generate mysql duplication [%s.%s] success", tbl.DB, tbl.Name)
	}

	tbl.Name = backName
	tbl.SqlTable = backTable

	return
}

func generateTx(tbl *SqlTable) (out []*tpl.BuildOutput, err error) {

	// genTpl, err := gen.Clone()
	// if err != nil {
	// 	return nil, err
	// }
	genTpl := genTx

	genTpl.AddImportFunc(tbl)

	for _, pkg := range []string{
		"context",
		"database/sql",
		"fmt",
		"github.com/jmoiron/sqlx",
	} {
		tbl.Import(pkg, "pkg")
	}

	tbl.Import(Config.SvcPkg, "svc_db")

	log.Printf("%#v", tbl)

	data, err := genTpl.Exec(tbl)
	if err != nil {
		return nil, err
	}

	out = append(out, &tpl.BuildOutput{
		File: filepath.Join(Config.CodePath, strings.ToLower(tbl.DB), strings.ToLower(tbl.DB)+"_tx.go"),
		Data: data,
	})

	log.Printf("generate mysql tx [%s] success", tbl.DB)

	return
}
