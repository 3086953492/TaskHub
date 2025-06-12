# 用户服务接口文档

## 认证相关接口

### 刷新令牌
`POST /auth/refresh`

**请求头**
- Authorization: Bearer <旧token>

**成功响应**
```json
{
  "token": "新JWT令牌"
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
  "token": "JWT令牌"
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
  "user": {
    "id": 1,
    "username": "testuser",
    "email": "test@example.com",
    "created_at": "2023-01-01T00:00:00Z",
    "role": "user"
  }
}
```

### 更新用户信息
`PUT /user/info`

**请求头**
- Authorization: Bearer <有效token>

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
  "msg": "用户信息更新成功"
}
```