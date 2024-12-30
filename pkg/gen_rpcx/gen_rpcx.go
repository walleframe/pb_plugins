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
package gen_rpcx

import (
	_ "embed"
	"errors"
	"strings"

	"path/filepath"

	"github.com/walleframe/pb_plugins/pkg/tpl"
	"github.com/walleframe/pb_plugins/pkg/utils"
)

//go:embed rpcx.tpl
var mainTpl string

//go:embed stub.tpl
var stubTpl string

//go:embed gin.tpl
var ginTpl string

type GenerateEnv struct {
	PkgServer   string
	PkgClient   string
	PkgProtocol string
	PkgCtrl     string
	PkgStub     string
	PkgGin      string
	StubPath    string
	GinPath     string
}

var Config = &GenerateEnv{
	PkgServer:   "github.com/smallnest/rpcx/server",
	PkgClient:   "github.com/smallnest/rpcx/client",
	PkgProtocol: "github.com/smallnest/rpcx/protocol",
	PkgCtrl:     "",
	PkgStub:     "",
	PkgGin:      "",
	StubPath:    "pkg/gen/rpcx/",
	GinPath:     "pkg/gen/gin/",
}

func ParseConfigFromEnv() {
	// 如果环境变量设置了值, 读取作为为默认值. 优先使用传递的参数
	utils.GetEnvString("RPCX_PKG_SERVER", &Config.PkgServer)
	utils.GetEnvString("RPCX_PKG_CLIENT", &Config.PkgClient)
	utils.GetEnvString("RPCX_PKG_PROTOCOL", &Config.PkgProtocol)
	utils.GetEnvString("RPCX_PKG_CTRL", &Config.PkgCtrl)
	utils.GetEnvString("RPCX_PKG_STUB", &Config.PkgStub)
	utils.GetEnvString("RPCX_PKG_GIN", &Config.PkgGin)
	utils.GetEnvString("RPCX_STUB_PATH", &Config.StubPath)
	utils.GetEnvString("RPCX_GIN_PATH", &Config.GinPath)
}

var genService *tpl.GoTemplate
var genStub *tpl.GoTemplate
var genGin *tpl.GoTemplate

func InitTemplate() (err error) {
	genService = tpl.NewTemplate("rpcx")
	genService.Funcs(map[string]interface{}{})
	err = genService.Parse(mainTpl)
	if err != nil {
		return err
	}

	genStub = tpl.NewTemplate("rpcx_stub")
	genStub.Funcs(map[string]interface{}{})
	err = genStub.Parse(stubTpl)
	if err != nil {
		return err
	}

	genGin = tpl.NewTemplate("rpcx_gin")
	genGin.Funcs(map[string]interface{}{})
	err = genGin.Parse(ginTpl)
	if err != nil {
		return err
	}

	return nil
}

type File struct {
	tpl.GoObject

	RPCFile string

	PkgCtrl   string
	PkgGin    string
	Package   string
	GoPackage string

	Services []*Service

	Stub *Service
}

type Service struct {
	StubPkg  string
	StubCtrl string
	StubGin  string
	Name     string
	Methods  []*Method
}

type Method struct {
	Name    string
	RQ      string
	RS      string
	AppInfo bool
	AppRQ   bool
	Doc     tpl.Commets
}

func Generate(svc *File) (out []*tpl.BuildOutput, err error) {
	if len(svc.Services) < 1 {
		return nil, nil
	}
	svc.Version = "0.0.1"

	if svc.PkgCtrl == "" {
		svc.PkgCtrl = Config.PkgCtrl
	}
	if svc.PkgCtrl == "" {
		return nil, errors.New("pkg ctrl is empty")
	}
	if svc.RPCFile == "" {
		return nil, errors.New("rpc file is empty")
	}
	if svc.PkgGin == "" {
		svc.PkgGin = Config.PkgGin
	}

	svc.PkgCtrl = filepath.Base(svc.PkgCtrl)
	svc.PkgGin = filepath.Base(svc.PkgGin)

	o1, err := generateRPCX(svc)
	if err != nil {
		return nil, err
	}
	out = append(out, o1...)

	o2, err := generateStub(svc)
	if err != nil {
		return nil, err
	}
	out = append(out, o2...)

	if len(svc.PkgGin) > 0 {
		o3, err := generateGIN(svc)
		if err != nil {
			return nil, err
		}
		out = append(out, o3...)
	}

	return
}

func generateRPCX(svc *File) (out []*tpl.BuildOutput, err error) {
	genTpl := genService

	genTpl.AddImportFunc(svc)

	for _, pkg := range []string{
		"context",
	} {
		svc.Import(pkg, "pkg")
	}
	svc.ImportAlias("rpcx_server", Config.PkgServer, "rpcx")
	svc.ImportAlias("rpcx_client", Config.PkgClient, "rpcx")
	svc.Import(Config.PkgCtrl, "ctrl")
	svc.Prepare("protocol", Config.PkgProtocol)

	data, err := genTpl.Exec(svc)
	if err != nil {
		return nil, err
	}

	out = append(out, &tpl.BuildOutput{
		File: svc.RPCFile,
		Data: data,
	})
	return
}

func generateStub(svc *File) (out []*tpl.BuildOutput, err error) {
	genTpl := genStub

	genTpl.AddImportFunc(svc)

	for _, pkg := range []string{
		"context",
	} {
		svc.Import(pkg, "pkg")
	}
	// svc.Import(Config.PkgServer, "rpcx")
	svc.ImportAlias("rpcx_client", Config.PkgClient, "rpcx")
	svc.Import(Config.PkgStub, "stub")
	svc.Import(svc.GoPackage, "current")
	pkgName := filepath.Base(svc.GoPackage)

	for _, v := range svc.Services {
		v.StubCtrl = filepath.Base(Config.PkgStub)
		svc.Stub = v

		backup := v.Methods
		v.Methods = make([]*Method, 0, len(backup))
		for _, method := range backup {
			m := *method
			if !strings.Contains(m.RQ, ".") {
				m.RQ = pkgName + "." + m.RQ
			}
			if !strings.Contains(m.RS, ".") {
				m.RS = pkgName + "." + m.RS
			}
			v.Methods = append(v.Methods, &m)
		}

		data, err := genTpl.Exec(svc)
		v.Methods = backup
		if err != nil {
			return nil, err
		}

		out = append(out, &tpl.BuildOutput{
			File: filepath.Join(Config.StubPath, v.StubPkg, v.StubPkg+".stub.go"),
			Data: data,
		})
	}
	svc.Remove(svc.GoPackage)
	return
}

func generateGIN(svc *File) (out []*tpl.BuildOutput, err error) {
	genTpl := genGin

	genTpl.AddImportFunc(svc)

	for _, pkg := range []string{
		"context",
		"github.com/gin-gonic/gin",
	} {
		svc.Import(pkg, "pkg")
	}

	svc.Import(Config.PkgGin, "gin")
	svc.Import(svc.GoPackage, "current")
	pkgName := filepath.Base(svc.GoPackage)

	for _, v := range svc.Services {
		v.StubCtrl = filepath.Base(Config.PkgGin)
		svc.Stub = v

		backup := v.Methods
		v.Methods = make([]*Method, 0, len(backup))
		for _, method := range backup {
			m := *method
			if !strings.Contains(m.RQ, ".") {
				m.RQ = pkgName + "." + m.RQ
			}
			if !strings.Contains(m.RS, ".") {
				m.RS = pkgName + "." + m.RS
			}
			m.AppInfo = method.AppInfo
			m.AppRQ = method.AppRQ
			v.Methods = append(v.Methods, &m)
		}

		data, err := genTpl.Exec(svc)
		v.Methods = backup
		if err != nil {
			return nil, err
		}

		out = append(out, &tpl.BuildOutput{
			File: filepath.Join(Config.GinPath, v.StubGin, v.StubGin+".go"),
			Data: data,
		})
	}
	svc.Remove(svc.GoPackage)
	return
}
