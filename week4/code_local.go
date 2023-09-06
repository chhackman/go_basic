package cache

import (
	"context"
	"errors"
	"fmt"
	"sync"
)

var (
	LocalErrCodeSendTooMany        = errors.New("发送验证码太频繁")
	LocalErrCodeVerifyTooManyTimes = errors.New("验证次数太多")
	LocalErrUnknowForCode          = errors.New("我也不知道发生什么，反正是和code有关")
)

type CodeCache interface {
	Set(ctx context.Context, biz, phone, code string) error
	Verify(ctx context.Context, biz, phone, inputCode string) (bool, error)
}

type CodeMemoryCache struct {
	cache map[string]string
	mu    sync.Mutex
}

func NewCodeMemoryCache() *CodeMemoryCache {
	return &CodeMemoryCache{
		cache: make(map[string]string),
	}
}

func (c *CodeMemoryCache) Set(ctx context.Context, biz, phone, code string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	key := c.key(biz, phone)
	if _, exists := c.cache[key]; exists {
		return ErrCodeSendTooMany
	}

	c.cache[key] = code
	return nil
}

func (c *CodeMemoryCache) Verify(ctx context.Context, biz, phone, inputCode string) (bool, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	key := c.key(biz, phone)
	storedCode, exists := c.cache[key]
	if !exists {
		return false, ErrCodeVerifyTooManyTimes
	}

	if storedCode == inputCode {
		delete(c.cache, key)
		return true, nil
	}

	return false, nil
}

func (c *CodeMemoryCache) key(biz, phone string) string {
	return fmt.Sprintf("phone_code:%s:%s", biz, phone)
}
