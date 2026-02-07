# EnterpriseHub

EnterpriseHub 是一个基于 **Go 语言** 构建的企业级分布式微服务电商系统。项目采用 **gRPC** 进行服务间通信，**Gin** 作为 HTTP 网关框架，结合 **MySQL**、**Redis**、**RabbitMQ** 等主流中间件，实现了用户、认证、订单、支付、通知等核心电商业务模块。

本项目致力于展示现代云原生架构下的高性能后端开发实践，特别是在**系统解耦**、**高并发处理**和**服务治理**方面的设计思路。

## 🛠 技术栈亮点 (Tech Stack)

本项目后端技术选型遵循**高性能**、**强类型**和**易维护**的原则：

### 核心框架
- **Go (Golang)**: 1.25+，利用其原生并发特性（Goroutines & Channels）处理高吞吐量请求。
- **Gin**: 轻量级、高性能的 HTTP Web 框架，作为 API 网关的底层支撑。
- **gRPC + Protobuf**: 高性能 RPC 框架，用于微服务内部的高效通信，拥有比 RESTful 更小的传输体积和更快的解析速度。

### 数据存储
- **MySQL**: 核心业务关系型数据库，使用 **GORM** 进行 ORM 映射，支持事务和复杂查询。
- **Redis**: 高性能 Key-Value 缓存，用于缓存热点数据（如用户信息、商品详情）和分布式锁实现。

### 消息队列与异步处理
- **RabbitMQ**: AMQP 协议消息队列。用于核心业务解耦（如：订单创建后异步通知、削峰填谷），确保系统在流量高峰下的稳定性。

### 服务治理与工具
- **JWT (JSON Web Tokens)**: 无状态认证机制，保障微服务间的安全鉴权。
- **Zap**: Uber 开源的高性能日志库，结合 **Lumberjack** 实现日志切割与归档。
- **Viper / Godotenv**: 配置管理，支持多环境配置热加载。
- **Docker**: 容器化部署支持（Makefile 脚本集成），简化开发与运维流程。

## 🏗 系统设计亮点 (System Design)

### 1. 微服务架构拆分 (Microservices Architecture)
系统按业务领域（DDD 思想）拆分为多个独立服务，降低了模块间的耦合度，各服务可独立部署与扩展：
- **Gateway Service**: 统一流量入口，负责路由分发、鉴权、限流与熔断。
- **User Service**: 用户身份管理与个人信息维护。
- **Auth Service**: 负责颁发与校验 JWT 令牌。
- **Order Service**: 订单全生命周期管理（创建、支付、取消、超时处理）。
- **Payment Service**: 支付网关对接与状态管理。
- **Notify Service**: 异步通知服务（邮件/短信），通过 MQ 消费事件。

### 2. 清洁架构 (Clean Architecture)
每个微服务内部均严格遵循分层架构原则，确保代码的可测试性与可维护性：
- **Handler Layer (Transport)**: 处理 HTTP/gRPC 请求与响应转换。
- **Service Layer (Business Logic)**: 纯粹的业务逻辑层，不依赖具体底层实现。
- **Repository Layer (Data Access)**: 负责数据库、缓存的数据持久化操作。

### 3. 高并发与性能优化
- **多级缓存策略**: 引入 Redis 缓存层，减轻 MySQL 压力。采用 "Cache Aside" 模式保证数据最终一致性。
- **异步解耦**: 关键链路（如“下单成功”）采用事件驱动模型。订单服务仅需发送消息到 RabbitMQ，通知服务异步消费，极大降低了接口响应时间（RT）。
- **连接池管理**: 数据库与 Redis 均配置了连接池，复用 TCP 连接，避免频繁握手带来的性能开销。

### 4. 统一网关设计
- 所有的外部请求必须通过 **Gateway-Service**。
- 网关集成了 **JWT Middleware**，在请求到达具体微服务前完成身份校验，实现了统一的认证拦截，不仅提高了安全性，也避免了每个微服务重复实现鉴权逻辑。

## 🚀 快速开始 (Quick Start)

### 依赖环境
- Go 1.25+
- MySQL 8.0+
- Redis 6.0+
- RabbitMQ 3.8+

### 运行服务
项目包含 `Makefile` 便于快速启动各个微服务：

```bash
# 启动网关
make gateway

# 启动用户服务
make user_service

# 启动订单服务
make order_service

# ... 其他服务参考 Makefile
```
