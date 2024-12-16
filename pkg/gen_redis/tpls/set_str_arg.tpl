{{ $Name := .Name}} {{$mem := .TypeSet.BaseType }} {{$farg := .TypeSet.Args}}
	{{- if $mem -}}
		{{- if eq $mem.Type "string" -}}
			mem
		{{- else -}}
			rdconv.{{$mem.RedisFunc}}ToString(mem)
		{{- end -}}
	{{- else if $farg -}}
		Merge{{$Name}}Member({{- range $i,$arg := $farg -}}
			{{$arg.ArgName}},
		{{- end -}})
	{{- end -}}