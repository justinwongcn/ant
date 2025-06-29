// Package handlers 提供HTTP处理器的接口层实现
package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/justinwongcn/ant/internal/application/dto"
	"github.com/justinwongcn/ant/internal/application/services"
)

// HTTPHandler HTTP请求处理器
type HTTPHandler struct {
	appService services.WebServerService
}

// NewHTTPHandler 创建新的HTTP处理器
func NewHTTPHandler(appService services.WebServerService) *HTTPHandler {
	return &HTTPHandler{
		appService: appService,
	}
}

// HandleServers 处理服务器相关请求
func (h *HTTPHandler) HandleServers(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.handleListServers(w, r)
	case http.MethodPost:
		h.handleCreateServer(w, r)
	default:
		h.writeError(w, http.StatusMethodNotAllowed, "方法不允许")
	}
}

// HandleServerByID 处理特定服务器的请求
func (h *HTTPHandler) HandleServerByID(w http.ResponseWriter, r *http.Request) {
	// 从URL路径中提取服务器ID
	path := strings.TrimPrefix(r.URL.Path, "/api/servers/")
	parts := strings.Split(path, "/")
	if len(parts) == 0 || parts[0] == "" {
		h.writeError(w, http.StatusBadRequest, "无效的服务器ID")
		return
	}

	serverID := parts[0]

	switch r.Method {
	case http.MethodGet:
		h.handleGetServer(w, r, serverID)
	case http.MethodDelete:
		h.handleDeleteServer(w, r, serverID)
	default:
		h.writeError(w, http.StatusMethodNotAllowed, "方法不允许")
	}
}

// HandleRoutes 处理路由相关请求
func (h *HTTPHandler) HandleRoutes(w http.ResponseWriter, r *http.Request) {
	serverID := h.extractServerID(r.URL.Path)
	if serverID == "" {
		h.writeError(w, http.StatusBadRequest, "无效的服务器ID")
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.handleGetRoutes(w, r, serverID)
	case http.MethodPost:
		h.handleRegisterRoute(w, r, serverID)
	default:
		h.writeError(w, http.StatusMethodNotAllowed, "方法不允许")
	}
}

// HandleMiddlewares 处理中间件相关请求
func (h *HTTPHandler) HandleMiddlewares(w http.ResponseWriter, r *http.Request) {
	serverID := h.extractServerID(r.URL.Path)
	if serverID == "" {
		h.writeError(w, http.StatusBadRequest, "无效的服务器ID")
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.handleGetMiddlewares(w, r, serverID)
	case http.MethodPost:
		h.handleAddMiddleware(w, r, serverID)
	default:
		h.writeError(w, http.StatusMethodNotAllowed, "方法不允许")
	}
}

// HandleStartServer 处理启动服务器请求
func (h *HTTPHandler) HandleStartServer(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.writeError(w, http.StatusMethodNotAllowed, "方法不允许")
		return
	}

	serverID := h.extractServerID(r.URL.Path)
	if serverID == "" {
		h.writeError(w, http.StatusBadRequest, "无效的服务器ID")
		return
	}

	req := &dto.StartServerRequest{
		ServerID: serverID,
	}

	response, err := h.appService.StartServer(r.Context(), req)
	if err != nil {
		h.writeError(w, http.StatusInternalServerError, fmt.Sprintf("启动服务器失败: %v", err))
		return
	}

	h.writeJSON(w, http.StatusOK, response)
}

// HandleStopServer 处理停止服务器请求
func (h *HTTPHandler) HandleStopServer(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.writeError(w, http.StatusMethodNotAllowed, "方法不允许")
		return
	}

	serverID := h.extractServerID(r.URL.Path)
	if serverID == "" {
		h.writeError(w, http.StatusBadRequest, "无效的服务器ID")
		return
	}

	req := &dto.StopServerRequest{
		ServerID: serverID,
	}

	response, err := h.appService.StopServer(r.Context(), req)
	if err != nil {
		h.writeError(w, http.StatusInternalServerError, fmt.Sprintf("停止服务器失败: %v", err))
		return
	}

	h.writeJSON(w, http.StatusOK, response)
}

// HandleServerStats 处理服务器统计信息请求
func (h *HTTPHandler) HandleServerStats(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.writeError(w, http.StatusMethodNotAllowed, "方法不允许")
		return
	}

	serverID := h.extractServerID(r.URL.Path)
	if serverID == "" {
		h.writeError(w, http.StatusBadRequest, "无效的服务器ID")
		return
	}

	response, err := h.appService.GetServerStats(r.Context(), serverID)
	if err != nil {
		h.writeError(w, http.StatusInternalServerError, fmt.Sprintf("获取服务器统计信息失败: %v", err))
		return
	}

	h.writeJSON(w, http.StatusOK, response)
}

// HandleHealth 处理健康检查请求
func (h *HTTPHandler) HandleHealth(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.writeError(w, http.StatusMethodNotAllowed, "方法不允许")
		return
	}

	response := map[string]interface{}{
		"status":    "healthy",
		"timestamp": "2025-06-28T02:20:41Z",
		"service":   "ant-web-framework",
	}

	h.writeJSON(w, http.StatusOK, response)
}

// HandleReady 处理就绪检查请求
func (h *HTTPHandler) HandleReady(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.writeError(w, http.StatusMethodNotAllowed, "方法不允许")
		return
	}

	response := map[string]interface{}{
		"status":    "ready",
		"timestamp": "2025-06-28T02:20:41Z",
	}

	h.writeJSON(w, http.StatusOK, response)
}

// HandleRoot 处理根路径请求
func (h *HTTPHandler) HandleRoot(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.writeError(w, http.StatusMethodNotAllowed, "方法不允许")
		return
	}

	response := map[string]interface{}{
		"message": "欢迎使用Ant Web框架",
		"version": "2.0.0-ddd",
		"docs":    "/api/docs",
		"health":  "/health",
	}

	h.writeJSON(w, http.StatusOK, response)
}

// handleListServers 处理列出服务器请求
func (h *HTTPHandler) handleListServers(w http.ResponseWriter, r *http.Request) {
	// 解析查询参数
	query := r.URL.Query()
	req := &dto.ListServersRequest{
		State:  query.Get("state"),
		Limit:  h.parseIntParam(query.Get("limit"), 10),
		Offset: h.parseIntParam(query.Get("offset"), 0),
	}

	response, err := h.appService.ListServers(r.Context(), req)
	if err != nil {
		h.writeError(w, http.StatusInternalServerError, fmt.Sprintf("获取服务器列表失败: %v", err))
		return
	}

	h.writeJSON(w, http.StatusOK, response)
}

// handleCreateServer 处理创建服务器请求
func (h *HTTPHandler) handleCreateServer(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateWebServerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, fmt.Sprintf("解析请求失败: %v", err))
		return
	}

	response, err := h.appService.CreateServer(r.Context(), &req)
	if err != nil {
		h.writeError(w, http.StatusInternalServerError, fmt.Sprintf("创建服务器失败: %v", err))
		return
	}

	h.writeJSON(w, http.StatusCreated, response)
}

// handleGetServer 处理获取服务器请求
func (h *HTTPHandler) handleGetServer(w http.ResponseWriter, r *http.Request, serverID string) {
	req := &dto.GetServerRequest{
		ServerID: serverID,
	}

	response, err := h.appService.GetServer(r.Context(), req)
	if err != nil {
		h.writeError(w, http.StatusNotFound, fmt.Sprintf("获取服务器失败: %v", err))
		return
	}

	h.writeJSON(w, http.StatusOK, response)
}

// handleDeleteServer 处理删除服务器请求
func (h *HTTPHandler) handleDeleteServer(w http.ResponseWriter, r *http.Request, serverID string) {
	// 这里应该调用删除服务器的应用服务方法
	// 目前简化实现
	h.writeJSON(w, http.StatusOK, map[string]string{
		"message":  "服务器删除功能待实现",
		"serverID": serverID,
	})
}

// handleGetRoutes 处理获取路由请求
func (h *HTTPHandler) handleGetRoutes(w http.ResponseWriter, r *http.Request, serverID string) {
	query := r.URL.Query()
	req := &dto.GetRoutesRequest{
		ServerID: serverID,
		Method:   query.Get("method"),
		Limit:    h.parseIntParam(query.Get("limit"), 10),
		Offset:   h.parseIntParam(query.Get("offset"), 0),
	}

	if enabledStr := query.Get("enabled"); enabledStr != "" {
		enabled := enabledStr == "true"
		req.Enabled = &enabled
	}

	response, err := h.appService.GetRoutes(r.Context(), req)
	if err != nil {
		h.writeError(w, http.StatusInternalServerError, fmt.Sprintf("获取路由列表失败: %v", err))
		return
	}

	h.writeJSON(w, http.StatusOK, response)
}

// handleRegisterRoute 处理注册路由请求
func (h *HTTPHandler) handleRegisterRoute(w http.ResponseWriter, r *http.Request, serverID string) {
	var req dto.RegisterRouteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, fmt.Sprintf("解析请求失败: %v", err))
		return
	}

	req.ServerID = serverID

	response, err := h.appService.RegisterRoute(r.Context(), &req)
	if err != nil {
		h.writeError(w, http.StatusInternalServerError, fmt.Sprintf("注册路由失败: %v", err))
		return
	}

	h.writeJSON(w, http.StatusCreated, response)
}

// handleGetMiddlewares 处理获取中间件请求
func (h *HTTPHandler) handleGetMiddlewares(w http.ResponseWriter, r *http.Request, serverID string) {
	query := r.URL.Query()
	req := &dto.GetMiddlewaresRequest{
		ServerID: serverID,
		Type:     query.Get("type"),
		Limit:    h.parseIntParam(query.Get("limit"), 10),
		Offset:   h.parseIntParam(query.Get("offset"), 0),
	}

	if enabledStr := query.Get("enabled"); enabledStr != "" {
		enabled := enabledStr == "true"
		req.Enabled = &enabled
	}

	response, err := h.appService.GetMiddlewares(r.Context(), req)
	if err != nil {
		h.writeError(w, http.StatusInternalServerError, fmt.Sprintf("获取中间件列表失败: %v", err))
		return
	}

	h.writeJSON(w, http.StatusOK, response)
}

// handleAddMiddleware 处理添加中间件请求
func (h *HTTPHandler) handleAddMiddleware(w http.ResponseWriter, r *http.Request, serverID string) {
	var req dto.AddMiddlewareRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, fmt.Sprintf("解析请求失败: %v", err))
		return
	}

	req.ServerID = serverID

	response, err := h.appService.AddMiddleware(r.Context(), &req)
	if err != nil {
		h.writeError(w, http.StatusInternalServerError, fmt.Sprintf("添加中间件失败: %v", err))
		return
	}

	h.writeJSON(w, http.StatusCreated, response)
}

// 工具方法

// extractServerID 从URL路径中提取服务器ID
func (h *HTTPHandler) extractServerID(path string) string {
	// 例如: /api/servers/123/routes -> 123
	parts := strings.Split(path, "/")
	if len(parts) >= 4 && parts[1] == "api" && parts[2] == "servers" {
		return parts[3]
	}
	return ""
}

// parseIntParam 解析整数参数
func (h *HTTPHandler) parseIntParam(param string, defaultValue int) int {
	if param == "" {
		return defaultValue
	}

	if value, err := strconv.Atoi(param); err == nil {
		return value
	}

	return defaultValue
}

// writeJSON 写入JSON响应
func (h *HTTPHandler) writeJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "编码响应失败", http.StatusInternalServerError)
	}
}

// writeError 写入错误响应
func (h *HTTPHandler) writeError(w http.ResponseWriter, statusCode int, message string) {
	response := map[string]interface{}{
		"error":   true,
		"message": message,
		"code":    statusCode,
	}

	h.writeJSON(w, statusCode, response)
}
