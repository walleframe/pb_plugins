package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"go.uber.org/multierr"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/reflect/protoreflect"

	"github.com/walleframe/pb_plugins/pkg/gen_zerolog"
)

const genGoDocURL = "https://github.com/walleframe/pb_plugins/blob/main/cmd/protoc-gen-zerolog/readme.org"

var (
	Version = "0.0.1"
)

func main() {
	if len(os.Args) == 2 && os.Args[1] == "--version" {
		fmt.Fprintf(os.Stdout, "%v %v\n", filepath.Base(os.Args[0]), Version)
		os.Exit(0)
	}
	if len(os.Args) == 2 && os.Args[1] == "--help" {
		fmt.Fprintf(os.Stdout, "See "+genGoDocURL+" for usage information.\n")
		os.Exit(0)
	}

	// 从环境变量中读取配置
	gen_zerolog.ParseConfigFromEnv()

	// // 外部命令行参数(比环境变量优先级高)
	// var flags pflag.FlagSet
	// cfg := gen_zerolog.Config
	// _ = cfg
	// flags.StringVar(&cfg.SvcPkg, "pkg_svc", cfg.SvcPkg, "service package name")
	protogen.Options{
		//ParamFunc: flags.Set,
	}.Run(func(gen *protogen.Plugin) (err error) {
		err = gen_zerolog.InitTemplate()
		if err != nil {
			log.Println("gen_mysql InitTemplate error: ", err)
			return
		}

		for _, f := range gen.Files {
			if !f.Generate {
				continue
			}
			err = genZerologFile(gen, f)
			if err != nil {
				return err
			}
		}
		//gen.SupportedFeatures = uint64(pluginpb.CodeGeneratorResponse_FEATURE_PROTO3_OPTIONAL)
		return
	})
}

func genZerologFile(gen *protogen.Plugin, file *protogen.File) (err error) {
	filename := file.GeneratedFilenamePrefix + ".zlog.go"

	zlog := &gen_zerolog.ZlogMessages{
		File: filename,
	}
	zlog.ToolName = "protoc-gen-zerolog"
	zlog.Package = string(file.GoPackageName)
	// 扫描要生成的消息
	genZerologObject(file, zlog)
	// generate
	data, err := gen_zerolog.Generate(zlog)
	if err != nil {
		err = multierr.Append(err, fmt.Errorf("generate zerolog file %s error: %v", filename, err))
		return
	}
	// write response
	for _, d := range data {
		file := gen.NewGeneratedFile(d.File, "")
		file.Write(d.Data)
	}
	return
}

func genZerologObject(file *protogen.File, zlog *gen_zerolog.ZlogMessages) {
	for _, m := range file.Messages {
		genZerologMessage(m, zlog)

	}
}

func genZerologMessage(m *protogen.Message, zlog *gen_zerolog.ZlogMessages) {
	if m.Desc.IsMapEntry() {
		return
	}
	msg := &gen_zerolog.Message{
		Name: m.GoIdent.GoName,
		Type: string(m.GoIdent.GoName),
	}
	for _, f := range m.Fields {
		genZerologField(f, msg, zlog)
	}
	zlog.Msgs = append(zlog.Msgs, msg)

	for _, sub := range m.Messages {
		genZerologMessage(sub, zlog)
	}
}

func genZerologField(f *protogen.Field, msg *gen_zerolog.Message, zlog *gen_zerolog.ZlogMessages) {
	if f.Desc.IsWeak() {
		return
	}
	field := &gen_zerolog.Field{
		Key:  string(f.Desc.Name()),
		Name: f.GoName,
	}
	field.Value = func(obj string) string {
		return fmt.Sprintf("%s.%s", obj, field.Name)
	}
	if f.Desc.IsMap() {
		field.Func = "Object"
		field.Map = &gen_zerolog.MapField{}
		// Type: string(f.Desc.MapKey().Name()),
		//KeyFunc: getFieldFunc(f.Message.Fields[0]),
		field.Map.Type = getMapFieldType(f)

		field.Map.KeyFunc = getMapKeyFunc(f.Message.Fields[0], zlog)

		field.Map.ValFunc = getFieldFunc(f.Message.Fields[1], field, zlog)
	} else {
		field.Func = getFieldFunc(f, field, zlog)
	}
	msg.Fields = append(msg.Fields, field)
}

func getMapFieldType(f *protogen.Field) string {
	key := f.Message.Fields[0]
	value := f.Message.Fields[1]
	typ := "map[" + getFieldGoType(key.Desc.Kind()) + "]"
	if value.Desc.Kind() == protoreflect.MessageKind {
		if value.Message.GoIdent.GoImportPath == f.Parent.GoIdent.GoImportPath {
			typ += "*" + value.Message.GoIdent.GoName
		} else {
			typ += "*" + filepath.Base(string(value.Message.GoIdent.GoImportPath)) + "." + value.Message.GoIdent.GoName
		}
	} else {
		typ += getFieldGoType(value.Desc.Kind())
	}
	return typ
}

func getMapKeyFunc(f *protogen.Field, zlog *gen_zerolog.ZlogMessages) func(k string) string {

	if f.Desc.Kind() != protoreflect.StringKind {
		zlog.Import("strconv", "FormatInt/FormatUint")
	}
	switch f.Desc.Kind() {
	case protoreflect.StringKind:
		return func(k string) string {
			return k
		}
	case protoreflect.EnumKind:
		return func(k string) string {
			return fmt.Sprintf("%s.String()", k)
		}
	case protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Sfixed32Kind:
		return func(k string) string {
			return fmt.Sprintf("strconv.FormatInt(int64(%s), 10)", k)
		}
	case protoreflect.Uint32Kind, protoreflect.Fixed32Kind:
		return func(k string) string {
			return fmt.Sprintf("strconv.FormatUint(uint64(%s), 10)", k)
		}
	case protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Sfixed64Kind:
		return func(k string) string {
			return fmt.Sprintf("strconv.FormatInt(%s, 10)", k)
		}
	case protoreflect.Uint64Kind, protoreflect.Fixed64Kind:
		return func(k string) string {
			return fmt.Sprintf("strconv.FormatUint(%s, 10)", k)
		}
	case protoreflect.FloatKind, protoreflect.DoubleKind, protoreflect.BytesKind, protoreflect.MessageKind, protoreflect.GroupKind, protoreflect.BoolKind:
		// protobuf 不支持这些类型作为map的key
		return nil
	}
	return nil
}

func getFieldFunc(f *protogen.Field, fcfg *gen_zerolog.Field, zlog *gen_zerolog.ZlogMessages) string {
	if fcfg.Map != nil {
		fcfg.Value = func(obj string) string {
			return obj
		}
	}
	if f.Desc.IsList() {
		switch f.Desc.Kind() {
		case protoreflect.StringKind:
			return "Strs"
		case protoreflect.BytesKind:
			fcfg.Value = func(obj string) string {
				zlog.Import("encoding/base64", "EncodeToString")
				return fmt.Sprintf(`func()(vals []string){
	vals = make([]string, 0, len(%[1]s.%[2]s))
	for _,v := range %[1]s.%[2]s {
		vals = append(vals, base64.StdEncoding.EncodeToString(v))
	}
	return vals
}()`, obj, fcfg.Name)
			}
			return "Strs"
		case protoreflect.BoolKind:
			return "Bools"
		case protoreflect.EnumKind:
			fcfg.Value = func(obj string) string {
				return fmt.Sprintf(`func()(vals []string){
	vals = make([]string, 0, len(%[1]s.%[2]s))
	for _,v := range %[1]s.%[2]s {
		vals = append(vals, v.String())
	}
	return vals
}()`, obj, fcfg.Name)
			}
			return "Strs"
		case protoreflect.MessageKind, protoreflect.GroupKind:
			GoName := f.Message.GoIdent.GoName
			if f.Message.GoIdent.GoImportPath != f.Parent.GoIdent.GoImportPath {
				GoName = filepath.Base(string(f.Message.GoIdent.GoImportPath)) + "." + GoName
			}
			//log.Println("--> ", GoName, f.Message.GoIdent.GoImportPath, f.Parent.GoIdent.GoImportPath)
			zlog.Import(string(f.Message.GoIdent.GoImportPath), f.Message.GoIdent.GoName)
			fcfg.Value = func(obj string) string {
				return fmt.Sprintf(`%[3]sArray(%[1]s.%[2]s)`, obj, fcfg.Name, GoName)
			}
			return "Array"
		case protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Sfixed32Kind:
			return "Ints32"
		case protoreflect.Uint32Kind, protoreflect.Fixed32Kind:
			return "Uints32"
		case protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Sfixed64Kind:
			return "Ints64"
		case protoreflect.Uint64Kind, protoreflect.Fixed64Kind:
			return "Uints64"
		case protoreflect.FloatKind:
			return "Floats32"
		case protoreflect.DoubleKind:
			return "Floats64"
		}
		return ""
	}
	switch f.Desc.Kind() {
	case protoreflect.StringKind:
		return "Str"
	case protoreflect.BytesKind:
		if fcfg.Map != nil {
			fcfg.Value = func(obj string) string {
				zlog.Import("encoding/base64", "EncodeToString")
				return fmt.Sprintf("base64.StdEncoding.EncodeToString(%s)", obj)
			}
		} else {
			fcfg.Value = func(obj string) string {
				zlog.Import("encoding/base64", "EncodeToString")
				return fmt.Sprintf("base64.StdEncoding.EncodeToString(%s.%s)", obj, fcfg.Name)
			}
		}
		return "Str"
	case protoreflect.BoolKind:
		return "Bool"
	case protoreflect.EnumKind:
		return "Stringer"
	case protoreflect.MessageKind, protoreflect.GroupKind:
		if f.Oneof != nil {
			fcfg.Value = func(obj string) string {
				return fmt.Sprintf("%s.Get%s()", obj, fcfg.Name)
			}
		}
		return "Object"
	case protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Sfixed32Kind:
		return "Int32"
	case protoreflect.Uint32Kind, protoreflect.Fixed32Kind:
		return "Uint32"
	case protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Sfixed64Kind:
		return "Int64"
	case protoreflect.Uint64Kind, protoreflect.Fixed64Kind:
		return "Uint64"
	case protoreflect.FloatKind:
		return "Float32"
	case protoreflect.DoubleKind:
		return "Float64"
	}
	return ""
}

func getFieldGoType(kind protoreflect.Kind) string {
	switch kind {
	case protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Sfixed32Kind:
		return "int32"
	case protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Sfixed64Kind:
		return "int64"
	case protoreflect.Uint32Kind, protoreflect.Fixed32Kind:
		return "uint32"
	case protoreflect.Uint64Kind, protoreflect.Fixed64Kind:
		return "uint64"
	case protoreflect.FloatKind:
		return "float32"
	case protoreflect.DoubleKind:
		return "float64"
	case protoreflect.StringKind:
		return "string"
	case protoreflect.BytesKind:
		return "[]byte"
	case protoreflect.BoolKind:
		return "bool"
	}
	return ""
}
