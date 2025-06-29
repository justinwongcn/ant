// Package memory 提供基于内存的仓储实现
package memory

import (
	"context"
	"fmt"
	"sync"

	"github.com/justinwongcn/ant/internal/domain/shared/errors"
	"github.com/justinwongcn/ant/internal/domain/webserver/aggregates"
	"github.com/justinwongcn/ant/internal/domain/webserver/repositories"
)

// WebServerRepository 基于内存的Web服务器仓储实现
type WebServerRepository struct {
	servers map[aggregates.ServerID]*aggregates.WebServer
	mu      sync.RWMutex
}

// NewWebServerRepository 创建新的内存Web服务器仓储
func NewWebServerRepository() *WebServerRepository {
	return &WebServerRepository{
		servers: make(map[aggregates.ServerID]*aggregates.WebServer),
	}
}

// Save 保存Web服务器聚合
func (r *WebServerRepository) Save(ctx context.Context, server *aggregates.WebServer) error {
	if server == nil {
		return errors.NewValidationError("server", "服务器不能为空")
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	r.servers[server.ID()] = server
	return nil
}

// FindByID 根据ID查找Web服务器
func (r *WebServerRepository) FindByID(ctx context.Context, id aggregates.ServerID) (*aggregates.WebServer, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	server, exists := r.servers[id]
	if !exists {
		return nil, errors.NewDomainError("SERVER_NOT_FOUND", fmt.Sprintf("服务器不存在: %s", id))
	}

	return server, nil
}

// FindByName 根据名称查找Web服务器
func (r *WebServerRepository) FindByName(ctx context.Context, name string) (*aggregates.WebServer, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, server := range r.servers {
		if server.Name() == name {
			return server, nil
		}
	}

	return nil, errors.NewDomainError("SERVER_NOT_FOUND", fmt.Sprintf("服务器不存在: %s", name))
}

// FindAll 查找所有Web服务器
func (r *WebServerRepository) FindAll(ctx context.Context) ([]*aggregates.WebServer, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	servers := make([]*aggregates.WebServer, 0, len(r.servers))
	for _, server := range r.servers {
		servers = append(servers, server)
	}

	return servers, nil
}

// FindRunning 查找所有运行中的Web服务器
func (r *WebServerRepository) FindRunning(ctx context.Context) ([]*aggregates.WebServer, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var runningServers []*aggregates.WebServer
	for _, server := range r.servers {
		if server.IsRunning() {
			runningServers = append(runningServers, server)
		}
	}

	return runningServers, nil
}

// Remove 移除Web服务器
func (r *WebServerRepository) Remove(ctx context.Context, id aggregates.ServerID) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.servers[id]; !exists {
		return errors.NewDomainError("SERVER_NOT_FOUND", fmt.Sprintf("服务器不存在: %s", id))
	}

	delete(r.servers, id)
	return nil
}

// Exists 检查Web服务器是否存在
func (r *WebServerRepository) Exists(ctx context.Context, id aggregates.ServerID) (bool, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	_, exists := r.servers[id]
	return exists, nil
}

// Count 返回Web服务器总数
func (r *WebServerRepository) Count(ctx context.Context) (int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return len(r.servers), nil
}

// Clear 清空所有数据（用于测试）
func (r *WebServerRepository) Clear() {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.servers = make(map[aggregates.ServerID]*aggregates.WebServer)
}

// GetAll 获取所有服务器（用于调试）
func (r *WebServerRepository) GetAll() map[aggregates.ServerID]*aggregates.WebServer {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make(map[aggregates.ServerID]*aggregates.WebServer)
	for id, server := range r.servers {
		result[id] = server
	}
	return result
}

// 确保实现了接口
var _ repositories.WebServerRepository = (*WebServerRepository)(nil)
