// Code generated by {{.ToolName}} {{.Version}}. DO NOT EDIT.{{$svc := .Stub}}{{$Name := $svc.Name}}{{$methods := $svc.Methods}}
package {{$svc.StubPkg}}

$Import-Packages$


var (
	//global{{Title $Name}}OP atomic.Pointer[{{.Package}}.{{$Name}}XClient]
	xclient rpcx_client.XClient
)

func init() {
	{{$svc.StubCtrl}}.RegisterClient({{.Package}}.{{$Name}}ServerName, UpdateClient)
}

func UpdateClient(xc rpcx_client.XClient) {
	xclient = xc
}


{{range $i,$method := $methods}}
{{Doc $method.Doc}}
func {{$method.Name}}(ctx context.Context, rq *{{$method.RQ}}) (rs *{{$method.RS}}, err error){
	rs = &{{$method.RS}}{}
	err = xclient.Call(ctx, "{{$method.Name}}", rq, rs)
	return
}{{end}}