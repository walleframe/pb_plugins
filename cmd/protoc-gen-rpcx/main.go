package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"go.uber.org/multierr"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"

	"github.com/spf13/pflag"
	"github.com/walleframe/pb_plugins/extend/gin"
	"github.com/walleframe/pb_plugins/pkg/gen_rpcx"
	"github.com/walleframe/pb_plugins/pkg/utils"
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
	gen_rpcx.ParseConfigFromEnv()

	// 外部命令行参数(比环境变量优先级高)
	var flags pflag.FlagSet
	cfg := gen_rpcx.Config
	flags.StringVar(&cfg.PkgServer, "pkg_server", cfg.PkgServer, "server package name")
	flags.StringVar(&cfg.PkgClient, "pkg_client", cfg.PkgClient, "client package name")
	flags.StringVar(&cfg.PkgProtocol, "pkg_protocol", cfg.PkgProtocol, "protocol package name")
	flags.StringVar(&cfg.PkgCtrl, "pkg_ctrl", cfg.PkgCtrl, "controller package name")
	flags.StringVar(&cfg.PkgStub, "pkg_stub", cfg.PkgStub, "stub package name")
	flags.StringVar(&cfg.PkgGin, "pkg_gin", cfg.PkgGin, "gin package name")
	flags.StringVar(&cfg.StubPath, "stub_path", cfg.StubPath, "stub file path")
	flags.StringVar(&cfg.GinPath, "gin_path", cfg.GinPath, "gin file path")
	protogen.Options{
		ParamFunc: flags.Set,
	}.Run(func(gen *protogen.Plugin) (err error) {
		err = gen_rpcx.InitTemplate()
		if err != nil {
			log.Println("gen_mysql InitTemplate error: ", err)
			return
		}

		for _, f := range gen.Files {
			if !f.Generate {
				continue
			}
			err = genRPCXFile(gen, f)
			if err != nil {
				return err
			}
		}
		//gen.SupportedFeatures = uint64(pluginpb.CodeGeneratorResponse_FEATURE_PROTO3_OPTIONAL)
		return
	})
}

func genRPCXFile(gen *protogen.Plugin, file *protogen.File) (err error) {
	filename := file.GeneratedFilenamePrefix + ".rpcx.go"

	zlog := &gen_rpcx.File{
		RPCFile: filename,
		Package: string(file.GoPackageName),
	}
	zlog.ToolName = "protoc-gen-rpcx"
	zlog.Package = string(file.GoPackageName)
	zlog.GoPackage = string(file.GoImportPath)
	getOptString := func(opt protoreflect.ExtensionType, def string) (v string) {
		v = def
		proto.RangeExtensions(file.Proto.Options, func(et protoreflect.ExtensionType, a any) bool {
			if et.TypeDescriptor().FullName() == opt.TypeDescriptor().FullName() {
				v = a.(string)
				return false
			}
			return true
		})
		return
	}

	// 扫描要生成的服务
	genRPCXServcies(file, zlog, getOptString(gin.E_Group, ""))
	// generate
	data, err := gen_rpcx.Generate(zlog)
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

func genRPCXServcies(file *protogen.File, rpc *gen_rpcx.File, apiGroup string) (err error) {

	getApiConfig := func(method *protogen.Method) (apiMethod, apiPath string) {
		getOptString := func(opt protoreflect.ExtensionType) (v string) {
			proto.RangeExtensions(method.Desc.Options(), func(et protoreflect.ExtensionType, a any) bool {
				if et.TypeDescriptor().FullName() == opt.TypeDescriptor().FullName() {
					v = a.(string)
					return false
				}
				return true
			})
			return
		}
		get := getOptString(gin.E_Get)
		post := getOptString(gin.E_Post)
		put := getOptString(gin.E_Put)
		del := getOptString(gin.E_Delete)
		customs := proto.GetExtension(method.Desc.Options(), gin.E_Custom).([]*gin.ApiPath)
		count := 0
		if len(get) > 0 {
			count++
			apiMethod = "GET"
			apiPath = get
		}
		if len(post) > 0 {
			if count > 0 {
				err = multierr.Append(err, fmt.Errorf("method %s.%s has more than one http rule[%s,%s]", file.GoPackageName, method.GoName, post, apiMethod))
			}
			count++
			apiMethod = "POST"
			apiPath = post

		}
		if len(put) > 0 {
			if count > 0 {
				err = multierr.Append(err, fmt.Errorf("method %s.%s has more than one http rule[%s,%s]", file.GoPackageName, method.GoName, put, apiMethod))
			}
			count++
			apiMethod = "PUT"
			apiPath = put
		}
		if len(del) > 0 {
			if count > 0 {
				err = multierr.Append(err, fmt.Errorf("method %s.%s has more than one http rule[%s,%s]", file.GoPackageName, method.GoName, del, apiMethod))
			}
			count++
			apiMethod = "DELETE"
			apiPath = del
		}
		if len(customs) > 0 {
			if count > 0 {
				err = multierr.Append(err, fmt.Errorf("method %s.%s has more than one http rule[%s,%s]", file.GoPackageName, method.GoName, customs[0].Path, apiMethod))
			}
			count++
			apiMethod = customs[0].Method
			apiPath = customs[0].Path
		}
		apiPath = filepath.Clean(filepath.Join("/", apiGroup, apiPath))
		apiMethod = strings.ToUpper(apiMethod)
		apiPath = strings.ToLower(apiPath)
		if apiPath == "/" {
			apiPath = ""
			apiMethod = ""
		}
		return
	}
	for _, service := range file.Services {
		svc := &gen_rpcx.Service{
			Name:    service.GoName,
			StubPkg: utils.PascalToSnake(service.GoName) + "_client",
			StubGin: utils.PascalToSnake(service.GoName) + "_gin",
			Methods: make([]*gen_rpcx.Method, 0, len(service.Methods)),
		}
		for _, method := range service.Methods {
			m := &gen_rpcx.Method{
				Name: method.GoName,
			}
			if method.Input.GoIdent.GoImportPath != file.GoImportPath {
				m.RQ = filepath.Base(string(method.Input.GoIdent.GoImportPath)) + "." + method.Input.GoIdent.GoName
				rpc.Import(string(method.Input.GoIdent.GoImportPath), filepath.Base(string(method.Input.GoIdent.GoImportPath)))
			} else {
				m.RQ = method.Input.GoIdent.GoName
			}
			if method.Output.GoIdent.GoImportPath != file.GoImportPath {
				m.RS = filepath.Base(string(method.Output.GoIdent.GoImportPath)) + "." + method.Output.GoIdent.GoName
				rpc.Import(string(method.Output.GoIdent.GoImportPath), filepath.Base(string(method.Output.GoIdent.GoImportPath)))
			} else {
				m.RS = method.Output.GoIdent.GoName
			}
			for _, field := range method.Input.Fields {
				//log.Println("field.GoIdent.GoName: ", field.GoIdent.GoName)
				if field.GoName == "AppInfo" {
					m.AppInfo = true
					break
				}
			}
			if !m.AppInfo {
				if m.RQ == "comm.AppInfo" {
					m.AppRQ = true
				}
			}
			m.ApiMethod, m.ApiPath = getApiConfig(method)
			if !svc.GenApi && len(m.ApiPath) > 0 {
				svc.GenApi = true
			}
			svc.Methods = append(svc.Methods, m)
		}
		rpc.Services = append(rpc.Services, svc)
	}
	for _, svc := range rpc.Services {
		if !svc.GenApi {
			continue
		}
		for _, method := range svc.Methods {
			if len(method.ApiMethod) == 0 {
				err = multierr.Append(err, fmt.Errorf("method %s.%s.%s has no http rule", file.GoPackageName, svc.Name, method.Name))
			}
		}
	}
	return
}
