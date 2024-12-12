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
package tpl

import (
	"bytes"
	"fmt"
	"log"
	"strings"
	"text/template"

	"github.com/walleframe/pb_plugins/pkg/utils"
	"go.uber.org/multierr"
)

// const (
// 	ProtoGoPkg = "proto.gopkg"
// )

type GoObject struct {
	// innerImportPkgs 导入包
	innerImportPkgs []string
	// ToolName 工具名
	ToolName string
	// Version 版本
	Version string
}

func (obj *GoObject) Import(pkg, fun string) (_ string) {
	for _, v := range obj.innerImportPkgs {
		if v == pkg {
			return
		}
	}
	obj.innerImportPkgs = append(obj.innerImportPkgs, pkg)
	return
}

func (g *GoObject) customImport() string {
	buf := strings.Builder{}
	buf.Grow(len(g.innerImportPkgs) * 30)
	for _, v := range g.innerImportPkgs {
		buf.WriteByte('\n')
		buf.WriteByte('\t')
		buf.WriteString(`"`)
		buf.WriteString(v)
		buf.WriteByte('"')
	}
	return buf.String()
}

// TODO: fix ImportDepend
// func (g *GoObject) ImportDepend(prog *buildpb.FileDesc, depend map[string]*buildpb.FileDesc) (err error) {
// 	// import others
// 	for _, imp := range prog.Imports {
// 		dep, ok := depend[imp.File]
// 		if !ok {
// 			err = multierr.Append(err, fmt.Errorf("%s import %s, but not found %s", prog.File, imp.File, imp.File))
// 			continue
// 		}
// 		pkg, ok := dep.Options.GetStringCheck(ProtoGoPkg)
// 		if !ok {
// 			err = multierr.Append(err, fmt.Errorf("%s import %s, but %s not set '%s' option", prog.File, imp.File, imp.File, ProtoGoPkg))
// 			continue
// 		}
// 		g.Import(pkg, "import others")
// 	}
// 	return
// }

type GoTemplate struct {
	*template.Template
	// 代码格式化
	formatCode FormatFunc
}

func NewTemplate(name string, formatFunc ...FormatFunc) (tpl *GoTemplate) {
	tpl = &GoTemplate{}
	if len(formatFunc) > 0 {
		tpl.formatCode = formatFunc[0]
	} else {
		tpl.formatCode = GoFormat2
	}
	tpl.Template = template.New(name).Funcs(DefaultFuncMap)
	return
}

func (tpl *GoTemplate) Clone() (*GoTemplate, error) {
	inner, err := tpl.Template.Clone()
	return &GoTemplate{
		Template:   inner,
		formatCode: tpl.formatCode,
	}, err
}

func (tpl *GoTemplate) AddImportFunc(obj interface {
	Import(pkg, fun string) (_ string)
}) {
	tpl.Funcs(template.FuncMap{
		"Import": obj.Import,
	})
}

func (tpl *GoTemplate) PrepareTemplate(preTpl map[string]string) (err error) {
	// 解析前置的全部模板
	for k, v := range preTpl {
		ts := fmt.Sprintf(`{{define "%s"}} %s {{end}}`, k, strings.TrimSpace(v))
		err = tpl.Parse(ts)
		if err != nil {
			err = multierr.Append(err, fmt.Errorf("parse template %s failed:%+v", k, err.Error()))
		}
	}
	return
}

func (tpl *GoTemplate) Parse(v string) error {
	_, err := tpl.Template.Parse(v)
	return err
}

func (tpl *GoTemplate) Exec(obj interface{}) (data []byte, err error) {
	buf := bytes.NewBuffer(make([]byte, 0, 4096))
	err = tpl.Execute(buf, obj)
	if err != nil {
		return nil, err
	}
	if !bytes.Contains(buf.Bytes(), []byte("$Import-Packages$")) {
		return buf.Bytes(), nil
	}

	impObj, ok := obj.(interface {
		customImport() string
	})
	if !ok {
		return buf.Bytes(), nil
	}

	bdata := bytes.Replace(buf.Bytes(), []byte("$Import-Packages$"), []byte(fmt.Sprintf(`
import (
	%s
)`, impObj.customImport())), 1)

	data, err = tpl.formatCode(bdata)
	if err != nil {
		printWithLine(bdata)
		return nil, fmt.Errorf("format code failed,[%w]", err)
	}
	return
}

func printWithLine(data []byte) {
	for k, v := range bytes.Split(data, []byte{'\n'}) {
		log.Printf("%d\t%s\n", k+1, string(v))
	}
}

var DefaultFuncMap = template.FuncMap{}

func init() {
	DefaultFuncMap["Doc"] = func(doc *Commets) string {
		if doc == nil {
			return ""
		}
		buf := strings.Builder{}
		for _, v := range doc.Doc {
			buf.WriteString(v)
			buf.WriteByte('\n')
		}
		if len(doc.TailDoc) > 0 {
			buf.WriteString(doc.TailDoc)
			buf.WriteByte('\n')
		}
		return buf.String()
	}
	DefaultFuncMap["Title"] = utils.Title
	DefaultFuncMap["Comma"] = func(k int) string {
		if k == 0 {
			return ""
		}
		return ","
	}
	DefaultFuncMap["BackQuote"] = func(v string) string {
		return "`" + v + "`"
	}
}

type Commets struct {
	Doc     []string
	TailDoc string
}

type Import struct {
	Pkg string
	Fun string
}

type BuildOutput struct {
	File string
	Data []byte
}
