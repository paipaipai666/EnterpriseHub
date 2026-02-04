# EnterpriseHub API 文档

> EnterpriseHub 企业级微服务项目接口规范

---

## 一、接口通用规范

### 1.1 基础约定

| 项目 | 规范 |
|------|------|
| 接口风格 | RESTful |
| 数据格式 | JSON |
| 字符编码 | UTF-8 |
| 统一入口 | API Gateway (8080) |
| 认证方式 | JWT Bearer Token |

### 1.2 统一响应结构

所有接口返回统一格式：

```json
{
    "code": 0,
    "message": "success",
    "data": {}
}
```

| 字段 | 类型 | 说明 |
|------|------|------|
| code | int | 0 表示成功，非 0 表示错误 |
| message | string | 状态描述信息 |
| data | object | 响应数据 |

### 1.3 认证方式

除公开接口外，所有请求需携带 JWT Token：

```
Authorization: Bearer <token>
```

**公开接口（无需认证）**：
- `POST /api/v1/auth/login`
- `POST /api/v1/users/register`

---

## 二、用户服务 (user-service)

> 端口：8001 | 基础路径：`/api/v1/users`

### 2.1 用户注册

- **URL**: `POST /api/v1/users/register`
- **认证**: 否

**请求参数**：

```json
{
    "username": "testuser",
    "Password": "123456",
    "email": "test@example.com"
}
```

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| username | string | 是 | 用户名 |
| Password | string | 是 | 密码 |
| email | string | 是 | 邮箱 |

**响应示例** (200)：

```json
{
    "code": 0,
    "message": "success",
    "data": "user_id"
}
```

---

### 2.2 获取用户信息

- **URL**: `GET /api/v1/users/:id`
- **认证**: 是

**路径参数**：

| 参数 | 类型 | 说明 |
|------|------|------|
| id | string | 用户ID |

**响应示例** (200)：

```json
{
    "code": 0,
    "message": "success",
    "data": {
        "Id": "uuid-xxx",
        "Username": "testuser",
        "Password": "123456",
        "Email": "test@example.com"
    }
}
```

---

### 2.3 获取所有用户

- **URL**: `GET /api/v1/users`
- **认证**: 是

**响应示例** (200)：

```json
{
    "code": 0,
    "message": "success",
    "data": [
        {
            "Id": "uuid-xxx",
            "Username": "user1",
            "Password": "xxx",
            "Email": "user1@example.com"
        },
        {
            "Id": "uuid-yyy",
            "Username": "user2",
            "Password": "xxx",
            "Email": "user2@example.com"
        }
    ]
}
```

---

## 三、认证服务 (auth-service)

> 端口：8002 | 基础路径：`/api/v1/auth`

### 3.1 用户登录

- **URL**: `POST /api/v1/auth/login`
- **认证**: 否

**请求参数**：

```json
{
    "username": "testuser",
    "password": "123456"
}
```

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| username | string | 是 | 用户名 |
| password | string | 是 | 密码 |

**响应示例** (200)：

```json
{
    "code": 0,
    "message": "success",
    "data": {
        "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
    }
}
```

**错误响应** (401)：

```json
{
    "code": 401,
    "message": "failed",
    "data": "invalid credentials"
}
```

---

### 3.2 用户登出

- **URL**: `POST /api/v1/auth/logout`
- **认证**: 是

**响应示例** (200)：

```json
{
    "code": 0,
    "message": "success",
    "data": "你已安全退出"
}
```

---

### 3.3 验证 Token

- **URL**: `POST /api/v1/auth/verify`
- **认证**: 是

**响应示例** (200)：

```json
{
    "code": 0,
    "message": "success",
    "data": {
        "valid": true,
        "username": "testuser"
    }
}
```

---

## 四、订单服务 (order-service)

> 端口：10000 | 基础路径：`/api/v1/order`

### 4.1 创建订单

- **URL**: `POST /api/v1/order/create`
- **认证**: 是

**请求参数**：

```json
{
    "user_id": "uuid-xxx",
    "amount": 199.99
}
```

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| user_id | string | 是 | 用户ID |
| amount | float64 | 是 | 订单金额 |

**响应示例** (200)：

```json
{
    "code": 0,
    "message": "success",
    "data": "order-uuid-xxx"
}
```

**业务流程**：
1. 创建订单（状态：CREATED）
2. 发布 `order.created` 事件到 RabbitMQ
3. notify-service 消费事件并发送通知

---

### 4.2 获取订单

- **URL**: `GET /api/v1/order/get/:id`
- **认证**: 是

**路径参数**：

| 参数 | 类型 | 说明 |
|------|------|------|
| id | string | 订单ID |

**响应示例** (200)：

```json
{
    "code": 0,
    "message": "success",
    "data": {
        "Id": "order-uuid-xxx",
        "UserId": "user-uuid-xxx",
        "Amount": 199.99,
        "Status": "CREATED",
        "CreateAt": "2024-01-01T10:00:00Z",
        "UpdateAt": "2024-01-01T10:00:00Z"
    }
}
```

---

### 4.3 获取用户订单列表

- **URL**: `GET /api/v1/order/list/:user_id`
- **认证**: 是

**路径参数**：

| 参数 | 类型 | 说明 |
|------|------|------|
| user_id | string | 用户ID |

**响应示例** (200)：

```json
{
    "code": 0,
    "message": "success",
    "data": [
        {
            "Id": "order-uuid-1",
            "UserId": "user-uuid-xxx",
            "Amount": 199.99,
            "Status": "CREATED"
        },
        {
            "Id": "order-uuid-2",
            "UserId": "user-uuid-xxx",
            "Amount": 99.00,
            "Status": "PAID"
        }
    ]
}
```

---

### 4.4 取消订单

- **URL**: `PUT /api/v1/order/cancel/:id`
- **认证**: 是

**路径参数**：

| 参数 | 类型 | 说明 |
|------|------|------|
| id | string | 订单ID |

**响应示例** (200)：

```json
{
    "code": 0,
    "message": "success",
    "data": "订单 order-uuid-xxx 已取消！"
}
```

**业务规则**：
- 仅允许取消 `CREATED` 状态的订单

---

### 4.5 支付订单

- **URL**: `POST /api/v1/order/pay/:id`
- **认证**: 是

**路径参数**：

| 参数 | 类型 | 说明 |
|------|------|------|
| id | string | 订单ID |

**请求参数**：

```json
{
    "method": "ALIPAY"
}
```

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| method | string | 是 | 支付方式 |

**支付方式枚举**：

| 值 | 说明 |
|------|------|
| ALIPAY | 支付宝 |
| WECHAT | 微信支付 |
| BALANCE | 余额支付 |

**响应示例** (200)：

```json
{
    "code": 0,
    "message": "success",
    "data": "payment-uuid-xxx"
}
```

**业务流程**：
1. 调用 payment-service 创建支付
2. 处理支付结果
3. 更新订单状态为 `PAID`
4. 发布 `order.paid` 事件到 RabbitMQ

---

## 五、订单状态说明

| 状态 | 说明 | 可转换到 |
|------|------|---------|
| CREATED | 已创建 | PAID, CANCELLED |
| PAID | 已支付 | COMPLETED |
| COMPLETED | 已完成 | - |
| CANCELLED | 已取消 | - |

---

## 六、支付服务 (payment-service)

> 端口：11001 | 内部 gRPC 服务

支付服务不对外暴露 HTTP 接口，仅供内部服务调用。

### 6.1 gRPC 接口

**服务定义**：

```protobuf
service PaymentService {
    rpc CreatePayment(CreatePaymentRequest) returns (CreatePaymentResponse);
    rpc Pay(PayRequest) returns (PayResponse);
    rpc QueryPayment(QueryPaymentRequest) returns (QueryPaymentResponse);
}
```

---

### 6.2 CreatePayment

**请求**：

```json
{
    "order_id": "order-uuid-xxx",
    "user_id": "user-uuid-xxx",
    "amount": 199.99,
    "method": "PAYMENT_METHOD_ALIPAY"
}
```

**响应**：

```json
{
    "payment_id": "payment-uuid-xxx",
    "status": "PAYMENT_STATUS_CREATED",
    "created_at": 1704067200
}
```

---

### 6.3 Pay

**请求**：

```json
{
    "payment_id": "payment-uuid-xxx"
}
```

**响应**：

```json
{
    "status": "PAYMENT_STATUS_SUCCESS",
    "transaction_id": "txn-xxx",
    "paid_at": 1704067200
}
```

---

### 6.4 QueryPayment

**请求**：

```json
{
    "payment_id": "payment-uuid-xxx"
}
```

**响应**：

```json
{
    "payment_id": "payment-uuid-xxx",
    "order_id": "order-uuid-xxx",
    "user_id": "user-uuid-xxx",
    "amount": 199.99,
    "method": "ALIPAY",
    "status": "SUCCESS",
    "created_at": 1704067200,
    "paid_at": 1704067200
}
```

---

### 6.5 支付状态枚举

| 状态 | 说明 |
|------|------|
| CREATED | 已创建 |
| PROCESSING | 处理中 |
| SUCCESS | 成功 |
| FAILED | 失败 |
| TIMEOUT | 超时 |

---

## 七、通知服务 (notify-service)

> 消息队列消费者 | Exchange: `order_events`

通知服务通过 RabbitMQ 异步消费订单事件，不提供 HTTP 接口。

### 7.1 事件类型

| 事件类型 | 说明 | 路由键 |
|----------|------|--------|
| order.created | 订单创建 | order.created |
| order.paid | 订单支付 | order.paid |
| order.cancelled | 订单取消 | order.cancelled |

### 7.2 事件格式

**order.created**：

```json
{
    "event_type": "order.created",
    "order_id": "order-uuid-xxx",
    "user_id": "user-uuid-xxx",
    "amount": 199.99,
    "created_at": "2024-01-01T10:00:00Z"
}
```

**order.paid**：

```json
{
    "event_type": "order.paid",
    "order_id": "order-uuid-xxx",
    "user_id": "user-uuid-xxx",
    "amount": 199.99,
    "paid_at": "2024-01-01T10:30:00Z"
}
```

---

## 八、错误码说明

### 8.1 通用错误码

| code | 说明 |
|------|------|
| 0 | 成功 |
| 400 | 请求参数错误 |
| 401 | 未认证 |
| 403 | 无权限 |
| 404 | 资源不存在 |
| 500 | 系统内部错误 |

### 8.2 业务错误码

| code | 说明 |
|------|------|
| 10001 | 参数错误 |
| 10002 | 未认证 |
| 10003 | 无权限 |
| 20001 | 用户不存在 |
| 20002 | 用户已存在 |
| 30001 | 订单不存在 |
| 30002 | 订单状态不允许此操作 |
| 40001 | 支付创建失败 |
| 40002 | 支付处理失败 |
| 40003 | 重复支付 |

---

## 九、用户服务 gRPC

> 端口：8001 | 内部 gRPC 服务

### 9.1 gRPC 接口

**服务定义**：

```protobuf
service UserService {
    rpc GetUserById(GetUserRequest) returns (GetUserResponse);
    rpc GetUserByUsername(GetUserByUsernameRequest) returns (GetUserResponse);
}
```

---

### 9.2 GetUserById

**请求**：

```json
{
    "id": "user-uuid-xxx"
}
```

**响应**：

```json
{
    "id": "user-uuid-xxx",
    "username": "testuser",
    "email": "test@example.com"
}
```

---

### 9.3 GetUserByUsername

**请求**：

```json
{
    "username": "testuser"
}
```

**响应**：

```json
{
    "id": "user-uuid-xxx",
    "username": "testuser",
    "email": "test@example.com"
}
```

---

## 十、接口调用流程

### 10.1 用户登录 + 下单流程

```
1. POST /api/v1/auth/login
   └─> 获取 JWT Token

2. POST /api/v1/order/create (携带 Token)
   └─> 创建订单
   └─> 发布 order.created 事件
   └─> notify-service 消费事件

3. POST /api/v1/order/pay/:id (携带 Token)
   └─> 调用 payment-service gRPC
   └─> 支付成功
   └─> 更新订单状态为 PAID
   └─> 发布 order.paid 事件
   └─> notify-service 消费事件
```

### 10.2 消息流转

```
┌─────────────┐     order.created      ┌───────────────┐
│ order-service│ ────────────────────> │ order_events  │
│              │                        │   (exchange)  │
└─────────────┘                        └───────┬───────┘
                                                 │
                                                 │ order.#
                                                 ▼
┌─────────────┐                       ┌───────────────┐
│notify-service│ <──────────────────── │ notifications │
│              │                       │    (queue)    │
└─────────────┘                       └───────────────┘
```

---

## 十一、环境配置

### 11.1 服务端口

| 服务 | 端口 | 说明 |
|------|------|------|
| gateway-service | 8080 | API 网关入口 |
| user-service | 8001 | 用户服务 |
| auth-service | 8002 | 认证服务 |
| order-service | 10000 | 订单服务 |
| payment-service | 11001 | 支付服务 |

### 11.2 中间件

| 中间件 | 端口 | 说明 |
|--------|------|------|
| MySQL | 3306 | 主数据库 |
| Redis | 6379 | 缓存 |
| RabbitMQ | 5672 | 消息队列 |
| RabbitMQ Management | 15672 | 管理界面 |

---

## 十二、附录

### 12.1 JWT Payload 示例

```json
{
    "username": "testuser",
    "exp": 1704070800,
    "iat": 1704067200
}
```

### 12.2 Content-Type 要求

所有请求和响应的 Content-Type 均为 `application/json`。

### 12.3 请求超时

- HTTP 接口：30 秒
- gRPC 接口：30 秒
