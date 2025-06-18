# 用户服务接口文档

## 认证相关接口

### 刷新令牌
`POST /auth/refresh`

**请求头**
- Authorization:  <旧token>

**成功响应**
```json
{
  "code": 200,
  "msg": "刷新令牌成功",
  "data": {
    "token": "新JWT令牌"
  }
}
```

## 用户相关接口

### 用户登录
`POST /user/login`

**请求体**
```json
{
  "username": "用户名",
  "password": "密码"
}
```

**成功响应**
```json
{
  "code": 200,
  "msg": "登录成功",
  "data": {
    "user": {
      "id": 1,
      "username": "testuser",
      "email": "test@example.com",
      "role": "user",
      "nickname": "user",
      "avatar": "url"
    },
    "token": "JWT令牌"
  }
}
```

**错误响应**
```json
{
  "code": 400,
  "msg": "请求参数错误: 错误详情"
}
```
```json
{
  "code": 500,
  "msg": "登录失败",
  "data": {
    "token": ""
  }
}
```

### 用户注册
`POST /user/register`

**请求体**
```json
{
  "username": "用户名（3-15字符）",
  "email": "有效邮箱",
  "password": "密码（最少6位）",
  "nickname": "昵称（最多50字符）"
}
```

**成功响应**
```json
{
  "code": 200,
  "msg": "注册成功",
  "data": {
    "user": {
      "id": 1,
      "username": "testuser",
      "email": "test@example.com",
      "created_at": "2023-01-01T00:00:00Z",
      "role": "user"
    }
  }
}
```

**错误响应**
```json
{
  "code": 400,
  "msg": "请求参数错误: 错误详情"
}
```
```json
{
  "code": 500,
  "msg": "注册失败",
  "data": {
    "user": ""
  }
}
```

### 更新用户信息
`PUT /user/info`

**请求头**
- Authorization:  <有效token>

**请求体**
```json
{
  "username": "新用户名（可选）",
  "email": "新邮箱（可选）",
  "password": "新密码（可选）",
  "nickname": "新昵称（可选）",
  "avatar": "新头像URL（可选）"
}
```

**成功响应**
```json
{
  "code": 200,
  "msg": "用户信息更新成功"
}
```

**错误响应**
```json
{
  "code": 400,
  "msg": "请求参数错误: 错误详情"
}
```
```json
{
  "code": 500,
  "msg": "更新失败",
  "data": {
    "user": ""
  }
}
```

## 响应状态码说明

- **200**: 请求成功
- **400**: 请求参数错误
- **500**: 服务器内部错误

## 通用响应格式

所有接口都遵循统一的响应格式：

```json
{
  "code": 200,        // 响应状态码
  "msg": "操作结果描述", // 响应消息
  "data": {}          // 响应数据（可选）
}
```