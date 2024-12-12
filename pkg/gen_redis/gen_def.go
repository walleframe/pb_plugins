// Copyright 2024 aggronmagi <czy463@163.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package gen_redis

import (
	"strings"

	"github.com/walleframe/pb_plugins/pkg/gen_redis/keyarg"
	"github.com/walleframe/pb_plugins/pkg/tpl"
)

// RedisObject 生成redis消息的对象
type RedisObject struct {
	tpl.GoObject

	// 生成包名
	Package string
	// 操作名,消息名
	Name string
	// 消息体注释
	Doc *tpl.Commets
	// key 参数
	Args []keyarg.KeyArg
	// service package
	SvcPkg string
	// walle package
	WPbPkg string
	//
	KeySize int

	// redis 脚本
	Scripts []*RedisScript

	// key 相关接口
	TypeKeys bool
	// string 生成
	TypeString *RedisTypeString
	// hash 生成
	TypeHash *RedisTypeHash
	// set 生成
	TypeSet *RedisTypeSet
	// zset 生成
	TypeZSet *RedisTypeZSet
	// redis 分布式锁
	Lock bool
}

type RedisTypeString struct {
	Type     string
	Signed   bool
	Number   bool
	String   bool
	Protobuf bool // github.com/gogo/protobuf/proto
	WProto   bool // github.com/walleframe/walle/process/message
	Custom   bool
	Float    bool
}

type RedisTypeHash struct {
	HashObject  *RedisHashObject
	HashDynamic *RedisHashDynamic
}

type RedisGenType struct {
	// 字段名
	Name string
	// 字段类型
	Type string
	// Number
	Number bool
	//
	RedisFunc string
}

func (x *RedisGenType) IsFloat() bool {
	return strings.HasPrefix(x.Type, "float")
}

func (x *RedisGenType) IsInt() bool {
	return strings.Contains(x.Type, "int")
}

type RedisHashObject struct {
	Fields  []*RedisGenType
	Name    string
	Type    string
	HGetAll bool
}

type RedisHashDynamic struct {
	Field     *RedisGenType
	Value     *RedisGenType
	GenMap    bool
	FieldArgs []keyarg.KeyArg
	ValueArgs []keyarg.KeyArg
}

type RedisGenMsg struct {
	Type      string
	Marshal   func(objName string) string
	Unmarshal func(objName, paramName string) string
	New       string
}

type RedisTypeSet struct {
	BaseType *RedisGenType
	Message  *RedisGenMsg
}

type RedisTypeZSet struct {
	// score
	Score *RedisGenType
	// mem
	Member *RedisGenType
	// 拼接string做field
	Args []keyarg.KeyArg
	// 使用消息作为field
	Message *RedisGenMsg
}

// RedisScript redis脚本
type RedisScript struct {
	// 脚本操作名称
	Name string
	// 脚本数据
	Script string
	// 脚本参数
	Args []keyarg.KeyArg
	// 脚本输出
	Output []keyarg.KeyArg
	//
	TemplateName string
	//
	CommandName string
}
