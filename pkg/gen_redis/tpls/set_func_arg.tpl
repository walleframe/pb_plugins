{{ $Name := .Name}} {{$mem := .TypeSet.BaseType }} {{$farg := .TypeSet.Args}}
	{{- if $mem -}}
		mem {{$mem.Type}},
	{{- else if $farg -}}
		{{- range $i,$arg := $farg -}}
			{{$arg.ArgName}} {{$arg.ArgType}},
		{{- end -}}
	{{- end -}}