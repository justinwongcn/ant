// Package server 提供HTTP服务器的接口层实现
package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/justinwongcn/ant/internal/application/services"
	"github.com/justinwongcn/ant/internal/interfaces/http/handlers"
)

// HTTPServer HTTP服务器实现
type HTTPServer struct {
	server      *http.Server
	mux         *http.ServeMux
	appService  services.WebServerService
	httpHandler *handlers.HTTPHandler
	mu          sync.RWMutex
	running     bool
}

// Config HTTP服务器配置
type Config struct {
	Address         string
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	IdleTimeout     time.Duration
	ShutdownTimeout time.Duration
}

// DefaultConfig 返回默认配置
func DefaultConfig() *Config {
	return &Config{
		Address:         ":8080",
		ReadTimeout:     30 * time.Second,
		WriteTimeout:    30 * time.Second,
		IdleTimeout:     120 * time.Second,
		ShutdownTimeout: 30 * time.Second,
	}
}

// NewHTTPServer 创建新的HTTP服务器
func NewHTTPServer(config *Config, appService services.WebServerService) *HTTPServer {
	if config == nil {
		config = DefaultConfig()
	}

	mux := http.NewServeMux()
	httpHandler := handlers.NewHTTPHandler(appService)

	server := &http.Server{
		Addr:         config.Address,
		Handler:      mux,
		ReadTimeout:  config.ReadTimeout,
		WriteTimeout: config.WriteTimeout,
		IdleTimeout:  config.IdleTimeout,
	}

	httpServer := &HTTPServer{
		server:      server,
		mux:         mux,
		appService:  appService,
		httpHandler: httpHandler,
	}

	// 注册路由
	httpServer.registerRoutes()

	return httpServer
}

// registerRoutes 注册HTTP路由
func (s *HTTPServer) registerRoutes() {
	// API路由
	s.mux.HandleFunc("/api/servers", s.httpHandler.HandleServers)
	s.mux.HandleFunc("/api/servers/", s.httpHandler.HandleServerByID)
	s.mux.HandleFunc("/api/servers/{id}/routes", s.httpHandler.HandleRoutes)
	s.mux.HandleFunc("/api/servers/{id}/middlewares", s.httpHandler.HandleMiddlewares)
	s.mux.HandleFunc("/api/servers/{id}/start", s.httpHandler.HandleStartServer)
	s.mux.HandleFunc("/api/servers/{id}/stop", s.httpHandler.HandleStopServer)
	s.mux.HandleFunc("/api/servers/{id}/stats", s.httpHandler.HandleServerStats)

	// 健康检查
	s.mux.HandleFunc("/health", s.httpHandler.HandleHealth)
	s.mux.HandleFunc("/ready", s.httpHandler.HandleReady)

	// 根路径
	s.mux.HandleFunc("/", s.httpHandler.HandleRoot)
}

// Start 启动HTTP服务器
func (s *HTTPServer) Start() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.running {
		return fmt.Errorf("HTTP服务器已经在运行")
	}

	log.Printf("启动HTTP服务器，监听地址: %s", s.server.Addr)

	// 在goroutine中启动服务器
	go func() {
		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("HTTP服务器启动失败: %v", err)
		}
	}()

	s.running = true
	return nil
}

// Stop 停止HTTP服务器
func (s *HTTPServer) Stop(ctx context.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.running {
		return fmt.Errorf("HTTP服务器未运行")
	}

	log.Printf("停止HTTP服务器...")

	// 优雅关闭服务器
	if err := s.server.Shutdown(ctx); err != nil {
		return fmt.Errorf("停止HTTP服务器失败: %w", err)
	}

	s.running = false
	log.Printf("HTTP服务器已停止")
	return nil
}

// IsRunning 检查服务器是否运行中
func (s *HTTPServer) IsRunning() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.running
}

// Address 返回服务器地址
func (s *HTTPServer) Address() string {
	return s.server.Addr
}

// AddRoute 添加自定义路由
func (s *HTTPServer) AddRoute(pattern string, handler http.HandlerFunc) {
	s.mux.HandleFunc(pattern, handler)
}

// AddMiddleware 添加全局中间件
func (s *HTTPServer) AddMiddleware(middleware func(http.Handler) http.Handler) {
	s.server.Handler = middleware(s.server.Handler)
}

// HTTPServerInterface 定义HTTP服务器接口
type HTTPServerInterface interface {
	Start() error
	Stop(ctx context.Context) error
	IsRunning() bool
	Address() string
	AddRoute(pattern string, handler http.HandlerFunc)
	AddMiddleware(middleware func(http.Handler) http.Handler)
}

// 确保实现了接口
var _ HTTPServerInterface = (*HTTPServer)(nil)
