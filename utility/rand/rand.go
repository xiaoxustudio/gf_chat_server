package rand

import (
	"crypto/rand"
	"fmt"
	"sync"
)

var (
	ids = make(map[string]struct{})
	mu  sync.Mutex
)

// 生成唯一ID
func GenUniqueID() string {
	var b [16]byte
	if _, err := rand.Read(b[:]); err != nil {
		panic(err) // Handle error properly in production code
	}

	id := fmt.Sprintf("%x", b)
	mu.Lock()
	defer mu.Unlock()

	if _, exists := ids[id]; exists {
		return GenUniqueID()
	}

	ids[id] = struct{}{}
	return id
}

// GetID 生成指定位数的ID
func GetID(numDigits int) string {
	if numDigits == 0 {
		return ""
	}

	// 计算所需字节数
	numBytes := (numDigits + 1) / 2

	// 生成随机字节
	var b = make([]byte, numBytes)
	if _, err := rand.Read(b); err != nil {
		panic(err) // Handle error properly in production code
	}

	// 转换为十六进制字符串
	id := fmt.Sprintf("%x", b)[:numDigits]
	mu.Lock()
	defer mu.Unlock()

	// 检查ID是否唯一
	if _, exists := ids[id]; exists {
		return GetID(numDigits)
	}

	ids[id] = struct{}{}
	return id
}
