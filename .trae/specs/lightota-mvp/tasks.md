# Tasks

## Phase 1: 项目初始化和基础框架

- [ ] Task 1: 初始化 Go 项目结构，配置 Go modules
  - 创建项目基础目录结构 (cmd, internal, pkg 等)
  - 初始化 go.mod，添加依赖 (Gin, GORM, 等)

- [ ] Task 2: 配置处理模块
  - 实现配置结构定义
  - 支持配置文件和环境变量
  - 使用 Viper 库

- [ ] Task 3: 创建数据模型
  - 实现 tenants 租户模型
  - 实现 products 产品模型
  - 实现 devices 设备关联模型
  - 实现 firmwares 固件模型
  - 实现 upgrade_tasks 升级任务模型
  - 实现 device_upgrade_records 设备升级记录模型
  - 配置 GORM 数据库连接

- [ ] Task 4: API Key 认证中间件
  - 实现 API Key 认证中间件
  - 实现租户隔离逻辑
  - 返回 401 错误处理

## Phase 2: 基础 API 实现

- [ ] Task 5: 租户管理 API
  - 实现租户创建、列表、详情、更新 API

- [ ] Task 6: 产品管理 API
  - 实现产品创建、列表、详情、更新、删除 API
  - 租户隔离

- [ ] Task 7: 固件管理 API
  - 实现固件上传接口（文件存储 + MD5 计算）
  - 实现固件列表、详情、删除 API
  - 固件文件存储（本地文件系统）

## Phase 3: 第三方对接和升级任务

- [ ] Task 8: 第三方设备对接 HTTP 客户端
  - 实现获取单个设备信息
  - 实现批量筛选设备接口调用

- [ ] Task 9: 升级任务创建 API
  - 实现升级任务创建，支持三种方式：指定设备/全量/灰度
  - 灰度模式调用第三方接口筛选设备
  - 创建设备升级记录

- [ ] Task 10: 升级任务查询 API
  - 实现任务列表查询
  - 实现任务详情查询，包含统计信息

## Phase 4: 推送和状态跟踪

- [ ] Task 11: MQTT 推送模块
  - 实现 MQTT 客户端连接
  - 实现升级通知推送

- [ ] Task 12: HTTP 推送模块
  - 实现 HTTP 长连接管理
  - 实现升级通知推送

- [ ] Task 13: 设备端 API
  - 实现 check-update 接口
  - 固件下载接口
  - 升级结果上报接口
  - 更新设备升级记录状态

- [ ] Task 14: 任务进度统计
  - 实现任务进度统计逻辑
  - 计算各状态设备数量和完成百分比

- [ ] Task 15: 升级任务执行和限流
  - 实现按速率推送升级通知
  - 任务状态管理（创建、运行、暂停、完成、取消）

## Phase 5: 前端开发

- [ ] Task 16: 初始化 React 前端项目
  - 使用 Vite 创建 React + TypeScript 项目
  - 配置 Ant Design
  - 配置 Axios 和 Zustand

- [ ] Task 17: 前端页面开发
  - 租户列表和管理页面
  - 产品列表和管理页面
  - 固件管理和上传页面
  - 创建升级任务页面
  - 任务列表和详情页面

## Phase 6: 部署和验证

- [ ] Task 18: 创建 Docker Compose 配置
  - LightOTA 服务配置
  - SQLite 数据卷配置
  - 可选 Mosquitto 配置
  - 可选 PostgreSQL 配置

- [ ] Task 19: 端到端测试验证
  - 验证所有 API 正常工作
  - 验证完整升级流程
  - 验证 Docker Compose 一键部署

# Task Dependencies

- [Task 2] 依赖 [Task 1]
- [Task 3] 依赖 [Task 2]
- [Task 4] 依赖 [Task 3]
- [Task 5] 依赖 [Task 4]
- [Task 6] 依赖 [Task 4]
- [Task 7] 依赖 [Task 6]
- [Task 8] 依赖 [Task 3]
- [Task 9] 依赖 [Task 6] [Task 7] [Task 8]
- [Task 10] 依赖 [Task 9]
- [Task 11] 依赖 [Task 2]
- [Task 12] 依赖 [Task 2]
- [Task 13] 依赖 [Task 9]
- [Task 14] 依赖 [Task 13]
- [Task 15] 依赖 [Task 11] [Task 12] [Task 14]
- [Task 17] 依赖 [Task 16]
- [Task 18] 不依赖其他任务，可以并行
- [Task 19] 依赖所有任务完成
