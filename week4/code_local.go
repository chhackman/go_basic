package cache

import (
	"context"
	"errors"
	"fmt"
	"sync"
)

// 定义错误变量，用于表示不同的错误情况
var (
	LocallErrCodeSendTooMany       = errors.New("发送验证码太频繁")
	LocalErrCodeVerifyTooManyTimes = errors.New("验证次数太多")
	LocalErrUnknowForCode          = errors.New("我也不知道发生什么，反正是和code有关")
)

// CodeCache 接口定义了缓存操作的方法
type CodeCache interface {
	Set(ctx context.Context, biz, phone, code string) error
	Verify(ctx context.Context, biz, phone, inputCode string) (bool, error)
}

// CodeMemoryCache 是 CodeCache 接口的实现，基于本地内存缓存
type CodeMemoryCache struct {
	cache map[string]string // 使用 map 存储验证码
	mu    sync.Mutex        // 用于并发安全的互斥锁
}

// NewCodeMemoryCache 创建一个新的 CodeMemoryCache 实例
func NewCodeMemoryCache() *CodeMemoryCache {
	return &CodeMemoryCache{
		cache: make(map[string]string), // 初始化存储验证码的 map
	}
}

// Set 方法用于将验证码存入缓存
func (c *CodeMemoryCache) Set(ctx context.Context, biz, phone, code string) error {
	c.mu.Lock()         // 锁定互斥锁，保证并发安全
	defer c.mu.Unlock() // 在函数结束后解锁

	key := c.key(biz, phone) // 生成用于存储的键

	// 如果已经存在相同的键（即相同的业务和手机号），返回发送验证码太频繁的错误
	if _, exists := c.cache[key]; exists {
		return ErrCodeSendTooMany
	}

	// 否则，将验证码存入缓存
	c.cache[key] = code
	return nil
}

// Verify 方法用于验证输入的验证码是否正确
func (c *CodeMemoryCache) Verify(ctx context.Context, biz, phone, inputCode string) (bool, error) {
	c.mu.Lock()         // 锁定互斥锁，保证并发安全
	defer c.mu.Unlock() // 在函数结束后解锁

	key := c.key(biz, phone) // 生成用于存储的键

	// 获取存储的验证码
	storedCode, exists := c.cache[key]

	// 如果不存在相同的键（即没有发送验证码记录），返回验证次数太多的错误
	if !exists {
		return false, ErrCodeVerifyTooManyTimes
	}

	// 如果输入的验证码与存储的验证码匹配，删除缓存中的验证码并返回验证成功
	if storedCode == inputCode {
		delete(c.cache, key) // 从缓存中删除验证码
		return true, nil
	}

	// 否则，返回验证失败
	return false, nil
}

// key 方法用于生成存储的键
func (c *CodeMemoryCache) key(biz, phone string) string {
	return fmt.Sprintf("phone_code:%s:%s", biz, phone)
}
