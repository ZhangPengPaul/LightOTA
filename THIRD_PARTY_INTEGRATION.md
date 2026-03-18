# 第三方设备管理系统对接指南

本文档描述了如果您已有自己的设备管理系统，如何对接 LightOTA 实现 OTA 升级功能。

## 概述

LightOTA 采用松耦合设计，只负责固件存储、升级任务管理、推送通知。设备信息由您现有的设备管理系统维护，LightOTA 通过 HTTP API 从您的系统获取设备信息。

## 对接要求

您需要在您的设备管理系统中暴露两个 HTTP API 接口给 LightOTA 调用：

1. **获取单个设备信息** - 根据设备 ID 查询设备详情
2. **批量筛选设备** - 根据条件筛选设备（用于灰度升级）

## API 规范

### 1. 获取单个设备信息

**Endpoint:** `GET /api/v1/devices/{deviceId}`

**Headers:**
- `Authorization: Bearer {apiKey}` - LightOTA 会携带租户配置的 API Key

**Response:**

```json
{
  "code": 0,
  "data": {
    "deviceId": "string",
    "deviceName": "string",
    "currentVersion": "string",
    "productId": "string",
    "online": "boolean"
  }
}
```

**字段说明：**
| 字段 | 类型 | 说明 |
|------|------|------|
| `deviceId` | string | 设备唯一 ID（您这边的设备 ID）|
| `deviceName` | string | 设备名称 |
| `currentVersion` | string | 当前固件版本 |
| `productId` | string | 产品 ID（对应 LightOTA 中的产品）|
| `online` | boolean | 是否在线 |

---

### 2. 批量筛选设备

**Endpoint:** `POST /api/v1/devices/query`

**Headers:**
- `Authorization: Bearer {apiKey}`
- `Content-Type: application/json`

**Request Body:**

```json
{
  "productId": "string",
  "percent": "int",
  "excludeVersions": "string[]",
  "onlyOnline": "boolean",
  "limit": "int"
}
```

**字段说明：**
| 字段 | 类型 | 说明 |
|------|------|------|
| `productId` | string | 筛选指定产品 |
| `percent` | int | 灰度百分比，1-100。比如 10 表示随机选 10% 的设备 |
| `excludeVersions` | string[] | 排除指定版本的设备 |
| `onlyOnline` | boolean | 只返回在线设备 |
| `limit` | int | 最大返回数量 |

**Response:**

```json
{
  "code": 0,
  "data": {
    "total": "int",
    "selected": [
      {
        "deviceId": "string",
        "deviceName": "string",
        "currentVersion": "string",
        "productId": "string",
        "online": "boolean"
      }
    ]
  }
}
```

**Response 字段说明：**
| 字段 | 类型 | 说明 |
|------|------|------|
| `total` | int | 符合条件的设备总数 |
| `selected` | array | 选中的设备列表 |

选中的设备列表每个设备字段同「获取单个设备信息」。

---

## 配置步骤

### 1. 在 LightOTA 创建租户

1. 登录 LightOTA 管理后台
2. 创建租户
3. 复制生成的 **API Key**
4. 填写 **External Device API URL**（您的 API 地址，例如 `https://your-device-api.com`）

### 2. 在您的系统配置 API Key

在您的系统中配置 LightOTA 请求过来时要校验的 API Key，与 LightOTA 租户配置保持一致。

### 3. 创建产品和固件

1. 在 LightOTA 创建产品，产品 ID 需要和您系统中的产品 ID 对应
2. 上传固件版本

### 4. 创建升级任务

1. 在 LightOTA 创建升级任务，选择：
   - **Specified** - 指定设备 ID 列表
   - **All** - 全量升级该产品下所有设备
   - **Gray** - 按百分比灰度升级
2. LightOTA 会调用您的 API 获取设备列表
3. LightOTA 创建升级记录，开始按速率推送升级通知

---

## 工作流程

```
1. 用户在 LightOTA 创建升级任务
       ↓
2. LightOTA 调用您的 API 获取符合条件的设备列表
       ↓
3. LightOTA 为每个设备创建升级记录
       ↓
4. LightOTA 后台goroutine按推送速率处理
       ↓
5. LightOTA 通过 MQTT 或 HTTP SSE 推送升级通知给设备
       ↓
6. 设备收到通知后，调用 LightOTA check-update 接口确认更新
       ↓
7. 设备从 LightOTA 下载新固件
       ↓
8. 设备升级完成后上报结果给 LightOTA
       ↓
9. LightOTA 更新升级记录状态，您可以在管理后台查看进度
```

---

## 设备端流程

设备端仍然直接和 LightOTA 交互，不需要通过您的系统转发：

1. **检查更新** - `GET https://lightota-host/api/v1/ota/check-update?deviceId={xxx}&currentVersion={yyy}`
2. **下载固件** - 从 LightOTA 返回的 downloadUrl 下载
3. **上报结果** - `POST https://lightota-host/api/v1/ota/report-result` 上报升级结果

> 您不需要在您的设备端做任何修改，只需要暴露管理 API 给 LightOTA 即可。

---

## 安全说明

- LightOTA 每次请求都会携带 `Authorization: Bearer {apiKey}`，请务必校验
- API Key 要保存在安全地方，不要泄露
- 建议使用 HTTPS 加密传输
- LightOTA 会设置 30 秒超时，请确保您的 API 在超时内响应

---

## 示例代码

### Node.js Express 示例

```javascript
app.get('/api/v1/devices/:deviceId', (req, res) => {
  const auth = req.headers.authorization;
  if (auth !== `Bearer ${process.env.LIGHTOTA_API_KEY}`) {
    return res.status(401).json({ code: 401, message: 'unauthorized' });
  }
  
  const device = await db.devices.find(req.params.deviceId);
  res.json({
    code: 0,
    data: {
      deviceId: device.id,
      deviceName: device.name,
      currentVersion: device.currentVersion,
      productId: device.productId,
      online: device.online
    }
  });
});

app.post('/api/v1/devices/query', (req, res) => {
  const auth = req.headers.authorization;
  if (auth !== `Bearer ${process.env.LIGHTOTA_API_KEY}`) {
    return res.status(401).json({ code: 401, message: 'unauthorized' });
  }
  
  const { productId, percent, excludeVersions, onlyOnline, limit } = req.body;
  
  let query = db.devices.where({ productId });
  if (onlyOnline) {
    query = query.where({ online: true });
  }
  if (excludeVersions.length > 0) {
    query = query.whereNot({ currentVersion: excludeVersions });
  }
  
  let devices = await query.get();
  
  // 按百分比随机抽样
  if (percent < 100 && devices.length > 0) {
    const targetCount = Math.round(devices.length * percent / 100);
    devices = shuffle(devices).slice(0, targetCount);
  }
  
  if (limit) {
    devices = devices.slice(0, limit);
  }
  
  res.json({
    code: 0,
    data: {
      total: devices.length,
      selected: devices.map(d => ({
        deviceId: d.id,
        deviceName: d.name,
        currentVersion: d.currentVersion,
        productId: d.productId,
        online: d.online
      }))
    }
  });
});
```

### Go Gin 示例

```go
package main

import (
	"net/http"
	"math/rand"

	"github.com/gin-gonic/gin"
)

type Device struct {
	ID             string `json:"deviceId"`
	Name           string `json:"deviceName"`
	CurrentVersion string `json:"currentVersion"`
	ProductID      string `json:"productId"`
	Online         bool   `json:"online"`
}

type QueryRequest struct {
	ProductID       string   `json:"productId"`
	Percent         int      `json:"percent"`
	ExcludeVersions []string `json:"excludeVersions"`
	OnlyOnline      bool     `json:"onlyOnline"`
	Limit           int      `json:"limit"`
}

func main() {
	r := gin.Default()
	
	r.GET("/api/v1/devices/:deviceId", func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		expected := "Bearer " + os.Getenv("LIGHTOTA_API_KEY")
		if auth != expected {
			c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "unauthorized"})
			return
		}
		
		deviceID := c.Param("deviceId")
		device := findDevice(deviceID)
		
		c.JSON(http.StatusOK, gin.H{
			"code": 0,
			"data": device,
		})
	})
	
	r.POST("/api/v1/devices/query", func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		expected := "Bearer " + os.Getenv("LIGHTOTA_API_KEY")
		if auth != expected {
			c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "unauthorized"})
			return
		}
		
		var req QueryRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error()})
			return
		}
		
		devices := queryDevices(req.ProductID, req.OnlyOnline, req.ExcludeVersions)
		
		if req.Percent < 100 && len(devices) > 0 {
			targetCount := len(devices) * req.Percent / 100
			rand.Shuffle(len(devices), func(i, j int) {
				devices[i], devices[j] = devices[j], devices[i]
			})
			devices = devices[:targetCount]
		}
		
		if req.Limit > 0 && len(devices) > req.Limit {
			devices = devices[:req.Limit]
		}
		
		c.JSON(http.StatusOK, gin.H{
			"code": 0,
			"data": gin.H{
				"total": len(devices),
				"selected": devices,
			},
		})
	})
	
	r.Run(":8080")
}
```

---

## 故障排查

**Q: LightOTA 提示 "third party api returned error code"**  
A: 请检查您的 API 返回的 JSON 格式是否正确，`code` 字段是否为 0。

**Q: 创建任务后没有设备**  
A: 请检查：
- 产品 ID 是否一致（LightOTA 的 productId 需要和您系统的 productId 一致）
- API 返回的 selected 数组是否为空
- API 是否正确处理了 percent 参数

**Q: 网络超时**  
A: LightOTA 设置了 30 秒超时，请优化您的 API 查询速度，确保 30 秒内响应。

---

## 联系

如果有对接问题，请提交 GitHub Issue。
