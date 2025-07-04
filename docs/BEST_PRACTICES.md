# Ant Web框架最佳实践指南

本文档提供使用Ant Web框架开发应用的最佳实践建议，帮助您构建高质量、可维护的Web应用。

## 目录

- [项目结构组织](#项目结构组织)
- [DDD架构实践](#ddd架构实践)
- [错误处理](#错误处理)
- [性能优化](#性能优化)
- [安全性考虑](#安全性考虑)
- [测试策略](#测试策略)
- [代码规范](#代码规范)

## 项目结构组织

### 推荐的项目结构

```
your-project/
├── cmd/                           # 应用程序入口
│   └── server/
│       └── main.go
├── internal/                      # 内部代码（不对外暴露）
│   ├── domain/                    # 领域层
│   │   ├── user/                  # 用户领域
│   │   │   ├── entity.go          # 实体定义
│   │   │   ├── repository.go      # 仓储接口
│   │   │   ├── service.go         # 领域服务
│   │   │   └── value_object.go    # 值对象
│   │   └── shared/                # 共享领域概念
│   │       ├── events/            # 领域事件
│   │       └── errors/            # 领域错误
│   ├── application/               # 应用层
│   │   ├── commands/              # 命令对象
│   │   ├── queries/               # 查询对象
│   │   ├── dto/                   # 数据传输对象
│   │   ├── handlers/              # 命令/查询处理器
│   │   └── services/              # 应用服务
│   ├── infrastructure/            # 基础设施层
│   │   ├── database/              # 数据库相关
│   │   ├── cache/                 # 缓存实现
│   │   ├── external/              # 外部服务
│   │   └── repositories/          # 仓储实现
│   └── interfaces/                # 接口层
│       ├── http/                  # HTTP接口
│       │   ├── handlers/          # HTTP处理器
│       │   ├── middleware/        # 中间件
│       │   └── routes/            # 路由定义
│       └── grpc/                  # gRPC接口（可选）
├── pkg/                           # 公共库（可对外暴露）
│   ├── logger/                    # 日志工具
│   ├── config/                    # 配置管理
│   └── utils/                     # 工具函数
├── configs/                       # 配置文件
├── scripts/                       # 脚本文件
├── docs/                          # 文档
├── tests/                         # 测试文件
│   ├── integration/               # 集成测试
│   └── e2e/                       # 端到端测试
├── web/                           # 前端资源（可选）
│   ├── static/                    # 静态文件
│   └── templates/                 # 模板文件
├── go.mod
├── go.sum
├── Makefile                       # 构建脚本
├── Dockerfile                     # Docker配置
└── README.md
```

### 包命名规范

```go
// ✅ 好的包命名
package user          // 简洁明了
package repository    // 描述性强
package middleware    // 功能清晰

// ❌ 避免的包命名
package userService   // 避免驼峰命名
package utils         // 过于通用
package common        // 意义不明确
```

## DDD架构实践

### 领域层设计

#### 实体设计

```go
// ✅ 好的实体设计
package user

import (
    "time"
    "errors"
)

// User 用户实体
type User struct {
    id       UserID
    email    Email
    name     string
    status   UserStatus
    createdAt time.Time
    updatedAt time.Time
}

// NewUser 创建新用户（工厂方法）
func NewUser(email Email, name string) (*User, error) {
    if name == "" {
        return nil, errors.New("用户名不能为空")
    }
    
    return &User{
        id:        NewUserID(),
        email:     email,
        name:      name,
        status:    UserStatusActive,
        createdAt: time.Now(),
        updatedAt: time.Now(),
    }, nil
}

// ChangeEmail 修改邮箱（业务方法）
func (u *User) ChangeEmail(newEmail Email) error {
    if u.status != UserStatusActive {
        return errors.New("非活跃用户不能修改邮箱")
    }
    
    u.email = newEmail
    u.updatedAt = time.Now()
    return nil
}

// ID 获取用户ID
func (u *User) ID() UserID {
    return u.id
}
```

#### 值对象设计

```go
// ✅ 好的值对象设计
package user

import (
    "errors"
    "regexp"
)

// Email 邮箱值对象
type Email struct {
    value string
}

// NewEmail 创建邮箱值对象
func NewEmail(email string) (Email, error) {
    if email == "" {
        return Email{}, errors.New("邮箱不能为空")
    }
    
    emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
    if !emailRegex.MatchString(email) {
        return Email{}, errors.New("邮箱格式不正确")
    }
    
    return Email{value: email}, nil
}

// String 获取邮箱字符串值
func (e Email) String() string {
    return e.value
}

// Equals 比较两个邮箱是否相等
func (e Email) Equals(other Email) bool {
    return e.value == other.value
}
```

#### 仓储接口设计

```go
// ✅ 好的仓储接口设计
package user

import "context"

// Repository 用户仓储接口
type Repository interface {
    // Save 保存用户
    Save(ctx context.Context, user *User) error
    
    // FindByID 根据ID查找用户
    FindByID(ctx context.Context, id UserID) (*User, error)
    
    // FindByEmail 根据邮箱查找用户
    FindByEmail(ctx context.Context, email Email) (*User, error)
    
    // Delete 删除用户
    Delete(ctx context.Context, id UserID) error
    
    // List 列出用户（带分页）
    List(ctx context.Context, offset, limit int) ([]*User, error)
}
```

### 应用层设计

#### 命令和查询分离

```go
// ✅ 命令对象设计
package commands

// CreateUserCommand 创建用户命令
type CreateUserCommand struct {
    Email string `json:"email" validate:"required,email"`
    Name  string `json:"name" validate:"required,min=2,max=50"`
}

// UpdateUserCommand 更新用户命令
type UpdateUserCommand struct {
    ID    string `json:"id" validate:"required"`
    Name  string `json:"name" validate:"required,min=2,max=50"`
    Email string `json:"email" validate:"required,email"`
}
```

```go
// ✅ 查询对象设计
package queries

// GetUserQuery 获取用户查询
type GetUserQuery struct {
    ID string `json:"id" validate:"required"`
}

// ListUsersQuery 列出用户查询
type ListUsersQuery struct {
    Page     int    `json:"page" validate:"min=1"`
    PageSize int    `json:"page_size" validate:"min=1,max=100"`
    Email    string `json:"email,omitempty"`
    Status   string `json:"status,omitempty"`
}
```

#### 应用服务设计

```go
// ✅ 好的应用服务设计
package services

import (
    "context"
    "your-project/internal/domain/user"
    "your-project/internal/application/commands"
    "your-project/internal/application/dto"
)

// UserService 用户应用服务
type UserService struct {
    userRepo user.Repository
    eventBus EventBus
}

// NewUserService 创建用户服务
func NewUserService(userRepo user.Repository, eventBus EventBus) *UserService {
    return &UserService{
        userRepo: userRepo,
        eventBus: eventBus,
    }
}

// CreateUser 创建用户
func (s *UserService) CreateUser(ctx context.Context, cmd commands.CreateUserCommand) (*dto.UserResponse, error) {
    // 1. 验证命令
    if err := s.validateCommand(cmd); err != nil {
        return nil, err
    }
    
    // 2. 创建值对象
    email, err := user.NewEmail(cmd.Email)
    if err != nil {
        return nil, err
    }
    
    // 3. 检查邮箱是否已存在
    existingUser, _ := s.userRepo.FindByEmail(ctx, email)
    if existingUser != nil {
        return nil, errors.New("邮箱已存在")
    }
    
    // 4. 创建用户实体
    newUser, err := user.NewUser(email, cmd.Name)
    if err != nil {
        return nil, err
    }
    
    // 5. 保存用户
    if err := s.userRepo.Save(ctx, newUser); err != nil {
        return nil, err
    }
    
    // 6. 发布领域事件
    event := user.NewUserCreatedEvent(newUser.ID(), newUser.Email())
    s.eventBus.Publish(ctx, event)
    
    // 7. 返回响应
    return &dto.UserResponse{
        ID:    newUser.ID().String(),
        Email: newUser.Email().String(),
        Name:  newUser.Name(),
    }, nil
}
```

## 错误处理

### 错误类型设计

```go
// ✅ 好的错误设计
package errors

import "fmt"

// ErrorType 错误类型
type ErrorType string

const (
    ErrorTypeValidation   ErrorType = "VALIDATION"
    ErrorTypeNotFound     ErrorType = "NOT_FOUND"
    ErrorTypeConflict     ErrorType = "CONFLICT"
    ErrorTypeInternal     ErrorType = "INTERNAL"
    ErrorTypeUnauthorized ErrorType = "UNAUTHORIZED"
    ErrorTypeForbidden    ErrorType = "FORBIDDEN"
)

// AppError 应用错误
type AppError struct {
    Type    ErrorType `json:"type"`
    Message string    `json:"message"`
    Code    string    `json:"code,omitempty"`
    Details any       `json:"details,omitempty"`
}

func (e *AppError) Error() string {
    return e.Message
}

// NewValidationError 创建验证错误
func NewValidationError(message string, details any) *AppError {
    return &AppError{
        Type:    ErrorTypeValidation,
        Message: message,
        Details: details,
    }
}

// NewNotFoundError 创建未找到错误
func NewNotFoundError(resource string) *AppError {
    return &AppError{
        Type:    ErrorTypeNotFound,
        Message: fmt.Sprintf("%s not found", resource),
    }
}
```

### 错误处理中间件

```go
// ✅ 统一错误处理中间件
package middleware

import (
    "encoding/json"
    "net/http"
    "your-project/pkg/errors"
    "github.com/justinwongcn/ant"
)

// ErrorHandlerMiddleware 错误处理中间件
func ErrorHandlerMiddleware(next ant.HandleFunc) ant.HandleFunc {
    return func(ctx *ant.Context) {
        defer func() {
            if r := recover(); r != nil {
                handlePanic(ctx, r)
            }
        }()
        
        next(ctx)
        
        // 检查是否有错误需要处理
        if ctx.RespStatusCode >= 400 {
            handleError(ctx)
        }
    }
}

func handleError(ctx *ant.Context) {
    // 如果已经设置了响应数据，不再处理
    if len(ctx.RespData) > 0 {
        return
    }
    
    var response map[string]any
    
    switch ctx.RespStatusCode {
    case http.StatusNotFound:
        response = map[string]any{
            "error": map[string]any{
                "type":    "NOT_FOUND",
                "message": "资源未找到",
            },
        }
    case http.StatusUnauthorized:
        response = map[string]any{
            "error": map[string]any{
                "type":    "UNAUTHORIZED",
                "message": "未授权访问",
            },
        }
    case http.StatusForbidden:
        response = map[string]any{
            "error": map[string]any{
                "type":    "FORBIDDEN",
                "message": "禁止访问",
            },
        }
    default:
        response = map[string]any{
            "error": map[string]any{
                "type":    "INTERNAL",
                "message": "服务器内部错误",
            },
        }
    }
    
    ctx.Resp.Header().Set("Content-Type", "application/json")
    jsonData, _ := json.Marshal(response)
    ctx.RespData = jsonData
}

func handlePanic(ctx *ant.Context, r any) {
    ctx.RespStatusCode = http.StatusInternalServerError
    
    response := map[string]any{
        "error": map[string]any{
            "type":    "INTERNAL",
            "message": "服务器内部错误",
        },
    }
    
    ctx.Resp.Header().Set("Content-Type", "application/json")
    jsonData, _ := json.Marshal(response)
    ctx.RespData = jsonData
}
```

## 性能优化

### 数据库优化

```go
// ✅ 使用连接池
package database

import (
    "database/sql"
    "time"
)

func NewDB(dsn string) (*sql.DB, error) {
    db, err := sql.Open("mysql", dsn)
    if err != nil {
        return nil, err
    }
    
    // 设置连接池参数
    db.SetMaxOpenConns(25)                 // 最大打开连接数
    db.SetMaxIdleConns(5)                  // 最大空闲连接数
    db.SetConnMaxLifetime(5 * time.Minute) // 连接最大生存时间
    db.SetConnMaxIdleTime(1 * time.Minute) // 连接最大空闲时间
    
    return db, nil
}
```

### 缓存策略

```go
// ✅ 缓存装饰器模式
package repository

import (
    "context"
    "encoding/json"
    "fmt"
    "time"
    "your-project/internal/domain/user"
)

// CachedUserRepository 带缓存的用户仓储
type CachedUserRepository struct {
    repo  user.Repository
    cache Cache
    ttl   time.Duration
}

func NewCachedUserRepository(repo user.Repository, cache Cache) *CachedUserRepository {
    return &CachedUserRepository{
        repo:  repo,
        cache: cache,
        ttl:   10 * time.Minute,
    }
}

func (r *CachedUserRepository) FindByID(ctx context.Context, id user.UserID) (*user.User, error) {
    // 尝试从缓存获取
    cacheKey := fmt.Sprintf("user:%s", id.String())
    if data, err := r.cache.Get(ctx, cacheKey); err == nil {
        var u user.User
        if err := json.Unmarshal(data, &u); err == nil {
            return &u, nil
        }
    }
    
    // 从数据库获取
    u, err := r.repo.FindByID(ctx, id)
    if err != nil {
        return nil, err
    }
    
    // 存入缓存
    if data, err := json.Marshal(u); err == nil {
        r.cache.Set(ctx, cacheKey, data, r.ttl)
    }
    
    return u, nil
}
```

### 并发处理

```go
// ✅ 使用工作池模式
package worker

import (
    "context"
    "sync"
)

// WorkerPool 工作池
type WorkerPool struct {
    workerCount int
    jobQueue    chan Job
    wg          sync.WaitGroup
}

// Job 工作任务
type Job interface {
    Execute(ctx context.Context) error
}

// NewWorkerPool 创建工作池
func NewWorkerPool(workerCount, queueSize int) *WorkerPool {
    return &WorkerPool{
        workerCount: workerCount,
        jobQueue:    make(chan Job, queueSize),
    }
}

// Start 启动工作池
func (p *WorkerPool) Start(ctx context.Context) {
    for i := 0; i < p.workerCount; i++ {
        p.wg.Add(1)
        go p.worker(ctx)
    }
}

// Submit 提交任务
func (p *WorkerPool) Submit(job Job) {
    p.jobQueue <- job
}

// Stop 停止工作池
func (p *WorkerPool) Stop() {
    close(p.jobQueue)
    p.wg.Wait()
}

func (p *WorkerPool) worker(ctx context.Context) {
    defer p.wg.Done()
    
    for {
        select {
        case job, ok := <-p.jobQueue:
            if !ok {
                return
            }
            job.Execute(ctx)
        case <-ctx.Done():
            return
        }
    }
}
```

## 安全性考虑

### 输入验证

```go
// ✅ 使用验证库
package validation

import (
    "github.com/go-playground/validator/v10"
    "your-project/pkg/errors"
)

var validate = validator.New()

// ValidateStruct 验证结构体
func ValidateStruct(s interface{}) error {
    if err := validate.Struct(s); err != nil {
        var validationErrors []string
        for _, err := range err.(validator.ValidationErrors) {
            validationErrors = append(validationErrors,
                fmt.Sprintf("字段 %s 验证失败: %s", err.Field(), err.Tag()))
        }
        return errors.NewValidationError("输入验证失败", validationErrors)
    }
    return nil
}

// 自定义验证规则
func init() {
    validate.RegisterValidation("phone", validatePhone)
}

func validatePhone(fl validator.FieldLevel) bool {
    phone := fl.Field().String()
    phoneRegex := regexp.MustCompile(`^1[3-9]\d{9}$`)
    return phoneRegex.MatchString(phone)
}
```

### 认证和授权

```go
// ✅ JWT认证实现
package auth

import (
    "time"
    "github.com/golang-jwt/jwt/v5"
)

// JWTService JWT服务
type JWTService struct {
    secretKey []byte
    issuer    string
}

// Claims JWT声明
type Claims struct {
    UserID   string   `json:"user_id"`
    Username string   `json:"username"`
    Roles    []string `json:"roles"`
    jwt.RegisteredClaims
}

// GenerateToken 生成JWT令牌
func (j *JWTService) GenerateToken(userID, username string, roles []string) (string, error) {
    claims := Claims{
        UserID:   userID,
        Username: username,
        Roles:    roles,
        RegisteredClaims: jwt.RegisteredClaims{
            Issuer:    j.issuer,
            Subject:   userID,
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
            NotBefore: jwt.NewNumericDate(time.Now()),
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(j.secretKey)
}

// ValidateToken 验证JWT令牌
func (j *JWTService) ValidateToken(tokenString string) (*Claims, error) {
    token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
        return j.secretKey, nil
    })

    if err != nil {
        return nil, err
    }

    if claims, ok := token.Claims.(*Claims); ok && token.Valid {
        return claims, nil
    }

    return nil, jwt.ErrInvalidKey
}
```

### 安全中间件

```go
// ✅ 安全头中间件
package middleware

import (
    "github.com/justinwongcn/ant"
)

// SecurityHeadersMiddleware 安全头中间件
func SecurityHeadersMiddleware(next ant.HandleFunc) ant.HandleFunc {
    return func(ctx *ant.Context) {
        // 防止XSS攻击
        ctx.Resp.Header().Set("X-Content-Type-Options", "nosniff")
        ctx.Resp.Header().Set("X-Frame-Options", "DENY")
        ctx.Resp.Header().Set("X-XSS-Protection", "1; mode=block")

        // HTTPS相关
        ctx.Resp.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")

        // 内容安全策略
        ctx.Resp.Header().Set("Content-Security-Policy", "default-src 'self'")

        // 引用策略
        ctx.Resp.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")

        next(ctx)
    }
}

// RateLimitMiddleware 限流中间件
func RateLimitMiddleware(limiter RateLimiter) ant.Middleware {
    return func(next ant.HandleFunc) ant.HandleFunc {
        return func(ctx *ant.Context) {
            clientIP := getClientIP(ctx.Req)

            if !limiter.Allow(clientIP) {
                ctx.RespStatusCode = http.StatusTooManyRequests
                ctx.RespJSON(http.StatusTooManyRequests, map[string]string{
                    "error": "请求过于频繁",
                })
                return
            }

            next(ctx)
        }
    }
}
```

### 数据加密

```go
// ✅ 密码加密
package crypto

import (
    "crypto/rand"
    "crypto/subtle"
    "encoding/base64"
    "fmt"
    "golang.org/x/crypto/argon2"
)

// PasswordHasher 密码哈希器
type PasswordHasher struct {
    time    uint32
    memory  uint32
    threads uint8
    keyLen  uint32
}

// NewPasswordHasher 创建密码哈希器
func NewPasswordHasher() *PasswordHasher {
    return &PasswordHasher{
        time:    1,
        memory:  64 * 1024,
        threads: 4,
        keyLen:  32,
    }
}

// HashPassword 哈希密码
func (p *PasswordHasher) HashPassword(password string) (string, error) {
    salt := make([]byte, 16)
    if _, err := rand.Read(salt); err != nil {
        return "", err
    }

    hash := argon2.IDKey([]byte(password), salt, p.time, p.memory, p.threads, p.keyLen)

    b64Salt := base64.RawStdEncoding.EncodeToString(salt)
    b64Hash := base64.RawStdEncoding.EncodeToString(hash)

    encodedHash := fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
        argon2.Version, p.memory, p.time, p.threads, b64Salt, b64Hash)

    return encodedHash, nil
}

// VerifyPassword 验证密码
func (p *PasswordHasher) VerifyPassword(password, encodedHash string) (bool, error) {
    // 解析哈希值（简化实现）
    // 实际应用中需要完整的解析逻辑

    // 这里假设已经解析出salt和hash
    salt, hash, err := parseHash(encodedHash)
    if err != nil {
        return false, err
    }

    otherHash := argon2.IDKey([]byte(password), salt, p.time, p.memory, p.threads, p.keyLen)

    return subtle.ConstantTimeCompare(hash, otherHash) == 1, nil
}
```

## 测试策略

### 单元测试

```go
// ✅ 表格驱动测试
package user_test

import (
    "testing"
    "your-project/internal/domain/user"
)

func TestNewEmail(t *testing.T) {
    tests := []struct {
        name    string
        email   string
        wantErr bool
        errMsg  string
    }{
        {
            name:    "有效邮箱",
            email:   "test@example.com",
            wantErr: false,
        },
        {
            name:    "空邮箱",
            email:   "",
            wantErr: true,
            errMsg:  "邮箱不能为空",
        },
        {
            name:    "无效邮箱格式",
            email:   "invalid-email",
            wantErr: true,
            errMsg:  "邮箱格式不正确",
        },
        {
            name:    "缺少@符号",
            email:   "testexample.com",
            wantErr: true,
            errMsg:  "邮箱格式不正确",
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            email, err := user.NewEmail(tt.email)

            if tt.wantErr {
                if err == nil {
                    t.Errorf("NewEmail() 期望错误，但没有返回错误")
                    return
                }
                if err.Error() != tt.errMsg {
                    t.Errorf("NewEmail() 错误消息 = %v, 期望 %v", err.Error(), tt.errMsg)
                }
                return
            }

            if err != nil {
                t.Errorf("NewEmail() 意外错误 = %v", err)
                return
            }

            if email.String() != tt.email {
                t.Errorf("NewEmail() = %v, 期望 %v", email.String(), tt.email)
            }
        })
    }
}
```

### 集成测试

```go
// ✅ 集成测试示例
package integration_test

import (
    "bytes"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"
    "your-project/internal/interfaces/http/handlers"
    "github.com/justinwongcn/ant"
)

func TestUserHandler_CreateUser(t *testing.T) {
    // 设置测试环境
    server := ant.NewHTTPServer()
    handler := handlers.NewUserHandler(mockUserService)

    server.Handle("POST /users", handler.CreateUser)

    tests := []struct {
        name           string
        requestBody    map[string]interface{}
        expectedStatus int
        expectedBody   map[string]interface{}
    }{
        {
            name: "创建用户成功",
            requestBody: map[string]interface{}{
                "email": "test@example.com",
                "name":  "测试用户",
            },
            expectedStatus: http.StatusCreated,
            expectedBody: map[string]interface{}{
                "id":    "123",
                "email": "test@example.com",
                "name":  "测试用户",
            },
        },
        {
            name: "邮箱格式错误",
            requestBody: map[string]interface{}{
                "email": "invalid-email",
                "name":  "测试用户",
            },
            expectedStatus: http.StatusBadRequest,
            expectedBody: map[string]interface{}{
                "error": map[string]interface{}{
                    "type":    "VALIDATION",
                    "message": "输入验证失败",
                },
            },
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // 准备请求
            body, _ := json.Marshal(tt.requestBody)
            req := httptest.NewRequest("POST", "/users", bytes.NewBuffer(body))
            req.Header.Set("Content-Type", "application/json")

            // 执行请求
            w := httptest.NewRecorder()
            server.ServeHTTP(w, req)

            // 验证响应
            if w.Code != tt.expectedStatus {
                t.Errorf("状态码 = %v, 期望 %v", w.Code, tt.expectedStatus)
            }

            var response map[string]interface{}
            json.Unmarshal(w.Body.Bytes(), &response)

            // 验证响应体（简化比较）
            if tt.expectedBody != nil {
                // 这里应该进行更详细的比较
                // 可以使用testify/assert库简化断言
            }
        })
    }
}
```

## 代码规范

### 命名规范

```go
// ✅ 好的命名
var userCount int                    // 变量使用驼峰命名
const MaxRetryCount = 3             // 常量使用大写字母开头
type UserService struct{}           // 类型使用大写字母开头
func (s *UserService) CreateUser()  // 方法使用大写字母开头（公开）
func (s *UserService) validateUser() // 方法使用小写字母开头（私有）

// ❌ 避免的命名
var user_count int                  // 避免下划线
const maxRetryCount = 3            // 常量应该大写开头
type userService struct{}          // 公开类型应该大写开头
```

### 注释规范

```go
// ✅ 好的注释
// UserService 提供用户相关的业务逻辑处理
// 包括用户的创建、更新、删除和查询功能
type UserService struct {
    repo user.Repository
}

// CreateUser 创建新用户
// 参数:
//   - ctx: 上下文对象，用于控制请求生命周期
//   - cmd: 创建用户命令，包含用户基本信息
// 返回值:
//   - *dto.UserResponse: 创建成功的用户信息
//   - error: 创建过程中的错误
func (s *UserService) CreateUser(ctx context.Context, cmd commands.CreateUserCommand) (*dto.UserResponse, error) {
    // 验证邮箱格式
    email, err := user.NewEmail(cmd.Email)
    if err != nil {
        return nil, fmt.Errorf("邮箱格式错误: %w", err)
    }

    // 检查邮箱是否已存在
    if existingUser, _ := s.repo.FindByEmail(ctx, email); existingUser != nil {
        return nil, errors.New("邮箱已存在")
    }

    // ... 其他逻辑
}
```

### 错误处理规范

```go
// ✅ 好的错误处理
func (s *UserService) GetUser(ctx context.Context, id string) (*dto.UserResponse, error) {
    userID, err := user.ParseUserID(id)
    if err != nil {
        return nil, fmt.Errorf("解析用户ID失败: %w", err)
    }

    u, err := s.repo.FindByID(ctx, userID)
    if err != nil {
        if errors.Is(err, user.ErrUserNotFound) {
            return nil, errors.NewNotFoundError("用户")
        }
        return nil, fmt.Errorf("查询用户失败: %w", err)
    }

    return &dto.UserResponse{
        ID:    u.ID().String(),
        Email: u.Email().String(),
        Name:  u.Name(),
    }, nil
}

// ❌ 避免的错误处理
func (s *UserService) GetUser(ctx context.Context, id string) (*dto.UserResponse, error) {
    u, err := s.repo.FindByID(ctx, id) // 没有验证ID格式
    if err != nil {
        return nil, err // 直接返回底层错误，没有包装
    }

    return &dto.UserResponse{
        ID:    u.ID,
        Email: u.Email,
        Name:  u.Name,
    }, nil
}
```

这些最佳实践将帮助您构建高质量、可维护的Ant Web应用。记住，好的代码不仅要能工作，还要易于理解、测试和维护。
```
