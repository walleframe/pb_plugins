{{ $Name := .Name}} {{$typ := .TypeSet.BaseType }} {{$farg := .TypeSet.Args}}
{{if $farg }}
func Merge{{$Name}}Member({{- range $i,$arg := $farg -}}{{$arg.ArgName}} {{$arg.ArgType }}, {{- end -}}) string {
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
func Split{{$Name}}Member(val string)({{- range $i,$arg := $farg -}}{{$arg.ArgName}} {{$arg.ArgType }}, {{- end -}}err error) {
	items := strings.Split(val, ":")
	if len(items) != {{len $farg}} {
		err = errors.New("invalid {{$Name}} mem value")
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
func (x *x{{$Name}}) SAdd(ctx context.Context, {{GenTypeTemplate "set_func_arg" .}}) (bool, error) {
	n, err := x.rds.SAdd(ctx, x.key, {{GenTypeTemplate "set_str_arg" .}}).Result()
	if err != nil {
		return false, err
	}
	return n == 1, nil
}

func (x *x{{$Name}}) SCard(ctx context.Context) (int64, error) {
	return x.rds.SCard(ctx, x.key).Result()
}

func (x *x{{$Name}}) SRem(ctx context.Context, {{GenTypeTemplate "set_func_arg" .}}) (bool, error) {
	n, err := x.rds.SRem(ctx, x.key, {{GenTypeTemplate "set_str_arg" .}}).Result()
	if err != nil {
		return false, err
	}
	return n == 1, nil
}

func (x *x{{$Name}}) SIsMember(ctx context.Context, {{GenTypeTemplate "set_func_arg" .}}) (bool, error) {
	return x.rds.SIsMember(ctx, x.key, {{GenTypeTemplate "set_str_arg" .}}).Result()
}

func (x *x{{$Name}}) SPop(ctx context.Context) ({{GenTypeTemplate "set_func_arg" .}} err error) {
{{if $farg}}
	 v, err := x.rds.SPop(ctx, x.key).Result()
	 return Split{{$Name}}Member(v)
{{else if eq $typ.Type "string" -}}
	return x.rds.SPop(ctx, x.key).Result()
{{else -}}
	v, err := x.rds.SPop(ctx, x.key).Result()
	if err != nil {
		return
	}
	return rdconv.StringTo{{$typ.RedisFunc}}(v)
{{end -}}
}
func (x *x{{$Name}}) SRandMember(ctx context.Context) ({{GenTypeTemplate "set_func_arg" .}} err error) {
{{if $farg}}
	 v, err := x.rds.SRandMember(ctx, x.key).Result()
	 return Split{{$Name}}Member(v)
{{else if eq $typ.Type "string" -}}
	return x.rds.SRandMember(ctx, x.key).Result()
{{else -}}
	v, err := x.rds.SRandMember(ctx, x.key).Result()
	if err != nil {
		return
	}
	return rdconv.StringTo{{$typ.RedisFunc}}(v)
{{end -}}
}

{{if not $farg}}
func (x *x{{$Name}}) SRandMemberN(ctx context.Context, count int) (vals []{{$typ.Type}}, err error) {
{{if eq $typ.Type "string" -}}
	return x.rds.SRandMemberN(ctx, x.key, int64(count)).Result()
{{else -}}
	ret, err := x.rds.SRandMemberN(ctx, x.key, int64(count)).Result()
	if err != nil {
		return
	}
	for _, v := range ret {
		val, err := rdconv.StringTo{{$typ.RedisFunc}}(v)
		if err != nil {
			return nil, err
		}
		vals = append(vals, val)
	}
	return
{{end -}}
}

func (x *x{{$Name}}) SMembers(ctx context.Context, count int) (vals []{{$typ.Type}}, err error) {
{{if eq $typ.Type "string" -}}
	return x.rds.SMembers(ctx, x.key).Result()
{{else -}}
	ret, err := x.rds.SMembers(ctx, x.key).Result()
	if err != nil {
		return
	}
	for _, v := range ret {
		val, err := rdconv.StringTo{{$typ.RedisFunc}}(v)
		if err != nil {
			return nil, err
		}
		vals = append(vals, val)
	}
	return
{{end -}}
}

func (x *x{{$Name}}) SScan(ctx context.Context, match string, count int) (vals []{{$typ.Type}}, err error) {
	cursor := uint64(0)
	var ret []string
	for {
		ret, cursor, err = x.rds.SScan(ctx, x.key, cursor, match, int64(count)).Result()
		if err != nil {
			return nil, err
		}
{{if eq $typ.Type "string" -}}
		vals = append(vals, ret...)
{{else -}}
		for _, v := range ret {
			val, err := rdconv.StringTo{{$typ.RedisFunc}}(v)
			if err != nil {
				return nil, err
			}
			vals = append(vals, val)
		}
{{end -}}
		if cursor == 0 {
			break
		}
	}
	return
}
{{end}}

func (x *x{{$Name}}) SScanRange(ctx context.Context, match string, count int, filter func({{GenTypeTemplate "set_func_arg" .}}) bool) (err error) {
	cursor := uint64(0)
	var ret []string
	for {
		ret, cursor, err = x.rds.SScan(ctx, x.key, cursor, match, int64(count)).Result()
		if err != nil {
			return err
		}
		for _, v := range ret {
{{- if $farg}}
			{{range $i,$arg := $farg -}}{{$arg.ArgName}}, {{- end -}}err := Split{{$Name}}Member(v)
			if err != nil {
				return err
			}
	 		if !filter({{- range $i,$arg := $farg -}}{{$arg.ArgName}}, {{- end -}}) {
				return nil
			}
{{- else if eq $typ.Type "string" }}
			if !filter(v) {
				return nil
			}
{{- else }}
			val, err := rdconv.StringTo{{$typ.RedisFunc}}(v)
			if err != nil {
				return err
			}
			if !filter(val) {
				return nil
			}
{{end -}}
		}
		if cursor == 0 {
			break
		}
	}
	return
}