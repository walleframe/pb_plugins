package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"go.uber.org/multierr"
	"google.golang.org/protobuf/compiler/protogen"

	"github.com/spf13/pflag"
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
	flags.StringVar(&cfg.StubPath, "stub_path", cfg.StubPath, "stub file path")
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
	// 扫描要生成的服务
	genRPCXServcies(file, zlog)
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

func genRPCXServcies(file *protogen.File, rpc *gen_rpcx.File) {
	for _, service := range file.Services {
		svc := &gen_rpcx.Service{
			Name:    service.GoName,
			StubPkg: utils.PascalToSnake(service.GoName) + "_client",
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
			svc.Methods = append(svc.Methods, m)

		}
		rpc.Services = append(rpc.Services, svc)
	}
}
