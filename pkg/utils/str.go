package utils

import (
	"regexp"
	"strings"
)

var (
	// 使用正则表达式匹配大写字母, 用于将大驼峰命名转换为小写加下划线的形式
	reToSnake = regexp.MustCompile("([a-z])([A-Z])")
)

// PascalToSnake 将大驼峰命名转换为小写加下划线的形式
func PascalToSnake(pascalStr string) string {
	// 插入下划线并转换为小写
	snakeStr := reToSnake.ReplaceAllString(pascalStr, "${1}_${2}")
	return strings.ToLower(snakeStr)
}
