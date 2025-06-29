// Package commands 包含应用层的命令定义
package commands

import (
	"time"

	"github.com/google/uuid"
)

// Command 表示一个命令的基础接口
type Command interface {
	// CommandID 返回命令的唯一标识符
	CommandID() string
	// CommandType 返回命令类型
	CommandType() string
	// Timestamp 返回命令创建时间
	Timestamp() time.Time
	// Validate 验证命令的有效性
	Validate() error
}

// BaseCommand 提供命令的基础实现
type BaseCommand struct {
	commandID   string
	commandType string
	timestamp   time.Time
}

// NewBaseCommand 创建新的基础命令
func NewBaseCommand(commandType string) *BaseCommand {
	return &BaseCommand{
		commandID:   uuid.New().String(),
		commandType: commandType,
		timestamp:   time.Now(),
	}
}

func (c *BaseCommand) CommandID() string    { return c.commandID }
func (c *BaseCommand) CommandType() string  { return c.commandType }
func (c *BaseCommand) Timestamp() time.Time { return c.timestamp }

// CreateWebServerCommand 创建Web服务器命令
type CreateWebServerCommand struct {
	*BaseCommand
	Name    string
	Address string
}

// NewCreateWebServerCommand 创建新的创建Web服务器命令
func NewCreateWebServerCommand(name, address string) *CreateWebServerCommand {
	return &CreateWebServerCommand{
		BaseCommand: NewBaseCommand("CreateWebServer"),
		Name:        name,
		Address:     address,
	}
}

// Validate 验证创建Web服务器命令
func (c *CreateWebServerCommand) Validate() error {
	if c.Name == "" {
		return ErrInvalidCommandParameter("name", "服务器名称不能为空")
	}
	if c.Address == "" {
		return ErrInvalidCommandParameter("address", "服务器地址不能为空")
	}
	return nil
}

// StartServerCommand 启动服务器命令
type StartServerCommand struct {
	*BaseCommand
	ServerID string
}

// NewStartServerCommand 创建新的启动服务器命令
func NewStartServerCommand(serverID string) *StartServerCommand {
	return &StartServerCommand{
		BaseCommand: NewBaseCommand("StartServer"),
		ServerID:    serverID,
	}
}

// Validate 验证启动服务器命令
func (c *StartServerCommand) Validate() error {
	if c.ServerID == "" {
		return ErrInvalidCommandParameter("serverID", "服务器ID不能为空")
	}
	return nil
}

// StopServerCommand 停止服务器命令
type StopServerCommand struct {
	*BaseCommand
	ServerID string
}

// NewStopServerCommand 创建新的停止服务器命令
func NewStopServerCommand(serverID string) *StopServerCommand {
	return &StopServerCommand{
		BaseCommand: NewBaseCommand("StopServer"),
		ServerID:    serverID,
	}
}

// Validate 验证停止服务器命令
func (c *StopServerCommand) Validate() error {
	if c.ServerID == "" {
		return ErrInvalidCommandParameter("serverID", "服务器ID不能为空")
	}
	return nil
}

// RegisterRouteCommand 注册路由命令
type RegisterRouteCommand struct {
	*BaseCommand
	ServerID    string
	Method      string
	Path        string
	HandlerName string
	Name        string
	Description string
	Metadata    map[string]interface{}
}

// NewRegisterRouteCommand 创建新的注册路由命令
func NewRegisterRouteCommand(serverID, method, path, handlerName string) *RegisterRouteCommand {
	return &RegisterRouteCommand{
		BaseCommand: NewBaseCommand("RegisterRoute"),
		ServerID:    serverID,
		Method:      method,
		Path:        path,
		HandlerName: handlerName,
		Metadata:    make(map[string]interface{}),
	}
}

// WithName 设置路由名称
func (c *RegisterRouteCommand) WithName(name string) *RegisterRouteCommand {
	c.Name = name
	return c
}

// WithDescription 设置路由描述
func (c *RegisterRouteCommand) WithDescription(description string) *RegisterRouteCommand {
	c.Description = description
	return c
}

// WithMetadata 设置路由元数据
func (c *RegisterRouteCommand) WithMetadata(key string, value interface{}) *RegisterRouteCommand {
	c.Metadata[key] = value
	return c
}

// Validate 验证注册路由命令
func (c *RegisterRouteCommand) Validate() error {
	if c.ServerID == "" {
		return ErrInvalidCommandParameter("serverID", "服务器ID不能为空")
	}
	if c.Method == "" {
		return ErrInvalidCommandParameter("method", "HTTP方法不能为空")
	}
	if c.Path == "" {
		return ErrInvalidCommandParameter("path", "路径不能为空")
	}
	if c.HandlerName == "" {
		return ErrInvalidCommandParameter("handlerName", "处理器名称不能为空")
	}
	return nil
}

// AddMiddlewareCommand 添加中间件命令
type AddMiddlewareCommand struct {
	*BaseCommand
	ServerID    string
	Name        string
	Type        string
	Priority    int
	Description string
	Metadata    map[string]interface{}
	HandlerName string
}

// NewAddMiddlewareCommand 创建新的添加中间件命令
func NewAddMiddlewareCommand(serverID, name, middlewareType, handlerName string, priority int) *AddMiddlewareCommand {
	return &AddMiddlewareCommand{
		BaseCommand: NewBaseCommand("AddMiddleware"),
		ServerID:    serverID,
		Name:        name,
		Type:        middlewareType,
		Priority:    priority,
		HandlerName: handlerName,
		Metadata:    make(map[string]interface{}),
	}
}

// WithDescription 设置中间件描述
func (c *AddMiddlewareCommand) WithDescription(description string) *AddMiddlewareCommand {
	c.Description = description
	return c
}

// WithMetadata 设置中间件元数据
func (c *AddMiddlewareCommand) WithMetadata(key string, value interface{}) *AddMiddlewareCommand {
	c.Metadata[key] = value
	return c
}

// Validate 验证添加中间件命令
func (c *AddMiddlewareCommand) Validate() error {
	if c.ServerID == "" {
		return ErrInvalidCommandParameter("serverID", "服务器ID不能为空")
	}
	if c.Name == "" {
		return ErrInvalidCommandParameter("name", "中间件名称不能为空")
	}
	if c.Type == "" {
		return ErrInvalidCommandParameter("type", "中间件类型不能为空")
	}
	if c.HandlerName == "" {
		return ErrInvalidCommandParameter("handlerName", "处理器名称不能为空")
	}
	return nil
}

// RemoveRouteCommand 移除路由命令
type RemoveRouteCommand struct {
	*BaseCommand
	ServerID string
	Method   string
	Path     string
}

// NewRemoveRouteCommand 创建新的移除路由命令
func NewRemoveRouteCommand(serverID, method, path string) *RemoveRouteCommand {
	return &RemoveRouteCommand{
		BaseCommand: NewBaseCommand("RemoveRoute"),
		ServerID:    serverID,
		Method:      method,
		Path:        path,
	}
}

// Validate 验证移除路由命令
func (c *RemoveRouteCommand) Validate() error {
	if c.ServerID == "" {
		return ErrInvalidCommandParameter("serverID", "服务器ID不能为空")
	}
	if c.Method == "" {
		return ErrInvalidCommandParameter("method", "HTTP方法不能为空")
	}
	if c.Path == "" {
		return ErrInvalidCommandParameter("path", "路径不能为空")
	}
	return nil
}

// CommandValidationError 命令验证错误
type CommandValidationError struct {
	Field   string
	Message string
}

func (e *CommandValidationError) Error() string {
	return e.Message
}

// ErrInvalidCommandParameter 创建无效命令参数错误
func ErrInvalidCommandParameter(field, message string) *CommandValidationError {
	return &CommandValidationError{
		Field:   field,
		Message: message,
	}
}
