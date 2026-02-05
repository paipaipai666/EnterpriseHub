# EnterpriseHub API 文档

本文档按服务分类整理了系统中所有的 API 接口，包括 REST API 和 gRPC API。

---

## 目录

- [User Service](#user-service)
- [Auth Service](#auth-service)
- [Order Service](#order-service)
- [Payment Service](#payment-service)
- [Gateway Service](#gateway-service)

---

## User Service

用户服务，提供用户注册、查询等功能的 REST 和 gRPC 接口。

**服务端口:**
- HTTP: `8000`
- gRPC: `8001`

### REST API

| 方法 | 路径 | 描述 | 请求参数 |
|------|------|------|----------|
| POST | `/api/v1/users/register` | 用户注册 | JSON: `{username, password}` |
| GET | `/api/v1/users/get/:id` | 根据ID获取用户 | 路径参数: `id` |
| GET | `/api/v1/users/get_all` | 获取所有用户 | 无 |

**请求示例:**

```bash
# 用户注册
POST /api/v1/users/register
Content-Type: application/json

{
  "username": "testuser",
  "password": "password123"
}

# 根据ID获取用户
GET /api/v1/users/get/123

# 获取所有用户
GET /api/v1/users/get_all
```

**响应格式:**

```json
{
  "message": "success",
  "data": {}
}
```

### gRPC API

**服务定义:** `UserService`

| 方法 | 请求类型 | 响应类型 | 描述 |
|------|----------|----------|------|
| GetUserById | GetUserByIdRequest | UserResponse | 根据用户ID获取用户信息 |
| GetUserByUsername | GetUserByUsernameRequest | UserResponse | 根据用户名获取用户信息 |

**消息定义:**

```protobuf
// 根据用户ID查询请求
message GetUserByIdRequest {
  string id = 1;
}

// 根据用户名查询请求
message GetUserByUsernameRequest {
  string username = 1;
}

// 用户响应消息
message UserResponse {
  string id = 1;
  string username = 2;
  string password = 3; // 仅内部使用
}
```

---

## Auth Service

认证服务，提供用户登录、登出功能的 REST 接口。

**服务端口:** `9000`

### REST API

| 方法 | 路径 | 描述 | 请求参数 |
|------|------|------|----------|
| POST | `/api/v1/auth/login` | JWT登录 | JSON: `{username, password}` |
| DELETE | `/api/v1/auth/logout` | 退出登录 | 无 |

**请求示例:**

```bash
# JWT登录
POST /api/v1/auth/login
Content-Type: application/json

{
  "username": "testuser",
  "password": "password123"
}

# 退出登录
DELETE /api/v1/auth/logout
```

**响应示例:**

```json
// 登录成功
{
  "message": "success",
  "data": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}

// 退出成功
{
  "message": "success",
  "data": "你已安全退出"
}
```

### gRPC API

Auth Service 本身不提供 gRPC 服务，它通过 gRPC 客户端调用 User Service 的 gRPC 接口来完成用户认证。

---

## Order Service

订单服务，提供订单创建、查询、支付、取消等功能的 REST 接口。

**服务端口:** `10000`

### REST API

| 方法 | 路径 | 描述 | 请求参数 |
|------|------|------|----------|
| POST | `/api/v1/order/create` | 创建订单 | JSON: `{userId, amount}` |
| GET | `/api/v1/order/get/:id` | 根据ID获取订单 | 路径参数: `id` |
| GET | `/api/v1/order/list/:user_id` | 获取用户订单列表 | 路径参数: `user_id` |
| PUT | `/api/v1/order/cancel/:id` | 取消订单 | 路径参数: `id` |
| POST | `/api/v1/order/pay/:id` | 支付订单 | 路径参数: `id`, JSON: `{method}` |

**请求示例:**

```bash
# 创建订单
POST /api/v1/order/create
Content-Type: application/json

{
  "userId": "user123",
  "amount": 99.00
}

# 根据ID获取订单
GET /api/v1/order/get/ord_123

# 获取用户订单列表
GET /api/v1/order/list/user123

# 取消订单
PUT /api/v1/order/cancel/ord_123

# 支付订单
POST /api/v1/order/pay/ord_123
Content-Type: application/json

{
  "method": 1
}
```

**支付方式 (PaymentMethod):**

| 值 | 描述 |
|----|------|
| 0 | 未知 |
| 1 | 支付宝 |
| 2 | 微信支付 |
| 3 | 银行卡 |

### gRPC API

Order Service 本身不提供 gRPC 服务，它通过 gRPC 客户端调用 Payment Service 的 gRPC 接口来完成支付功能。

---

## Payment Service

支付服务，提供支付订单创建、执行支付、查询支付状态等功能的 gRPC 接口。

**服务端口:** `11001`

### REST API

Payment Service 不提供 REST API，仅通过 gRPC 接口对外提供服务。

### gRPC API

**服务定义:** `PaymentService`

| 方法 | 请求类型 | 响应类型 | 描述 |
|------|----------|----------|------|
| CreatePayment | CreatePaymentRequest | CreatePaymentResponse | 创建支付订单 |
| Pay | PayRequest | PayResponse | 执行支付操作 |
| QueryPayment | QueryPaymentRequest | QueryPaymentResponse | 查询支付状态 |

**消息定义:**

```protobuf
// 创建支付请求
message CreatePaymentRequest {
  string order_id = 1;    // 业务订单号
  double amount = 2;      // 支付金额（单位：分）
  PaymentMethod method = 3; // 支付方式
  string currency = 4;     // 币种（如 CNY）
  string user_id = 5;      // 支付用户
}

// 创建支付响应
message CreatePaymentResponse {
  string payment_id = 1;  // 支付单号
  PaymentStatus status = 2; // 初始状态（一般是 PENDING）
  int64 created_at = 3;    // 创建时间（时间戳）
}

// 支付请求
message PayRequest {
  string payment_id = 1;   // 支付单号
  string request_id = 2;   // 幂等ID（防重复支付）
}

// 支付响应
message PayResponse {
  PaymentStatus status = 1;        // 支付结果
  string transaction_id = 2;      // 第三方交易号（可为空）
  int64 paid_at = 3;              // 支付完成时间
}

// 查询支付请求
message QueryPaymentRequest {
  string payment_id = 1;
}

// 查询支付响应
message QueryPaymentResponse {
  string payment_id = 1;
  string order_id = 2;
  string user_id = 3;
  double amount = 4;
  string currency = 5;
  PaymentMethod method = 6;
  PaymentStatus status = 7;
  int64 created_at = 8;
  int64 paid_at = 9;
}

// 支付方式枚举
enum PaymentMethod {
  UNKNOWN = 0;
  ALIPAY = 1;
  WECHAT = 2;
  BANK_CARD = 3;
}

// 支付状态枚举
enum PaymentStatus {
  PAYMENT_STATUS_UNSPECIFIED = 0;
  PAYMENT_STATUS_PENDING = 1;
  PAYMENT_STATUS_SUCCESS = 2;
  PAYMENT_STATUS_FAILED = 3;
}
```

---

## Gateway Service

网关服务，提供统一的 API 入口，对所有服务进行路由分发。

**服务端口:** `8080`

### 路由规则

| 路径模式 | 目标服务 | 目标端口 |
|----------|----------|----------|
| `/api/v1/users/*path` | User Service | 8000 |
| `/api/v1/auth/*path` | Auth Service | 9000 |
| `/api/v1/order/*path` | Order Service | 10000 |

**注意:** Gateway Service 本身不提供业务 API，仅作为统一的 API 网关入口。所有请求需要通过 JWT 中间件认证。

---

## 通用说明

### 响应格式

所有 REST API 返回统一的 JSON 格式：

```json
{
  "message": "success",
  "data": {}
}
```

| 字段 | 类型 | 描述 |
|------|------|------|
| message | string | 状态消息: "success" 或 "failed" |
| data | any | 响应数据 |

### 认证

Gateway Service 使用 JWT 进行认证。请求需要在 Header 中携带 Token：

```bash
Authorization: Bearer <your_jwt_token>
```

### 服务间通信

- **User Service** ← **Auth Service**: gRPC 调用
- **User Service** ← **Order Service**: gRPC 调用
- **Order Service** ← **Payment Service**: gRPC 调用

---

## 端口汇总

| 服务 | HTTP 端口 | gRPC 端口 |
|------|----------|----------|
| User Service | 8000 | 8001 |
| Auth Service | 9000 | - |
| Order Service | 10000 | - |
| Payment Service | - | 11001 |
| Gateway Service | 8080 | - |
