# 🚀 TaskHub 项目说明

---

## 🔧 环境要求
- **Go 语言** : 1.20+  
- **MySQL 数据库** : 5.7+  
- **Nginx 服务器** : 1.18+  

---

## 🗄️ 数据库初始化脚本
### [taskhub.sql](server/taskhub.sql)
---
## ⚡ 快速启动

### 👤 用户服务
```powershell
# 进入服务目录
cd server\user_service

# 启动服务（默认端口 8081）
go run main.go
```

---

### 🔄 Nginx 配置
#### 步骤 1：配置 `nginx.conf`
```nginx
worker_processes 1;

events {
    worker_connections 1024;
}

http {
    include       mime.types;
    default_type  application/octet-stream;
    sendfile      on;
    keepalive_timeout  65;

    # 用户服务代理
    upstream user_server {
        server 127.0.0.1:8081;
    }

    # 图片服务代理
    upstream img_server {
        server 127.0.0.1:8082;
    }

    server {
        listen 8080;

        location /user {
            proxy_pass      http://user_server;
            proxy_set_header Host $host;
            proxy_set_header X-Real-IP $remote_addr;
        }

        location /auth {
            proxy_pass      http://user_server;
            proxy_set_header Host $host;
            proxy_set_header X-Real-IP $remote_addr;
        }
        
        location /img {
            proxy_pass      http://img_server;
            proxy_set_header Host $host;
            proxy_set_header X-Real-IP $remote_addr;
        }
    }
}
```

#### 步骤 2：启动 Nginx 服务
```powershell
nginx.exe
```

---

## 📚 接口文档
### [用户服务接口文档](server/user_service/API_DOCUMENTATION.md)
### [图片服务接口文档](server/img_service/API文档.md)