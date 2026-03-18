package handler

import (
	"net/http"
	"strconv"

	"github.com/ZhangPengPaul/LightOTA/internal/service"
	"github.com/gin-gonic/gin"
)

type ProductHandler struct {
	service *service.ProductService
}

func NewProductHandler(service *service.ProductService) *ProductHandler {
	return &ProductHandler{service: service}
}

func (h *ProductHandler) Register(r *gin.RouterGroup) {
	group := r.Group("/products")
	{
		group.POST("", h.create)
		group.GET("", h.list)
		group.GET("/:id", h.get)
		group.PUT("/:id", h.update)
		group.DELETE("/:id", h.delete)
	}
}

func (h *ProductHandler) create(c *gin.Context) {
	tenant := GetTenant(c)
	var req service.CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error()})
		return
	}

	product, err := h.service.Create(tenant.ID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "data": product})
}

func (h *ProductHandler) list(c *gin.Context) {
	tenant := GetTenant(c)
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

	products, total, err := h.service.List(tenant.ID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "data": gin.H{"list": products, "total": total}})
}

func (h *ProductHandler) get(c *gin.Context) {
	tenant := GetTenant(c)
	id := c.Param("id")
	product, err := h.service.GetByID(tenant.ID, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "Product not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "data": product})
}

func (h *ProductHandler) update(c *gin.Context) {
	tenant := GetTenant(c)
	id := c.Param("id")
	var req service.UpdateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error()})
		return
	}

	product, err := h.service.Update(tenant.ID, id, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "data": product})
}

func (h *ProductHandler) delete(c *gin.Context) {
	tenant := GetTenant(c)
	id := c.Param("id")
	err := h.service.Delete(tenant.ID, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "deleted"})
}
