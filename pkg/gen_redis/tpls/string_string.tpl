
// GetRange 获取字符串子串 [start,end]
func (x *x{{.Name}}) GetRange(ctx context.Context, start, end int64) (_ string, err error) {
	cmd := redis.NewStringCmd(ctx, "getrange", x.key, strconv.FormatInt(start, 10), strconv.FormatInt(end, 10))
	err = x.rds.Process(ctx, cmd)
	if err != nil {
		return
	}
	return cmd.Val(), nil
}

// SetRange 设置字符串子串
func (x *x{{.Name}}) SetRange(ctx context.Context, offset int64, value string) (_ int64, err error) {
	cmd := redis.NewIntCmd(ctx, "setrange", x.key, strconv.FormatInt(offset, 10), value) {{- Import "strconv" "strconv.Format"}}
	err = x.rds.Process(ctx, cmd)
	if err != nil {
		return
	}
	return cmd.Val(), nil
}

// Append 追加字符串
func (x *x{{.Name}}) Append(ctx context.Context, val string) (int64, error) {
	return x.rds.Append(ctx, x.key, val).Result()
}

// StrLen 获取字符串长度
func (x *x{{.Name}}) StrLen(ctx context.Context) (int64, error) {
	return x.rds.StrLen(ctx, x.key).Result()
}

// Get 获取字符串值
func (x *x{{.Name}}) Get(ctx context.Context) (string, error) {
	return x.rds.Get(ctx, x.key).Result()
}

// Set 设置字符串值
func (x *x{{.Name}}) Set(ctx context.Context, data string, expire time.Duration) error { {{- Import "time" "time.Duration"}}
	return x.rds.Set(ctx, x.key, data, expire).Err()
}

// SetNX 设置字符串值(仅当key不存在时)
func (x *x{{.Name}}) SetNX(ctx context.Context, data string, expire time.Duration) (bool, error) {
	return x.rds.SetNX(ctx, x.key, data, expire).Result()
}

// SetEx 设置字符串值并指定过期时间
func (x *x{{.Name}}) SetEx(ctx context.Context, data string, expire time.Duration) error {
	return x.rds.SetEx(ctx, x.key, data, expire).Err()
}