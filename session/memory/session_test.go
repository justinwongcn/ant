package memory

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestStoreNewStore(t *testing.T) {
	expiration := 30 * time.Minute
	store := NewStore(expiration)

	assert.NotNil(t, store)
	assert.NotNil(t, store.c)
	assert.Equal(t, expiration, store.expiration)
}

func TestStoreGenerate(t *testing.T) {
	store := NewStore(30 * time.Minute)
	ctx := context.Background()

	// 测试生成新会话
	sess, err := store.Generate(ctx, "test-id")
	assert.NoError(t, err)
	assert.NotNil(t, sess)
	assert.Equal(t, "test-id", sess.ID())

	// 验证会话已被存储
	storedSess, err := store.Get(ctx, "test-id")
	assert.NoError(t, err)
	assert.Equal(t, sess.ID(), storedSess.ID())
}

func TestStoreGet(t *testing.T) {
	store := NewStore(30 * time.Minute)
	ctx := context.Background()

	// 测试获取不存在的会话
	sess, err := store.Get(ctx, "non-existent")
	assert.Error(t, err)
	assert.Nil(t, sess)

	// 测试获取存在的会话
	store.Generate(ctx, "test-id")
	sess, err = store.Get(ctx, "test-id")
	assert.NoError(t, err)
	assert.NotNil(t, sess)
	assert.Equal(t, "test-id", sess.ID())
}

func TestStoreRefresh(t *testing.T) {
	store := NewStore(30 * time.Minute)
	ctx := context.Background()

	// 测试刷新不存在的会话
	err := store.Refresh(ctx, "non-existent")
	assert.Error(t, err)

	// 测试刷新存在的会话
	store.Generate(ctx, "test-id")
	err = store.Refresh(ctx, "test-id")
	assert.NoError(t, err)
}

func TestStoreRemove(t *testing.T) {
	store := NewStore(30 * time.Minute)
	ctx := context.Background()

	// 生成一个会话
	store.Generate(ctx, "test-id")

	// 测试删除会话
	err := store.Remove(ctx, "test-id")
	assert.NoError(t, err)

	// 验证会话已被删除
	sess, err := store.Get(ctx, "test-id")
	assert.Error(t, err)
	assert.Nil(t, sess)
}

func TestMemorySessionGet(t *testing.T) {
	sess := &memorySession{
		id:   "test-id",
		data: make(map[string]any),
	}
	ctx := context.Background()

	// 测试获取不存在的键
	val, err := sess.Get(ctx, "non-existent")
	assert.Error(t, err)
	assert.Equal(t, "", val)

	// 测试获取存在的键
	sess.data["test-key"] = "test-value"
	val, err = sess.Get(ctx, "test-key")
	assert.NoError(t, err)
	assert.Equal(t, "test-value", val)
}

func TestMemorySessionSet(t *testing.T) {
	sess := &memorySession{
		id:   "test-id",
		data: make(map[string]any),
	}
	ctx := context.Background()

	// 测试设置新键值对
	err := sess.Set(ctx, "test-key", "test-value")
	assert.NoError(t, err)
	assert.Equal(t, "test-value", sess.data["test-key"])

	// 测试更新已存在的键
	err = sess.Set(ctx, "test-key", "new-value")
	assert.NoError(t, err)
	assert.Equal(t, "new-value", sess.data["test-key"])
}

func TestMemorySessionID(t *testing.T) {
	sess := &memorySession{
		id:   "test-id",
		data: make(map[string]any),
	}

	assert.Equal(t, "test-id", sess.ID())
}

func TestSessionExpiration(t *testing.T) {
	// 使用较短的过期时间进行测试
	store := NewStore(100 * time.Millisecond)
	ctx := context.Background()

	// 生成会话
	sess, err := store.Generate(ctx, "test-id")
	assert.NoError(t, err)
	// 设置会话数据
	err = sess.Set(ctx, "test-key", "test-value")
	assert.NoError(t, err)

	// 等待会话过期
	time.Sleep(200 * time.Millisecond)

	// 验证会话已过期
	sess, err = store.Get(ctx, "test-id")
	assert.Error(t, err)
	assert.Nil(t, sess)
}

func TestConcurrentAccess(t *testing.T) {
	store := NewStore(30 * time.Minute)
	ctx := context.Background()

	// 生成会话
	sess, err := store.Generate(ctx, "test-id")
	assert.NoError(t, err)

	// 并发访问会话
	doneChan := make(chan bool)
	for i := 0; i < 10; i++ {
		go func(val int) {
			key := fmt.Sprintf("key-%d", val)
			err := sess.Set(ctx, key, val)
			assert.NoError(t, err)
			doneChan <- true
		}(i)
	}

	// 等待所有goroutine完成
	for i := 0; i < 10; i++ {
		<-doneChan
	}

	// 验证所有数据都被正确设置
	for i := 0; i < 10; i++ {
		key := fmt.Sprintf("key-%d", i)
		val, err := sess.Get(ctx, key)
		assert.NoError(t, err)
		assert.Equal(t, i, val)
	}
}
