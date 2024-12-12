{{ $Name := .Name}} {{$field := .TypeHash.HashDynamic.Field }} {{$value := .TypeHash.HashDynamic.Value }} {{$farg := .TypeHash.HashDynamic.FieldArgs}} {{$varg := .TypeHash.HashDynamic.ValueArgs}}
{{- if $value -}}
	{{- if eq $value.Type "string" -}}
		value
	{{- else -}}
		rdconv.{{$value.RedisFunc}}ToString(value)
	{{- end -}}
{{- else -}}
	Merge{{$Name}}Value({{- range $i,$arg := $varg -}}
		{{$arg.ArgName}},
	{{- end -}})
{{- end -}}
