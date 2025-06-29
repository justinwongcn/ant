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
	"github.com/justinwongcn/ant/internal/domain/webserver/valueobjects"
)

// RouteRepository 基于内存的路由仓储实现
type RouteRepository struct {
	// routes 存储路由，key为路由ID
	routes map[entities.RouteID]*entities.Route
	// serverRoutes 存储服务器到路由的映射，key为服务器ID
	serverRoutes map[aggregates.ServerID][]entities.RouteID
	mu           sync.RWMutex
}

// NewRouteRepository 创建新的内存路由仓储
func NewRouteRepository() *RouteRepository {
	return &RouteRepository{
		routes:       make(map[entities.RouteID]*entities.Route),
		serverRoutes: make(map[aggregates.ServerID][]entities.RouteID),
	}
}

// Save 保存路由实体
func (r *RouteRepository) Save(ctx context.Context, serverID aggregates.ServerID, route *entities.Route) error {
	if route == nil {
		return errors.NewValidationError("route", "路由不能为空")
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	// 保存路由
	r.routes[route.ID()] = route

	// 更新服务器路由映射
	routeIDs := r.serverRoutes[serverID]

	// 检查是否已存在
	found := false
	for _, id := range routeIDs {
		if id == route.ID() {
			found = true
			break
		}
	}

	if !found {
		r.serverRoutes[serverID] = append(routeIDs, route.ID())
	}

	return nil
}

// FindByID 根据ID查找路由
func (r *RouteRepository) FindByID(ctx context.Context, id entities.RouteID) (*entities.Route, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	route, exists := r.routes[id]
	if !exists {
		return nil, errors.NewDomainError("ROUTE_NOT_FOUND", fmt.Sprintf("路由不存在: %s", id))
	}

	return route, nil
}

// FindByPattern 根据模式查找路由
func (r *RouteRepository) FindByPattern(ctx context.Context, serverID aggregates.ServerID, pattern *valueobjects.URLPattern) (*entities.Route, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	routeIDs, exists := r.serverRoutes[serverID]
	if !exists {
		return nil, errors.NewDomainError("ROUTE_NOT_FOUND", "服务器没有路由")
	}

	for _, routeID := range routeIDs {
		route, exists := r.routes[routeID]
		if exists && route.Pattern().Equals(pattern) {
			return route, nil
		}
	}

	return nil, errors.NewDomainError("ROUTE_NOT_FOUND", fmt.Sprintf("路由不存在: %s", pattern.String()))
}

// FindByServerID 查找服务器的所有路由
func (r *RouteRepository) FindByServerID(ctx context.Context, serverID aggregates.ServerID) ([]*entities.Route, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	routeIDs, exists := r.serverRoutes[serverID]
	if !exists {
		return []*entities.Route{}, nil
	}

	routes := make([]*entities.Route, 0, len(routeIDs))
	for _, routeID := range routeIDs {
		if route, exists := r.routes[routeID]; exists {
			routes = append(routes, route)
		}
	}

	return routes, nil
}

// FindMatching 查找匹配的路由
func (r *RouteRepository) FindMatching(ctx context.Context, serverID aggregates.ServerID, method *valueobjects.HTTPMethod, path string) ([]*entities.Route, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	routeIDs, exists := r.serverRoutes[serverID]
	if !exists {
		return []*entities.Route{}, nil
	}

	var matchingRoutes []*entities.Route
	for _, routeID := range routeIDs {
		route, exists := r.routes[routeID]
		if exists && route.Matches(method, path) {
			matchingRoutes = append(matchingRoutes, route)
		}
	}

	return matchingRoutes, nil
}

// Remove 移除路由
func (r *RouteRepository) Remove(ctx context.Context, id entities.RouteID) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// 检查路由是否存在
	if _, exists := r.routes[id]; !exists {
		return errors.NewDomainError("ROUTE_NOT_FOUND", fmt.Sprintf("路由不存在: %s", id))
	}

	// 从路由映射中删除
	delete(r.routes, id)

	// 从服务器路由映射中删除
	for serverID, routeIDs := range r.serverRoutes {
		for i, routeID := range routeIDs {
			if routeID == id {
				// 删除该路由ID
				r.serverRoutes[serverID] = append(routeIDs[:i], routeIDs[i+1:]...)
				break
			}
		}
	}

	return nil
}

// RemoveByPattern 根据模式移除路由
func (r *RouteRepository) RemoveByPattern(ctx context.Context, serverID aggregates.ServerID, pattern *valueobjects.URLPattern) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	routeIDs, exists := r.serverRoutes[serverID]
	if !exists {
		return errors.NewDomainError("ROUTE_NOT_FOUND", "服务器没有路由")
	}

	for i, routeID := range routeIDs {
		route, exists := r.routes[routeID]
		if exists && route.Pattern().Equals(pattern) {
			// 删除路由
			delete(r.routes, routeID)
			// 从服务器路由映射中删除
			r.serverRoutes[serverID] = append(routeIDs[:i], routeIDs[i+1:]...)
			return nil
		}
	}

	return errors.NewDomainError("ROUTE_NOT_FOUND", fmt.Sprintf("路由不存在: %s", pattern.String()))
}

// RemoveByServerID 移除服务器的所有路由
func (r *RouteRepository) RemoveByServerID(ctx context.Context, serverID aggregates.ServerID) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	routeIDs, exists := r.serverRoutes[serverID]
	if !exists {
		return nil // 没有路由需要删除
	}

	// 删除所有路由
	for _, routeID := range routeIDs {
		delete(r.routes, routeID)
	}

	// 清空服务器路由映射
	delete(r.serverRoutes, serverID)

	return nil
}

// Exists 检查路由是否存在
func (r *RouteRepository) Exists(ctx context.Context, id entities.RouteID) (bool, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	_, exists := r.routes[id]
	return exists, nil
}

// ExistsByPattern 检查路由模式是否存在
func (r *RouteRepository) ExistsByPattern(ctx context.Context, serverID aggregates.ServerID, pattern *valueobjects.URLPattern) (bool, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	routeIDs, exists := r.serverRoutes[serverID]
	if !exists {
		return false, nil
	}

	for _, routeID := range routeIDs {
		route, exists := r.routes[routeID]
		if exists && route.Pattern().Equals(pattern) {
			return true, nil
		}
	}

	return false, nil
}

// Count 返回服务器的路由总数
func (r *RouteRepository) Count(ctx context.Context, serverID aggregates.ServerID) (int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	routeIDs, exists := r.serverRoutes[serverID]
	if !exists {
		return 0, nil
	}

	return len(routeIDs), nil
}

// Clear 清空所有数据（用于测试）
func (r *RouteRepository) Clear() {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.routes = make(map[entities.RouteID]*entities.Route)
	r.serverRoutes = make(map[aggregates.ServerID][]entities.RouteID)
}

// 确保实现了接口
var _ repositories.RouteRepository = (*RouteRepository)(nil)
