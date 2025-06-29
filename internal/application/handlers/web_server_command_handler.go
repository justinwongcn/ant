// Package handlers 包含应用层的命令和查询处理器
package handlers

import (
	"context"
	"fmt"
	"log"

	"github.com/justinwongcn/ant/internal/application/commands"
	"github.com/justinwongcn/ant/internal/application/dto"
	"github.com/justinwongcn/ant/internal/domain/shared/events"
	"github.com/justinwongcn/ant/internal/domain/webserver/aggregates"
	"github.com/justinwongcn/ant/internal/domain/webserver/entities"
	"github.com/justinwongcn/ant/internal/domain/webserver/repositories"
	"github.com/justinwongcn/ant/internal/domain/webserver/valueobjects"
)

// WebServerCommandHandler 处理Web服务器相关的命令
type WebServerCommandHandler struct {
	webServerRepo   repositories.WebServerRepository
	routeRepo       repositories.RouteRepository
	middlewareRepo  repositories.MiddlewareRepository
	eventPublisher  events.EventPublisher
	handlerRegistry HandlerRegistry
}

// HandlerRegistry 处理器注册表接口
type HandlerRegistry interface {
	// GetRouteHandler 获取路由处理器
	GetRouteHandler(name string) (entities.HandlerFunc, error)
	// GetMiddlewareHandler 获取中间件处理器
	GetMiddlewareHandler(name string) (entities.MiddlewareFunc, error)
}

// NewWebServerCommandHandler 创建新的Web服务器命令处理器
func NewWebServerCommandHandler(
	webServerRepo repositories.WebServerRepository,
	routeRepo repositories.RouteRepository,
	middlewareRepo repositories.MiddlewareRepository,
	eventPublisher events.EventPublisher,
	handlerRegistry HandlerRegistry,
) *WebServerCommandHandler {
	return &WebServerCommandHandler{
		webServerRepo:   webServerRepo,
		routeRepo:       routeRepo,
		middlewareRepo:  middlewareRepo,
		eventPublisher:  eventPublisher,
		handlerRegistry: handlerRegistry,
	}
}

// HandleCreateWebServer 处理创建Web服务器命令
func (h *WebServerCommandHandler) HandleCreateWebServer(ctx context.Context, cmd *commands.CreateWebServerCommand) (*dto.CreateWebServerResponse, error) {
	// 验证命令
	if err := cmd.Validate(); err != nil {
		return nil, fmt.Errorf("命令验证失败: %w", err)
	}

	// 创建Web服务器聚合
	server, err := aggregates.NewWebServer(cmd.Name, cmd.Address)
	if err != nil {
		return nil, fmt.Errorf("创建Web服务器失败: %w", err)
	}

	// 保存到仓储
	if err := h.webServerRepo.Save(ctx, server); err != nil {
		return nil, fmt.Errorf("保存Web服务器失败: %w", err)
	}

	// 发布领域事件
	domainEvents := server.GetEvents()
	for _, event := range domainEvents {
		if err := h.eventPublisher.Publish(ctx, event); err != nil {
			log.Printf("发布领域事件失败: %v", err)
		}
	}

	// 构建响应
	response := &dto.CreateWebServerResponse{
		ServerID:  server.ID().String(),
		Name:      server.Name(),
		Address:   server.Address(),
		State:     dto.ConvertServerStateToString(server.State()),
		CreatedAt: server.CreatedAt(),
	}

	return response, nil
}

// HandleStartServer 处理启动服务器命令
func (h *WebServerCommandHandler) HandleStartServer(ctx context.Context, cmd *commands.StartServerCommand) (*dto.StartServerResponse, error) {
	// 验证命令
	if err := cmd.Validate(); err != nil {
		return nil, fmt.Errorf("命令验证失败: %w", err)
	}

	// 获取服务器聚合
	serverID := aggregates.ServerID(cmd.ServerID)
	server, err := h.webServerRepo.FindByID(ctx, serverID)
	if err != nil {
		return nil, fmt.Errorf("获取Web服务器失败: %w", err)
	}

	// 启动服务器
	if err := server.Start(); err != nil {
		return nil, fmt.Errorf("启动服务器失败: %w", err)
	}

	// 保存更新
	if err := h.webServerRepo.Save(ctx, server); err != nil {
		return nil, fmt.Errorf("保存服务器状态失败: %w", err)
	}

	// 发布领域事件
	domainEvents := server.GetEvents()
	for _, event := range domainEvents {
		if err := h.eventPublisher.Publish(ctx, event); err != nil {
			log.Printf("发布领域事件失败: %v", err)
		}
	}

	// 构建响应
	response := &dto.StartServerResponse{
		ServerID:  server.ID().String(),
		State:     dto.ConvertServerStateToString(server.State()),
		StartedAt: server.StartedAt(),
		Message:   "服务器启动成功",
	}

	return response, nil
}

// HandleStopServer 处理停止服务器命令
func (h *WebServerCommandHandler) HandleStopServer(ctx context.Context, cmd *commands.StopServerCommand) (*dto.StopServerResponse, error) {
	// 验证命令
	if err := cmd.Validate(); err != nil {
		return nil, fmt.Errorf("命令验证失败: %w", err)
	}

	// 获取服务器聚合
	serverID := aggregates.ServerID(cmd.ServerID)
	server, err := h.webServerRepo.FindByID(ctx, serverID)
	if err != nil {
		return nil, fmt.Errorf("获取Web服务器失败: %w", err)
	}

	// 停止服务器
	if err := server.Stop(); err != nil {
		return nil, fmt.Errorf("停止服务器失败: %w", err)
	}

	// 保存更新
	if err := h.webServerRepo.Save(ctx, server); err != nil {
		return nil, fmt.Errorf("保存服务器状态失败: %w", err)
	}

	// 发布领域事件
	domainEvents := server.GetEvents()
	for _, event := range domainEvents {
		if err := h.eventPublisher.Publish(ctx, event); err != nil {
			log.Printf("发布领域事件失败: %v", err)
		}
	}

	// 构建响应
	response := &dto.StopServerResponse{
		ServerID:  server.ID().String(),
		State:     dto.ConvertServerStateToString(server.State()),
		StoppedAt: server.StoppedAt(),
		Message:   "服务器停止成功",
	}

	return response, nil
}

// HandleRegisterRoute 处理注册路由命令
func (h *WebServerCommandHandler) HandleRegisterRoute(ctx context.Context, cmd *commands.RegisterRouteCommand) (*dto.RegisterRouteResponse, error) {
	// 验证命令
	if err := cmd.Validate(); err != nil {
		return nil, fmt.Errorf("命令验证失败: %w", err)
	}

	// 获取服务器聚合
	serverID := aggregates.ServerID(cmd.ServerID)
	server, err := h.webServerRepo.FindByID(ctx, serverID)
	if err != nil {
		return nil, fmt.Errorf("获取Web服务器失败: %w", err)
	}

	// 创建HTTP方法值对象
	method, err := valueobjects.NewHTTPMethod(cmd.Method)
	if err != nil {
		return nil, fmt.Errorf("创建HTTP方法失败: %w", err)
	}

	// 创建URL模式值对象
	pattern, err := valueobjects.NewURLPattern(method, cmd.Path)
	if err != nil {
		return nil, fmt.Errorf("创建URL模式失败: %w", err)
	}

	// 获取处理器
	handler, err := h.handlerRegistry.GetRouteHandler(cmd.HandlerName)
	if err != nil {
		return nil, fmt.Errorf("获取路由处理器失败: %w", err)
	}

	// 创建路由实体
	route, err := entities.NewRoute(pattern, handler)
	if err != nil {
		return nil, fmt.Errorf("创建路由失败: %w", err)
	}

	// 设置路由属性
	if cmd.Name != "" {
		if err := route.SetName(cmd.Name); err != nil {
			return nil, fmt.Errorf("设置路由名称失败: %w", err)
		}
	}

	if cmd.Description != "" {
		route.SetDescription(cmd.Description)
	}

	for key, value := range cmd.Metadata {
		if err := route.SetMetadata(key, value); err != nil {
			return nil, fmt.Errorf("设置路由元数据失败: %w", err)
		}
	}

	// 注册路由到服务器
	if err := server.RegisterRoute(route); err != nil {
		return nil, fmt.Errorf("注册路由失败: %w", err)
	}

	// 保存路由到路由仓储
	if err := h.routeRepo.Save(ctx, serverID, route); err != nil {
		return nil, fmt.Errorf("保存路由失败: %w", err)
	}

	// 保存更新的服务器
	if err := h.webServerRepo.Save(ctx, server); err != nil {
		return nil, fmt.Errorf("保存服务器失败: %w", err)
	}

	// 发布领域事件
	domainEvents := server.GetEvents()
	for _, event := range domainEvents {
		if err := h.eventPublisher.Publish(ctx, event); err != nil {
			log.Printf("发布领域事件失败: %v", err)
		}
	}

	// 构建响应
	response := &dto.RegisterRouteResponse{
		RouteID:     route.ID().String(),
		ServerID:    server.ID().String(),
		Method:      method.Value(),
		Path:        pattern.Path(),
		Name:        route.Name(),
		Description: route.Description(),
		CreatedAt:   route.CreatedAt(),
		Message:     "路由注册成功",
	}

	return response, nil
}

// HandleAddMiddleware 处理添加中间件命令
func (h *WebServerCommandHandler) HandleAddMiddleware(ctx context.Context, cmd *commands.AddMiddlewareCommand) (*dto.AddMiddlewareResponse, error) {
	// 验证命令
	if err := cmd.Validate(); err != nil {
		return nil, fmt.Errorf("命令验证失败: %w", err)
	}

	// 获取服务器聚合
	serverID := aggregates.ServerID(cmd.ServerID)
	server, err := h.webServerRepo.FindByID(ctx, serverID)
	if err != nil {
		return nil, fmt.Errorf("获取Web服务器失败: %w", err)
	}

	// 获取中间件处理器
	middlewareFunc, err := h.handlerRegistry.GetMiddlewareHandler(cmd.HandlerName)
	if err != nil {
		return nil, fmt.Errorf("获取中间件处理器失败: %w", err)
	}

	// 转换中间件类型
	middlewareType := convertStringToMiddlewareType(cmd.Type)

	// 创建中间件实体
	middleware, err := entities.NewMiddleware(cmd.Name, middlewareType, middlewareFunc, cmd.Priority)
	if err != nil {
		return nil, fmt.Errorf("创建中间件失败: %w", err)
	}

	// 设置中间件属性
	if cmd.Description != "" {
		middleware.SetDescription(cmd.Description)
	}

	for key, value := range cmd.Metadata {
		if err := middleware.SetMetadata(key, value); err != nil {
			return nil, fmt.Errorf("设置中间件元数据失败: %w", err)
		}
	}

	// 添加中间件到服务器
	if err := server.AddMiddleware(middleware); err != nil {
		return nil, fmt.Errorf("添加中间件失败: %w", err)
	}

	// 保存中间件到中间件仓储
	if err := h.middlewareRepo.Save(ctx, serverID, middleware); err != nil {
		return nil, fmt.Errorf("保存中间件失败: %w", err)
	}

	// 保存更新的服务器
	if err := h.webServerRepo.Save(ctx, server); err != nil {
		return nil, fmt.Errorf("保存服务器失败: %w", err)
	}

	// 构建响应
	response := &dto.AddMiddlewareResponse{
		MiddlewareID: middleware.ID().String(),
		ServerID:     server.ID().String(),
		Name:         middleware.Name(),
		Type:         middleware.Type().String(),
		Priority:     middleware.Priority(),
		CreatedAt:    middleware.CreatedAt(),
		Message:      "中间件添加成功",
	}

	return response, nil
}

// HandleRemoveRoute 处理移除路由命令
func (h *WebServerCommandHandler) HandleRemoveRoute(ctx context.Context, cmd *commands.RemoveRouteCommand) error {
	// 验证命令
	if err := cmd.Validate(); err != nil {
		return fmt.Errorf("命令验证失败: %w", err)
	}

	// 获取服务器聚合
	serverID := aggregates.ServerID(cmd.ServerID)
	server, err := h.webServerRepo.FindByID(ctx, serverID)
	if err != nil {
		return fmt.Errorf("获取Web服务器失败: %w", err)
	}

	// 创建HTTP方法值对象
	method, err := valueobjects.NewHTTPMethod(cmd.Method)
	if err != nil {
		return fmt.Errorf("创建HTTP方法失败: %w", err)
	}

	// 创建URL模式值对象
	pattern, err := valueobjects.NewURLPattern(method, cmd.Path)
	if err != nil {
		return fmt.Errorf("创建URL模式失败: %w", err)
	}

	// 从服务器移除路由
	if err := server.UnregisterRoute(pattern); err != nil {
		return fmt.Errorf("移除路由失败: %w", err)
	}

	// 保存更新
	if err := h.webServerRepo.Save(ctx, server); err != nil {
		return fmt.Errorf("保存服务器失败: %w", err)
	}

	return nil
}

// convertStringToMiddlewareType 将字符串转换为中间件类型
func convertStringToMiddlewareType(typeStr string) entities.MiddlewareType {
	switch typeStr {
	case "auth":
		return entities.MiddlewareTypeAuth
	case "logging":
		return entities.MiddlewareTypeLogging
	case "recovery":
		return entities.MiddlewareTypeRecovery
	case "ratelimit":
		return entities.MiddlewareTypeRateLimit
	case "cors":
		return entities.MiddlewareTypeCORS
	case "compression":
		return entities.MiddlewareTypeCompression
	case "cache":
		return entities.MiddlewareTypeCache
	default:
		return entities.MiddlewareTypeGeneral
	}
}
