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
package gen_redis

import (
	"embed"
	_ "embed"
	"io/fs"

	"bytes"
	"fmt"
	"log"
	"path/filepath"
	"strings"

	"github.com/walleframe/pb_plugins/pkg/tpl"
	"github.com/walleframe/pb_plugins/pkg/utils"
)

type GenerateEnv struct {
	ProtobufPackage string
	WProtoPackage   string
	ServicePackage  string
	UtilPkg         string
	CodePath        string
}

var Config = &GenerateEnv{
	ProtobufPackage: "", // github.com/gogo/protobuf/proto
	WProtoPackage:   "github.com/walleframe/walle/process/message",
	ServicePackage:  "github.com/walleframe/svc_redis",
	UtilPkg:         "github.com/walleframe/walle/util",
	CodePath:        "pkg/gen/redisop/",
}

var envFlag = struct {
}{}

func init() {

	utils.GetEnvString("WREDIS_PB_PKG", &Config.ProtobufPackage)
	utils.GetEnvString("WREDIS_WPB_PKG", &Config.WProtoPackage)
	utils.GetEnvString("WREDIS_SVC_PKG", &Config.ServicePackage)
	utils.GetEnvString("WREDIS_UTIL_PKG", &Config.UtilPkg)
	utils.GetEnvString("WREDIS_OPCODE_PATH", &Config.CodePath)
}

//go:embed redis.tpl
var redisTpl string

//go:embed tpls/*.tpl
var redisOpTpls embed.FS

var gen *tpl.GoTemplate

func InitTemplate() (err error) {
	gen = tpl.NewTemplate("redis")
	gen.Funcs(map[string]interface{}{
		// 循环调用模板函数
		"GenTypeTemplate": func(typeTplName string, obj *RedisObject) (string, error) {

			buf := &bytes.Buffer{}
			dst := gen.Lookup(typeTplName)
			if dst == nil {
				return "", fmt.Errorf("%v not found", typeTplName)
			}

			err = dst.Execute(buf, obj)
			if err != nil {
				log.Println(err)
				return "", err
			}

			//		log.Println(buf.String())

			return buf.String(), nil
		},
		"GenScriptTemplate": func(obj *RedisObject, script *RedisScript) (string, error) {
			buf := &bytes.Buffer{}
			dst := gen.Lookup(script.TemplateName)
			if dst == nil {
				return "", fmt.Errorf("%v not found", script.TemplateName)
			}

			err = dst.Execute(buf, map[string]interface{}{
				"Obj":    obj,
				"Script": script,
			})
			if err != nil {
				log.Println(err)
				return "", err
			}

			return buf.String(), nil
		},
	})

	// 解析前置的全部模板
	fs.WalkDir(redisOpTpls, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		data, err := redisOpTpls.ReadFile(path)
		if err != nil {
			return err
		}

		k := strings.TrimSuffix(path, filepath.Ext(path))
		k = strings.TrimPrefix(k, "tpls/")

		ts := fmt.Sprintf(`{{define "%s"}} %s {{end}}`, k, strings.TrimSpace(string(data)))
		err = gen.Parse(ts)
		if err != nil {
			err = fmt.Errorf("parse template %s failed:%+v", k, err.Error())
			return err
		}

		return nil
	})

	err = gen.Parse(redisTpl)
	if err != nil {
		return err
	}

	return nil
}

func Generate(tbl *RedisObject) (out []*tpl.BuildOutput, err error) {
	tbl.Version = "0.0.1"

	// genTpl, err := gen.Clone()
	// if err != nil {
	// 	return nil, err
	// }
	genTpl := gen

	genTpl.AddImportFunc(tbl)

	for _, pkg := range []string{
		"context",
		"github.com/redis/go-redis/v9",
		"github.com/walleframe/walle/util",
		"github.com/walleframe/walle/util/rdconv",
	} {
		tbl.Import(pkg, "pkg")
	}

	tbl.Import(Config.UtilPkg, "util.Builder")
	tbl.Import(filepath.Join(Config.UtilPkg, "rdconv"), "AnyToString/Int64ToString/...")
	tbl.Import(Config.ServicePackage, "RegisterDBName/GetDBLink")

	for _, arg := range tbl.Args {
		for _, pkg := range arg.Imports() {
			if pkg == "" {
				continue
			}
			tbl.Import(pkg, "keyargv")
		}
	}

	data, err := genTpl.Exec(tbl)
	if err != nil {
		return nil, err
	}

	out = append(out, &tpl.BuildOutput{
		File: filepath.Join(Config.CodePath, strings.ToLower(tbl.Package), strings.ToLower(tbl.Name)+".redisop.go"),
		Data: data,
	})

	log.Println(tbl.Package, tbl.Name, "generate redis code success")

	return
}
