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

	log.Printf("msg:%s redis.pb:%t redis.json:%t env:%t usePb:%t\n", msg.Desc.Name(),
		getOptBool(redis.E_Pb, false), getOptBool(redis.E_Json, false), gen_redis.Config.MsgProtocol == "proto", usePb,
	)

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
	case "zset":
	case "":
		if !obj.TypeKeys {
			err = multierr.Append(err, fmt.Errorf("msg:%s redis type config invalid", msg.Desc.Name()))
			return
		}
	default:
		err = multierr.Append(err, fmt.Errorf("msg:%s redis type not support [%s]", msg.Desc.Name(), redisType))
		return
	}
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
		opt.Type = fieldType.Kind().String()
	case protoreflect.Uint32Kind, protoreflect.Uint64Kind, protoreflect.Fixed32Kind, protoreflect.Fixed64Kind:
		opt.Number = true
		opt.Type = fieldType.Kind().String()
	case protoreflect.FloatKind, protoreflect.DoubleKind:
		opt.Float = true
		if fieldType.Kind() == protoreflect.FloatKind {
			opt.Type = "float32"
		} else {
			opt.Type = "float64"
		}
	case protoreflect.StringKind:
		opt.String = true
		opt.Type = fieldType.Kind().String()
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
				return fmt.Errorf("redis-hash type message field [%s.%s] is not basic type.", msg.Desc.Name(), field.Desc.Name())
			}
			field := &gen_redis.RedisGenType{
				Name:      string(field.Desc.Name()),
				Type:      field.Desc.Kind().String(),
				Number:    isNumber(field.Desc.Kind()),
				RedisFunc: utils.Title(field.Desc.Kind().String()),
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
				return fmt.Errorf("redis-hash type message field [%s.%s] is not message type.", msg.Desc.Name(), fieldType.Name())
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
				Type:      fieldType.Kind().String(),
				Number:    isNumber(fieldType.Kind()),
				RedisFunc: utils.Title(fieldType.Kind().String()),
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
				return fmt.Errorf("redis-hash type message field [%s.%s] is not basic type.", msg.Desc.Name(), fieldType.Name())
			}
			dynamic.Field = &gen_redis.RedisGenType{
				Name:      utils.PascalToSnake(msg.Fields[0].GoName),
				Type:      fieldType.Kind().String(),
				Number:    isNumber(fieldType.Kind()),
				RedisFunc: utils.Title(fieldType.Kind().String()),
			}
		}

	case 2:
		fieldType1 := msg.Fields[0].Desc
		fieldType2 := msg.Fields[1].Desc
		if fieldType1.IsList() || fieldType1.IsMap() || fieldType1.Kind() == protoreflect.MessageKind {
			return fmt.Errorf("redis-hash type message field [%s.%s] is not basic type.", msg.Desc.Name(), fieldType1.Name())
		}
		dynamic := &gen_redis.RedisHashDynamic{
			GenMap: true,
		}
		hash.HashDynamic = dynamic
		dynamic.Field = &gen_redis.RedisGenType{
			Name:      utils.PascalToSnake(msg.Fields[0].GoName),
			Type:      fieldType1.Kind().String(),
			Number:    isNumber(fieldType1.Kind()),
			RedisFunc: utils.Title(fieldType1.Kind().String()),
		}
		if fieldType2.IsList() || fieldType2.IsMap() {
			return fmt.Errorf("redis-hash type message value [%s.%s] not support array or map type.", msg.Desc.Name(), fieldType2.Name())
		}
		dynamic.Value = &gen_redis.RedisGenType{
			Name:      utils.PascalToSnake(msg.Fields[1].GoName),
			Type:      fieldType2.Kind().String(),
			Number:    isNumber(fieldType2.Kind()),
			RedisFunc: utils.Title(fieldType2.Kind().String()),
		}
		if fieldType2.Kind() == protoreflect.MessageKind {
			dynamic.Value.Type = filepath.Base(string(msg.Fields[1].Message.GoIdent.GoImportPath)) + "." + msg.Fields[1].Message.GoIdent.GoName

			log.Printf("msgFrom:%s msgTo:%s usePB:%t\n", msg.Desc.Name(), dynamic.Value.Type, usePb)
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

// func analyseTypeSet(msg *buildpb.MsgDesc, obj *RedisObject) (err error) {
// 	// 最多只能有一个字段
// 	if len(msg.Fields) > 1 {
// 		return errors.New("redis-set type fields too many. max 1 fields")
// 	}
// 	opt := &RedisTypeSet{}
// 	obj.TypeSet = opt

// 	// 无数据设置,查看是否生成通用接口
// 	if len(msg.Fields) < 1 {
// 		// 生成protobuf接口
// 		if msg.Options.GetOptionBool(options.RedisOpProtobuf) && len(envFlag.ProtobufPackage) > 0 {
// 			obj.Import(envFlag.ProtobufPackage, "Marshal/Unmarshal")
// 			opt.Message = &RedisGenMsg{
// 				Type: "proto.Message",
// 				Marshal: func(objName string) string {
// 					return fmt.Sprintf("proto.Marshal(%s)", objName)
// 				},
// 				Unmarshal: func(objName string, paramName string) string {
// 					return fmt.Sprintf("proto.Unmarshal(%s,%s)", objName, paramName)
// 				},
// 				New: "",
// 			}
// 			return
// 		}
// 		// 生成walle message 接口
// 		if msg.Options.GetOptionBool(options.RedisOpWalleMsg) && len(envFlag.WProtoPackage) > 0 {
// 			obj.Import(envFlag.WProtoPackage, "MarshalObject/UnmarshalObject")
// 			opt.Message = &RedisGenMsg{
// 				Type: "message.Message",
// 				Marshal: func(objName string) string {
// 					return fmt.Sprintf("%s.MarshalObject()", objName)
// 				},
// 				Unmarshal: func(objName string, paramName string) string {
// 					return fmt.Sprintf("%s.UnmarshalObject(%s)", objName, paramName)
// 				},
// 				New: "",
// 			}
// 			return
// 		}
// 		// 无任何设置,直接生成string类型接口
// 		opt.BaseType = &RedisGenType{
// 			Name:      "",
// 			Type:      "string",
// 			Number:    false,
// 			RedisFunc: "String",
// 		}
// 		return
// 	}

// 	fieldType := msg.Fields[0].Type

// 	switch fieldType.Type {
// 	case buildpb.FieldType_BaseType:
// 		switch fieldType.KeyBase {
// 		// case buildpb.BaseTypeDesc_Binary:
// 		// 	return errors.New("redis-string type generation not support binary basic type.")
// 		case buildpb.BaseTypeDesc_Bool:
// 			return errors.New("redis-set type generation not support bool basic type")
// 		default:
// 			opt.BaseType = &RedisGenType{
// 				Name:      fieldType.Key,
// 				Type:      fieldType.Key,
// 				Number:    false,
// 				RedisFunc: fieldType.KeyBase.String(),
// 			}
// 		}
// 	case buildpb.FieldType_CustomType:
// 		opt.Message = &RedisGenMsg{
// 			Type: "*" + keyType(fieldType.Key),
// 			Marshal: func(objName string) string {
// 				return fmt.Sprintf("%s.MarshalObject()", objName)
// 			},
// 			Unmarshal: func(objName string, paramName string) string {
// 				return fmt.Sprintf("%s.UnmarshalObject(%s)", objName, paramName)
// 			},
// 			New: "&" + keyType(fieldType.Key) + "{}",
// 		}
// 		if msg.Options.GetOptionBool(options.RedisOpProtobuf) && len(envFlag.ProtobufPackage) > 0 {
// 			obj.Import(envFlag.ProtobufPackage, "Marshal/Unmarshal")
// 			opt.Message.Marshal = func(objName string) string {
// 				return fmt.Sprintf("proto.Marshal(%s)", objName)
// 			}
// 			opt.Message.Unmarshal = func(objName string, paramName string) string {
// 				return fmt.Sprintf("proto.Unmarshal(%s,%s)", objName, paramName)
// 			}
// 		}
// 	default:
// 		return errors.New("redis-set type generation not support array or map type")
// 	}

// 	return
// }

// func analyseTypeZSet(msg *buildpb.MsgDesc, obj *RedisObject) (err error) {
// 	// 只支持1个或者2个字段
// 	fieldCount := len(msg.Fields)
// 	if fieldCount < 1 || fieldCount > 2 {
// 		return errors.New("redis-zset type fields invalid. only support 1 or 2 fields")
// 	}
// 	opt := &RedisTypeZSet{}
// 	obj.TypeZSet = opt

// 	fieldType := msg.Fields[0].Type

// 	switch fieldType.Type {
// 	case buildpb.FieldType_BaseType:
// 		switch fieldType.KeyBase {
// 		case buildpb.BaseTypeDesc_Bool:
// 			return errors.New("redis-zset type member not support bool basic type")
// 		default:
// 			if fieldType.KeyBase == buildpb.BaseTypeDesc_String && msg.HasOption(options.RedisOpMatchMember) {
// 				// 拼接string做value
// 				opt.Args, err = keyarg.MatchGoTypes(msg.Options.GetString(options.RedisOpMatchMember, ""), nil)
// 				if err != nil {
// 					return fmt.Errorf("redis-zset analyse redis.member failed.%v", err)
// 				}
// 			} else {
// 				opt.Member = &RedisGenType{
// 					Name:      fieldType.Key,
// 					Type:      fieldType.Key,
// 					Number:    false,
// 					RedisFunc: fieldType.KeyBase.String(),
// 				}
// 			}
// 		}
// 	case buildpb.FieldType_CustomType:
// 		opt.Message = &RedisGenMsg{
// 			Type: "*" + keyType(fieldType.Key),
// 			Marshal: func(objName string) string {
// 				return fmt.Sprintf("%s.MarshalObject()", objName)
// 			},
// 			Unmarshal: func(objName string, paramName string) string {
// 				return fmt.Sprintf("%s.UnmarshalObject(%s)", objName, paramName)
// 			},
// 			New: "&" + keyType(fieldType.Key) + "{}",
// 		}
// 		if msg.Options.GetOptionBool(options.RedisOpProtobuf) && len(envFlag.ProtobufPackage) > 0 {
// 			obj.Import(envFlag.ProtobufPackage, "Marshal/Unmarshal")
// 			opt.Message.Marshal = func(objName string) string {
// 				return fmt.Sprintf("proto.Marshal(%s)", objName)
// 			}
// 			opt.Message.Unmarshal = func(objName string, paramName string) string {
// 				return fmt.Sprintf("proto.Unmarshal(%s,%s)", objName, paramName)
// 			}
// 		}
// 	default:
// 		return errors.New("redis-zset type generation not support array or map type")
// 	}

// 	if fieldCount < 2 {
// 		opt.Score = &RedisGenType{
// 			Name:      "score",
// 			Type:      "float64",
// 			Number:    true,
// 			RedisFunc: "Float64",
// 		}
// 		return
// 	}

// 	scoreType := msg.Fields[1].Type
// 	if scoreType.Type != buildpb.FieldType_BaseType {
// 		return errors.New("redis-zset type score only support signed int or float type")
// 	}

// 	if !strings.HasPrefix(scoreType.Key, "int") && !strings.HasPrefix(scoreType.Key, "float") {
// 		return errors.New("redis-zset type score only support signed int or float type")
// 	}

// 	opt.Score = &RedisGenType{
// 		Name:      "score",
// 		Type:      scoreType.Key,
// 		Number:    false,
// 		RedisFunc: scoreType.KeyBase.String(),
// 	}

// 	return
// }

// func analyseScript(msg *buildpb.MsgDesc, obj *RedisObject) (err error) {
// 	for optKey := range msg.Options.Options {
// 		if !strings.HasPrefix(optKey, options.RedisScriptPrefix) {
// 			continue
// 		}
// 		if !strings.HasSuffix(optKey, options.RedisScriptSuffixScript) {
// 			continue
// 		}
// 		scriptName := strings.TrimSuffix(strings.TrimPrefix(optKey, options.RedisScriptPrefix), options.RedisScriptSuffixScript)
// 		if strings.Contains(scriptName, ".") {
// 			return fmt.Errorf("MessageName:%s define redis script failed. script name [%s] invalid", msg.Name, scriptName)
// 		}
// 		scriptData := msg.GetString(options.RedisScriptPrefix+scriptName+options.RedisScriptSuffixScript, "")
// 		scriptArgv := msg.GetString(options.RedisScriptPrefix+scriptName+options.RedisScriptSuffixInput, "")
// 		scriptReply := msg.GetString(options.RedisScriptPrefix+scriptName+options.RedisScriptSuffixReply, "")
// 		// log.Println(options.RedisScriptPrefix+scriptName+options.RedisScriptSuffixScript,
// 		// 	options.RedisScriptPrefix+scriptName+options.RedisScriptSuffixInput,
// 		// 	options.RedisScriptPrefix+scriptName+options.RedisScriptSuffixReply)

// 		if scriptData == "" {
// 			return fmt.Errorf("MessageName:%s redis script [%s] data empty", msg.Name, scriptName)
// 		}
// 		if scriptArgv == "" {
// 			return fmt.Errorf("MessageName:%s redis script [%s] argv empty", msg.Name, scriptName)
// 		}
// 		if scriptReply == "" {
// 			return fmt.Errorf("MessageName:%s redis script [%s] reply empty", msg.Name, scriptName)
// 		}
// 		argv, err := keyarg.MatchGoTypes(scriptArgv, nil)
// 		if err != nil {
// 			return fmt.Errorf("MessageName:%s redis script [%s] argv invalid. %+v", msg.Name, scriptName, err)
// 		}

// 		reply, err := keyarg.MatchGoTypes(scriptReply, nil)
// 		if err != nil {
// 			return fmt.Errorf("MessageName:%s redis script [%s] reply invalid. %+v", msg.Name, scriptName, err)
// 		}
// 		if len(reply) < 1 {
// 			return fmt.Errorf("MessageName:%s redis script [%s] reply must >= 1", msg.Name, scriptName)
// 		}

// 		script := &RedisScript{
// 			Name:         scriptName,
// 			Script:       scriptData,
// 			Args:         argv,
// 			Output:       reply,
// 			TemplateName: "script_return_mul",
// 			CommandName:  "",
// 		}
// 		if len(reply) == 1 {
// 			script.TemplateName = "script_return_1"
// 			switch reply[0].ArgType() {
// 			case "bool":
// 				script.CommandName = "NewBoolCmd"
// 			case "float32", "float64":
// 				script.CommandName = "NewFloatCmd"
// 			case "string":
// 				script.CommandName = "NewStringCmd"
// 			default:
// 				script.CommandName = "NewIntCmd"
// 			}
// 		}

// 		obj.Scripts = append(obj.Scripts, script)

// 	}

// 	return
// }

func isNumber(typ protoreflect.Kind) bool {
	switch typ {
	case protoreflect.StringKind, protoreflect.BytesKind, protoreflect.BoolKind, protoreflect.FloatKind, protoreflect.DoubleKind:
		return false
	default:
		return true
	}
}

func checkHashMatchModeArg(obj *gen_redis.RedisHashDynamic) (err error) {
	if obj == nil {
		return
	}
	if obj.FieldArgs == nil || obj.ValueArgs == nil {
		return
	}
	checks := make(map[string]struct{})
	for _, v := range obj.FieldArgs {
		checks[v.ArgName()] = struct{}{}
	}
	for _, v := range obj.ValueArgs {
		if _, ok := checks[v.ArgName()]; ok {
			return fmt.Errorf("redis-hash match field named repeated[%s]", v.ArgName())
		}
	}
	return
}
