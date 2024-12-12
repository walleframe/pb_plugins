{{ $Name := .Name}} {{$mem := .TypeZSet.Member }} {{$farg := .TypeZSet.Args}}
	{{- if $mem -}}
		mem {{$mem.Type}},
	{{- else if $farg -}}
		{{- range $i,$arg := $farg -}}
			{{$arg.ArgName}} {{$arg.ArgType}},
		{{- end -}}
	{{- end -}}