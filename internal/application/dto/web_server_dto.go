// Package dto 包含应用层的数据传输对象
package dto

import (
	"time"

	"github.com/justinwongcn/ant/internal/domain/webserver/aggregates"
)

// CreateWebServerRequest 创建Web服务器的请求
type CreateWebServerRequest struct {
	Name    string `json:"name" validate:"required,min=1,max=100"`
	Address string `json:"address" validate:"required"`
}

// CreateWebServerResponse 创建Web服务器的响应
type CreateWebServerResponse struct {
	ServerID  string    `json:"server_id"`
	Name      string    `json:"name"`
	Address   string    `json:"address"`
	State     string    `json:"state"`
	CreatedAt time.Time `json:"created_at"`
}

// StartServerRequest 启动服务器的请求
type StartServerRequest struct {
	ServerID string `json:"server_id" validate:"required"`
}

// StartServerResponse 启动服务器的响应
type StartServerResponse struct {
	ServerID  string     `json:"server_id"`
	State     string     `json:"state"`
	StartedAt *time.Time `json:"started_at"`
	Message   string     `json:"message"`
}

// StopServerRequest 停止服务器的请求
type StopServerRequest struct {
	ServerID string `json:"server_id" validate:"required"`
}

// StopServerResponse 停止服务器的响应
type StopServerResponse struct {
	ServerID  string     `json:"server_id"`
	State     string     `json:"state"`
	StoppedAt *time.Time `json:"stopped_at"`
	Message   string     `json:"message"`
}

// RegisterRouteRequest 注册路由的请求
type RegisterRouteRequest struct {
	ServerID    string                 `json:"server_id" validate:"required"`
	Method      string                 `json:"method" validate:"required"`
	Path        string                 `json:"path" validate:"required"`
	HandlerName string                 `json:"handler_name" validate:"required"`
	Name        string                 `json:"name,omitempty"`
	Description string                 `json:"description,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

// RegisterRouteResponse 注册路由的响应
type RegisterRouteResponse struct {
	RouteID     string    `json:"route_id"`
	ServerID    string    `json:"server_id"`
	Method      string    `json:"method"`
	Path        string    `json:"path"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	Message     string    `json:"message"`
}

// AddMiddlewareRequest 添加中间件的请求
type AddMiddlewareRequest struct {
	ServerID    string                 `json:"server_id" validate:"required"`
	Name        string                 `json:"name" validate:"required"`
	Type        string                 `json:"type" validate:"required"`
	Priority    int                    `json:"priority"`
	Description string                 `json:"description,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
	HandlerName string                 `json:"handler_name" validate:"required"`
}

// AddMiddlewareResponse 添加中间件的响应
type AddMiddlewareResponse struct {
	MiddlewareID string    `json:"middleware_id"`
	ServerID     string    `json:"server_id"`
	Name         string    `json:"name"`
	Type         string    `json:"type"`
	Priority     int       `json:"priority"`
	CreatedAt    time.Time `json:"created_at"`
	Message      string    `json:"message"`
}

// GetServerRequest 获取服务器信息的请求
type GetServerRequest struct {
	ServerID string `json:"server_id" validate:"required"`
}

// GetServerResponse 获取服务器信息的响应
type GetServerResponse struct {
	ServerID        string     `json:"server_id"`
	Name            string     `json:"name"`
	Address         string     `json:"address"`
	State           string     `json:"state"`
	RouteCount      int        `json:"route_count"`
	MiddlewareCount int        `json:"middleware_count"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
	StartedAt       *time.Time `json:"started_at"`
	StoppedAt       *time.Time `json:"stopped_at"`
}

// ListServersRequest 列出服务器的请求
type ListServersRequest struct {
	State  string `json:"state,omitempty"`
	Limit  int    `json:"limit,omitempty"`
	Offset int    `json:"offset,omitempty"`
}

// ListServersResponse 列出服务器的响应
type ListServersResponse struct {
	Servers []GetServerResponse `json:"servers"`
	Total   int                 `json:"total"`
	Limit   int                 `json:"limit"`
	Offset  int                 `json:"offset"`
}

// RouteInfo 路由信息
type RouteInfo struct {
	RouteID     string                 `json:"route_id"`
	Method      string                 `json:"method"`
	Path        string                 `json:"path"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Priority    int                    `json:"priority"`
	Enabled     bool                   `json:"enabled"`
	Metadata    map[string]interface{} `json:"metadata"`
	CreatedAt   time.Time              `json:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at"`
}

// MiddlewareInfo 中间件信息
type MiddlewareInfo struct {
	MiddlewareID string                 `json:"middleware_id"`
	Name         string                 `json:"name"`
	Type         string                 `json:"type"`
	Priority     int                    `json:"priority"`
	Description  string                 `json:"description"`
	Enabled      bool                   `json:"enabled"`
	Metadata     map[string]interface{} `json:"metadata"`
	CreatedAt    time.Time              `json:"created_at"`
	UpdatedAt    time.Time              `json:"updated_at"`
}

// GetRoutesRequest 获取路由列表的请求
type GetRoutesRequest struct {
	ServerID string `json:"server_id" validate:"required"`
	Method   string `json:"method,omitempty"`
	Enabled  *bool  `json:"enabled,omitempty"`
	Limit    int    `json:"limit,omitempty"`
	Offset   int    `json:"offset,omitempty"`
}

// GetRoutesResponse 获取路由列表的响应
type GetRoutesResponse struct {
	Routes []RouteInfo `json:"routes"`
	Total  int         `json:"total"`
	Limit  int         `json:"limit"`
	Offset int         `json:"offset"`
}

// GetMiddlewaresRequest 获取中间件列表的请求
type GetMiddlewaresRequest struct {
	ServerID string `json:"server_id" validate:"required"`
	Type     string `json:"type,omitempty"`
	Enabled  *bool  `json:"enabled,omitempty"`
	Limit    int    `json:"limit,omitempty"`
	Offset   int    `json:"offset,omitempty"`
}

// GetMiddlewaresResponse 获取中间件列表的响应
type GetMiddlewaresResponse struct {
	Middlewares []MiddlewareInfo `json:"middlewares"`
	Total       int              `json:"total"`
	Limit       int              `json:"limit"`
	Offset      int              `json:"offset"`
}

// ServerStatsResponse 服务器统计信息响应
type ServerStatsResponse struct {
	ServerID        string `json:"server_id"`
	ServerName      string `json:"server_name"`
	State           string `json:"state"`
	RouteCount      int    `json:"route_count"`
	MiddlewareCount int    `json:"middleware_count"`
	Uptime          string `json:"uptime,omitempty"`
	CreatedAt       string `json:"created_at"`
	StartedAt       string `json:"started_at,omitempty"`
}

// ConvertServerStateToString 将服务器状态转换为字符串
func ConvertServerStateToString(state aggregates.ServerState) string {
	return state.String()
}

// ConvertStringToServerState 将字符串转换为服务器状态
func ConvertStringToServerState(state string) aggregates.ServerState {
	switch state {
	case "stopped":
		return aggregates.ServerStateStopped
	case "starting":
		return aggregates.ServerStateStarting
	case "running":
		return aggregates.ServerStateRunning
	case "stopping":
		return aggregates.ServerStateStopping
	default:
		return aggregates.ServerStateStopped
	}
}
