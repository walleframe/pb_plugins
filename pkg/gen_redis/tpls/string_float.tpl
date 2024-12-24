
func (x *x{{.Name}}) Get(ctx context.Context) ({{.TypeString.Type}}, error) {
	data,err := x.rds.Get(ctx, x.key).Result()
	if err != nil {
		return 0,err
	}
	val,err := strconv.ParseFloat(data, 64)
	if err != nil {
		return 0,err
	}
	return {{.TypeString.Type}}(val), nil
}

func (x *x{{.Name}}) IncrBy(ctx context.Context, val int) (_ {{.TypeString.Type}},err error) {
	cmd := redis.NewFloatCmd(ctx, "incrbyfloat", x.key, strconv.FormatInt(int64(val), 10)) {{Import "strconv" "strconv.FormatInt"}}
	err = x.rds.Process(ctx, cmd)
	if err != nil {
		return
	}
	return {{.TypeString.Type}}(cmd.Val()), nil
}

func (x *x{{.Name}}) Set(ctx context.Context, val {{.TypeString.Type}}, expire time.Duration) error { {{- Import "time" "time.Duration"}}
	return x.rds.Set(ctx, x.key, rdconv.Float64ToString(float64(val)), expire).Err()
}

func (x *x{{.Name}}) SetNX(ctx context.Context, val {{.TypeString.Type}}, expire time.Duration) (bool, error) {
	return x.rds.SetNX(ctx, x.key, rdconv.Float64ToString(float64(val)), expire).Result()
}

func (x *x{{.Name}}) SetEx(ctx context.Context, val {{.TypeString.Type}}, expire time.Duration) error {
	return x.rds.SetEx(ctx, x.key, rdconv.Float64ToString(float64(val)), expire).Err()
} {{UsePackage "rdconv" "Float64ToString"}}