package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/spf13/pflag"
	"go.uber.org/multierr"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"

	"github.com/walleframe/pb_plugins/extend/mysql"
	"github.com/walleframe/pb_plugins/pkg/gen_mysql"
	"github.com/walleframe/pb_plugins/pkg/tpl"
	"github.com/walleframe/pb_plugins/pkg/utils"
)

const genGoDocURL = "https://protobuf.dev/reference/go/go-generated"
const grpcDocURL = "https://grpc.io/docs/languages/go/quickstart/#regenerate-grpc-code"

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

	_ = mysql.E_DbName
	// 从环境变量中读取配置
	gen_mysql.ParseConfigFromEnv()

	// 外部命令行参数(比环境变量优先级高)
	var flags pflag.FlagSet
	cfg := gen_mysql.Config
	flags.StringVar(&cfg.SvcPkg, "pkg_svc", cfg.SvcPkg, "service package name")
	flags.StringVar(&cfg.UtilPkg, "pkg_util", cfg.UtilPkg, "util package name")
	flags.StringVar(&cfg.CodePath, "code_path", cfg.CodePath, "code path")
	flags.StringVar(&cfg.Charset, "charset", cfg.Charset, "charset")
	flags.StringVar(&cfg.Collate, "collate", cfg.Collate, "collate")
	protogen.Options{
		ParamFunc: flags.Set,
	}.Run(func(gen *protogen.Plugin) (err error) {

		log.Println(gen_mysql.Config)
		err = gen_mysql.InitTemplate()
		if err != nil {
			log.Println("gen_mysql InitTemplate error: ", err)
			return
		}

		for _, f := range gen.Files {
			if !f.Generate {
				continue
			}
			err = genDatabase(gen, f)
		}
		//gen.SupportedFeatures = uint64(pluginpb.CodeGeneratorResponse_FEATURE_PROTO3_OPTIONAL)
		return
	})
}

func genDatabase(gen *protogen.Plugin, file *protogen.File) (err error) {
	genDatabase := proto.GetExtension(file.Proto.Options, mysql.E_DbName)
	if genDatabase == nil {
		return
	}

	dbName := genDatabase.(string)

	getOptString := func(opt protoreflect.ExtensionType, def string) (v string) {
		v = def
		proto.RangeExtensions(file.Proto.Options, func(et protoreflect.ExtensionType, a any) bool {
			if opt.TypeDescriptor().FullName() == et.TypeDescriptor().FullName() {
				log.Println("get ext value:", a.(string), "def:", def)
				v = a.(string)
				return false
			}
			log.Println("want:", opt.TypeDescriptor().FullName(), " range:", et.TypeDescriptor().FullName(), a)
			return true
		})
		return
	}

	getOptBool := func(opt protoreflect.ExtensionType, def bool) (v bool) {
		v = def
		proto.RangeExtensions(file.Proto.Options, func(et protoreflect.ExtensionType, a any) bool {
			if et.TypeDescriptor().FullName() == opt.TypeDescriptor().FullName() {
				v = a.(bool)
				return false
			}
			return true
		})
		return
	}

	table := &gen_mysql.SqlTable{
		DB:        dbName,
		Charset:   getOptString(mysql.E_DbCharset, gen_mysql.Config.Charset),
		Collate:   getOptString(mysql.E_DbCollate, gen_mysql.Config.Collate),
		DisableEx: getOptBool(mysql.E_DisableEx, false),
	}
	table.ToolName = "protoc-gen-mysql"

	for _, msg := range file.Messages {
		err = genDBTable(gen, table, msg)
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

func genDBTable(gen *protogen.Plugin, table *gen_mysql.SqlTable, msg *protogen.Message) (err error) {
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

	// 指定忽略消息
	if getOptBool(mysql.E_Ignore, false) {
		return
	}

	// 复制有一份新的表结构
	table = table.Clone()

	// 大驼峰转小写加下划线
	name := pascalToSnake(string(msg.Desc.Name()))

	// 表引擎
	table.Engine = getOptString(mysql.E_Engine, "InnoDB")
	// 表名
	table.SqlTable = getOptString(mysql.E_TblName, "tbl_"+name)
	// 生成操作结构体名
	table.Name = string(msg.Desc.Name())
	// 对应的Go结构体名
	table.Struct = filepath.Base(string(msg.GoIdent.GoImportPath)) + "." + msg.GoIdent.GoName
	table.Import(string(msg.GoIdent.GoImportPath), "msg")
	//
	table.GenUpdate = getOptBool(mysql.E_Update, true)
	table.GenUpsert = getOptBool(mysql.E_Upsert, true)
	table.DisableEx = getOptBool(mysql.E_TblNoex, table.DisableEx)
	if !table.DisableEx {
		table.GenEx = getOptBool(mysql.E_GenEx, false)
	}
	table.CustomOptions = getOptString(mysql.E_TblOpt, "")

	// primary keys
	primaryKeys := getOptString(mysql.E_Pks, "")
	// index keys
	indexKeys := getOptString(mysql.E_Index, "")
	// unique keys
	uniqueKeys := getOptString(mysql.E_Unique, "")

	// 解析字段
	for _, field := range msg.Fields {
		if field.Desc.IsWeak() {
			continue
		}
		// 不支持枚举类型
		if field.Desc.Kind() == protoreflect.EnumKind {
			err = multierr.Append(err, fmt.Errorf("field %s is enum type, not supported", field.Desc.Name()))
			continue
		}
		err = multierr.Append(err, parseMysqlField(table, field))
	}
	// 解析主键
	err = multierr.Append(err, checkAndSetPrimaryKey(table, primaryKeys))
	// 解析索引
	err = multierr.Append(err, checkAndSetIndex(table, indexKeys, false))
	err = multierr.Append(err, checkAndSetIndex(table, uniqueKeys, true))
	if err != nil {
		return err
	}
	_ = primaryKeys
	_ = indexKeys
	_ = uniqueKeys

	// generate
	data, err := gen_mysql.Generate(table)
	if err != nil {
		err = multierr.Append(err, fmt.Errorf("generate table %s error: %v", table.SqlTable, err))
		return
	}
	// write response
	for _, d := range data {
		file := gen.NewGeneratedFile(d.File, "")
		file.Write(d.Data)
	}
	return
}

// parseMysqlField 解析字段
func parseMysqlField(table *gen_mysql.SqlTable, field *protogen.Field) (err error) {
	col := &gen_mysql.SqlColumn{
		Name:   pascalToSnake(field.GoName), // 大驼峰转小写加下划线
		GoType: field.Desc.Kind().String(),
	}
	if field.Desc.IsList() {
		col.GoType = "[]" + col.GoType
	}
	if field.Desc.IsMap() {
		col.GoType = "map[" + field.Message.Fields[0].Desc.Kind().String() + "]" +
			field.Message.Fields[1].Desc.Kind().String()
	}
	// log.Println("type:", col.GoType)
	// 字段名 log.Println("x:", field.Desc.Name(), field.Desc.FullName())
	// field.Desc.IsMap()
	if field.Comments.Trailing != "" || field.Comments.Leading != "" {
		col.Doc = &tpl.Commets{
			Doc:     []string{string(field.Comments.Leading)},
			TailDoc: string(field.Comments.Trailing),
		}
	}

	table.InnerAllColumns = append(table.InnerAllColumns, col)
	getOptString := func(opt protoreflect.ExtensionType, def string) (v string) {
		v = def
		proto.RangeExtensions(field.Desc.Options(), func(et protoreflect.ExtensionType, a any) bool {
			if et.TypeDescriptor().FullName() == opt.TypeDescriptor().FullName() {
				v = a.(string)
				return false
			}
			return true
		})
		return
	}
	getOptBool := func(opt protoreflect.ExtensionType, def bool) (v bool) {
		v = def
		proto.RangeExtensions(field.Desc.Options(), func(et protoreflect.ExtensionType, a any) bool {
			if et.TypeDescriptor().FullName() == opt.TypeDescriptor().FullName() {
				v = a.(bool)
				return false
			}
			return true
		})
		return
	}
	getOptInt := func(opt protoreflect.ExtensionType, def int32) (v int32) {
		v = def
		proto.RangeExtensions(field.Desc.Options(), func(et protoreflect.ExtensionType, a any) bool {
			if et.TypeDescriptor().FullName() == opt.TypeDescriptor().FullName() {
				v = a.(int32)
				return false
			}
			return true
		})
		return
	}
	//
	col.Unmarshal = utils.Title(col.GoType)
	if field.Desc.IsList() {
		col.Unmarshal = fmt.Sprintf("Slice[%s]", strings.TrimPrefix(col.GoType, "[]"))
	} else if field.Desc.IsMap() {
		col.Unmarshal = fmt.Sprintf("Map[%s,%s]",
			field.Message.Fields[0].Desc.Kind().String(),
			field.Message.Fields[1].Desc.Kind().String(),
		)
	} else {
		switch field.Desc.Kind() {
		case protoreflect.MessageKind, protoreflect.GroupKind:
			col.Unmarshal = fmt.Sprintf("Object[%s]", strings.TrimPrefix(col.GoType, "*"))
		}
	}
	// 对本字段,自定义序列化和反序列化函数
	setCustom := getOptBool(mysql.E_Custom, false)

	if setCustom {
		col.Marshal = col.Unmarshal
	}

	// 自定义字段设置,全部都需要手动写
	col.SqlType = getOptString(mysql.E_Column, "")
	if len(col.SqlType) > 0 {
		if strings.Contains(strings.ToLower(col.SqlType), "auto_increment") {
			table.AutoIncr = col
		}
		return
	}
	defaultSqlType := ""
	defaultSqlValue := ""
	if field.Desc.IsList() || field.Desc.IsMap() {
		defaultSqlType = fmt.Sprintf("varchar(%d)",
			getOptInt(mysql.E_Size, 64),
		)
		defaultSqlValue = "''"
	} else {
		switch field.Desc.Kind() {
		case protoreflect.BoolKind:
			defaultSqlType = "tinyint"
			defaultSqlValue = "0"
		case protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Sfixed32Kind:
			defaultSqlType = "int"
			defaultSqlValue = "0"
		case protoreflect.Uint32Kind, protoreflect.Fixed32Kind:
			defaultSqlType = "int unsigned"
			defaultSqlValue = "0"
		case protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Sfixed64Kind:
			defaultSqlType = "bigint"
			defaultSqlValue = "0"
		case protoreflect.Uint64Kind, protoreflect.Fixed64Kind:
			defaultSqlType = "bigint unsigned"
			defaultSqlValue = "0"
		case protoreflect.FloatKind:
			defaultSqlType = "float"
			defaultSqlValue = "0"
		case protoreflect.DoubleKind:
			defaultSqlType = "double"
			defaultSqlValue = "0"
		case protoreflect.StringKind:
			defaultSqlType = fmt.Sprintf("varchar(%d)",
				getOptInt(mysql.E_Size, 64),
			)
			defaultSqlValue = "''"
		case protoreflect.BytesKind:
			defaultSqlType = "blob"
			defaultSqlValue = "''"
		case protoreflect.MessageKind, protoreflect.GroupKind:
			defaultSqlType = fmt.Sprintf("varchar(%d)",
				getOptInt(mysql.E_Size, 64),
			)
			defaultSqlValue = "''"
		}
	}
	col.SqlType = getOptString(mysql.E_Type, defaultSqlType) + " not null default " + defaultSqlValue
	if getOptBool(mysql.E_Increment, false) {
		table.AutoIncr = col
		col.SqlType += " auto_increment"
	}
	return
}

// primaryKeys
func checkAndSetPrimaryKey(table *gen_mysql.SqlTable, primaryKeyConfig string) (err error) {
	// 解析主键
	primaryKeys := strings.Split(primaryKeyConfig, ",")
	primaryMap := make(map[string]struct{}, len(primaryKeys))
	for _, pk := range primaryKeys {
		pk = strings.TrimSpace(pk)
		if len(pk) == 0 {
			continue
		}
		var col *gen_mysql.SqlColumn
		for _, c := range table.InnerAllColumns {
			if c.Name == pk {
				col = c
				break
			}
		}
		if col == nil {
			err = multierr.Append(err, fmt.Errorf("primary key %s not found in table %s", pk, table.SqlTable))
			continue
		}
		if _, ok := primaryMap[pk]; ok {
			err = multierr.Append(err, fmt.Errorf("primary key %s duplicate in table %s", pk, table.SqlTable))
			continue
		}
		table.PrimaryKey = append(table.PrimaryKey, col)
		primaryMap[pk] = struct{}{}
	}
	for _, c := range table.InnerAllColumns {
		if _, ok := primaryMap[c.Name]; !ok {
			table.Columns = append(table.Columns, c)
			continue
		}
	}
	return
}

// name(column,...);name(column,...) to gen_mysql.SqlIndex
func checkAndSetIndex(table *gen_mysql.SqlTable, indexConfig string, uniqueFlagSet bool) (err error) {
	indexes := strings.Split(indexConfig, ";")
	for _, index := range indexes {
		index = strings.TrimSpace(index)
		if len(index) == 0 {
			continue
		}
		// index 的格式是 name(column,...), 先去name,然后判断(),最后再解析column名称列表
		name := ""
		if !strings.Contains(index, "(") || !strings.HasSuffix(index, ")") {
			err = multierr.Append(err, fmt.Errorf("index %s format error", index))
			continue
		}
		name = index[:strings.Index(index, "(")]
		columns := strings.Split(strings.TrimSuffix(strings.TrimPrefix(index, name+"("), ")"), ",")
		name = strings.TrimSpace(name)
		idx := &gen_mysql.SqlIndex{
			Name:     name,
			Columns:  make([]*gen_mysql.SqlColumn, 0, len(columns)),
			IsUnique: uniqueFlagSet,
			IdxName:  "idx_" + table.SqlTable + "_" + name,
		}
		table.Index = append(table.Index, idx)
		for _, col := range columns {
			col = strings.TrimSpace(col)
			if len(col) == 0 {
				continue
			}
			var column *gen_mysql.SqlColumn
			for _, c := range table.InnerAllColumns {
				if c.Name == col {
					column = c
					break
				}
			}
			if column == nil {
				err = multierr.Append(err, fmt.Errorf("index %s column %s not found in table %s",
					col, name, table.SqlTable,
				))
				continue
			}
			idx.Columns = append(idx.Columns, column)
		}

	}
	return
}

var (
	// 使用正则表达式匹配大写字母, 用于将大驼峰命名转换为小写加下划线的形式
	reToSnake = regexp.MustCompile("([a-z])([A-Z])")
)

// pascalToSnake 将大驼峰命名转换为小写加下划线的形式
func pascalToSnake(pascalStr string) string {
	// 插入下划线并转换为小写
	snakeStr := reToSnake.ReplaceAllString(pascalStr, "${1}_${2}")
	return strings.ToLower(snakeStr)
}
