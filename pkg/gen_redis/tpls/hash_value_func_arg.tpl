{{ $Name := .Name}} {{$field := .TypeHash.HashDynamic.Field }} {{$value := .TypeHash.HashDynamic.Value }} {{$farg := .TypeHash.HashDynamic.FieldArgs}} {{$varg := .TypeHash.HashDynamic.ValueArgs}}
{{- if $value -}}
	value {{ if $value.MarshalPkg }}*{{end}}{{$value.Type}},
{{- else -}}
	{{- range $i,$arg := $varg -}}
		{{$arg.ArgName}} {{$arg.ArgType }},
	{{- end -}}
{{- end -}}