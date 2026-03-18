# LightOTA

LightOTA 是一个轻量级、易于集成的嵌入式设备 OTA (Over-the-Air) 升级服务。如果你已经有自己的设备管理系统，只需要暴露简单的 API 就能快速获得完整的固件升级功能。

## 🎯 项目目标

很多 IoT 项目都需要固件升级功能，但从零开发一套稳定可靠的 OTA 系统比较耗时，而且多数情况下你已经有了自己的设备管理系统，只缺一个固件管理和升级任务调度的模块。

LightOTA 的设计目标：

- **轻量简单** - 单个二进制文件，依赖少，部署简单
- **易于集成** - 松耦合设计，不强制你把设备数据迁移过来
- **功能完整** - 支持全量升级、灰度升级、进度追踪、速率控制
- **多租户** - 支持多租户隔离，适合 SaaS 场景

## ✨ 核心功能

| 功能 | 说明 |
|------|------|
| 固件管理 | 多版本管理，支持查看变更日志和发布说明 |
| 多种升级模式 | 指定设备 / 全量升级 / 按百分比灰度升级 |
| 速率控制 | 可配置每秒推送数量，避免服务器突刺 |
| 进度追踪 | 实时查看每个设备的升级状态和结果 |
| 第三方集成 | 对接你现有的设备管理系统，无需重复造轮子 |
| 租户隔离 | 严格的数据隔离，支持多租户 SaaS 部署 |
| Web 管理后台 | 基于 React 的现代 UI，操作简单直观 |
| MQTT 通知 | 可选支持 MQTT 主动推送升级通知 |
| HTTP 长连接 | 支持 HTTPS SSE 长连接推送 |

## 🏗️ 架构设计

```
┌─────────────────────────────────────────────────────────────┐
│  Your Device Management System                               │
│  ──────────────────                                          │
│  维护设备信息，暴露两个 API 给 LightOTA 调用                  │
└─────────────────────────────────────────────────────────────┘
                              ↓ HTTP API 获取设备列表
┌─────────────────────────────────────────────────────────────┐
│  LightOTA Server (Go + Gin + GORM)                           │
│  ─────────────────────────                                   │
│  • 固件存储管理                                              │
│  • 升级任务创建和调度                                        │
│  • 按速率推送升级通知                                        │
│  • 升级进度记录                                              │
└─────────────────────────────────────────────────────────────┘
                              ↓ MQTT / SSE 通知
┌─────────────────────────────────────────────────────────────┐
│  Embedded Device                                             │
│  ─────────────────                                           │
│  1. 检查更新: GET /api/v1/ota/check-update                   │
│  2. 下载固件: 从 LightOTA 直接下载                           │
│  3. 上报结果: POST /api/v1/ota/report-result                │
└─────────────────────────────────────────────────────────────┘
```

**设计特点：**

- 设备信息由你管理，LightOTA 只负责升级相关
- 设备端直接和 LightOTA 交互，不经过你的系统转发
- 松耦合设计，对接成本极低，只需要两个 API

## 🚀 快速开始

### 使用 Docker Compose 启动（推荐）

```bash
git clone https://github.com/ZhangPengPaul/LightOTA.git
cd LightOTA
docker-compose -f deploy/docker-compose.yml up -d
```

然后访问 `http://localhost:8080` 即可。

### 本地编译运行

**要求:**
- Go 1.24+
- Node.js 18+ (编译前端)

```bash
# 克隆项目
git clone https://github.com/ZhangPengPaul/LightOTA.git
cd LightOTA

# 编译前端
cd web/react-admin
npm install
npm run build
cd ../..

# 编译后端
go build -o lightota cmd/server/main.go

# 运行
./lightota
```

### 环境变量

| 变量 | 说明 | 默认值 |
|------|------|--------|
| `DB_TYPE` | 数据库类型: sqlite/postgres | `sqlite` |
| `DB_PATH` | SQLite 数据库文件路径 | `./lightota.db` |
| `DB_HOST` | PostgreSQL 主机 | |
| `DB_PORT` | PostgreSQL 端口 | `5432` |
| `DB_USER` | PostgreSQL 用户名 | |
| `DB_PASSWORD` | PostgreSQL 密码 | |
| `DB_NAME` | PostgreSQL 数据库名 | |
| `FIRMWARE_STORAGE_PATH` | 固件存储目录 | `./firmwares` |
| `MQTT_ENABLE` | 是否启用 MQTT | `false` |
| `MQTT_BROKER` | MQTT 地址 | |
| `MQTT_USERNAME` | MQTT 用户名 | |
| `MQTT_PASSWORD` | MQTT 密码 | |
| `MQTT_CLIENT_ID` | MQTT 客户端 ID | |
| `SERVER_PORT` | 服务端口 | `8080` |

## 🔌 第三方对接

如果你已经有自己的设备管理系统，请阅读 [第三方对接指南](./THIRD_PARTY_INTEGRATION.md)，里面有详细的 API 规范和示例代码。

**对接只需要两步:**
1. 在你的系统暴露两个 API（获取单个设备、批量筛选设备）
2. 在 LightOTA 租户配置里填写 API 地址和 API Key

## 📖 使用流程

1. **登录管理后台** - 默认使用 API Key 认证
2. **创建租户** - 每个租户数据隔离
3. **创建产品** - 对应你的设备产品型号
4. **上传固件** - 填写版本号、版本代码、变更日志
5. **创建升级任务** - 选择升级模式：
   - **Specified** - 指定设备 ID 列表升级
   - **All** - 升级该产品下所有设备
   - **Gray** - 按百分比灰度升级
6. **查看进度** - 在表格中查看每个设备的升级状态
7. **设备端** - 设备调用 check-update 接口检查更新，有更新就下载升级

## 设备端 API

设备端需要调用三个接口：

### 1. 检查更新

```
GET /api/v1/ota/check-update?deviceId={deviceId}&currentVersion={currentVersion}
```

**响应示例:**
```json
{
  "code": 0,
  "data": {
    "has_update": true,
    "version": "v1.2.3",
    "version_code": 123,
    "download_url": "https://your-lightota-host/api/v1/ota/download/{firmwareId}?deviceId={deviceId}",
    "md5": "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx",
    "file_size": 123456
  }
}
```

如果没有更新，`has_update` 为 `false`。

### 2. 下载固件

直接 GET 请求 `download_url` 下载固件文件，会返回二进制内容。

### 3. 上报升级结果

```
POST /api/v1/ota/report-result
Content-Type: application/json

{
  "device_id": "xxxx",
  "task_id": "xxxx",
  "success": true,
  "error_message": ""
}
```

## 📦 项目结构

```
LightOTA/
├── cmd/server/          # 主程序入口
├── internal/
│   ├── config/          # 配置加载
│   ├── handler/         # HTTP 处理器
│   ├── model/           # 数据模型定义
│   ├── repository/      # 数据访问层
│   ├── service/         # 业务逻辑
│   ├── mqtt/            # MQTT 客户端
│   └── httpsse/         # SSE 推送
├── web/react-admin/     # React 管理后台
├── deploy/              # Docker 部署文件
└── THIRD_PARTY_INTEGRATION.md  # 第三方对接文档
```

## 🛠️ 技术栈

- **后端**: Go + Gin + GORM
- **数据库**: SQLite (默认) / PostgreSQL (可选)
- **前端**: React + TypeScript + Vite
- **消息推送**: MQTT / HTTP SSE
- **部署**: Docker / 二进制直接运行

## 🤔 什么时候选择 LightOTA

✅ **适合使用:**
- 你已经有设备管理系统，只需要 OTA 升级功能
- 你想要快速集成，不想从零开发
- 你需要灰度升级、进度追踪等高级功能
- 你需要多租户 SaaS 部署

❌ **不适合:**
- 你想要一个包含设备管理的完整 IoT 平台（LightOTA 只做 OTA 部分）
- 你需要非常复杂的分组、标签等高级设备筛选（这些应该在你的设备管理系统实现）

## 📄 许可证

MIT License - 详见 [LICENSE](./LICENSE) 文件。

## 🙏 贡献

欢迎提交 Issue 和 Pull Request！
