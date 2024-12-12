{{ $Name := .Name}} {{$field := .TypeHash.HashDynamic.Field }} {{$value := .TypeHash.HashDynamic.Value }} {{$farg := .TypeHash.HashDynamic.FieldArgs}} {{$varg := .TypeHash.HashDynamic.ValueArgs}}
{{- if $field -}}
	{{- if eq $field.Type "string" -}}
		field
	{{- else -}}
		rdconv.{{$field.RedisFunc}}ToString(field)
	{{- end -}}
{{- else -}}
	Merge{{$Name}}Field({{- range $i,$arg := $farg -}}
		{{$arg.ArgName}},
	{{- end -}})
{{- end -}}