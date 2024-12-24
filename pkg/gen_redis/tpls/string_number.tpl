
func (x *x{{.Name}}) Incr(ctx context.Context) ({{.TypeString.Type}}, error) {
	n,err := x.rds.Incr(ctx, x.key).Result()
	return {{.TypeString.Type}}(n), err
}

func (x *x{{.Name}}) IncrBy(ctx context.Context, val int) (_ {{.TypeString.Type}},err error) {
	cmd := redis.NewIntCmd(ctx, "incrby", x.key, strconv.FormatInt(int64(val), 10)) {{- Import "strconv" "strconv.Format"}}
	err = x.rds.Process(ctx, cmd)
	if err != nil {
		return
	}
	return {{.TypeString.Type}}(cmd.Val()), nil
}

func (x *x{{.Name}}) Decr(ctx context.Context) ({{.TypeString.Type}}, error) {
	n,err := x.rds.Decr(ctx, x.key).Result()
	return {{.TypeString.Type}}(n), err
}

func (x *x{{.Name}}) DecrBy(ctx context.Context, val int) (_ {{.TypeString.Type}}, err error) {
	cmd := redis.NewIntCmd(ctx, "decrby", x.key, strconv.FormatInt(int64(val), 10)) {{- Import "strconv" "strconv.FormatInt"}}
	err = x.rds.Process(ctx, cmd)
	if err != nil {
		return
	}
	return {{.TypeString.Type}}(cmd.Val()), nil
}

func (x *x{{.Name}}) Get(ctx context.Context) ({{.TypeString.Type}}, error) {
	data,err := x.rds.Get(ctx, x.key).Result()
	if err != nil {
		return 0,err
	}
	val,err := strconv.Parse{{if .TypeString.Signed}}Int{{else}}Uint{{end}}(data, 10, 64)
	if err != nil {
		return 0,err
	}
	return {{.TypeString.Type}}(val), nil
}

func (x *x{{.Name}}) Set(ctx context.Context, val {{.TypeString.Type}}, expire time.Duration) error { {{- Import "time" "time.Duration"}}
	return x.rds.Set(ctx, x.key, strconv.Format{{if .TypeString.Signed}}Int(int64{{else}}Uint(uint64{{end}}(val), 10), expire).Err()
}

func (x *x{{.Name}}) SetNX(ctx context.Context, val {{.TypeString.Type}}, expire time.Duration) (bool, error) {
	return x.rds.SetNX(ctx, x.key, strconv.Format{{if .TypeString.Signed}}Int(int64{{else}}Uint(uint64{{end}}(val), 10), expire).Result()
}

func (x *x{{.Name}}) SetEx(ctx context.Context, val {{.TypeString.Type}}, expire time.Duration) error {
	return x.rds.SetEx(ctx, x.key, strconv.Format{{if .TypeString.Signed}}Int(int64{{else}}Uint(uint64{{end}}(val), 10), expire).Err()
}