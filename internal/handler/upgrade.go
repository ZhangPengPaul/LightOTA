package handler

import (
	"net/http"
	"strconv"

	"github.com/ZhangPengPaul/LightOTA/internal/model"
	"github.com/ZhangPengPaul/LightOTA/internal/service"
	"github.com/gin-gonic/gin"
)

type UpgradeHandler struct {
	service *service.UpgradeService
}

func NewUpgradeHandler(service *service.UpgradeService) *UpgradeHandler {
	return &UpgradeHandler{service: service}
}

func (h *UpgradeHandler) Register(r *gin.RouterGroup) {
	group := r.Group("/upgrade")
	{
		group.POST("/task", h.createTask)
		group.GET("/tasks", h.listTasks)
		group.GET("/task/:id", h.getTask)
	}
}

func (h *UpgradeHandler) createTask(c *gin.Context) {
	tenant := GetTenant(c)
	var req service.CreateUpgradeTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error()})
		return
	}

	task, err := h.service.CreateTask(tenant.ID, &req, tenant)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
		return
	}

	err = h.service.StartTask(task)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
		return
	}

	go h.processPendingTasks(task)

	c.JSON(http.StatusOK, gin.H{"code": 0, "data": gin.H{"taskId": task.ID}})
}

func (h *UpgradeHandler) listTasks(c *gin.Context) {
	tenant := GetTenant(c)
	productID := c.Query("product_id")
	limit := 10
	offset := 0
	if queryLimit := c.Query("limit"); queryLimit != "" {
		if l, err := strconv.Atoi(queryLimit); err == nil && l > 0 {
			limit = l
		}
	}
	if queryOffset := c.Query("offset"); queryOffset != "" {
		if o, err := strconv.Atoi(queryOffset); err == nil && o >= 0 {
			offset = o
		}
	}

	var pid *string
	if productID != "" {
		pid = &productID
	}

	tasks, total, err := h.service.ListTasks(tenant.ID, pid, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "data": gin.H{"list": tasks, "total": total}})
}

func (h *UpgradeHandler) getTask(c *gin.Context) {
	tenant := GetTenant(c)
	id := c.Param("id")
	task, counts, err := h.service.GetTask(tenant.ID, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "Task not found"})
		return
	}

	total := 0
	for _, count := range counts {
		total += count
	}
	success := counts[model.DeviceStatusSuccess]
	failed := counts[model.DeviceStatusFailed]
	pending := 0
	for status, count := range counts {
		if status != model.DeviceStatusSuccess && status != model.DeviceStatusFailed {
			pending += count
		}
	}

	var percent int
	if total > 0 {
		percent = (success + failed) * 100 / total
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "data": gin.H{
		"task":          task,
		"total":         total,
		"successCount":  success,
		"failedCount":   failed,
		"pendingCount":  pending,
		"percent":       percent,
		"statusCounts":  counts,
	}})
}

func (h *UpgradeHandler) processPendingTasks(task *model.UpgradeTask) {
	for {
		if task.Status != model.TaskStatusRunning {
			break
		}

		records, err := h.service.FindPendingRecords(task.ID, task.PushRate)
		if err != nil || len(records) == 0 {
			break
		}

		for _, record := range records {
			firmware, err := h.service.GetFirmwareByID(task.TenantID, task.FirmwareID)
			if err != nil {
				continue
			}

			device, err := h.service.GetDeviceByID(record.DeviceID)
			if err != nil {
				continue
			}

			record.NewVersion = firmware.Version
			h.service.PushNotification(&record, firmware, device)
		}

		if len(records) < task.PushRate {
			break
		}
	}

	stats, err := h.service.GetStats(task.ID)
	if err != nil {
		return
	}

	total := 0
	pending := 0
	for status, count := range stats {
		total += count
		if status == model.DeviceStatusPending || status == model.DeviceStatusNotified {
			pending += count
		}
	}

	if pending == 0 && total > 0 {
		h.service.CompleteTask(task)
	}
}
