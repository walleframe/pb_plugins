
func (x *x{{.Name}}) Set(ctx context.Context, msg any, expire time.Duration) error {
	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	return x.rds.Set(ctx, x.key, util.BytesToString(data), expire).Err()
}

func (x *x{{.Name}}) SetNX(ctx context.Context, msg any, expire time.Duration) error {
	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	return x.rds.SetNX(ctx, x.key, util.BytesToString(data), expire).Err()
}

func (x *x{{.Name}}) SetEx(ctx context.Context, msg any, expire time.Duration) error {
	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	return x.rds.SetEx(ctx, x.key, util.BytesToString(data), expire).Err()
}

func (x *x{{.Name}}) Get(ctx context.Context, msg any) error {
	data, err := x.rds.Get(ctx, x.key).Result()
	if err != nil {
		return err
	}
	err = json.Unmarshal(util.StringToBytes(data), msg)
	if err != nil {
		return err
	}
	return nil
} {{UsePackage "util" "StringToBytes"}}{{UsePackage "json" "Marshal/Unmarshal"}}