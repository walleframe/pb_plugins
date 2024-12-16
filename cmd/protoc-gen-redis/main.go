package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/pflag"
	"go.uber.org/multierr"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"

	"github.com/walleframe/pb_plugins/extend/redis"
	"github.com/walleframe/pb_plugins/pkg/gen_redis"
	"github.com/walleframe/pb_plugins/pkg/gen_redis/keyarg"
	"github.com/walleframe/pb_plugins/pkg/tpl"
	"github.com/walleframe/pb_plugins/pkg/utils"

	lua "github.com/yuin/gopher-lua"
)

const genGoDocURL = "https://github.com/walleframe/pb_plugins/blob/main/cmd/protoc-gen-redis/readme.org"

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
	gen_redis.ParseConfigFromEnv()

	// 外部命令行参数(比环境变量优先级高)
	var flags pflag.FlagSet
	cfg := gen_redis.Config
	flags.StringVar(&cfg.CodePath, "code_path", cfg.CodePath, "code path")
	flags.StringVar(&cfg.PkgSvc, "pkg_svc", cfg.PkgSvc, "service package name")
	flags.StringVar(&cfg.PkgUtil, "pkg_util", cfg.PkgUtil, "util package name")
	//.flags.StringVar(&cfg.PkgWPB, "pkg_wproto", cfg.PkgWPB, "wproto messsage package")
	flags.StringVar(&cfg.PkgPB, "pkg_protobuf", cfg.PkgPB, "ptotobuf package")
	flags.StringVar(&cfg.PkgTime, "pkg_time", cfg.PkgTime, "time package")
	flags.StringVar(&cfg.MsgProtocol, "msg_protocol", cfg.MsgProtocol, "message protocol")
	flags.BoolVar(&cfg.OpAnotherFile, "op_another", cfg.OpAnotherFile, "operate message in another file")
	flags.BoolVar(&cfg.CheckLuaScript, "check_lua", cfg.CheckLuaScript, "check lua script")
	protogen.Options{
		ParamFunc: flags.Set,
	}.Run(func(gen *protogen.Plugin) (err error) {
		err = gen_redis.InitTemplate()
		if err != nil {
			log.Println("gen_redis InitTemplate error: ", err)
			return
		}

		for _, f := range gen.Files {
			if !f.Generate {
				continue
			}
			err = genRedisFile(gen, f)
			if err != nil {
				return
			}
		}
		//gen.SupportedFeatures = uint64(pluginpb.CodeGeneratorResponse_FEATURE_PROTO3_OPTIONAL)
		return
	})
}

func genRedisFile(gen *protogen.Plugin, file *protogen.File) (err error) {
	for _, msg := range file.Messages {
		err = genRedisOperation(gen, msg)
		if err != nil {
			errs := multierr.Errors(err)
			log.Printf("generate file:%s msg:%s failed:\n", file.Desc.Path(), msg.Desc.Name())
			for _, e := range errs {
				log.Println("\t", e)
			}

			err = errors.New("generate failed")
			return
		}
	}

	return
}

func genRedisOperation(gen *protogen.Plugin, msg *protogen.Message) (err error) {
	getOptString := func(opt protoreflect.ExtensionType, def string) (v string) {
		v = def
		proto.RangeExtensions(msg.Desc.Options(), func(et protoreflect.ExtensionType, a any) bool {
			if et.TypeDescriptor().FullName() == opt.TypeDescriptor().FullName() {
				v = a.(string)
				return false
			}
			return true
		})
		return
	}

	// int32
	getOptInt := func(opt protoreflect.ExtensionType, def int32) (v int32) {
		v = def
		proto.RangeExtensions(msg.Desc.Options(), func(et protoreflect.ExtensionType, a any) bool {
			if et.TypeDescriptor().FullName() == opt.TypeDescriptor().FullName() {
				v = a.(int32)
				return false
			}
			return true
		})
		return
	}
	// bool
	getOptBool := func(opt protoreflect.ExtensionType, def bool) (v bool) {
		v = def
		proto.RangeExtensions(msg.Desc.Options(), func(et protoreflect.ExtensionType, a any) bool {
			if et.TypeDescriptor().FullName() == opt.TypeDescriptor().FullName() {
				v = a.(bool)
				return false
			}
			return true
		})
		return
	}
	//
	hasOpt := func(opt protoreflect.ExtensionType) (v bool) {
		proto.RangeExtensions(msg.Desc.Options(), func(et protoreflect.ExtensionType, a any) bool {
			if et.TypeDescriptor().FullName() == opt.TypeDescriptor().FullName() {
				v = true
				return false
			}
			return true
		})
		return
	}

	// key
	redisKey := getOptString(redis.E_Key, "")
	if redisKey == "" {
		// 没有设置key,不生成
		return
	}
	redisKeyArgs, err := keyarg.MatchKey(redisKey, nil)
	if err != nil {
		err = multierr.Append(err, fmt.Errorf("parse key %s in msg:%s failed:%+v", redisKey, msg.Desc.Name(), err))
		return
	}

	// 禁止同时设置 redis.pb 和 redis.json
	if getOptBool(redis.E_Json, false) && getOptBool(redis.E_Pb, false) {
		err = multierr.Append(err, fmt.Errorf("msg:%s redis.pb and redis.json can't be set at the same time", msg.Desc.Name()))
		return
	}

	usePb := getOptBool(redis.E_Pb, false)
	if !usePb && hasOpt(redis.E_Json) {
		usePb = !getOptBool(redis.E_Json, false)
	} else if gen_redis.Config.MsgProtocol == "proto" {
		// 默认使用proto
		usePb = true
	}

	// log.Printf("msg:%s redis.pb:%t redis.json:%t env:%t usePb:%t\n", msg.Desc.Name(),
	// 	getOptBool(redis.E_Pb, false), getOptBool(redis.E_Json, false), gen_redis.Config.MsgProtocol == "proto", usePb,
	// )

	redisType := strings.ToLower(strings.TrimSpace(getOptString(redis.E_Type, "")))

	//
	fname := msg.Desc.ParentFile().Path()
	obj := &gen_redis.RedisObject{
		Package: getOptString(redis.E_OpPackage, string(strings.TrimSuffix(filepath.Base(fname), filepath.Ext(fname)))),
		Name:    msg.GoIdent.GoName,
		Doc: &tpl.Commets{
			Doc:     []string{string(msg.Comments.Leading)},
			TailDoc: string(msg.Comments.Trailing),
		},
		Args:     redisKeyArgs,
		KeySize:  int(getOptInt(redis.E_Size, 64)),
		TypeKeys: true,
		// Scripts:    []*gen_redis.RedisScript{},
		// TypeString: &gen_redis.RedisTypeString{},
		// TypeHash:   &gen_redis.RedisTypeHash{},
		// TypeSet:    &gen_redis.RedisTypeSet{},
		// TypeZSet:   &gen_redis.RedisTypeZSet{},
		// Lock:       false,
	}
	obj.ToolName = "protoc-gen-redis"

	// redisField := getOptString(redis.E_Field, "")
	// redisValue := getOptString(redis.E_Value, "")
	// redisMember := getOptString(redis.E_Member, "")
	// redisScripts := proto.GetExtension(msg.Desc.Options(), redis.E_Script).([]*redis.RedisScript)

	if strings.HasPrefix(redisType, "!") {
		obj.TypeKeys = false
		redisType = strings.TrimSpace(strings.TrimPrefix(redisType, "!"))
	}
	switch redisType {
	case "lock":
		obj.Lock = true
	case "string":
		err = multierr.Append(err, analyseTypeString(msg, obj, usePb))
	case "hash":
		err = multierr.Append(err, analyseTypeHash(msg, obj, usePb))
	case "set":
		err = multierr.Append(err, analyseTypeSet(msg, obj, usePb))
	case "zset":
		err = multierr.Append(err, analyseTypeZSet(msg, obj, usePb))
	case "":
		if !obj.TypeKeys {
			err = multierr.Append(err, fmt.Errorf("msg:%s redis type config invalid", msg.Desc.Name()))
			return
		}
	default:
		err = multierr.Append(err, fmt.Errorf("msg:%s redis type not support [%s]", msg.Desc.Name(), redisType))
		return
	}
	err = multierr.Append(err, analyseRedisScripts(msg, obj))
	if err != nil {
		return
	}
	// generate
	data, err := gen_redis.Generate(obj)
	if err != nil {
		err = multierr.Append(err, fmt.Errorf("generate redis %s error: %v", obj.Name, err))
		return
	}
	// write response
	for _, d := range data {
		file := gen.NewGeneratedFile(d.File, "")
		file.Write(d.Data)
	}
	return
}

func analyseTypeString(msg *protogen.Message, obj *gen_redis.RedisObject, usePb bool) (err error) {
	getOptBool := func(opt protoreflect.ExtensionType, def bool) (v bool) {
		v = def
		proto.RangeExtensions(msg.Desc.Options(), func(et protoreflect.ExtensionType, a any) bool {
			if et.TypeDescriptor().FullName() == opt.TypeDescriptor().FullName() {
				v = a.(bool)
				return false
			}
			return true
		})
		return
	}
	if !gen_redis.Config.OpAnotherFile && !getOptBool(redis.E_OpField, false) {
		// 未设置op_field,操作当前结构体
		obj.TypeString = &gen_redis.RedisTypeString{
			Custom:   true,
			Protobuf: usePb,
			Json:     !usePb,
			//WProto:   !gen_redis.Config.UseProtobuf,
			// 对应的Go结构体名
			Type: filepath.Base(string(msg.GoIdent.GoImportPath)) + "." + msg.GoIdent.GoName,
		}
		obj.Import(string(msg.GoIdent.GoImportPath), "msg")

		return
	}
	// 最多只能有一个字段
	if len(msg.Fields) > 1 {
		return errors.New("set redis.op_field=true ,redis-string type fields too many. max 1 fields")
	}

	opt := &gen_redis.RedisTypeString{}
	obj.TypeString = opt
	// 无数据设置,查看是否生成通用接口
	if len(msg.Fields) < 1 {
		// 生成protobuf接口
		if usePb {
			obj.Import(gen_redis.Config.PkgPB, "Marshal/Unmarshal")
			opt.Protobuf = true
			return
		}
		obj.Import(gen_redis.Config.PkgJson, "Marshal/Unmarshal")
		opt.Json = true
		//obj.Import(gen_redis.Config.PkgWPB, "MarshalObject/UnmarshalObject")
		//opt.WProto = true
		return
	}
	fieldType := msg.Fields[0].Desc
	if fieldType.IsList() || fieldType.IsMap() {
		return fmt.Errorf("msg:%s redis-string type generation not support array or map type", msg.Desc.Name())
	}

	switch fieldType.Kind() {
	case protoreflect.MessageKind, protoreflect.GroupKind:
		opt.Custom = true
		msg := msg.Fields[0].Message
		opt.Type = filepath.Base(string(msg.GoIdent.GoImportPath)) + "." + msg.GoIdent.GoName
		obj.Import(string(msg.GoIdent.GoImportPath), "msg")
		// 自定义类型
		if usePb {
			obj.Import(gen_redis.Config.PkgPB, "Marshal/Unmarshal")
			opt.Protobuf = true
		} else {
			// obj.Import(gen_redis.Config.PkgWPB, "MarshalObject/UnmarshalObject")
			// opt.WProto = true
			obj.Import(gen_redis.Config.PkgJson, "Marshal/Unmarshal")
			opt.Json = true
		}
	case protoreflect.Int32Kind, protoreflect.Int64Kind, protoreflect.Sint32Kind, protoreflect.Sint64Kind, protoreflect.Sfixed32Kind, protoreflect.Sfixed64Kind:
		opt.Number = true
		opt.Signed = true
		opt.Type = getFieldGoType(fieldType.Kind())
	case protoreflect.Uint32Kind, protoreflect.Uint64Kind, protoreflect.Fixed32Kind, protoreflect.Fixed64Kind:
		opt.Number = true
		opt.Type = getFieldGoType(fieldType.Kind())
	case protoreflect.FloatKind, protoreflect.DoubleKind:
		opt.Float = true
		if fieldType.Kind() == protoreflect.FloatKind {
			opt.Type = "float32"
		} else {
			opt.Type = "float64"
		}
	case protoreflect.StringKind:
		opt.String = true
		opt.Type = getFieldGoType(fieldType.Kind())
	default:
		return fmt.Errorf("msg:%s redis-string type generation not support %s basic type",
			msg.Desc.Name(), fieldType.Kind(),
		)
	}
	return
}

func analyseTypeHash(msg *protogen.Message, obj *gen_redis.RedisObject, usePb bool) (err error) {
	hash := &gen_redis.RedisTypeHash{}
	obj.TypeHash = hash
	getOptBool := func(opt protoreflect.ExtensionType, def bool) (v bool) {
		v = def
		proto.RangeExtensions(msg.Desc.Options(), func(et protoreflect.ExtensionType, a any) bool {
			if et.TypeDescriptor().FullName() == opt.TypeDescriptor().FullName() {
				v = a.(bool)
				return false
			}
			return true
		})
		return
	}
	getOptString := func(opt protoreflect.ExtensionType, def string) (v string) {
		v = def
		proto.RangeExtensions(msg.Desc.Options(), func(et protoreflect.ExtensionType, a any) bool {
			if et.TypeDescriptor().FullName() == opt.TypeDescriptor().FullName() {
				v = a.(string)
				return false
			}
			return true
		})
		return
	}

	opField := getOptString(redis.E_Field, "")
	opValue := getOptString(redis.E_Value, "")
	if opField != "" && opValue != "" {
		dynamic := &gen_redis.RedisHashDynamic{
			GenMap: false,
		}
		hash.HashDynamic = dynamic
		// 拼接string做field
		dynamic.FieldArgs, err = keyarg.MatchGoTypes(opField, nil)
		if err != nil {
			return fmt.Errorf("redis-hash analyse redis.field failed.%v", err)
		}
		// 拼接string做value
		dynamic.ValueArgs, err = keyarg.MatchGoTypes(opValue, nil)
		if err != nil {
			return fmt.Errorf("redis-hash analyse redis.value failed.%v", err)
		}
		// 检测两组参数内是否有重复的参数名称
		err = keyarg.CheckArgNameConflict(dynamic.FieldArgs, dynamic.ValueArgs)
		if err != nil {
			return fmt.Errorf("redis-hash analyse redis.field and redis.value conflict.%v", err)
		}
		return
	}
	// 分析hash object,展开所有字段.字段名作为field,字段值作为value
	hashObject := func(msg *protogen.Message) (err error) {
		hash.HashObject = &gen_redis.RedisHashObject{
			Name:    utils.Title(msg.GoIdent.GoName),
			Type:    filepath.Base(string(msg.GoIdent.GoImportPath)) + "." + msg.GoIdent.GoName,
			Fields:  make([]*gen_redis.RedisGenType, 0, len(msg.Fields)),
			HGetAll: true,
		}
		obj.Import(string(msg.GoIdent.GoImportPath), "msg")
		for _, field := range msg.Fields {
			if field.Desc.IsList() || field.Desc.IsMap() || field.Desc.Kind() == protoreflect.MessageKind ||
				field.Desc.Kind() == protoreflect.GroupKind {
				return fmt.Errorf("redis-hash type message field [%s.%s] is not basic type", msg.Desc.Name(), field.Desc.Name())
			}
			field := &gen_redis.RedisGenType{
				Name:      string(field.Desc.Name()),
				Type:      getFieldGoType(field.Desc.Kind()),
				Number:    isNumber(field.Desc.Kind()),
				RedisFunc: getRedisFunc(field.Desc.Kind()),
			}
			hash.HashObject.Fields = append(hash.HashObject.Fields, field)
		}
		return nil
	}
	// 未指定操作另一个文件
	if !gen_redis.Config.OpAnotherFile {
		if !getOptBool(redis.E_OpField, false) {
			// 并且未设置op_field,操作当前结构体
			hashObject(msg)
			return
		}
	}

	switch len(msg.Fields) {
	case 1:
		//field := msg.Fields[0]
		fieldType := msg.Fields[0].Desc
		switch {
		case opField == "" && opValue == "":
			// 未设置field和value,必须是message类型
			if fieldType.IsList() || fieldType.IsMap() || fieldType.Kind() != protoreflect.MessageKind {
				return fmt.Errorf("redis-hash type message field [%s.%s] is not message type", msg.Desc.Name(), fieldType.Name())
			}
			// 展开所有字段.字段名作为field,字段值作为value
			err = hashObject(msg.Fields[0].Message)
			if err != nil {
				return err
			}
			obj.Import(string(msg.Fields[0].Message.GoIdent.GoImportPath), "msg")
		case opField != "":
			dynamic := &gen_redis.RedisHashDynamic{
				GenMap: false,
			}
			hash.HashDynamic = dynamic
			// 拼接string做field
			dynamic.FieldArgs, err = keyarg.MatchGoTypes(opField, nil)
			if err != nil {
				return fmt.Errorf("redis-hash analyse redis.field failed.%v", err)
			}
			dynamic.Value = &gen_redis.RedisGenType{
				Name:      utils.PascalToSnake(msg.Fields[0].GoName),
				Type:      getFieldGoType(fieldType.Kind()),
				Number:    isNumber(fieldType.Kind()),
				RedisFunc: getRedisFunc(fieldType.Kind()),
			}
		case opValue != "":
			dynamic := &gen_redis.RedisHashDynamic{}
			hash.HashDynamic = dynamic
			// 拼接string做value
			dynamic.ValueArgs, err = keyarg.MatchGoTypes(opValue, nil)
			if err != nil {
				return fmt.Errorf("redis-hash analyse redis.value failed.%v", err)
			}
			if fieldType.IsList() || fieldType.IsMap() || fieldType.Kind() == protoreflect.MessageKind {
				return fmt.Errorf("redis-hash type message field [%s.%s] is not basic type", msg.Desc.Name(), fieldType.Name())
			}
			dynamic.Field = &gen_redis.RedisGenType{
				Name:      utils.PascalToSnake(msg.Fields[0].GoName),
				Type:      getFieldGoType(fieldType.Kind()),
				Number:    isNumber(fieldType.Kind()),
				RedisFunc: getRedisFunc(fieldType.Kind()),
			}
		}

	case 2:
		fieldType1 := msg.Fields[0].Desc
		fieldType2 := msg.Fields[1].Desc
		if fieldType1.IsList() || fieldType1.IsMap() || fieldType1.Kind() == protoreflect.MessageKind {
			return fmt.Errorf("redis-hash type message field [%s.%s] is not basic type", msg.Desc.Name(), fieldType1.Name())
		}
		dynamic := &gen_redis.RedisHashDynamic{
			GenMap: true,
		}
		hash.HashDynamic = dynamic
		dynamic.Field = &gen_redis.RedisGenType{
			Name:      utils.PascalToSnake(msg.Fields[0].GoName),
			Type:      getFieldGoType(fieldType1.Kind()),
			Number:    isNumber(fieldType1.Kind()),
			RedisFunc: getRedisFunc(fieldType1.Kind()),
		}
		if fieldType2.IsList() || fieldType2.IsMap() {
			return fmt.Errorf("redis-hash type message value [%s.%s] not support array or map type", msg.Desc.Name(), fieldType2.Name())
		}
		dynamic.Value = &gen_redis.RedisGenType{
			Name:      utils.PascalToSnake(msg.Fields[1].GoName),
			Type:      getFieldGoType(fieldType2.Kind()),
			Number:    isNumber(fieldType2.Kind()),
			RedisFunc: getRedisFunc(fieldType2.Kind()),
		}
		if fieldType2.Kind() == protoreflect.MessageKind {
			dynamic.Value.Type = filepath.Base(string(msg.Fields[1].Message.GoIdent.GoImportPath)) + "." + msg.Fields[1].Message.GoIdent.GoName

			//log.Printf("msgFrom:%s msgTo:%s usePB:%t\n", msg.Desc.Name(), dynamic.Value.Type, usePb)
			if usePb {
				obj.Import(gen_redis.Config.PkgPB, "Marshal/Unmarshal")
				dynamic.Value.MarshalPkg = "proto"
			} else {
				obj.Import(gen_redis.Config.PkgJson, "Marshal/Unmarshal")
				dynamic.Value.MarshalPkg = "json"
			}
		}

	case 3: // NOTE: hval,hfields 生成需要过滤,有需要再改吧. 先禁掉功能了.
		return fmt.Errorf("redis-hash type fields count 3 not support now,need modify hvals/hfields/range functions")
	// 	err = hashObject(msg.Fields[0].Type, obj.TypeHash)
	// 	if err != nil {
	// 		return err
	// 	}
	// 	err = hashDynamic(msg.Fields[1].Type, msg.Fields[2].Type, obj.TypeHash)
	// 	if err != nil {
	// 		return err
	// 	}
	// 	// 同时存在object和动态字段,禁止使用hgetall
	// 	obj.TypeHash.HashObject.HGetAll = false
	default:
		return fmt.Errorf("redis-hash type fields count %d not support", len(msg.Fields))
	}

	return
}

func analyseTypeSet(msg *protogen.Message, obj *gen_redis.RedisObject, usePb bool) (err error) {
	opt := &gen_redis.RedisTypeSet{}
	obj.TypeSet = opt

	getOptBool := func(opt protoreflect.ExtensionType, def bool) (v bool) {
		v = def
		proto.RangeExtensions(msg.Desc.Options(), func(et protoreflect.ExtensionType, a any) bool {
			if et.TypeDescriptor().FullName() == opt.TypeDescriptor().FullName() {
				v = a.(bool)
				return false
			}
			return true
		})
		return
	}

	getOptString := func(opt protoreflect.ExtensionType, def string) (v string) {
		v = def
		proto.RangeExtensions(msg.Desc.Options(), func(et protoreflect.ExtensionType, a any) bool {
			if et.TypeDescriptor().FullName() == opt.TypeDescriptor().FullName() {
				v = a.(string)
				return false
			}
			return true
		})
		return
	}
	// 设置此选项后将忽略其他配置选项,拼接member.
	opMember := getOptString(redis.E_Member, "")
	if opMember != "" {
		// 拼接string做member
		opt.Args, err = keyarg.MatchGoTypes(opMember, nil)
		if err != nil {
			return fmt.Errorf("redis-set analyse redis.member msg:%s failed.%v", msg.Desc.Name(), err)
		}
		// log.Println("set custom member:", opt.Args)
		return
	}

	// 未指定操作另一个文件, 并且未设置op_field,操作当前结构体
	if !gen_redis.Config.OpAnotherFile {
		if !getOptBool(redis.E_OpField, false) {
			pkgName := "proto"
			if usePb {
				obj.Import(gen_redis.Config.PkgPB, "Marshal/Unmarshal")
			} else {
				obj.Import(gen_redis.Config.PkgJson, "Marshal/Unmarshal")
				pkgName = "json"
			}
			// 并且未设置op_field,操作当前结构体
			typ := filepath.Base(string(msg.GoIdent.GoImportPath)) + "." + msg.GoIdent.GoName
			opt.Message = &gen_redis.RedisGenMsg{
				Type: "*" + typ,
				New:  "&" + typ + "{}",
				Marshal: func(objName string) string {
					return fmt.Sprintf("%s.Marshal(%s)", pkgName, objName)
				},
				Unmarshal: func(objName, paramName string) string {
					return fmt.Sprintf("%s.Unmarshal(%s,%s)", pkgName, paramName, objName)
				},
			}
			obj.Import(string(msg.GoIdent.GoImportPath), "msg")
			// log.Println("set current custom message:", opt.Message)
			return
		}
	}
	switch len(msg.Fields) {
	case 0:
		// 无数据设置,生成通用接口
		pkgName := "proto"
		typ := "proto.Message"
		if usePb {
			obj.Import(gen_redis.Config.PkgPB, "Marshal/Unmarshal")
		} else {
			obj.Import(gen_redis.Config.PkgJson, "Marshal/Unmarshal")
			pkgName = "json"
			typ = "any"
		}
		opt.Message = &gen_redis.RedisGenMsg{
			Type: typ,
			New:  "",
			Marshal: func(objName string) string {
				return fmt.Sprintf("%s.Marshal(%s)", pkgName, objName)
			},
			Unmarshal: func(objName, paramName string) string {
				return fmt.Sprintf("%s.Unmarshal(%s,%s)", pkgName, paramName, objName)
			},
		}
		//log.Println("set common message:", opt.Message)
		return
	case 1:
		fieldType := msg.Fields[0].Desc
		if fieldType.Kind() == protoreflect.BoolKind || fieldType.IsList() || fieldType.IsMap() {
			return fmt.Errorf("redis-set type member not support bool/Array/Map type, msg:%s", string(msg.Desc.Name()))
		}
		// 自定义类型
		if fieldType.Kind() == protoreflect.MessageKind {
			dstMsg := msg.Fields[0].Message
			pkgName := "proto"
			if usePb {
				obj.Import(gen_redis.Config.PkgPB, "Marshal/Unmarshal")
			} else {
				obj.Import(gen_redis.Config.PkgJson, "Marshal/Unmarshal")
				pkgName = "json"
			}
			typ := filepath.Base(string(dstMsg.GoIdent.GoImportPath)) + "." + dstMsg.GoIdent.GoName
			opt.Message = &gen_redis.RedisGenMsg{
				Type: "*" + typ,
				New:  "&" + typ + "{}",
				Marshal: func(objName string) string {
					return fmt.Sprintf("%s.Marshal(%s)", pkgName, objName)
				},
				Unmarshal: func(objName, paramName string) string {
					return fmt.Sprintf("%s.Unmarshal(%s,%s)", pkgName, paramName, objName)
				},
			}
			obj.Import(string(dstMsg.GoIdent.GoImportPath), "msg")
			// log.Println("set custom message:", opt.Message)
			return
		}
		field := msg.Fields[0]
		opt.BaseType = &gen_redis.RedisGenType{
			Name:      string(field.Desc.Name()),
			Type:      getFieldGoType(field.Desc.Kind()),
			Number:    isNumber(field.Desc.Kind()),
			RedisFunc: getRedisFunc(field.Desc.Kind()),
		}
		//log.Println("set base type:", opt.BaseType)
		return
	}
	return fmt.Errorf("redis-set type fields count %d not support,msg:%s", len(msg.Fields), string(msg.Desc.Name()))
}
func analyseTypeZSet(msg *protogen.Message, obj *gen_redis.RedisObject, usePb bool) (err error) {
	opt := &gen_redis.RedisTypeZSet{}
	obj.TypeZSet = opt

	getOptBool := func(opt protoreflect.ExtensionType, def bool) (v bool) {
		v = def
		proto.RangeExtensions(msg.Desc.Options(), func(et protoreflect.ExtensionType, a any) bool {
			if et.TypeDescriptor().FullName() == opt.TypeDescriptor().FullName() {
				v = a.(bool)
				return false
			}
			return true
		})
		return
	}

	getOptString := func(opt protoreflect.ExtensionType, def string) (v string) {
		v = def
		proto.RangeExtensions(msg.Desc.Options(), func(et protoreflect.ExtensionType, a any) bool {
			if et.TypeDescriptor().FullName() == opt.TypeDescriptor().FullName() {
				v = a.(string)
				return false
			}
			return true
		})
		return
	}
	if !gen_redis.Config.OpAnotherFile && !getOptBool(redis.E_OpField, false) {
		// 未指定操作另一个文件, 并且未设置op_field,操作当前结构体, 默认score是float64
		opt.Score = &gen_redis.RedisGenType{
			Name:      "score",
			Type:      "float64",
			Number:    true,
			RedisFunc: "Float64",
		}
		pkgName := "proto"
		if usePb {
			obj.Import(gen_redis.Config.PkgPB, "Marshal/Unmarshal")
		} else {
			obj.Import(gen_redis.Config.PkgJson, "Marshal/Unmarshal")
			pkgName = "json"
		}
		typ := filepath.Base(string(msg.GoIdent.GoImportPath)) + "." + msg.GoIdent.GoName
		opt.Message = &gen_redis.RedisGenMsg{
			Type: "*" + typ,
			New:  "&" + typ + "{}",
			Marshal: func(objName string) string {
				return fmt.Sprintf("%s.Marshal(%s)", pkgName, objName)
			},
			Unmarshal: func(objName, paramName string) string {
				return fmt.Sprintf("%s.Unmarshal(%s,%s)", pkgName, paramName, objName)
			},
		}
		return
	}
	// 设置此选项后将忽略其他配置选项,拼接member.
	opMember := getOptString(redis.E_Member, "")
	if opMember != "" {
		// 拼接string做member
		opt.Args, err = keyarg.MatchGoTypes(opMember, nil)
		if err != nil {
			return fmt.Errorf("redis-zset analyse redis.member msg:%s failed.%v", msg.Desc.Name(), err)
		}
		switch len(msg.Fields) {
		case 0:
			opt.Score = &gen_redis.RedisGenType{
				Name:      "score",
				Type:      "float64",
				Number:    true,
				RedisFunc: "Float64",
			}
			return
		case 1:
			field := msg.Fields[0]
			// field必须是 有符号数值类型或者float,double
			if field.Desc.IsList() || field.Desc.IsMap() || field.Desc.Kind() == protoreflect.MessageKind ||
				field.Desc.Kind() == protoreflect.GroupKind || field.Desc.Kind() == protoreflect.BoolKind ||
				field.Desc.Kind() == protoreflect.StringKind || field.Desc.Kind() == protoreflect.BytesKind {
				return fmt.Errorf("redis-zset type message score define [%s.%s] is not support type", msg.Desc.Name(), field.Desc.Name())
			}
			// 无符号类型过滤
			if field.Desc.Kind() == protoreflect.Uint32Kind || field.Desc.Kind() == protoreflect.Uint64Kind ||
				field.Desc.Kind() == protoreflect.Fixed32Kind || field.Desc.Kind() == protoreflect.Fixed64Kind {
				return fmt.Errorf("redis-zset type message score define [%s.%s] is not support unsigned type", msg.Desc.Name(), field.Desc.Name())
			}
			opt.Score = &gen_redis.RedisGenType{
				Name:      string(field.Desc.Name()),
				Type:      getFieldGoType(field.Desc.Kind()),
				Number:    isNumber(field.Desc.Kind()),
				RedisFunc: getRedisFunc(field.Desc.Kind()),
			}
			return
		default:
			return fmt.Errorf("redis-zset type fields count %d not support", len(msg.Fields))
		}
	}
	//
	if len(msg.Fields) < 1 {
		return fmt.Errorf("redis-zset type fields count %d not support", len(msg.Fields))
	}

	member := msg.Fields[0]
	// member
	if member.Desc.IsList() || member.Desc.IsMap() ||
		member.Desc.Kind() == protoreflect.GroupKind || member.Desc.Kind() == protoreflect.BoolKind {
		return fmt.Errorf(
			`redis-zset type message member define [%s.%s] is not support type`,
			msg.Desc.Name(),
			member.Desc.Name(),
		)
	}
	if member.Desc.Kind() == protoreflect.MessageKind {
		dstMsg := member.Message
		pkgName := "proto"
		if usePb {
			obj.Import(gen_redis.Config.PkgPB, "Marshal/Unmarshal")
		} else {
			obj.Import(gen_redis.Config.PkgJson, "Marshal/Unmarshal")
			pkgName = "json"
		}
		typ := filepath.Base(string(dstMsg.GoIdent.GoImportPath)) + "." + dstMsg.GoIdent.GoName
		opt.Message = &gen_redis.RedisGenMsg{
			Type: "*" + typ,
			New:  "&" + typ + "{}",
			Marshal: func(objName string) string {
				return fmt.Sprintf("%s.Marshal(%s)", pkgName, objName)
			},
			Unmarshal: func(objName, paramName string) string {
				return fmt.Sprintf("%s.Unmarshal(%s,%s)", pkgName, paramName, objName)
			},
		}
	} else {
		opt.Member = &gen_redis.RedisGenType{
			Name:      string(member.Desc.Name()),
			Type:      getFieldGoType(member.Desc.Kind()),
			Number:    isNumber(member.Desc.Kind()),
			RedisFunc: getRedisFunc(member.Desc.Kind()),
		}
	}

	switch len(msg.Fields) {
	case 1:
		opt.Score = &gen_redis.RedisGenType{
			Name:      "score",
			Type:      "float64",
			Number:    true,
			RedisFunc: "Float64",
		}
		return
	case 2:
		score := msg.Fields[1]
		if score.Desc.IsList() || score.Desc.IsMap() || score.Desc.Kind() == protoreflect.MessageKind ||
			score.Desc.Kind() == protoreflect.GroupKind || score.Desc.Kind() == protoreflect.BoolKind ||
			score.Desc.Kind() == protoreflect.StringKind || score.Desc.Kind() == protoreflect.BytesKind {
			return fmt.Errorf(
				`redis-zset type message score define [%s.%s] is not support type`,
				msg.Desc.Name(),
				score.Desc.Name(),
			)
		}
		// 无符号类型过滤
		if score.Desc.Kind() == protoreflect.Uint32Kind || score.Desc.Kind() == protoreflect.Uint64Kind ||
			score.Desc.Kind() == protoreflect.Fixed32Kind || score.Desc.Kind() == protoreflect.Fixed64Kind {
			return fmt.Errorf(
				`redis-zset type message score define [%s.%s] is not support unsigned type`,
				msg.Desc.Name(),
				score.Desc.Name(),
			)
		}
		opt.Score = &gen_redis.RedisGenType{
			Name:      string(score.Desc.Name()),
			Type:      getFieldGoType(score.Desc.Kind()),
			Number:    isNumber(score.Desc.Kind()),
			RedisFunc: getRedisFunc(score.Desc.Kind()),
		}
	}

	return
}

func analyseRedisScripts(msg *protogen.Message, obj *gen_redis.RedisObject) (err error) {
	redisScripts := proto.GetExtension(msg.Desc.Options(), redis.E_Script).([]*redis.RedisScript)
	if len(redisScripts) < 1 {
		return
	}
	scriptNameMap := make(map[string]int, len(redisScripts))
	var LuaState *lua.LState
	for k, script := range redisScripts {
		if script == nil {
			continue
		}
		// 检查是否有重复的脚本名称
		if last, ok := scriptNameMap[script.Name]; ok {
			err = multierr.Append(err, fmt.Errorf("msg:%s redis script name %s duplicate,cur:%d last:%d",
				msg.Desc.Name(), script.Name, k, last,
			))
		}
		scriptNameMap[script.Name] = k
		// 检查脚本内容是否为空
		if script.Lua == "" {
			err = multierr.Append(err, fmt.Errorf("msg:%s redis script name %s script content is empty",
				msg.Desc.Name(), script.Name,
			))
		}
		// 检查脚本参数是否为空
		argv, err := keyarg.MatchGoTypes(script.Argv, nil)
		if err != nil {
			err = multierr.Append(err, fmt.Errorf("msg:%s redis script name %s script argv error:%v",
				msg.Desc.Name(), script.Name, err,
			))
		}
		reply, err := keyarg.MatchGoTypes(script.Reply, nil)
		if err != nil {
			err = multierr.Append(err, fmt.Errorf("msg:%s redis script name %s script reply error:%v",
				msg.Desc.Name(), script.Name, err,
			))
		}
		if err == nil && len(reply) < 1 {
			err = fmt.Errorf("msg:%s redis script name %s script reply empty,must >= 1",
				msg.Desc.Name(), script.Name,
			)
		}

		// 检测lua基础语法
		if gen_redis.Config.CheckLuaScript && len(script.Lua) > 0 {
			if LuaState == nil {
				LuaState = lua.NewState()
			}
			if _, err := LuaState.LoadString(script.Lua); err != nil {
				err = multierr.Append(err, fmt.Errorf("msg:%s redis script name %s script lua error:%v",
					msg.Desc.Name(), script.Name, err,
				))
			}
		}

		if err != nil {
			continue
		}

		err = keyarg.CheckArgNameConflict(argv, reply)
		if err != nil {
			err = multierr.Append(err, fmt.Errorf("msg:%s redis script name %s script argv and reply conflict:%v",
				msg.Desc.Name(), script.Name, err,
			))
			continue
		}

		sobj := &gen_redis.RedisScript{
			Name:         script.Name,
			Script:       script.Lua,
			Args:         argv,
			Output:       reply,
			TemplateName: "script_return_mul",
			CommandName:  "",
		}
		if len(reply) == 1 {
			sobj.TemplateName = "script_return_1"
			switch reply[0].ArgType() {
			case "bool":
				sobj.CommandName = "NewBoolCmd"
			case "float32", "float64":
				sobj.CommandName = "NewFloatCmd"
			case "string":
				sobj.CommandName = "NewStringCmd"
			default:
				sobj.CommandName = "NewIntCmd"
			}
		}

		obj.Scripts = append(obj.Scripts, sobj)
	}

	return
}

// func initRedisLuaCheckState() *lua.LState {
// 	L := lua.NewState()
// 	for _, lib := range []struct {
// 		libName string
// 		libFunc lua.LGFunction
// 	}{
// 		{lua.TabLibName, lua.OpenTable},
// 		{lua.StringLibName, lua.OpenString},
// 		{lua.MathLibName, lua.OpenMath},
// 	} {
// 		L.Push(L.NewFunction(lib.libFunc))
// 		L.Push(lua.LString(lib.libName))
// 		L.Call(1, 0)
// 	}
// 	L.SetGlobal("KEYS", L.NewTable())
// 	L.SetGlobal("ARGV", L.NewTable())
// 	rds := L.NewTable()
// 	// call
// 	rds.RawSet(lua.LString("call"), L.NewFunction(func(l *lua.LState) int {
// 		return 0
// 	}))
// 	// pcall
// 	rds.RawSet(lua.LString("pcall"), L.NewFunction(func(l *lua.LState) int {
// 		return 0
// 	}))

// 	L.SetGlobal("redis", rds)
// 	return L
// }

func isNumber(typ protoreflect.Kind) bool {
	switch typ {
	case protoreflect.StringKind, protoreflect.BytesKind, protoreflect.BoolKind, protoreflect.FloatKind, protoreflect.DoubleKind:
		return false
	default:
		return true
	}
}

func getRedisFunc(kind protoreflect.Kind) string {
	switch kind {
	case protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Sfixed32Kind:
		return "Int32"
	case protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Sfixed64Kind:
		return "Int64"
	case protoreflect.Uint32Kind, protoreflect.Fixed32Kind:
		return "Uint32"
	case protoreflect.Uint64Kind, protoreflect.Fixed64Kind:
		return "Uint64"
	case protoreflect.FloatKind:
		return "Float32"
	case protoreflect.DoubleKind:
		return "Float64"
	case protoreflect.StringKind:
		return "String"
	case protoreflect.BytesKind:
		return "Binary"
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
	}
	return ""
}
