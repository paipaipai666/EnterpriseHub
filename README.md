# EnterpriseHub —— Go 企业级微服务项目

## 一、项目简介

**EnterpriseHub** 是一个基于 **Golang** 构建的企业级微服务项目示例，采用当前主流的微服务架构设计思想，覆盖用户管理、认证授权、订单处理、支付模拟、异步通知等典型业务场景。

本项目的目标并非功能堆砌，而是**完整还原企业真实微服务项目的工程结构、服务拆分方式与治理思路**，适合作为：

- Go 后端 / 微服务学习项目
- 实习 / 校招简历项目
- 微服务架构实践参考

------

## 二、整体架构设计

### 1. 架构概览

项目采用 **API Gateway + 多业务微服务** 的经典架构：

- 对外统一通过 Gateway 提供 HTTP 接口
- 服务内部使用 gRPC 通信
- 通过消息队列实现异步解耦
- 一服务一数据库，避免强耦合

### 2. 服务划分

| 服务名          | 职责说明                       |
| --------------- | ------------------------------ |
| gateway-service | API 网关、鉴权、路由转发、限流 |
| user-service    | 用户注册、用户信息管理         |
| auth-service    | 登录认证、JWT 生成与校验       |
| order-service   | 订单创建、状态流转             |
| payment-service | 支付模拟、支付回调             |
| notify-service  | 消息消费、通知处理             |

------

## 三、技术栈选型

### 1. 核心技术

- **语言**：Go 1.22
- **Web 框架**：Gin
- **RPC 通信**：gRPC + Protobuf
- **数据库**：MySQL 8
- **ORM**：GORM
- **缓存**：Redis
- **消息队列**：RabbitMQ
- **服务发现**：Consul
- **日志**：Zap
- **容器化**：Docker & Docker Compose

### 2. 通信方式

- 外部通信：HTTP / JSON
- 内部通信：gRPC
- 异步通信：RabbitMQ

------

## 四、项目目录结构

```text
enterprisehub/
├── gateway-service/     # API 网关服务
├── user-service/        # 用户服务
├── auth-service/        # 认证服务
├── order-service/       # 订单服务
├── payment-service/     # 支付服务
├── notify-service/      # 通知服务
├── proto/               # gRPC 协议定义
├── deploy/              # Docker & 中间件配置
│   └── docker-compose.yml
├── scripts/             # 辅助脚本
└── README.md
```

### 单个服务内部结构（示例：user-service）

```text
user-service/
├── cmd/main.go           # 程序入口
├── internal/
│   ├── api/              # HTTP Handler
│   ├── service/          # 业务逻辑层
│   ├── repository/       # 数据访问层
│   ├── model/            # 数据模型
│   ├── middleware/       # 中间件
│   └── config/           # 配置加载
├── proto/                # gRPC 接口
├── go.mod
```

------

## 五、核心业务流程

### 1. 用户注册与登录

1. 客户端通过 Gateway 调用注册接口
2. Gateway 转发请求至 user-service
3. 登录时由 auth-service 调用 user-service 校验用户
4. auth-service 生成 JWT 返回给客户端

### 2. 下单与支付流程

1. 客户端携带 JWT 调用下单接口
2. order-service 创建订单（状态：CREATED）
3. 调用 payment-service 进行支付模拟
4. 支付完成后更新订单状态（PAID）
5. 发送消息至 MQ
6. notify-service 消费消息并处理通知逻辑

------

## 六、设计原则与工程实践

- **一服务一数据库**，避免跨服务数据强依赖
- **接口隔离**，通过 gRPC 明确服务边界
- **业务与基础设施解耦**（HTTP / RPC / MQ 分层）
- **统一错误码与日志规范**
- **面向接口编程，便于扩展与测试**

------

## 七、运行方式（开发环境）

### 1. 启动基础中间件

```bash
docker-compose up -d
```

### 2. 启动各微服务

```bash
cd user-service
go run cmd/main.go
```

其他服务启动方式相同。

------

## 八、项目亮点（简历可用）

- 基于 **Go + Gin + gRPC** 构建多微服务系统
- 实现 API Gateway + JWT 鉴权架构
- 使用 RabbitMQ 实现服务解耦与异步处理
- 采用 Docker Compose 实现多服务本地编排
- 具备清晰的服务边界与企业级工程结构

------

## 九、后续可扩展方向

- 引入 Kubernetes 进行容器编排
- 接入 Prometheus + Grafana 实现监控
- 引入 OpenTelemetry 实现链路追踪
- 增加 CI/CD 自动化部署流程

------

## 十、作者说明

本项目为个人学习与实践项目，重点在于**工程能力、架构理解与微服务设计思路**，而非业务复杂度。适合对 Go 微服务与企业级后端开发感兴趣的学习者参考与扩展。
