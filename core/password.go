package core

import (
	"time"
	"math/rand"
)

const PasswordLength = 256

// Password 结构定义
type Password [PasswordLength]byte

func init() {
	rand.Seed(time.Now().Unix())
}

// RandPassword 随机密码
func RandPassword() *Password {
	// 随机生成一个由  0~255 组成的 byte 数组
	intArr := rand.Perm(PasswordLength)
	password := &Password{}
	for i, v := range intArr {
		password[i] = byte(v)
		if i == v {
			return RandPassword()
		}
	}
	return password
}