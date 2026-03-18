package handler

import (
	"net/http"
	"strconv"

	"github.com/ZhangPengPaul/LightOTA/internal/service"
	"github.com/gin-gonic/gin"
)

type TenantHandler struct {
	service *service.TenantService
}

func NewTenantHandler(service *service.TenantService) *TenantHandler {
	return &TenantHandler{service: service}
}

func (h *TenantHandler) Register(r *gin.RouterGroup) {
	group := r.Group("/tenants")
	{
		group.POST("", h.create)
		group.GET("", h.list)
		group.GET("/:id", h.get)
		group.PUT("/:id", h.update)
	}
}

func (h *TenantHandler) create(c *gin.Context) {
	var req service.CreateTenantRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error()})
		return
	}

	tenant, err := h.service.Create(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "data": tenant})
}

func (h *TenantHandler) list(c *gin.Context) {
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

	tenants, total, err := h.service.List(limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "data": gin.H{"list": tenants, "total": total}})
}

func (h *TenantHandler) get(c *gin.Context) {
	id := c.Param("id")
	tenant, err := h.service.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "Tenant not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "data": tenant})
}

func (h *TenantHandler) update(c *gin.Context) {
	id := c.Param("id")
	var req service.UpdateTenantRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error()})
		return
	}

	tenant, err := h.service.Update(id, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "data": tenant})
}
