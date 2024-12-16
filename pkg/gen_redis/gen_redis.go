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

	"github.com/walleframe/pb_plugins/pkg/gen_redis/keyarg"
	"github.com/walleframe/pb_plugins/pkg/tpl"
	"github.com/walleframe/pb_plugins/pkg/utils"
)

type GenerateEnv struct {
	PkgPB       string
	PkgJson     string
	PkgSvc      string
	PkgUtil     string
	PkgTime     string
	CodePath    string
	MsgProtocol string
	//PkgWPB        string
	OpAnotherFile  bool
	CheckLuaScript bool
}

var Config = &GenerateEnv{
	PkgJson:     "encoding/json",
	PkgPB:       "google.golang.org/protobuf/proto", // github.com/gogo/protobuf/proto
	PkgSvc:      "github.com/walleframe/svc_redis",
	PkgUtil:     "github.com/walleframe/walle/util",
	PkgTime:     "github.com/walleframe/walle/util/wtime",
	CodePath:    "pkg/gen/redisop/",
	MsgProtocol: "proto",
	//PkgWPB:        "github.com/walleframe/walle/process/message",
	OpAnotherFile:  false,
	CheckLuaScript: true,
}

func ParseConfigFromEnv() {
	utils.GetEnvString("REDIS_PKG_PB", &Config.PkgPB)
	utils.GetEnvString("REDIS_PKG_SVC", &Config.PkgSvc)
	utils.GetEnvString("REDIS_PKG_UTIL", &Config.PkgUtil)
	utils.GetEnvString("REDIS_PKG_TIME", &Config.PkgTime)
	utils.GetEnvString("REDIS_OPCODE_PATH", &Config.CodePath)
	utils.GetEnvString("REDIS_MSG_PROTOCOL", &Config.MsgProtocol)
	// utils.GetEnvString("REDIS_PKG_WPB", &Config.PkgWPB)
	utils.GetEnvBool("REDIS_OP_ANOTHER", &Config.OpAnotherFile)
	utils.GetEnvBool("REDIS_CHECK_LUA", &Config.CheckLuaScript)
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
		"Import": func(pkg, alias string) string {
			return fmt.Sprintf("import %s \"%s\"", alias, pkg)
		},
		"UsePackage": func(pkg, alias string) string {
			return fmt.Sprintf("import %s \"%s\"", alias, pkg)
		},
	})

	// 解析前置的全部模板
	err = fs.WalkDir(redisOpTpls, ".", func(path string, d fs.DirEntry, err error) error {
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

		log.Println("regiester key:", k)

		ts := fmt.Sprintf(`{{define "%s"}} %s {{end}}`, k, strings.TrimSpace(string(data)))
		err = gen.Parse(ts)
		if err != nil {
			err = fmt.Errorf("parse template %s failed:%+v", k, err.Error())
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}

	err = gen.Parse(redisTpl)
	if err != nil {
		return err
	}

	keyarg.WTimePkg = Config.PkgTime

	switch Config.MsgProtocol {
	case "proto":
		if len(Config.PkgPB) == 0 {
			return fmt.Errorf("USE_PROTOBUF is true, but REDIS_PKG_PB is empty")
		}
	case "json":
		if len(Config.PkgJson) == 0 {
			return fmt.Errorf("USE_PROTOBUF is false, but REDIS_PKG_WPB is empty")
		}
	default:
		return fmt.Errorf("unknown message protocol %s", Config.MsgProtocol)
	}

	log.Println("init redis template success")

	return nil
}

func Generate(obj *RedisObject) (out []*tpl.BuildOutput, err error) {
	obj.Version = "0.0.1"
	obj.SvcPkg = filepath.Base(Config.PkgSvc)
	//obj.WPbPkg = filepath.Base(Config.PkgWPB)

	// genTpl, err := gen.Clone()
	// if err != nil {
	// 	return nil, err
	// }
	genTpl := gen

	genTpl.AddImportFunc(obj)

	for _, pkg := range []string{
		"context",
		"sync/atomic",
		"github.com/redis/go-redis/v9",
	} {
		obj.Import(pkg, "pkg")
	}

	obj.Prepare("rdconv", filepath.Join(Config.PkgUtil, "rdconv"))
	obj.Prepare("wtime", Config.PkgTime)
	obj.Prepare("util", Config.PkgUtil)

	obj.Prepare("proto", Config.PkgPB)
	obj.Prepare("json", Config.PkgJson)

	// obj.Import(Config.PkgUtil, "util.Builder")
	obj.Import(Config.PkgSvc, "RegisterDBName/GetDBLink")

	for _, arg := range obj.Args {
		for _, pkg := range arg.Imports() {
			if pkg == "" {
				continue
			}
			obj.Import(pkg, "keyargv")
		}
	}

	data, err := genTpl.Exec(obj)
	if err != nil {
		return nil, err
	}

	out = append(out, &tpl.BuildOutput{
		File: filepath.Join(Config.CodePath, strings.ToLower(obj.Package), strings.ToLower(obj.Name)+".redisop.go"),
		Data: data,
	})

	log.Println(obj.Package, obj.Name, "generate redis code success")

	return
}
