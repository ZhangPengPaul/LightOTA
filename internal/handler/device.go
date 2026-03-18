package handler

import (
	"io"
	"net/http"
	"strconv"

	"github.com/ZhangPengPaul/LightOTA/internal/model"
	"github.com/ZhangPengPaul/LightOTA/internal/service"
	"github.com/gin-gonic/gin"
)

type DeviceHandler struct {
	service *service.UpgradeService
}

func NewDeviceHandler(service *service.UpgradeService) *DeviceHandler {
	return &DeviceHandler{service: service}
}

func (h *DeviceHandler) Register(r *gin.RouterGroup) {
	r.GET("/check-update", h.checkUpdate)
	r.GET("/download/:firmwareId", h.download)
	r.POST("/report-result", h.reportResult)
	r.GET("/events/:deviceId", h.subscribeEvents)
}

type CheckUpdateRequest struct {
	DeviceID      string `form:"deviceId" binding:"required"`
	CurrentVersion string `form:"currentVersion" binding:"required"`
}

func (h *DeviceHandler) checkUpdate(c *gin.Context) {
	var req CheckUpdateRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error()})
		return
	}

	device, err := h.service.GetDeviceByID(req.DeviceID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "Device not found"})
		return
	}

	firmware, err := h.service.GetLatestActiveFirmware(device.TenantID, device.ProductID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
		return
	}

	// If firmware version name matches what device reports, no update
	if firmware.Version == req.CurrentVersion {
		c.JSON(http.StatusOK, gin.H{"code": 0, "hasUpdate": false})
		return
	}

	scheme := "http"
	if c.Request.TLS != nil {
		scheme = "https"
	}
	downloadURL := scheme + "://" + c.Request.Host + "/api/v1/ota/download/" + firmware.ID + "?deviceId=" + req.DeviceID

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"hasUpdate": true,
		"data": gin.H{
			"version":    firmware.Version,
			"changelog":  firmware.Changelog,
			"downloadUrl": downloadURL,
			"md5":        firmware.MD5,
			"fileSize":   firmware.FileSize,
		},
	})
}

func (h *DeviceHandler) download(c *gin.Context) {
	firmwareID := c.Param("firmwareId")
	deviceID := c.Query("deviceId")
	if deviceID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "deviceId is required"})
		return
	}

	device, err := h.service.GetDeviceByID(deviceID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "Device not found"})
		return
	}

	file, size, err := h.service.GetFirmwareFile(device.TenantID, firmwareID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "Firmware not found"})
		return
	}
	defer file.Close()

	c.Writer.Header().Set("Content-Disposition", "attachment; filename=firmware.bin")
	c.Writer.Header().Set("Content-Length", strconv.FormatInt(size, 10))
	c.Data(http.StatusOK, "application/octet-stream", nil)
	io.Copy(c.Writer, file)
}

type ReportResultRequest struct {
	DeviceID  string                 `json:"deviceId" binding:"required"`
	TaskID    string                 `json:"taskId" binding:"required"`
	Status    model.DeviceUpgradeStatus `json:"status" binding:"required"`
	ErrorMsg  string                 `json:"errorMessage"`
}

func (h *DeviceHandler) reportResult(c *gin.Context) {
	var req ReportResultRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error()})
		return
	}

	// We need to find the record
	// Since service has the method, we can call it directly
	record, err := h.service.FindPendingByTaskAndDevice(req.TaskID, req.DeviceID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "Record not found"})
		return
	}

	err = h.service.UpdateDeviceRecordStatus(record.ID, req.Status, req.ErrorMsg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "ok"})
}

func (h *DeviceHandler) subscribeEvents(c *gin.Context) {
	deviceID := c.Param("deviceId")
	h.service.HandleSSE(c.Writer, c.Request, deviceID)
}
