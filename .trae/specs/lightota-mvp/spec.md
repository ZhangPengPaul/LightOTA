# LightOTA MVP 开发规格

## Why
打造一个轻量级、可对接、自托管的独立 OTA 升级平台，专注只做固件升级。解决已有设备管理系统需要 OTA 功能但不想使用重量级平台的问题，通过松耦合设计，保持功能完整的同时，追求极致轻量，部署简单，对接方便。

## What Changes
- 完整实现 LightOTA 后端 API 服务，包含 P0 所有核心模块
- 实现 React 后台管理前端页面
- 支持 Docker Compose 一键部署
- **包含的 P0 功能模块**:
  - 租户认证：API Key 认证，租户隔离
  - 产品管理：产品创建、编辑、列表
  - 固件管理：固件上传、版本管理、MD5 校验
  - 第三方对接：对接第三方设备管理 API
  - 升级任务创建：支持指定设备/全量/灰度三种方式
  - MQTT 推送：内置 Mosquitto 推送升级通知
  - HTTP 推送：支持 HTTP 长连接推送
  - 状态跟踪：实时跟踪每个设备升级状态
  - 基础监控：任务进度统计、成功率展示
  - 管理前端：React 后台管理页面
  - Docker 部署：Docker Compose 一键部署

## Impact
- Affected specs: 这是初始 MVP 版本，无现有规格
- Affected code:
  - `cmd/server/` - 主服务入口
  - `internal/config/` - 配置处理
  - `internal/model/` - 数据模型
  - `internal/handler/` - API 处理器
  - `internal/service/` - 业务逻辑
  - `internal/repository/` - 数据访问层
  - `internal/mqtt/` - MQTT 推送
  - `internal/util/` - 工具函数
  - `web/react-admin/` - React 管理后台
  - `deploy/` - 部署配置

## ADDED Requirements

### Requirement: 租户认证
系统 SHALL 支持多租户隔离，每个租户通过 API Key 进行认证。

#### Scenario: 正确 API Key 访问
- **WHEN** 请求携带正确的 API Key 在 Authorization header 中
- **THEN** 请求被允许访问，返回正常数据

#### Scenario: 错误 API Key 访问
- **WHEN** 请求携带错误或缺失 API Key
- **THEN** 返回 401 未授权错误

### Requirement: 产品管理
系统 SHALL 支持产品的创建、编辑、查询和列表功能。

#### Scenario: 创建产品
- **WHEN** 租户创建新产品
- **THEN** 产品被创建，返回产品信息

#### Scenario: 查询产品列表
- **WHEN** 租户查询产品列表
- **THEN** 返回该租户下所有产品列表

### Requirement: 固件管理
系统 SHALL 支持固件上传、版本管理和 MD5 校验。

#### Scenario: 上传固件
- **WHEN** 用户上传固件文件
- **THEN** 计算 MD5 校验值，存储固件文件，保存固件信息到数据库

#### Scenario: 版本管理
- **WHEN** 用户查看固件列表
- **THEN** 返回该产品下所有固件版本，按版本号排序

### Requirement: 第三方设备对接
系统 SHALL 支持通过 HTTP API 对接第三方设备管理系统。

#### Scenario: 获取单个设备信息
- **WHEN** LightOTA 需要获取设备信息
- **THEN** LightOTA 调用第三方 API 获取设备信息，包含当前版本和在线状态

#### Scenario: 批量筛选设备（灰度）
- **WHEN** 创建灰度升级任务
- **THEN** LightOTA 调用第三方筛选接口获取指定百分比的设备列表

### Requirement: 升级任务创建
系统 SHALL 支持三种升级方式：指定设备、全量升级、灰度百分比升级。

#### Scenario: 创建指定设备升级
- **WHEN** 用户选择指定设备创建任务
- **THEN** 任务创建成功，只包含指定设备

#### Scenario: 创建灰度升级
- **WHEN** 用户指定灰度百分比创建任务
- **THEN** 调用第三方接口筛选对应百分比设备，创建任务

### Requirement: MQTT 推送
系统 SHALL 支持通过 MQTT 推送升级通知给设备。

#### Scenario: 推送升级通知
- **WHEN** 升级任务需要通知设备升级
- **THEN** 通过 MQTT 发送升级通知消息到对应设备主题

### Requirement: HTTP 推送
系统 SHALL 支持通过 HTTP 长连接推送升级通知。

#### Scenario: HTTP 推送通知
- **WHEN** 设备保持 HTTP 长连接
- **THEN** 有升级任务时通过连接推送通知给设备

### Requirement: 状态跟踪
系统 SHALL 实时跟踪每个设备的升级过程状态。

#### Scenario: 更新状态
- **WHEN** 设备升级状态发生变化
- **THEN** 系统更新记录，任务进度同步更新

#### Scenario: 查询任务状态
- **WHEN** 用户查询任务进度
- **THEN** 返回任务总体进度和各状态统计（成功、失败、待处理）

### Requirement: 基础监控
系统 SHALL 展示任务进度统计和成功率。

#### Scenario: 任务统计
- **WHEN** 用户查看任务详情
- **THEN** 显示目标设备总数、成功数、失败数、待处理数和完成百分比

### Requirement: 设备端检查更新
系统 SHALL 提供设备端 API 检查是否有可用更新。

#### Scenario: 有更新
- **WHEN** 设备调用 check-update 接口且有新版本
- **THEN** 返回更新信息，包含版本、下载链接、MD5、文件大小

#### Scenario: 无更新
- **WHEN** 设备调用 check-update 接口且无新版本
- **THEN** 返回无更新状态

### Requirement: React 管理后台
系统 SHALL 提供 React 后台管理页面用于管理租户、产品、固件和升级任务。

#### Scenario: 管理界面访问
- **WHEN** 用户访问管理界面
- **THEN** 可以进行租户、产品、固件的增删改查操作，可以创建升级任务和查看任务进度

### Requirement: Docker Compose 部署
系统 SHALL 提供 Docker Compose 配置文件支持一键部署。

#### Scenario: 一键部署
- **WHEN** 用户执行 docker-compose up
- **THEN** LightOTA 服务启动，包含数据库和可选的 MQTT 服务
