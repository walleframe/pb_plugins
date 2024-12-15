{{ $Name := .Name}}
func (x *x{{$Name}}) Set{{Title .TypeHash.HashObject.Name}}(ctx context.Context, obj *{{.TypeHash.HashObject.Type}}) (err error) {
	n, err := x.rds.HSet(ctx, x.key, {{ range $i,$field := .TypeHash.HashObject.Fields }}
		_{{$Name}}_{{Title $field.Name}}, rdconv.{{$field.RedisFunc}}ToString(obj.{{Title $field.Name}}),{{end}}
	).Result()
	if err != nil {
		return err
	}
	if n != {{len .TypeHash.HashObject.Fields}} {
		return errors.New("set {{Title .TypeHash.HashObject.Name}} failed")
	}
	return
}

{{if .TypeHash.HashObject.HGetAll }}
func (x *x{{$Name}}) Get{{Title .TypeHash.HashObject.Name}}(ctx context.Context) (*{{.TypeHash.HashObject.Type}}, error) {
	ret, err := x.rds.HGetAll(ctx, x.key).Result()
	if err != nil {
		return nil, err
	}
	obj := &{{.TypeHash.HashObject.Type}}{}
{{ range $i,$field := .TypeHash.HashObject.Fields }}
	if val, ok := ret[_{{$Name}}_{{Title $field.Name}}]; ok {
{{if eq $field.Type "string"}}
		obj.{{Title $field.Name}} = val
{{ else -}}
		obj.{{Title $field.Name}}, err = rdconv.StringTo{{$field.RedisFunc}}(val)
		if err != nil {
			return nil, fmt.Errorf("parse {{$Name}}.{{Title $field.Name}} failed,%w", err)
		}
{{ end -}}
	}
{{end}}
	return obj, nil
}
{{end}}

func (x *x{{$Name}}) MGet{{Title .TypeHash.HashObject.Name}}(ctx context.Context) (*{{.TypeHash.HashObject.Type}}, error) {
	ret, err := x.rds.HMGet(ctx, x.key, {{ range $i,$field := .TypeHash.HashObject.Fields }} _{{$Name}}_{{Title $field.Name}},{{end}}).Result()
	if err != nil {
		return nil, err
	}
	obj := &{{.TypeHash.HashObject.Type}}{}
{{ range $i,$field := .TypeHash.HashObject.Fields }}
	if len(ret) > {{$i}} && ret[{{$i}}] != nil {
		obj.{{Title $field.Name}}, err = rdconv.AnyTo{{$field.RedisFunc}}(ret[ {{$i}} ])
		if err != nil {
			return nil, fmt.Errorf("parse {{$Name}}.{{Title $field.Name}} failed,%w", err)
		}
	}
{{end}}
	return obj, nil
}

{{range $i,$field := .TypeHash.HashObject.Fields}}
func (x *x{{$Name}}) Get{{Title $field.Name}}(ctx context.Context)(_ {{$field.Type}},err error) {
{{if eq $field.Type "string" -}}
	return x.rds.HGet(ctx, x.key, _{{$Name}}_{{Title $field.Name}}).Result()
{{ else -}}
	val, err := x.rds.HGet(ctx, x.key, _{{$Name}}_{{Title $field.Name}}).Result()
	if err != nil {
		return
	}
	return rdconv.StringTo{{$field.RedisFunc}}(val)
{{end}}
}
func (x *x{{$Name}}) Set{{Title $field.Name}}(ctx context.Context, val {{$field.Type}}) (err error) {
	n, err := x.rds.HSet(ctx, x.key,  _{{$Name}}_{{Title $field.Name}}, rdconv.{{$field.RedisFunc}}ToString(val)).Result()
	if err != nil {
		return err
	}
	if n != 1 {
		return errors.New("set {{$Name}}.{{Title $field.Name}} failed")
	}
	return nil
}
{{if $field.Number}}
func (x *x{{$Name}}) IncrBy{{Title $field.Name}}(ctx context.Context, incr int) ({{$field.Type}}, error) {
	num, err := x.rds.HIncrBy(ctx, x.key,  _{{$Name}}_{{Title $field.Name}}, int64(incr)).Result()
	if err != nil {
		return 0, err
	}
	return {{$field.Type}}(num), nil
}
{{end}}
{{end}}{{UsePackage "rdconv" "Float64ToString"}}