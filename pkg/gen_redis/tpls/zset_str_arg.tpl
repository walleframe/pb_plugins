{{ $Name := .Name}} {{$mem := .TypeZSet.Member }} {{$farg := .TypeZSet.Args}}
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