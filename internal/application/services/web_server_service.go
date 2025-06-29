// Package services 包含应用层的服务实现
package services

import (
	"context"

	"github.com/justinwongcn/ant/internal/application/commands"
	"github.com/justinwongcn/ant/internal/application/dto"
	"github.com/justinwongcn/ant/internal/application/handlers"
	"github.com/justinwongcn/ant/internal/application/queries"
)

// WebServerApplicationService Web服务器应用服务
type WebServerApplicationService struct {
	commandHandler *handlers.WebServerCommandHandler
	queryHandler   *handlers.WebServerQueryHandler
}

// NewWebServerApplicationService 创建新的Web服务器应用服务
func NewWebServerApplicationService(
	commandHandler *handlers.WebServerCommandHandler,
	queryHandler *handlers.WebServerQueryHandler,
) *WebServerApplicationService {
	return &WebServerApplicationService{
		commandHandler: commandHandler,
		queryHandler:   queryHandler,
	}
}

// CreateServer 创建Web服务器
func (s *WebServerApplicationService) CreateServer(ctx context.Context, req *dto.CreateWebServerRequest) (*dto.CreateWebServerResponse, error) {
	cmd := commands.NewCreateWebServerCommand(req.Name, req.Address)
	return s.commandHandler.HandleCreateWebServer(ctx, cmd)
}

// StartServer 启动服务器
func (s *WebServerApplicationService) StartServer(ctx context.Context, req *dto.StartServerRequest) (*dto.StartServerResponse, error) {
	cmd := commands.NewStartServerCommand(req.ServerID)
	return s.commandHandler.HandleStartServer(ctx, cmd)
}

// StopServer 停止服务器
func (s *WebServerApplicationService) StopServer(ctx context.Context, req *dto.StopServerRequest) (*dto.StopServerResponse, error) {
	cmd := commands.NewStopServerCommand(req.ServerID)
	return s.commandHandler.HandleStopServer(ctx, cmd)
}

// RegisterRoute 注册路由
func (s *WebServerApplicationService) RegisterRoute(ctx context.Context, req *dto.RegisterRouteRequest) (*dto.RegisterRouteResponse, error) {
	cmd := commands.NewRegisterRouteCommand(req.ServerID, req.Method, req.Path, req.HandlerName)

	if req.Name != "" {
		cmd.WithName(req.Name)
	}

	if req.Description != "" {
		cmd.WithDescription(req.Description)
	}

	for key, value := range req.Metadata {
		cmd.WithMetadata(key, value)
	}

	return s.commandHandler.HandleRegisterRoute(ctx, cmd)
}

// AddMiddleware 添加中间件
func (s *WebServerApplicationService) AddMiddleware(ctx context.Context, req *dto.AddMiddlewareRequest) (*dto.AddMiddlewareResponse, error) {
	cmd := commands.NewAddMiddlewareCommand(req.ServerID, req.Name, req.Type, req.HandlerName, req.Priority)

	if req.Description != "" {
		cmd.WithDescription(req.Description)
	}

	for key, value := range req.Metadata {
		cmd.WithMetadata(key, value)
	}

	return s.commandHandler.HandleAddMiddleware(ctx, cmd)
}

// RemoveRoute 移除路由
func (s *WebServerApplicationService) RemoveRoute(ctx context.Context, serverID, method, path string) error {
	cmd := commands.NewRemoveRouteCommand(serverID, method, path)
	return s.commandHandler.HandleRemoveRoute(ctx, cmd)
}

// GetServer 获取服务器信息
func (s *WebServerApplicationService) GetServer(ctx context.Context, req *dto.GetServerRequest) (*dto.GetServerResponse, error) {
	query := queries.NewGetServerQuery(req.ServerID)
	return s.queryHandler.HandleGetServer(ctx, query)
}

// ListServers 列出服务器
func (s *WebServerApplicationService) ListServers(ctx context.Context, req *dto.ListServersRequest) (*dto.ListServersResponse, error) {
	query := queries.NewListServersQuery()

	if req.State != "" {
		query.WithState(req.State)
	}

	if req.Limit > 0 || req.Offset > 0 {
		limit := req.Limit
		if limit <= 0 {
			limit = 10 // 默认限制
		}
		query.WithPagination(limit, req.Offset)
	}

	return s.queryHandler.HandleListServers(ctx, query)
}

// GetRoutes 获取路由列表
func (s *WebServerApplicationService) GetRoutes(ctx context.Context, req *dto.GetRoutesRequest) (*dto.GetRoutesResponse, error) {
	query := queries.NewGetRoutesQuery(req.ServerID)

	if req.Method != "" {
		query.WithMethod(req.Method)
	}

	if req.Enabled != nil {
		query.WithEnabled(*req.Enabled)
	}

	if req.Limit > 0 || req.Offset > 0 {
		limit := req.Limit
		if limit <= 0 {
			limit = 10 // 默认限制
		}
		query.WithPagination(limit, req.Offset)
	}

	return s.queryHandler.HandleGetRoutes(ctx, query)
}

// GetMiddlewares 获取中间件列表
func (s *WebServerApplicationService) GetMiddlewares(ctx context.Context, req *dto.GetMiddlewaresRequest) (*dto.GetMiddlewaresResponse, error) {
	query := queries.NewGetMiddlewaresQuery(req.ServerID)

	if req.Type != "" {
		query.WithType(req.Type)
	}

	if req.Enabled != nil {
		query.WithEnabled(*req.Enabled)
	}

	if req.Limit > 0 || req.Offset > 0 {
		limit := req.Limit
		if limit <= 0 {
			limit = 10 // 默认限制
		}
		query.WithPagination(limit, req.Offset)
	}

	return s.queryHandler.HandleGetMiddlewares(ctx, query)
}

// GetServerStats 获取服务器统计信息
func (s *WebServerApplicationService) GetServerStats(ctx context.Context, serverID string) (*dto.ServerStatsResponse, error) {
	query := queries.NewGetServerStatsQuery(serverID)
	return s.queryHandler.HandleGetServerStats(ctx, query)
}

// FindMatchingRoute 查找匹配的路由
func (s *WebServerApplicationService) FindMatchingRoute(ctx context.Context, serverID, method, path string) (*dto.RouteInfo, map[string]string, error) {
	query := queries.NewFindMatchingRouteQuery(serverID, method, path)
	return s.queryHandler.HandleFindMatchingRoute(ctx, query)
}

// WebServerService 定义Web服务器应用服务接口
type WebServerService interface {
	// 命令操作
	CreateServer(ctx context.Context, req *dto.CreateWebServerRequest) (*dto.CreateWebServerResponse, error)
	StartServer(ctx context.Context, req *dto.StartServerRequest) (*dto.StartServerResponse, error)
	StopServer(ctx context.Context, req *dto.StopServerRequest) (*dto.StopServerResponse, error)
	RegisterRoute(ctx context.Context, req *dto.RegisterRouteRequest) (*dto.RegisterRouteResponse, error)
	AddMiddleware(ctx context.Context, req *dto.AddMiddlewareRequest) (*dto.AddMiddlewareResponse, error)
	RemoveRoute(ctx context.Context, serverID, method, path string) error

	// 查询操作
	GetServer(ctx context.Context, req *dto.GetServerRequest) (*dto.GetServerResponse, error)
	ListServers(ctx context.Context, req *dto.ListServersRequest) (*dto.ListServersResponse, error)
	GetRoutes(ctx context.Context, req *dto.GetRoutesRequest) (*dto.GetRoutesResponse, error)
	GetMiddlewares(ctx context.Context, req *dto.GetMiddlewaresRequest) (*dto.GetMiddlewaresResponse, error)
	GetServerStats(ctx context.Context, serverID string) (*dto.ServerStatsResponse, error)
	FindMatchingRoute(ctx context.Context, serverID, method, path string) (*dto.RouteInfo, map[string]string, error)
}

// 确保 WebServerApplicationService 实现了 WebServerService 接口
var _ WebServerService = (*WebServerApplicationService)(nil)

// ApplicationServiceManager 应用服务管理器
type ApplicationServiceManager struct {
	webServerService WebServerService
}

// NewApplicationServiceManager 创建新的应用服务管理器
func NewApplicationServiceManager(webServerService WebServerService) *ApplicationServiceManager {
	return &ApplicationServiceManager{
		webServerService: webServerService,
	}
}

// WebServer 返回Web服务器应用服务
func (m *ApplicationServiceManager) WebServer() WebServerService {
	return m.webServerService
}

// ServiceManager 定义服务管理器接口
type ServiceManager interface {
	WebServer() WebServerService
}
