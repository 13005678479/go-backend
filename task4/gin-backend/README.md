# Blog Backend

基于 Go + Gin + GORM 的个人博客系统后端

## 项目结构

```
blogV2/
├── cmd/                    # 应用程序入口
│   └── api/               # API 服务入口
│       └── main.go        # 主程序入口
├── config/                 # 配置文件
│   └── config.go          # 配置结构体
├── internal/               # 内部包（不对外暴露）
│   ├── controllers/       # 控制器层
│   │   ├── user_controller.go
│   │   ├── post_controller.go
│   │   └── comment_controller.go
│   ├── models/           # 数据模型
│   │   └── blog.go
│   └── services/         # 业务逻辑层
│       ├── user_service.go
│       └── blog_service.go
├── middleware/            # 中间件
│   └── auth.go           # 认证中间件
├── pkg/                  # 可复用的公共包
│   ├── database/         # 数据库相关
│   │   └── connection.go
│   └── utils/            # 工具函数
│       ├── logger.go
│       └── jwt.go
├── router/               # 路由配置
│   └── router.go
├── go.mod
├── go.sum
└── README.md
```

## 功能特性

- ✅ 用户注册/登录
- ✅ JWT 认证
- ✅ 文章 CRUD 操作
- ✅ 评论功能
- ✅ 权限控制
- ✅ 错误处理
- ✅ 日志记录

## 快速开始

### 环境要求

- Go 1.25+
- MySQL 5.7+

### 安装依赖

```bash
go mod tidy
```

### 配置环境变量

```bash
# 服务器配置
export SERVER_PORT=8080
export SERVER_MODE=debug

# 数据库配置
export DB_HOST=localhost
export DB_PORT=3306
export DB_USER=root
export DB_PASSWORD=your_password
export DB_NAME=blog

# JWT配置
export JWT_SECRET=your-secret-key
export JWT_EXPIRE=24
```

### 运行应用

```bash
go run cmd/api/main.go
```

## API 文档

### 用户相关

- `POST /api/register` - 用户注册
- `POST /api/login` - 用户登录
- `GET /api/user/me` - 获取当前用户信息（需要认证）

### 文章相关

- `GET /api/posts` - 获取文章列表（公开）
- `GET /api/posts/:id` - 获取文章详情（公开）
- `POST /api/posts` - 创建文章（需要认证）
- `PUT /api/posts/:id` - 更新文章（需要认证）
- `DELETE /api/posts/:id` - 删除文章（需要认证）

### 评论相关

- `POST /api/posts/:id/comments` - 创建评论（需要认证）
- `DELETE /api/comments/:id` - 删除评论（需要认证）

## 开发说明

项目采用分层架构：

- **Controller**: 处理 HTTP 请求和响应
- **Service**: 业务逻辑处理
- **Model**: 数据模型定义
- **Middleware**: 中间件处理
- **Router**: 路由配置

遵循 Go 语言最佳实践，代码结构清晰，易于维护和扩展。