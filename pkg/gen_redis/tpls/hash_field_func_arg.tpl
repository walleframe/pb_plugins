{{ $Name := .Name}} {{$field := .TypeHash.HashDynamic.Field }} {{$value := .TypeHash.HashDynamic.Value }} {{$farg := .TypeHash.HashDynamic.FieldArgs}} {{$varg := .TypeHash.HashDynamic.ValueArgs}}
{{- if $field -}}
	field {{$field.Type}},
{{- else -}}
	{{- range $i,$arg := $farg -}}
		{{$arg.ArgName}} {{$arg.ArgType }},
	{{- end -}}
{{- end -}}