# EnterpriseHub 接口文档（API Specification）

> 本文档描述 **EnterpriseHub 企业级微服务项目** 对外 HTTP 接口与核心内部 gRPC 接口规范，用于开发、联调与测试。

------

## 一、接口通用规范

### 1. 基础约定

- 接口风格：RESTful
- 数据格式：JSON
- 字符编码：UTF-8
- 统一入口：API Gateway

### 2. 统一返回结构

```json
{
  "code": 0,
  "message": "success",
  "data": {}
}
```

| 字段    | 类型   | 说明                      |
| ------- | ------ | ------------------------- |
| code    | int    | 0 表示成功，非 0 表示错误 |
| message | string | 错误或成功信息            |
| data    | object | 返回数据                  |

------

### 3. 统一错误码示例

| code  | 含义         |
| ----- | ------------ |
| 200   | 成功         |
| 10001 | 参数错误     |
| 10002 | 未认证       |
| 10003 | 无权限       |
| 20001 | 用户不存在   |
| 30001 | 订单不存在   |
| 50000 | 系统内部错误 |

------

### 4. 认证方式

- 使用 **JWT Bearer Token**
- Header 示例：

```
Authorization: Bearer <token>
```

------

## 二、Gateway 接口（对外）

> 所有客户端请求均通过 Gateway 转发

------

## 三、用户服务（user-service）

### 1. 用户注册

- **URL**：`POST /api/users/register`
- **鉴权**：否

#### 请求参数

```json
{
  "username": "testuser",
  "password": "123456",
  "email": "test@test.com"
}
```

#### 响应示例

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "user_id": 1
  }
}
```

------

### 2. 获取用户信息

- **URL**：`GET /api/users/{id}`
- **鉴权**：是

#### 响应示例

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": 1,
    "username": "testuser",
    "email": "test@test.com",
    "status": 1
  }
}
```

------

### 3. 更新用户信息

- **URL**：`PUT /api/users/{id}`
- **鉴权**：是

#### 请求参数

```json
{
  "email": "new@test.com",
  "status": 1
}
```

------

## 四、认证服务（auth-service）

### 1. 用户登录

- **URL**：`POST /api/auth/login`
- **鉴权**：否

#### 请求参数

```json
{
  "username": "testuser",
  "password": "123456"
}
```

#### 响应示例

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "token": "jwt-token"
  }
}
```

------

### 2. Token 校验（内部）

- **URL**：`POST /api/auth/verify`
- **鉴权**：是

------

## 五、订单服务（order-service）

### 1. 创建订单

- **URL**：`POST /api/orders`
- **鉴权**：是

#### 请求参数

```json
{
  "amount": 199.99
}
```

#### 响应示例

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "order_id": 10001,
    "status": "CREATED"
  }
}
```

------

### 2. 查询订单

- **URL**：`GET /api/orders/{id}`
- **鉴权**：是

------

### 3. 支付订单

- **URL**：`POST /api/orders/{id}/pay`
- **鉴权**：是

------

## 六、支付服务（payment-service）

### 1. 支付处理（内部）

- **协议**：gRPC
- **说明**：模拟第三方支付，支付成功后回调订单服务

------

## 七、通知服务（notify-service）

### 1. 通知消息消费

- **方式**：RabbitMQ
- **事件类型**：
  - ORDER_CREATED
  - ORDER_PAID

------

## 八、gRPC 接口（核心示例）

### UserService

```proto
service UserService {
  rpc GetUserById(GetUserRequest) returns (GetUserResponse);
}
```

### OrderService

```proto
service OrderService {
  rpc UpdateOrderStatus(UpdateOrderRequest) returns (UpdateOrderResponse);
}
```

------

## 九、接口调用流程示例

### 登录 + 下单流程

1. 调用 `/api/auth/login` 获取 JWT
2. 携带 JWT 调用 `/api/orders`
3. order-service 创建订单
4. payment-service 处理支付
5. notify-service 发送通知

------

## 十、说明

- 本接口文档用于学习与项目实践
- 可根据业务复杂度进行裁剪或扩展
- 所有接口均可通过 Gateway 统一接入

------

> 建议：后续可基于本接口文档生成 Swagger / OpenAPI 文档