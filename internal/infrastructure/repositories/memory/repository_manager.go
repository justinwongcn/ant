// Package memory 提供基于内存的仓储管理器实现
package memory

import (
	"context"
	"sync"

	"github.com/justinwongcn/ant/internal/domain/webserver/repositories"
)

// RepositoryManager 基于内存的仓储管理器实现
type RepositoryManager struct {
	webServerRepo  *WebServerRepository
	routeRepo      *RouteRepository
	middlewareRepo *MiddlewareRepository
	queryRepo      *QueryRepository
	mu             sync.RWMutex
}

// NewRepositoryManager 创建新的内存仓储管理器
func NewRepositoryManager() *RepositoryManager {
	webServerRepo := NewWebServerRepository()
	routeRepo := NewRouteRepository()
	middlewareRepo := NewMiddlewareRepository()
	queryRepo := NewQueryRepository(webServerRepo, routeRepo, middlewareRepo)

	return &RepositoryManager{
		webServerRepo:  webServerRepo,
		routeRepo:      routeRepo,
		middlewareRepo: middlewareRepo,
		queryRepo:      queryRepo,
	}
}

// WebServer 返回Web服务器仓储
func (m *RepositoryManager) WebServer() repositories.WebServerRepository {
	return m.webServerRepo
}

// Route 返回路由仓储
func (m *RepositoryManager) Route() repositories.RouteRepository {
	return m.routeRepo
}

// Middleware 返回中间件仓储
func (m *RepositoryManager) Middleware() repositories.MiddlewareRepository {
	return m.middlewareRepo
}

// Query 返回查询仓储
func (m *RepositoryManager) Query() repositories.WebServerQueryRepository {
	return m.queryRepo
}

// BeginTransaction 开始新事务（内存实现中返回自身）
func (m *RepositoryManager) BeginTransaction(ctx context.Context) (repositories.Transaction, error) {
	return &Transaction{
		manager: m,
	}, nil
}

// Clear 清空所有数据（用于测试）
func (m *RepositoryManager) Clear() {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.webServerRepo.Clear()
	m.routeRepo.Clear()
	m.middlewareRepo.Clear()
}

// Transaction 内存事务实现（简化版本）
type Transaction struct {
	manager    *RepositoryManager
	committed  bool
	rolledBack bool
}

// Commit 提交事务
func (t *Transaction) Commit() error {
	if t.rolledBack {
		return repositories.ErrTransactionRolledBack
	}
	if t.committed {
		return repositories.ErrTransactionAlreadyCommitted
	}

	t.committed = true
	return nil
}

// Rollback 回滚事务
func (t *Transaction) Rollback() error {
	if t.committed {
		return repositories.ErrTransactionAlreadyCommitted
	}
	if t.rolledBack {
		return repositories.ErrTransactionRolledBack
	}

	t.rolledBack = true
	return nil
}

// WebServer 返回事务中的Web服务器仓储
func (t *Transaction) WebServer() repositories.WebServerRepository {
	return t.manager.webServerRepo
}

// Route 返回事务中的路由仓储
func (t *Transaction) Route() repositories.RouteRepository {
	return t.manager.routeRepo
}

// Middleware 返回事务中的中间件仓储
func (t *Transaction) Middleware() repositories.MiddlewareRepository {
	return t.manager.middlewareRepo
}

// 确保实现了接口
var (
	_ repositories.RepositoryManager = (*RepositoryManager)(nil)
	_ repositories.Transaction       = (*Transaction)(nil)
)
