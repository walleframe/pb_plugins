{{ $Name := .Name}} {{$field := .TypeHash.HashDynamic.Field }} {{$value := .TypeHash.HashDynamic.Value }} {{$farg := .TypeHash.HashDynamic.FieldArgs}} {{$varg := .TypeHash.HashDynamic.ValueArgs}}

{{if $farg }}
func Merge{{$Name}}Field({{- range $i,$arg := $farg -}}{{$arg.ArgName}} {{$arg.ArgType }}, {{- end -}}) string {
	buf := util.Builder{}
	buf.Grow({{.KeySize}})
{{- range $i,$arg := $farg -}}
{{- if gt $i 0 -}}
	buf.WriteByte(':')
{{- end }}
	{{$arg.FormatCode "buf"}}
{{end -}}
	return buf.String()
}
func Split{{$Name}}Field(val string)({{- range $i,$arg := $farg -}}{{$arg.ArgName}} {{$arg.ArgType }}, {{- end -}}err error) {
	items := strings.Split(val, ":")
	if len(items) != {{len $farg}} {
		err = errors.New("invalid {{$Name}} field value")
		return
	}
{{ range $i,$arg := $farg -}}
{{- if eq $arg.ArgType "string" -}}
	{{- $arg.ArgName}} = items[{{$i}}]
{{ else -}}
	{{$arg.ArgName}}, err = rdconv.StringTo{{Title $arg.ArgType}}(items[{{$i}}])
	if err != nil {
		return
	}
{{ end -}}
{{end -}}
	return
}
{{end}}

{{if $varg }}
func Merge{{$Name}}Value({{- range $i,$arg := $varg -}}{{$arg.ArgName}} {{$arg.ArgType}},{{- end -}}) string {
	buf := util.Builder{}
	buf.Grow({{.KeySize}})
{{- range $i,$arg := $varg -}}
{{- if gt $i 0 -}}
	buf.WriteByte(':')
{{- end }}
	{{$arg.FormatCode "buf"}}
{{end -}}
	return buf.String()
}
func Split{{$Name}}Value(val string)({{- range $i,$arg := $varg -}}{{$arg.ArgName}} {{$arg.ArgType }}, {{- end -}}err error) {
	items := strings.Split(val, ":")
	if len(items) != {{len $varg}} {
		err = errors.New("invalid {{$Name}} field value")
		return
	}
{{ range $i,$arg := $varg -}}
{{- if eq $arg.ArgType "string" -}}
	{{$arg.ArgName}} = items[{{$i}}]
{{ else -}}
	{{$arg.ArgName}}, err = rdconv.StringTo{{Title $arg.ArgType}}(items[{{$i}}])
	if err != nil {
		return
	}
{{ end -}}
{{end -}}
	return
}
{{end}}

func (x *x{{.Name}}) GetField(ctx context.Context, {{GenTypeTemplate "hash_field_func_arg" .}}) ({{GenTypeTemplate "hash_value_func_arg" .}} err error) {
{{- if $value }}
{{- if eq $value.Type "string"}}
	return x.rds.HGet(ctx, x.key, {{GenTypeTemplate "hash_filed_str_arg" .}}).Result()
{{- else }}
	v, err := x.rds.HGet(ctx, x.key, {{GenTypeTemplate "hash_filed_str_arg" .}}).Result()
	if err != nil {
		return
	}{{- if $value.MarshalPkg }}
	value = &{{$value.Type}}{}
	err = {{$value.MarshalPkg}}.Unmarshal(util.StringToBytes(v), value)
	if err != nil {
		return
	}
	return{{else}}
	return rdconv.StringTo{{$value.RedisFunc}}(v){{end}}
{{- end}}
{{- else }}
	v, err := x.rds.HGet(ctx, x.key, {{GenTypeTemplate "hash_filed_str_arg" .}}).Result()
	if err != nil {
		return
	}
	return Split{{$Name}}Value(v)
{{ end -}}
}
func (x *x{{.Name}}) SetField(ctx context.Context, {{GenTypeTemplate "hash_field_func_arg" .}} {{GenTypeTemplate "hash_value_func_arg" .}}) (err error) {
{{- if $value}}{{if $value.MarshalPkg }}
	data, err := {{ $value.MarshalPkg }}.Marshal(value)
	if err != nil {
		return err
	} {{- end }}{{end}}
	num, err := x.rds.HSet(ctx, x.key, {{GenTypeTemplate "hash_filed_str_arg" .}}, {{GenTypeTemplate "hash_value_str_arg" .}}).Result()
	if err != nil {
		return err
	}
	if num != 1 {
		return errors.New("set field failed")
	}
	return nil
}

{{if $field }}
func (x *x{{.Name}}) HKeys(ctx context.Context) (vals []{{$field.Type}}, err error) {
{{if eq $field.Type "string" -}}
	return x.rds.HKeys(ctx, x.key).Result()
{{- else -}}
	ret, err := x.rds.HKeys(ctx, x.key).Result()
	if err != nil {
		return
	}
	for _, v := range ret {
		key,err := rdconv.StringTo{{$field.RedisFunc}}(v)
		if err != nil {
			return nil, err
		}
		vals = append(vals, key)
	}
	return
{{- end}}
}


func (x *x{{.Name}}) HKeysRange(ctx context.Context, filter func({{$field.Type}}) bool)(error) {
	ret, err := x.rds.HKeys(ctx, x.key).Result()
	if err != nil {
		return err
	}
	for _, v := range ret {
{{if eq $field.Type "string" -}}
		if !filter(v) {
			return nil
		}
{{ else -}}
		key,err := rdconv.StringTo{{$field.RedisFunc}}(v)
		if err != nil {
			return err
		}
		if !filter(key) {
			return nil
		}
{{ end -}}
	}
	return nil
}

{{else}}

func (x *x{{.Name}}) HKeysRange(ctx context.Context, filter func({{GenTypeTemplate "hash_field_func_arg" .}})bool) (err error) {
	ret, err := x.rds.HKeys(ctx, x.key).Result()
	if err != nil {
		return
	}
	for _, v := range ret {
		{{ range $i,$arg := $farg -}} {{$arg.ArgName}}, {{end}}err := Split{{$Name}}Field(v)
		if err != nil {
			return err
		}
		if !filter({{ range $i,$arg := $farg -}} {{$arg.ArgName}}, {{end}}) {
			return nil
		}
	}
	return nil
}

{{end}}

{{if $value}}
func (x *x{{.Name}}) HVals(ctx context.Context) (vals []{{if $value.MarshalPkg}}*{{end}}{{$value.Type}}, err error) {
{{- if $value.MarshalPkg}}
	ret, err := x.rds.HVals(ctx, x.key).Result()
	if err != nil {
		return
	}
	for _, v := range ret {
		val := &{{$value.Type}}{}
		err = {{ $value.MarshalPkg }}.Unmarshal(util.StringToBytes(v), val)
		if err != nil {
			return nil, err
		}
		vals = append(vals, val)
	}
	return
{{- else if eq $value.Type "string"}}
	return x.rds.HVals(ctx, x.key).Result()
{{- else }}
	ret, err := x.rds.HVals(ctx, x.key).Result()
	if err != nil {
		return
	}
	for _, v := range ret {
		val,err := rdconv.StringTo{{$value.RedisFunc}}(v)
		if err != nil {
			return nil, err
		}
		vals = append(vals, val)
	}
	return
{{- end}}
}

func (x *x{{.Name}}) HValsRange(ctx context.Context, filter func({{if $value.MarshalPkg}}*{{end}}{{$value.Type}}) bool)(error) {
	ret, err := x.rds.HVals(ctx, x.key).Result()
	if err != nil {
		return err
	}
	for _, v := range ret {
{{ if $value.MarshalPkg}}
		val :=&{{$value.Type}}{}
		err = {{ $value.MarshalPkg }}.Unmarshal(util.StringToBytes(v), val)
		if err != nil {
			return err
		}
		if !filter(val) {
			return nil
		}
{{else if eq $value.Type "string" }}
		if !filter(v) {
			return nil
		}
{{ else -}}
		val,err := rdconv.StringTo{{$value.RedisFunc}}(v)
		if err != nil {
			return err
		}
		if !filter(val) {
			return nil
		}
{{ end -}}
	}
	return nil
}
{{else}}
func (x *x{{.Name}}) HValsRange(ctx context.Context, filter func({{GenTypeTemplate "hash_value_func_arg" .}})bool)(err error) {
	ret, err := x.rds.HVals(ctx, x.key).Result()
	if err != nil {
		return
	}
	for _, v := range ret {
		{{ range $i,$arg := $varg -}} {{$arg.ArgName}}, {{end}}err := Split{{$Name}}Value(v)
		if err != nil {
			return err
		}
		if !filter({{ range $i,$arg := $varg -}} {{$arg.ArgName}}, {{end}}) {
			return nil
		}
	}
	return nil
}

{{end}}

func (x *x{{.Name}}) HExists(ctx context.Context, {{GenTypeTemplate "hash_field_func_arg" .}}) (bool, error) {
	return x.rds.HExists(ctx, x.key, {{GenTypeTemplate "hash_filed_str_arg" .}}).Result()
}

func (x *x{{.Name}}) HDel(ctx context.Context, {{GenTypeTemplate "hash_field_func_arg" .}}) (bool, error) {
	n, err := x.rds.HDel(ctx, x.key, {{GenTypeTemplate "hash_filed_str_arg" .}}).Result()
	if err != nil {
		return false, err
	}
	return n == 1, nil
}

func (x *x{{.Name}}) HLen(ctx context.Context) (count int64, err error) {
	return x.rds.HLen(ctx, x.key).Result()
}
{{if $field }}
	func (x *x{{.Name}}) HRandField(ctx context.Context, count int) (vals []{{$field.Type}}, err error) {
	{{if eq $field.Type "string" -}}
		return x.rds.HRandField(ctx, x.key, count).Result()
	{{- else -}}
		ret, err := x.rds.HRandField(ctx, x.key, count).Result()
		if err != nil {
			return
		}
		for _, v := range ret {
			key, err := rdconv.StringTo{{$field.RedisFunc}}(v)
			if err != nil {
				return nil, err
			}
			vals = append(vals, key)
		}
		return
	{{- end}}
	}
{{else}}
	func (x *x{{.Name}}) HRandFieldRange(ctx context.Context, count int,filter func({{GenTypeTemplate "hash_field_func_arg" .}})bool) (err error) {
		ret, err := x.rds.HRandField(ctx, x.key, count).Result()
		if err != nil {
			return
		}
		for _, v := range ret {
			{{ range $i,$arg := $farg -}} {{$arg.ArgName}}, {{end}}err := Split{{$Name}}Field(v)
			if err != nil {
				return err
			}
			if !filter({{ range $i,$arg := $farg -}} {{$arg.ArgName}}, {{end}}) {
				return nil
			}
		}
		return
	}
{{end}}

{{if .TypeHash.HashDynamic.GenMap }}
func (x *x{{.Name}}) HRandFieldWithValues(ctx context.Context, count int) (vals map[{{$field.Type}}]{{if $value.MarshalPkg}}*{{end}}{{$value.Type}}, err error) {
	ret, err := x.rds.HRandFieldWithValues(ctx, x.key, count).Result()
	if err != nil {
		return
	}
	vals = make(map[{{$field.Type}}]{{if $value.MarshalPkg}}*{{end}}{{$value.Type}}, len(ret))
	for _, v := range ret {
{{if eq $field.Type "string" -}}
		key := v.Key
{{ else -}}
		key, err := rdconv.StringTo{{$field.RedisFunc}}(v.Key)
		if err != nil {
			return nil, err
		}
{{end -}}
{{ if $value.MarshalPkg }}
		value := &{{$value.Type}}{}
		err = {{$value.MarshalPkg}}.Unmarshal(util.StringToBytes(v.Value), value)
		if err != nil {
			return nil, err
		}
		vals[key] = value
{{ else if eq $value.Type "string" -}}
		vals[key] = v.Value
{{else -}}
		val, err := rdconv.StringTo{{$value.RedisFunc}}(v.Value)
		if err != nil {
			return nil, err
		}
		vals[key] = val
{{end -}}
	}
	return
}
{{end}}

func (x *x{{.Name}}) HRandFieldWithValuesRange(ctx context.Context, count int, filter func({{GenTypeTemplate "hash_field_func_arg" .}} {{GenTypeTemplate "hash_value_func_arg" .}})bool) (err error) {
	ret, err := x.rds.HRandFieldWithValues(ctx, x.key, count).Result()
	if err != nil {
		return
	}
	for _, v := range ret {
{{if $field}}
{{if eq $field.Type "string" -}}
		key := v.Key
{{ else -}}
		key, err := rdconv.StringTo{{$field.RedisFunc}}(v.Key)
		if err != nil {
			return  err
		}
{{end -}}
{{else}}
	{{ range $i,$arg := $farg -}} {{$arg.ArgName}}, {{end}}err := Split{{$Name}}Field(v.Key)
	if err != nil {
		return err
	}
{{end}}

{{if $value}}
{{ if $value.MarshalPkg }}
		val := &{{$value.Type}}{}
		err = {{$value.MarshalPkg }}.Unmarshal(util.StringToBytes(v.Value), val)
		if err != nil {
			return err
		}
{{else if eq $value.Type "string" -}}
		val := v.Value
{{else -}}
		val, err := rdconv.StringTo{{$value.RedisFunc}}(v.Value)
		if err != nil {
			return err
		}
{{end -}}
{{else}}
		{{ range $i,$arg := $varg -}} {{$arg.ArgName}}, {{end}}err := Split{{$Name}}Value(v.Value)
		if err != nil {
			return err
		}
{{end}}

		if !filter({{if $field}} key, {{else}} {{ range $i,$arg := $farg -}} {{$arg.ArgName}}, {{end}}{{end -}}
			{{- if $value}} val {{else}} {{ range $i,$arg := $varg -}} {{$arg.ArgName}}, {{end}} {{end}}) {
			return nil
		}

	}
	return
}
{{if .TypeHash.HashDynamic.GenMap }}
func (x *x{{.Name}}) HScan(ctx context.Context, match string, count int) (vals map[{{$field.Type}}]{{if $value.MarshalPkg}}*{{end}}{{$value.Type}}, err error) {
	cursor := uint64(0)
	vals = make(map[{{$field.Type}}]{{if $value.MarshalPkg}}*{{end}}{{$value.Type}})
	var kvs []string
	for {
		kvs, cursor, err = x.rds.HScan(ctx, x.key, cursor, match, int64(count)).Result()
		if err != nil {
			return nil, err
		}
		for k := 0; k < len(kvs); k += 2 {
{{if eq $field.Type "string" -}}
			key := kvs[k]
{{ else -}}
			key, err := rdconv.StringTo{{$field.RedisFunc}}(kvs[k])
			if err != nil {
				return nil, err
			}
{{end -}}
{{ if $value.MarshalPkg }}
			val := &{{$value.Type}}{}
			err = {{$value.MarshalPkg }}.Unmarshal(util.StringToBytes(kvs[k+1]), val)
			if err != nil {
				return nil, err
			}
{{else if eq $value.Type "string" -}}
			val := kvs[k+1]
{{else -}}
			val, err := rdconv.StringTo{{$value.RedisFunc}}(kvs[k+1])
			if err != nil {
				return nil, err
			}
{{end -}}
			vals[key] = val
		}
		if cursor == 0 {
			break
		}
	}

	return
}
{{end}}
func (x *x{{.Name}}) HScanRange(ctx context.Context, match string, count int, filter func({{GenTypeTemplate "hash_field_func_arg" .}} {{GenTypeTemplate "hash_value_func_arg" .}})bool) (err error) {
	cursor := uint64(0)
	var kvs []string
	for {
		kvs, cursor, err = x.rds.HScan(ctx, x.key, cursor, match, int64(count)).Result()
		if err != nil {
			return err
		}
		for k := 0; k < len(kvs); k += 2 {
{{if $field}}
{{if eq $field.Type "string" -}}
			key := kvs[k]
{{ else -}}
			key, err := rdconv.StringTo{{$field.RedisFunc}}(kvs[k])
			if err != nil {
				return err
			}
{{end -}}
{{else -}}
	{{ range $i,$arg := $farg -}} {{$arg.ArgName}}, {{end}}err := Split{{$Name}}Field(kvs[k])
	if err != nil {
		return err
	}
{{- end}}
{{if $value}}
{{ if $value.MarshalPkg }}
			val := &{{$value.Type}}{}
			err = {{$value.MarshalPkg }}.Unmarshal(util.StringToBytes(kvs[k+1]), val)
			if err != nil {
				return err
			}
{{else if eq $value.Type "string" -}}
			val := kvs[k+1]
{{else -}}
			val, err := rdconv.StringTo{{$value.RedisFunc}}(kvs[k+1])
			if err != nil {
				return err
			}
{{end -}}
{{else}}
			{{ range $i,$arg := $varg -}} {{$arg.ArgName}}, {{end}}err := Split{{$Name}}Value(kvs[k+1])
			if err != nil {
				return err
			}
{{end}}
			if !filter({{if $field}} key, {{else}} {{ range $i,$arg := $farg -}} {{$arg.ArgName}}, {{end}}{{end -}}
				{{- if $value}} val {{else}} {{ range $i,$arg := $varg -}} {{$arg.ArgName}}, {{end}} {{end}}) {
				return nil
			}
		}
		if cursor == 0 {
			break
		}
	}
	return
}{{UsePackage "rdconv" "ToString/FronString"}}