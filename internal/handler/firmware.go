package handler

import (
	"io"
	"net/http"
	"strconv"

	"github.com/ZhangPengPaul/LightOTA/internal/service"
	"github.com/gin-gonic/gin"
)

type FirmwareHandler struct {
	service *service.FirmwareService
}

func NewFirmwareHandler(service *service.FirmwareService) *FirmwareHandler {
	return &FirmwareHandler{service: service}
}

func (h *FirmwareHandler) Register(r *gin.RouterGroup) {
	group := r.Group("/firmwares")
	{
		group.POST("", h.create)
		group.GET("", h.list)
		group.GET("/:id", h.get)
		group.DELETE("/:id", h.delete)
		group.GET("/:id/download", h.download)
	}
}

func (h *FirmwareHandler) create(c *gin.Context) {
	tenant := GetTenant(c)

	var req service.CreateFirmwareRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error()})
		return
	}

	fileHeader, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "file is required"})
		return
	}

	file, err := fileHeader.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
		return
	}
	defer file.Close()

	firmware, err := h.service.Create(tenant.ID, &req, file, fileHeader.Filename)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "data": firmware})
}

func (h *FirmwareHandler) list(c *gin.Context) {
	tenant := GetTenant(c)
	productID := c.Query("product_id")
	if productID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "product_id is required"})
		return
	}

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

	firmwares, total, err := h.service.List(tenant.ID, productID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "data": gin.H{"list": firmwares, "total": total}})
}

func (h *FirmwareHandler) get(c *gin.Context) {
	tenant := GetTenant(c)
	id := c.Param("id")
	firmware, err := h.service.GetByID(tenant.ID, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "Firmware not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "data": firmware})
}

func (h *FirmwareHandler) delete(c *gin.Context) {
	tenant := GetTenant(c)
	id := c.Param("id")
	err := h.service.Delete(tenant.ID, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "deleted"})
}

func (h *FirmwareHandler) download(c *gin.Context) {
	tenant := GetTenant(c)
	id := c.Param("id")
	file, size, err := h.service.GetFile(tenant.ID, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "Firmware not found"})
		return
	}
	defer file.Close()

	firmware, err := h.service.GetByID(tenant.ID, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "Firmware not found"})
		return
	}

	filename := firmware.Version + "_" + firmware.ID + ".bin"
	c.Writer.Header().Set("Content-Disposition", "attachment; filename="+filename)
	c.Writer.Header().Set("Content-Length", strconv.FormatInt(size, 10))
	c.Data(http.StatusOK, "application/octet-stream", nil)
	io.Copy(c.Writer, file)
}
