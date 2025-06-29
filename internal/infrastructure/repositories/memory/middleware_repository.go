// Package memory 提供基于内存的仓储实现
package memory

import (
	"context"
	"fmt"
	"sync"

	"github.com/justinwongcn/ant/internal/domain/shared/errors"
	"github.com/justinwongcn/ant/internal/domain/webserver/aggregates"
	"github.com/justinwongcn/ant/internal/domain/webserver/entities"
	"github.com/justinwongcn/ant/internal/domain/webserver/repositories"
)

// MiddlewareRepository 基于内存的中间件仓储实现
type MiddlewareRepository struct {
	// middlewares 存储中间件，key为中间件ID
	middlewares map[entities.MiddlewareID]*entities.Middleware
	// serverMiddlewares 存储服务器到中间件的映射，key为服务器ID
	serverMiddlewares map[aggregates.ServerID][]entities.MiddlewareID
	mu                sync.RWMutex
}

// NewMiddlewareRepository 创建新的内存中间件仓储
func NewMiddlewareRepository() *MiddlewareRepository {
	return &MiddlewareRepository{
		middlewares:       make(map[entities.MiddlewareID]*entities.Middleware),
		serverMiddlewares: make(map[aggregates.ServerID][]entities.MiddlewareID),
	}
}

// Save 保存中间件实体
func (r *MiddlewareRepository) Save(ctx context.Context, serverID aggregates.ServerID, middleware *entities.Middleware) error {
	if middleware == nil {
		return errors.NewValidationError("middleware", "中间件不能为空")
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	// 保存中间件
	r.middlewares[middleware.ID()] = middleware

	// 更新服务器中间件映射
	middlewareIDs := r.serverMiddlewares[serverID]

	// 检查是否已存在
	found := false
	for _, id := range middlewareIDs {
		if id == middleware.ID() {
			found = true
			break
		}
	}

	if !found {
		r.serverMiddlewares[serverID] = append(middlewareIDs, middleware.ID())
	}

	return nil
}

// FindByID 根据ID查找中间件
func (r *MiddlewareRepository) FindByID(ctx context.Context, id entities.MiddlewareID) (*entities.Middleware, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	middleware, exists := r.middlewares[id]
	if !exists {
		return nil, errors.NewDomainError("MIDDLEWARE_NOT_FOUND", fmt.Sprintf("中间件不存在: %s", id))
	}

	return middleware, nil
}

// FindByName 根据名称查找中间件
func (r *MiddlewareRepository) FindByName(ctx context.Context, serverID aggregates.ServerID, name string) (*entities.Middleware, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	middlewareIDs, exists := r.serverMiddlewares[serverID]
	if !exists {
		return nil, errors.NewDomainError("MIDDLEWARE_NOT_FOUND", "服务器没有中间件")
	}

	for _, middlewareID := range middlewareIDs {
		middleware, exists := r.middlewares[middlewareID]
		if exists && middleware.Name() == name {
			return middleware, nil
		}
	}

	return nil, errors.NewDomainError("MIDDLEWARE_NOT_FOUND", fmt.Sprintf("中间件不存在: %s", name))
}

// FindByServerID 查找服务器的所有中间件
func (r *MiddlewareRepository) FindByServerID(ctx context.Context, serverID aggregates.ServerID) ([]*entities.Middleware, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	middlewareIDs, exists := r.serverMiddlewares[serverID]
	if !exists {
		return []*entities.Middleware{}, nil
	}

	middlewares := make([]*entities.Middleware, 0, len(middlewareIDs))
	for _, middlewareID := range middlewareIDs {
		if middleware, exists := r.middlewares[middlewareID]; exists {
			middlewares = append(middlewares, middleware)
		}
	}

	return middlewares, nil
}

// FindByType 根据类型查找中间件
func (r *MiddlewareRepository) FindByType(ctx context.Context, serverID aggregates.ServerID, middlewareType entities.MiddlewareType) ([]*entities.Middleware, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	middlewareIDs, exists := r.serverMiddlewares[serverID]
	if !exists {
		return []*entities.Middleware{}, nil
	}

	var typedMiddlewares []*entities.Middleware
	for _, middlewareID := range middlewareIDs {
		middleware, exists := r.middlewares[middlewareID]
		if exists && middleware.Type() == middlewareType {
			typedMiddlewares = append(typedMiddlewares, middleware)
		}
	}

	return typedMiddlewares, nil
}

// FindEnabled 查找所有启用的中间件
func (r *MiddlewareRepository) FindEnabled(ctx context.Context, serverID aggregates.ServerID) ([]*entities.Middleware, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	middlewareIDs, exists := r.serverMiddlewares[serverID]
	if !exists {
		return []*entities.Middleware{}, nil
	}

	var enabledMiddlewares []*entities.Middleware
	for _, middlewareID := range middlewareIDs {
		middleware, exists := r.middlewares[middlewareID]
		if exists && middleware.IsEnabled() {
			enabledMiddlewares = append(enabledMiddlewares, middleware)
		}
	}

	return enabledMiddlewares, nil
}

// Remove 移除中间件
func (r *MiddlewareRepository) Remove(ctx context.Context, id entities.MiddlewareID) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// 检查中间件是否存在
	if _, exists := r.middlewares[id]; !exists {
		return errors.NewDomainError("MIDDLEWARE_NOT_FOUND", fmt.Sprintf("中间件不存在: %s", id))
	}

	// 从中间件映射中删除
	delete(r.middlewares, id)

	// 从服务器中间件映射中删除
	for serverID, middlewareIDs := range r.serverMiddlewares {
		for i, middlewareID := range middlewareIDs {
			if middlewareID == id {
				// 删除该中间件ID
				r.serverMiddlewares[serverID] = append(middlewareIDs[:i], middlewareIDs[i+1:]...)
				break
			}
		}
	}

	return nil
}

// RemoveByName 根据名称移除中间件
func (r *MiddlewareRepository) RemoveByName(ctx context.Context, serverID aggregates.ServerID, name string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	middlewareIDs, exists := r.serverMiddlewares[serverID]
	if !exists {
		return errors.NewDomainError("MIDDLEWARE_NOT_FOUND", "服务器没有中间件")
	}

	for i, middlewareID := range middlewareIDs {
		middleware, exists := r.middlewares[middlewareID]
		if exists && middleware.Name() == name {
			// 删除中间件
			delete(r.middlewares, middlewareID)
			// 从服务器中间件映射中删除
			r.serverMiddlewares[serverID] = append(middlewareIDs[:i], middlewareIDs[i+1:]...)
			return nil
		}
	}

	return errors.NewDomainError("MIDDLEWARE_NOT_FOUND", fmt.Sprintf("中间件不存在: %s", name))
}

// RemoveByServerID 移除服务器的所有中间件
func (r *MiddlewareRepository) RemoveByServerID(ctx context.Context, serverID aggregates.ServerID) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	middlewareIDs, exists := r.serverMiddlewares[serverID]
	if !exists {
		return nil // 没有中间件需要删除
	}

	// 删除所有中间件
	for _, middlewareID := range middlewareIDs {
		delete(r.middlewares, middlewareID)
	}

	// 清空服务器中间件映射
	delete(r.serverMiddlewares, serverID)

	return nil
}

// Exists 检查中间件是否存在
func (r *MiddlewareRepository) Exists(ctx context.Context, id entities.MiddlewareID) (bool, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	_, exists := r.middlewares[id]
	return exists, nil
}

// ExistsByName 检查中间件名称是否存在
func (r *MiddlewareRepository) ExistsByName(ctx context.Context, serverID aggregates.ServerID, name string) (bool, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	middlewareIDs, exists := r.serverMiddlewares[serverID]
	if !exists {
		return false, nil
	}

	for _, middlewareID := range middlewareIDs {
		middleware, exists := r.middlewares[middlewareID]
		if exists && middleware.Name() == name {
			return true, nil
		}
	}

	return false, nil
}

// Count 返回服务器的中间件总数
func (r *MiddlewareRepository) Count(ctx context.Context, serverID aggregates.ServerID) (int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	middlewareIDs, exists := r.serverMiddlewares[serverID]
	if !exists {
		return 0, nil
	}

	return len(middlewareIDs), nil
}

// Clear 清空所有数据（用于测试）
func (r *MiddlewareRepository) Clear() {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.middlewares = make(map[entities.MiddlewareID]*entities.Middleware)
	r.serverMiddlewares = make(map[aggregates.ServerID][]entities.MiddlewareID)
}

// 确保实现了接口
var _ repositories.MiddlewareRepository = (*MiddlewareRepository)(nil)
