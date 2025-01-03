 {{ $Name := .Name}} {{$msg := .TypeSet.Message }}
func (x *x{{$Name}}) SAdd(ctx context.Context, val {{$msg.Type}}) (bool, error) {
	data, err := {{call $msg.Marshal "val"}}
	if err != nil {
		return false, err
	}
	n, err := x.rds.SAdd(ctx, x.key, util.BytesToString(data)).Result()
	if err != nil {
		return false, err
	}
	return n == 1, nil
}

func (x *x{{$Name}}) SCard(ctx context.Context) (int64, error) {
	return x.rds.SCard(ctx, x.key).Result()
}

func (x *x{{$Name}}) SRem(ctx context.Context, val {{$msg.Type}}) (bool, error) {
	data, err := {{call $msg.Marshal "val"}}
	if err != nil {
		return false, err
	}
	n, err := x.rds.SRem(ctx, x.key, util.BytesToString(data)).Result()
	if err != nil {
		return false, err
	}
	return n == 1, nil
}

func (x *x{{$Name}}) SIsMember(ctx context.Context, val {{$msg.Type}}) (bool, error) {
	data, err := {{call $msg.Marshal "val"}}
	if err != nil {
		return false, err
	}
	return x.rds.SIsMember(ctx, x.key, util.BytesToString(data)).Result()
}
{{if $msg.New }}
func (x *x{{$Name}}) SPop(ctx context.Context) (_ {{$msg.Type}}, err error) {
	v, err := x.rds.SPop(ctx, x.key).Result()
	if err != nil {
		return
	}
	val := {{$msg.New}}
	err = {{call $msg.Unmarshal "val" "util.StringToBytes(v)"}}
	if err != nil {
		return
	}
	return val, nil
}
func (x *x{{$Name}}) SRandMember(ctx context.Context) (_ {{$msg.Type}}, err error) {
	v, err := x.rds.SRandMember(ctx, x.key).Result()
	if err != nil {
		return
	}
	val := {{$msg.New}}
	err = {{call $msg.Unmarshal "val" "util.StringToBytes(v)"}}
	if err != nil {
		return
	}
	return val, nil
}

func (x *x{{$Name}}) SRandMemberN(ctx context.Context, count int) (vals []{{$msg.Type}}, err error) {
	ret, err := x.rds.SRandMemberN(ctx, x.key, int64(count)).Result()
	if err != nil {
		return
	}
	for _, v := range ret {
		val := {{$msg.New}}
		err = {{call $msg.Unmarshal "val" "util.StringToBytes(v)"}}
		if err != nil {
			return nil, err
		}
		vals = append(vals, val)
	}
	return
}

func (x *x{{$Name}}) SRandMemberNRange(ctx context.Context, count int, filter func({{$msg.Type}})bool) (err error) {
	ret, err := x.rds.SRandMemberN(ctx, x.key, int64(count)).Result()
	if err != nil {
		return
	}
	for _, v := range ret {
		val := {{$msg.New}}
		err = {{call $msg.Unmarshal "val" "util.StringToBytes(v)"}}
		if err != nil {
			return err
		}
		if !filter(val) {
			return nil
		}
	}
	return
}

func (x *x{{$Name}}) SMembers(ctx context.Context) (vals []{{$msg.Type}}, err error) {
	ret, err := x.rds.SMembers(ctx, x.key).Result()
	if err != nil {
		return
	}
	for _, v := range ret {
		val := {{$msg.New}}
		err = {{call $msg.Unmarshal "val" "util.StringToBytes(v)"}}
		if err != nil {
			return nil,err
		}
		vals = append(vals, val)
	}
	return
}

func (x *x{{$Name}}) SScan(ctx context.Context, match string, count int) (vals []{{$msg.Type}}, err error) {
	cursor := uint64(0)
	var ret []string
	for {
		ret, cursor, err = x.rds.SScan(ctx, x.key, cursor, match, int64(count)).Result()
		if err != nil {
			return nil, err
		}
		for _, v := range ret {
			val := {{$msg.New}}
			err = {{call $msg.Unmarshal "val" "util.StringToBytes(v)"}}
			if err != nil {
				return nil,err
			}
			vals = append(vals, val)
		}
		if cursor == 0 {
			break
		}
	}
	return
}

func (x *x{{$Name}}) SScanRange(ctx context.Context, match string, count int, filter func({{$msg.Type}}) bool) (err error) {
	cursor := uint64(0)
	var ret []string
	for {
		ret, cursor, err = x.rds.SScan(ctx, x.key, cursor, match, int64(count)).Result()
		if err != nil {
			return err
		}
		for _, v := range ret {
			val := {{$msg.New}}
			err = {{call $msg.Unmarshal "val" "util.StringToBytes(v)"}}
			if err != nil {
				return err
			}
			if !filter(val) {
				return nil
			}
		}
		if cursor == 0 {
			break
		}
	}
	return
}

{{else}}

func (x *x{{$Name}}) SPop(ctx context.Context,val {{$msg.Type}})(err error) {
	v, err := x.rds.SPop(ctx, x.key).Result()
	if err != nil {
		return
	}
	err = {{call $msg.Unmarshal "val" "util.StringToBytes(v)"}}
	if err != nil {
		return
	}
	return nil
}
func (x *x{{$Name}}) SRandMember(ctx context.Context,val {{$msg.Type}})(err error) {
	v, err := x.rds.SRandMember(ctx, x.key).Result()
	if err != nil {
		return
	}
	err = {{call $msg.Unmarshal "val" "util.StringToBytes(v)"}}
	if err != nil {
		return
	}
	return nil
}


func (x *x{{$Name}}) SRandMemberN(ctx context.Context, count int, newFunc func(){{$msg.Type}}) (vals []{{$msg.Type}}, err error) {
	ret, err := x.rds.SRandMemberN(ctx, x.key, int64(count)).Result()
	if err != nil {
		return
	}
	for _, v := range ret {
		val := newFunc()
		err = {{call $msg.Unmarshal "val" "util.StringToBytes(v)"}}
		if err != nil {
			return nil, err
		}
		vals = append(vals, val)
	}
	return
}

func (x *x{{$Name}}) SRandMemberNRange(ctx context.Context, count int, newFunc func(){{$msg.Type}}, filter func({{$msg.Type}})bool) (err error) {
	ret, err := x.rds.SRandMemberN(ctx, x.key, int64(count)).Result()
	if err != nil {
		return
	}
	for _, v := range ret {
		val := newFunc()
		err = {{call $msg.Unmarshal "val" "util.StringToBytes(v)"}}
		if err != nil {
			return err
		}
		if !filter(val) {
			return nil
		}
	}
	return
}

func (x *x{{$Name}}) SMembers(ctx context.Context, count int, newFunc func(){{$msg.Type}}) (vals []{{$msg.Type}}, err error) {
	ret, err := x.rds.SMembers(ctx, x.key, int64(count)).Result()
	if err != nil {
		return
	}
	for _, v := range ret {
		val := newFunc()
		err = {{call $msg.Unmarshal "val" "util.StringToBytes(v)"}}
		if err != nil {
			return nil,err
		}
		vals = append(vals, val)
	}
	return
}

func (x *x{{$Name}}) SScan(ctx context.Context, match string, count int, newFunc func(){{$msg.Type}}) (vals []{{$msg.Type}}, err error) {
	cursor := uint64(0)
	var ret []string
	for {
		ret, cursor, err = x.rds.SScan(ctx, x.key, cursor, match, int64(count)).Result()
		if err != nil {
			return nil, err
		}
		for _, v := range ret {
			val := newFunc()
			err = {{call $msg.Unmarshal "val" "util.StringToBytes(v)"}}
			if err != nil {
				return nil,err
			}
			vals = append(vals, val)
		}
		if cursor == 0 {
			break
		}
	}
	return
}

func (x *x{{$Name}}) SScanRange(ctx context.Context, match string, count int, newFunc func(){{$msg.Type}}, filter func({{$msg.Type}}) bool) (err error) {
	cursor := uint64(0)
	var ret []string
	for {
		ret, cursor, err = x.rds.SScan(ctx, x.key, cursor, match, int64(count)).Result()
		if err != nil {
			return err
		}
		for _, v := range ret {
			val := newFunc()
			err = {{call $msg.Unmarshal "val" "util.StringToBytes(v)"}}
			if err != nil {
				return err
			}
			if !filter(val) {
				return nil
			}
		}
		if cursor == 0 {
			break
		}
	}
	return
}
{{end}}