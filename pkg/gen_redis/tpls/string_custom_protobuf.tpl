
func (x *x{{.Name}}) Set(ctx context.Context, pb *{{.TypeString.Type}}, expire time.Duration) error {
	data, err := proto.Marshal(pb)
	if err != nil {
		return err
	}
	return x.rds.Set(ctx, x.key, util.BytesToString(data), expire).Err()
}

func (x *x{{.Name}}) SetNX(ctx context.Context, pb *{{.TypeString.Type}}, expire time.Duration) error {
	data, err := proto.Marshal(pb)
	if err != nil {
		return err
	}
	return x.rds.SetNX(ctx, x.key, util.BytesToString(data), expire).Err()
}

func (x *x{{.Name}}) SetEx(ctx context.Context, pb *{{.TypeString.Type}}, expire time.Duration) error {
	data, err := proto.Marshal(pb)
	if err != nil {
		return err
	}
	return x.rds.SetEx(ctx, x.key, util.BytesToString(data), expire).Err()
}

func (x *x{{.Name}}) Get(ctx context.Context, pb *{{.TypeString.Type}}) error {
	data, err := x.rds.Get(ctx, x.key).Result()
	if err != nil {
		return err
	}
	err = proto.Unmarshal(util.StringToBytes(data), pb)
	if err != nil {
		return err
	}
	return nil
} {{Import "github.com/walleframe/walle/util" "StringToBytes"}}