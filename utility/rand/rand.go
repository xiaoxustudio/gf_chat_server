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
