# 图片服务 API 接口文档

## 概述

图片服务提供图片上传和访问功能，支持JPEG、PNG、GIF格式的图片文件。

**基础URL**: `http://your-domain/img`

## 接口列表

### 1. 图片上传接口

#### 接口信息
- **URL**: `/img/`
- **方法**: `POST`
- **认证**: 需要认证（通过 AuthMiddleware）
- **内容类型**: `multipart/form-data`

#### 请求参数

| 参数名 | 类型 | 必填 | 描述 |
|--------|------|------|------|
| file | File | 是 | 要上传的图片文件 |

#### 支持的文件类型

- `image/jpeg` - JPEG格式图片
- `image/png` - PNG格式图片  
- `image/gif` - GIF格式图片

#### 响应格式

**成功响应 (200)**
```json
{
    "code": 200,
    "msg": "图片上传成功",
    "data": {
        "path": "图片存储路径"
    }
}
```

**错误响应**

| HTTP状态码 | 描述 | 响应示例 |
|------------|------|----------|
| 400 | 请求参数错误或文件类型不支持 | `{"code": 400, "msg": "只允许上传JPEG, PNG或GIF图片", "data": {"path": ""}}` |
| 500 | 服务器内部错误 | `{"code": 500, "msg": "图片上传失败", "data": {"path": ""}}` |
| 401 | 未认证 | `{"error": "认证失败"}` |

#### 请求示例

**cURL 示例**
```bash
curl -X POST "http://your-domain/img/" \
  -H "Authorization:  your-token" \
  -F "file=@/path/to/your/image.jpg"
```

**JavaScript 示例**
```javascript
const formData = new FormData();
formData.append('file', file);

fetch('/img/', {
    method: 'POST',
    headers: {
        'Authorization': ' your-token'
    },
    body: formData
})
.then(response => response.json())
.then(data => {
    if (data.code === 200) {
        console.log('上传成功，图片路径:', data.data.path);
        // 显示上传的图片
        showUploadedImage(data.data.path);
    } else {
        console.error('上传失败:', data.msg);
    }
});
```

#### 响应示例

**成功响应**
```json
{
    "code": 200,
    "msg": "图片上传成功",
    "data": {
        "path": "/uploads/2024/01/15/image_123456.jpg"
    }
}
```

**错误响应**
```json
{
    "code": 400,
    "msg": "只允许上传JPEG, PNG或GIF图片",
    "data": {
        "path": ""
    }
}
```

### 2. 图片访问接口

#### 接口信息
- **URL**: `/img/{filepath}`
- **方法**: `GET`
- **认证**: 不需要
- **描述**: 静态文件服务，用于访问已上传的图片

#### 请求参数

| 参数名 | 类型 | 必填 | 描述 |
|--------|------|------|------|
| filepath | String | 是 | 图片文件路径（在URL路径中） |

#### 响应格式

- **成功**: 直接返回图片文件内容，Content-Type根据文件类型设置
- **失败**: 返回404状态码

#### 请求示例

```bash
# 直接在浏览器中访问或使用以下cURL命令
curl "http://your-domain/img/uploads/2024/01/15/image_123456.jpg"
```

#### HTML中使用示例

```html
<img src="/img/uploads/2024/01/15/image_123456.jpg" alt="上传的图片">
```

## 错误码说明

| 错误码 | 描述 |
|--------|------|
| 400 | 请求参数错误、文件格式不支持或文件读取失败 |
| 401 | 认证失败（仅上传接口） |
| 404 | 文件不存在（仅访问接口） |
| 500 | 服务器内部错误 |

## 注意事项

1. **文件大小限制**: 请确认服务器配置的文件大小限制
2. **认证要求**: 上传接口需要有效的认证token
3. **文件类型**: 只支持JPEG、PNG、GIF三种图片格式
4. **路径保存**: 上传成功后返回的path可直接用于访问接口
5. **静态文件**: 图片文件存储在 `./img` 目录下

## 完整使用流程示例

### 1. 上传图片
```javascript
// 选择文件
const fileInput = document.getElementById('fileInput');
const file = fileInput.files[0];

// 创建FormData
const formData = new FormData();
formData.append('file', file);

// 上传请求
fetch('/img/', {
    method: 'POST',
    headers: {
        'Authorization': ' your-auth-token'
    },
    body: formData
})
.then(response => response.json())
.then(data => {
    if (data.code === 200) {
        console.log('上传成功，图片路径:', data.data.path);
        // 显示上传的图片
        showUploadedImage(data.data.path);
    } else {
        console.error('上传失败:', data.msg);
    }
});
```

### 2. 显示图片
```javascript
function showUploadedImage(path) {
    const img = document.createElement('img');
    img.src = `/img${path}`;  // 注意路径拼接
    img.alt = '上传的图片';
    document.body.appendChild(img);
}
``` 