{{ $Name := .Name}} {{$msg := .TypeZSet.Message }} {{$score := .TypeZSet.Score }}
func (x *x{{$Name}}) ZCard(ctx context.Context) (int64, error) {
	cmd := redis.NewIntCmd(ctx, "zcard", x.key)
	x.rds.Process(ctx, cmd)
	return cmd.Result()
}

func (x *x{{$Name}}) ZAdd(ctx context.Context, mem {{$msg.Type}}, score {{$score.Type}}) error {
	data, err := {{call $msg.Marshal "mem"}}
	if err != nil {
		return err
	}
	cmd := redis.NewIntCmd(ctx, "zadd", x.key, util.BytesToString(data), rdconv.{{$score.RedisFunc}}ToString(score))
	return x.rds.Process(ctx, cmd)
}

func (x *x{{$Name}}) ZAddNX(ctx context.Context, mem {{$msg.Type}}, score {{$score.Type}}) error {
	data, err := {{call $msg.Marshal "mem"}}
	if err != nil {
		return err
	}
	cmd := redis.NewIntCmd(ctx, "zadd", x.key, "nx", util.BytesToString(data), rdconv.{{$score.RedisFunc}}ToString(score))
	return x.rds.Process(ctx, cmd)
}

func (x *x{{$Name}}) ZAddXX(ctx context.Context, mem {{$msg.Type}}, score {{$score.Type}}) error {
	data, err := {{call $msg.Marshal "mem"}}
	if err != nil {
		return err
	}
	cmd := redis.NewIntCmd(ctx, "zadd", x.key, "xx", util.BytesToString(data), rdconv.{{$score.RedisFunc}}ToString(score))
	return x.rds.Process(ctx, cmd)
}

func (x *x{{$Name}}) ZAddLT(ctx context.Context, mem {{$msg.Type}}, score {{$score.Type}}) error {
	data, err := {{call $msg.Marshal "mem"}}
	if err != nil {
		return err
	}
	cmd := redis.NewIntCmd(ctx, "zadd", x.key, "lt", util.BytesToString(data), rdconv.{{$score.RedisFunc}}ToString(score))
	return x.rds.Process(ctx, cmd)
}

func (x *x{{$Name}}) ZAddGT(ctx context.Context, mem {{$msg.Type}}, score {{$score.Type}}) error {
	data, err := {{call $msg.Marshal "mem"}}
	if err != nil {
		return err
	}
	cmd := redis.NewIntCmd(ctx, "zadd", x.key, "gt", util.BytesToString(data), rdconv.{{$score.RedisFunc}}ToString(score))
	return x.rds.Process(ctx, cmd)
}

func (x *x{{$Name}}) ZAdds(ctx context.Context, vals map[{{$msg.Type}}]{{$score.Type}}) error {
	args := make([]interface{}, 2, 2+len(vals)*2)
	args[0] = "zadd"
	args[1] = x.key
	for k, v := range vals {
		data, err := {{call $msg.Marshal "k"}}
		if err != nil {
			return err
		}
		args = append(args, util.BytesToString(data))
		args = append(args, rdconv.{{$score.RedisFunc}}ToString(v))
	}
	cmd := redis.NewIntCmd(ctx, args...)
	return x.rds.Process(ctx, cmd)
}

func (x *x{{$Name}}) ZRem(ctx context.Context, mem {{$msg.Type}}) error {
	data, err := {{call $msg.Marshal "mem"}}
	if err != nil {
		return err
	}
	cmd := redis.NewIntCmd(ctx, "zrem", x.key, util.BytesToString(data))
	return x.rds.Process(ctx, cmd)
}

func (x *x{{$Name}}) ZIncrBy(ctx context.Context, increment {{$score.Type}}, mem {{$msg.Type}}) (_ {{$score.Type}}, err error) {
	data, err := {{call $msg.Marshal "mem"}}
	if err != nil {
		return
	}
	cmd := redis.{{if $score.IsFloat}}NewFloatCmd{{else}}NewIntCmd{{end}}(ctx, "zincrby", x.key, rdconv.{{$score.RedisFunc}}ToString(increment), util.BytesToString(data))
	err = x.rds.Process(ctx, cmd)
	if err != nil {
		return
	}
	return {{$score.Type}}(cmd.Val()), nil
}

func (x *x{{$Name}}) ZScore(ctx context.Context, mem {{$msg.Type}}) (_ {{$score.Type}}, err error) {
	data, err := {{call $msg.Marshal "mem"}}
	if err != nil {
		return
	}
	cmd := redis.{{if $score.IsFloat}}NewFloatCmd{{else}}NewIntCmd{{end}}(ctx, "zscore", x.key, util.BytesToString(data))
	err = x.rds.Process(ctx, cmd)
	if err != nil {
		return
	}
	return {{$score.Type}}(cmd.Val()), nil
}

func (x *x{{$Name}}) ZRank(ctx context.Context, mem {{$msg.Type}}) (_ int64, err error) {
	data, err := {{call $msg.Marshal "mem"}}
	if err != nil {
		return
	}
	cmd := redis.NewIntCmd(ctx, "zrank", x.key, util.BytesToString(data))
	err = x.rds.Process(ctx, cmd)
	if err != nil {
		return
	}
	return cmd.Val(), nil
}

func (x *x{{$Name}}) ZRankWithScore(ctx context.Context, mem {{$msg.Type}}) (rank int64, score {{$score.Type}}, err error) {
	data, err := {{call $msg.Marshal "mem"}}
	if err != nil {
		return
	}
	cmd := redis.NewRankWithScoreCmd(ctx, "zrank", x.key, util.BytesToString(data), "withscore")
	err = x.rds.Process(ctx, cmd)
	if err != nil {
		return
	}
	rank = cmd.Val().Rank
	score = {{$score.Type}}(cmd.Val().Score)
	return
}

func (x *x{{$Name}}) ZRevRank(ctx context.Context, mem {{$msg.Type}}) (_ int64, err error) {
	data, err := {{call $msg.Marshal "mem"}}
	if err != nil {
		return
	}
	cmd := redis.NewIntCmd(ctx, "zrevrank", x.key, util.BytesToString(data))
	err = x.rds.Process(ctx, cmd)
	if err != nil {
		return
	}
	return cmd.Val(), nil
}

func (x *x{{$Name}}) ZRevRankWithScore(ctx context.Context, mem {{$msg.Type}}) (rank int64, score {{$score.Type}}, err error) {
	data, err := {{call $msg.Marshal "mem"}}
	if err != nil {
		return
	}
	cmd := redis.NewRankWithScoreCmd(ctx, "zrevrank", x.key, util.BytesToString(data), "withscore")
	err = x.rds.Process(ctx, cmd)
	if err != nil {
		return
	}
	rank = cmd.Val().Rank
	score = {{$score.Type}}(cmd.Val().Score)
	return
}

func (x *x{{$Name}}) parseMemberSliceCmd(cmd *redis.StringSliceCmd) (vals []{{$msg.Type}}, err error) {
	for _, v := range cmd.Val() {
		val := {{$msg.New}}
		err = {{call $msg.Unmarshal "val" "util.StringToBytes(v)"}}
		if err != nil {
			return nil, err
		}
		vals = append(vals, val)
	}
	return
}

func (x *x{{$Name}}) ZRange(ctx context.Context, start, stop int64) (vals []{{$msg.Type}}, err error) {
	cmd := redis.NewStringSliceCmd(ctx, "zrange", x.key, rdconv.Int64ToString(start), rdconv.Int64ToString(stop))
	err = x.rds.Process(ctx, cmd)
	if err != nil {
		return
	}

	return x.parseMemberSliceCmd(cmd)
}

func (x *x{{$Name}}) ZRangeByScore(ctx context.Context, start, stop {{$score.Type}}) (vals []{{$msg.Type}}, err error) {
	cmd := redis.NewStringSliceCmd(ctx, "zrange", x.key, rdconv.{{$score.RedisFunc}}ToString(start), rdconv.{{$score.RedisFunc}}ToString(stop), "byscore")
	err = x.rds.Process(ctx, cmd)
	if err != nil {
		return
	}

	return x.parseMemberSliceCmd(cmd)
}

func (x *x{{$Name}}) ZRevRange(ctx context.Context, start, stop int64) (vals []{{$msg.Type}}, err error) {
	cmd := redis.NewStringSliceCmd(ctx, "zrange", x.key, rdconv.Int64ToString(stop), rdconv.Int64ToString(start), "rev")
	err = x.rds.Process(ctx, cmd)
	if err != nil {
		return
	}

	return x.parseMemberSliceCmd(cmd)
}

func (x *x{{$Name}}) ZRevRangeByScore(ctx context.Context, start, stop {{$score.Type}}) (vals []{{$msg.Type}}, err error) {
	cmd := redis.NewStringSliceCmd(ctx, "zrange", x.key, rdconv.{{$score.RedisFunc}}ToString(stop), rdconv.{{$score.RedisFunc}}ToString(start), "byscore", "rev")
	err = x.rds.Process(ctx, cmd)
	if err != nil {
		return
	}

	return x.parseMemberSliceCmd(cmd)
}

func (x *x{{$Name}}) rangeMemberSliceCmd(cmd *redis.StringSliceCmd, f func({{$msg.Type}}) bool) (err error) {
	for _, v := range cmd.Val() {
		val := {{$msg.New}}
		err = {{call $msg.Unmarshal "val" "util.StringToBytes(v)"}}
		if err != nil {
			return err
		}
		if !f(val) {
			return nil
		}
	}
	return
}

func (x *x{{$Name}}) ZRangeF(ctx context.Context, start, stop int64, f func({{$msg.Type}}) bool) (err error) {
	cmd := redis.NewStringSliceCmd(ctx, "zrange", x.key, rdconv.Int64ToString(start), rdconv.Int64ToString(stop))
	err = x.rds.Process(ctx, cmd)
	if err != nil {
		return
	}

	return x.rangeMemberSliceCmd(cmd, f)
}

func (x *x{{$Name}}) ZRangeByScoreF(ctx context.Context, start, stop {{$score.Type}}, f func({{$msg.Type}}) bool) (err error) {
	cmd := redis.NewStringSliceCmd(ctx, "zrange", x.key, rdconv.{{$score.RedisFunc}}ToString(start), rdconv.{{$score.RedisFunc}}ToString(stop), "byscore")
	err = x.rds.Process(ctx, cmd)
	if err != nil {
		return
	}

	return x.rangeMemberSliceCmd(cmd, f)
}

func (x *x{{$Name}}) ZRevRangeF(ctx context.Context, start, stop int64, f func({{$msg.Type}}) bool) (err error) {
	cmd := redis.NewStringSliceCmd(ctx, "zrange", x.key, rdconv.Int64ToString(stop), rdconv.Int64ToString(start), "rev")
	err = x.rds.Process(ctx, cmd)
	if err != nil {
		return
	}

	return x.rangeMemberSliceCmd(cmd, f)
}

func (x *x{{$Name}}) ZRevRangeByScoreF(ctx context.Context, start, stop {{$score.Type}}, f func({{$msg.Type}}) bool) (err error) {
	cmd := redis.NewStringSliceCmd(ctx, "zrange", x.key, rdconv.{{$score.RedisFunc}}ToString(stop), rdconv.{{$score.RedisFunc}}ToString(start), "byscore", "rev")
	err = x.rds.Process(ctx, cmd)
	if err != nil {
		return
	}

	return x.rangeMemberSliceCmd(cmd, f)
}

func (x *x{{$Name}}) parseZSliceCmd(cmd *redis.ZSliceCmd) (vals map[{{$msg.Type}}]{{$score.Type}}, err error) {
	vals = make(map[{{$msg.Type}}]{{$score.Type}})
	for _, v := range cmd.Val() {
		str, err := rdconv.AnyToString(v.Member)
		if err != nil {
			return nil, err
		}
		val := {{$msg.New}}
		err = {{call $msg.Unmarshal "val" "util.StringToBytes(str)"}}
		if err != nil {
			return nil, err
		}
		vals[val] = {{$score.Type}}(v.Score)
	}
	return
}

func (x *x{{$Name}}) ZRangeWithScores(ctx context.Context, start, stop int64) (vals map[{{$msg.Type}}]{{$score.Type}}, err error) {
	cmd := redis.NewZSliceCmd(ctx, "zrange", x.key, rdconv.Int64ToString(start), rdconv.Int64ToString(stop), "withscores")
	err = x.rds.Process(ctx, cmd)
	if err != nil {
		return
	}

	return x.parseZSliceCmd(cmd)
}

func (x *x{{$Name}}) ZRangeByScoreWithScores(ctx context.Context, start, stop {{$score.Type}}) (vals map[{{$msg.Type}}]{{$score.Type}}, err error) {
	cmd := redis.NewZSliceCmd(ctx, "zrange", x.key, rdconv.{{$score.RedisFunc}}ToString(start), rdconv.{{$score.RedisFunc}}ToString(stop), "byscore", "withscores")
	err = x.rds.Process(ctx, cmd)
	if err != nil {
		return
	}

	return x.parseZSliceCmd(cmd)
}
func (x *x{{$Name}}) ZRevRangeWithScores(ctx context.Context, start, stop int64) (vals map[{{$msg.Type}}]{{$score.Type}}, err error) {
	cmd := redis.NewZSliceCmd(ctx, "zrange", x.key, rdconv.Int64ToString(stop), rdconv.Int64ToString(start), "rev", "withscores")
	err = x.rds.Process(ctx, cmd)
	if err != nil {
		return
	}

	return x.parseZSliceCmd(cmd)
}

func (x *x{{$Name}}) ZRevRangeByScoreWithScores(ctx context.Context, start, stop {{$score.Type}}) (vals map[{{$msg.Type}}]{{$score.Type}}, err error) {
	cmd := redis.NewZSliceCmd(ctx, "zrange", x.key, rdconv.{{$score.RedisFunc}}ToString(stop), rdconv.{{$score.RedisFunc}}ToString(start), "byscore", "rev", "withscores")
	err = x.rds.Process(ctx, cmd)
	if err != nil {
		return
	}

	return x.parseZSliceCmd(cmd)
}

func (x *x{{$Name}}) rangeZSliceCmd(cmd *redis.ZSliceCmd, f func({{$msg.Type}}, {{$score.Type}}) bool) (err error) {
	for _, v := range cmd.Val() {
		str, err := rdconv.AnyToString(v.Member)
		if err != nil {
			return err
		}
		val := {{$msg.New}}
		err = {{call $msg.Unmarshal "val" "util.StringToBytes(str)"}}
		if err != nil {
			return err
		}
		if !f(val, {{$score.Type}}(v.Score)) {
			return nil
		}
	}
	return
}

func (x *x{{$Name}}) ZRangeWithScoresF(ctx context.Context, start, stop int64, f func({{$msg.Type}}, {{$score.Type}}) bool) (err error) {
	cmd := redis.NewZSliceCmd(ctx, "zrange", x.key, rdconv.Int64ToString(start), rdconv.Int64ToString(stop), "withscores")
	err = x.rds.Process(ctx, cmd)
	if err != nil {
		return
	}

	return x.rangeZSliceCmd(cmd, f)
}

func (x *x{{$Name}}) ZRangeByScoreWithScoresF(ctx context.Context, start, stop {{$score.Type}}, f func({{$msg.Type}}, {{$score.Type}}) bool) (err error) {
	cmd := redis.NewZSliceCmd(ctx, "zrange", x.key, rdconv.{{$score.RedisFunc}}ToString(start), rdconv.{{$score.RedisFunc}}ToString(stop), "byscore", "withscores")
	err = x.rds.Process(ctx, cmd)
	if err != nil {
		return
	}

	return x.rangeZSliceCmd(cmd, f)
}
func (x *x{{$Name}}) ZRevRangeWithScoresF(ctx context.Context, start, stop int64, f func({{$msg.Type}}, {{$score.Type}}) bool) (err error) {
	cmd := redis.NewZSliceCmd(ctx, "zrange", x.key, rdconv.Int64ToString(stop), rdconv.Int64ToString(start), "rev", "withscores")
	err = x.rds.Process(ctx, cmd)
	if err != nil {
		return
	}

	return x.rangeZSliceCmd(cmd, f)
}

func (x *x{{$Name}}) ZRevRangeByScoreWithScoresF(ctx context.Context, start, stop {{$score.Type}}, f func({{$msg.Type}}, {{$score.Type}}) bool) (err error) {
	cmd := redis.NewZSliceCmd(ctx, "zrange", x.key, rdconv.{{$score.RedisFunc}}ToString(stop), rdconv.{{$score.RedisFunc}}ToString(start), "byscore", "rev", "withscores")
	err = x.rds.Process(ctx, cmd)
	if err != nil {
		return
	}

	return x.rangeZSliceCmd(cmd, f)
}

func (x *x{{$Name}}) ZPopMin(ctx context.Context, count int64) (_ map[{{$msg.Type}}]{{$score.Type}}, err error) {
	cmd := x.rds.ZPopMin(ctx, x.key, int64(count))
	if cmd.Err() != nil {
		return nil, cmd.Err()
	}
	return x.parseZSliceCmd(cmd)
}

func (x *x{{$Name}}) ZPopMinF(ctx context.Context, count int64, f func({{$msg.Type}}, {{$score.Type}}) bool) (err error) {
	cmd := x.rds.ZPopMin(ctx, x.key, int64(count))
	if cmd.Err() != nil {
		return cmd.Err()
	}
	return x.rangeZSliceCmd(cmd, f)
}

func (x *x{{$Name}}) ZPopMax(ctx context.Context, count int64) (_ map[{{$msg.Type}}]{{$score.Type}}, err error) {
	cmd := x.rds.ZPopMax(ctx, x.key, int64(count))
	if cmd.Err() != nil {
		return nil, cmd.Err()
	}
	return x.parseZSliceCmd(cmd)
}

func (x *x{{$Name}}) ZPopMaxF(ctx context.Context, count int64, f func({{$msg.Type}}, {{$score.Type}}) bool) (err error) {
	cmd := x.rds.ZPopMax(ctx, x.key, int64(count))
	if cmd.Err() != nil {
		return cmd.Err()
	}
	return x.rangeZSliceCmd(cmd, f)
}

func (x *x{{$Name}}) ZPopGTScore(ctx context.Context, limitScore {{$score.Type}}, count int64) (vals []{{$msg.Type}}, err error) {
	cmd := redis.NewStringSliceCmd(ctx, "evalsha", {{.SvcPkg}}.ZPopMaxValue.Hash, "1", x.key, rdconv.{{$score.RedisFunc}}ToString(limitScore), rdconv.Int64ToString(count))
	err = x.rds.Process(ctx, cmd)
	if err != nil {
		if !redis.HasErrorPrefix(err, "NOSCRIPT") {
			return
		}
		cmd = redis.NewStringSliceCmd(ctx, "eval", {{.SvcPkg}}.ZPopMaxValue.Script, "1", x.key, rdconv.{{$score.RedisFunc}}ToString(limitScore), rdconv.Int64ToString(count))
		err = x.rds.Process(ctx, cmd)
		if err != nil {
			return
		}
	}
	return x.parseMemberSliceCmd(cmd)
}

func (x *x{{$Name}}) ZPopGTScoreF(ctx context.Context, limitScore {{$score.Type}}, count int64, f func({{$msg.Type}}) bool) (err error) {
	cmd := redis.NewStringSliceCmd(ctx, "evalsha", {{.SvcPkg}}.ZPopMaxValue.Hash, "1", x.key, rdconv.{{$score.RedisFunc}}ToString(limitScore), rdconv.Int64ToString(count))
	err = x.rds.Process(ctx, cmd)
	if err != nil {
		if !redis.HasErrorPrefix(err, "NOSCRIPT") {
			return
		}
		cmd = redis.NewStringSliceCmd(ctx, "eval", {{.SvcPkg}}.ZPopMaxValue.Script, "1", x.key, rdconv.{{$score.RedisFunc}}ToString(limitScore), rdconv.Int64ToString(count))
		err = x.rds.Process(ctx, cmd)
		if err != nil {
			return
		}
	}
	//return cmd.Val(), nil
	return x.rangeMemberSliceCmd(cmd, f)
}

func (x *x{{$Name}}) ZPopGTScoreWithScores(ctx context.Context, limitScore {{$score.Type}}, count int64) (vals map[{{$msg.Type}}]{{$score.Type}}, err error) {
	cmd := redis.NewZSliceCmd(ctx, "evalsha", {{.SvcPkg}}.ZPopMaxValueWithScore.Hash, "1", x.key, rdconv.{{$score.RedisFunc}}ToString(limitScore), rdconv.Int64ToString(count))
	err = x.rds.Process(ctx, cmd)
	if err != nil {
		if !redis.HasErrorPrefix(err, "NOSCRIPT") {
			return
		}
		cmd = redis.NewZSliceCmd(ctx, "eval", {{.SvcPkg}}.ZPopMaxValueWithScore.Script, "1", x.key, rdconv.{{$score.RedisFunc}}ToString(limitScore), rdconv.Int64ToString(count)) {{- UsePackage "rdconv" "ToString/FronString"}}
		err = x.rds.Process(ctx, cmd)
		if err != nil {
			return
		}
	}
	//return cmd.Val(), nil
	return x.parseZSliceCmd(cmd)
}

func (x *x{{$Name}}) ZPopGTScoreWithScoresF(ctx context.Context, limitScore {{$score.Type}}, count int64, f func({{$msg.Type}}, {{$score.Type}}) bool) (err error) {
	cmd := redis.NewZSliceCmd(ctx, "evalsha", {{.SvcPkg}}.ZPopMaxValueWithScore.Hash, "1", x.key, rdconv.{{$score.RedisFunc}}ToString(limitScore), rdconv.Int64ToString(count))
	err = x.rds.Process(ctx, cmd)
	if err != nil {
		if !redis.HasErrorPrefix(err, "NOSCRIPT") {
			return
		}
		cmd = redis.NewZSliceCmd(ctx, "eval", {{.SvcPkg}}.ZPopMaxValueWithScore.Script, "1", x.key, rdconv.{{$score.RedisFunc}}ToString(limitScore), rdconv.Int64ToString(count))
		err = x.rds.Process(ctx, cmd)
		if err != nil {
			return
		}
	}
	//return cmd.Val(), nil
	return x.rangeZSliceCmd(cmd, f)
}