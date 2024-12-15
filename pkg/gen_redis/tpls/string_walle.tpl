
func (x *x{{.Name}}) Set(ctx context.Context, pb {{.WPbPkg}}.Message, expire time.Duration) error {
	data, err := pb.MarshalObject()
	if err != nil {
		return err
	}
	return x.rds.Set(ctx, x.key, util.BytesToString(data), expire).Err()
}

func (x *x{{.Name}}) SetNX(ctx context.Context, pb {{.WPbPkg}}.Message, expire time.Duration) error {
	data, err := pb.MarshalObject()
	if err != nil {
		return err
	}
	return x.rds.SetNX(ctx, x.key, util.BytesToString(data), expire).Err()
}

func (x *x{{.Name}}) SetEx(ctx context.Context, pb {{.WPbPkg}}.Message, expire time.Duration) error {
	data, err := pb.MarshalObject()
	if err != nil {
		return err
	}
	return x.rds.SetEx(ctx, x.key, util.BytesToString(data), expire).Err()
}

func (x *x{{.Name}}) Get(ctx context.Context, pb {{.WPbPkg}}.Message) error {
	data, err := x.rds.Get(ctx, x.key).Result()
	if err != nil {
		return err
	}
	err = pb.UnmarshalObject(util.StringToBytes(data))
	if err != nil {
		return err
	}
	return nil
} {{UsePackage "util" "StringToBytes"}}