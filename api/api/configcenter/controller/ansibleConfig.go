package controller

import (
	"strconv"

	"dodevops-api/api/configcenter/model"
	"dodevops-api/api/configcenter/service"
	"dodevops-api/common/result"

	"github.com/gin-gonic/gin"
)

type AnsibleConfigController struct {
	service *service.AnsibleConfigService
}

func NewAnsibleConfigController() *AnsibleConfigController {
	return &AnsibleConfigController{
		service: service.NewAnsibleConfigService(),
	}
}

func (c *AnsibleConfigController) Create(ctx *gin.Context) {
	var req model.AnsibleConfigRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		result.Failed(ctx, 400, err.Error())
		return
	}

	config, err := c.service.Create(req)
	if err != nil {
		result.Failed(ctx, 500, err.Error())
		return
	}
	result.Success(ctx, config)
}

func (c *AnsibleConfigController) Update(ctx *gin.Context) {
	id, ok := parseIDParam(ctx)
	if !ok {
		return
	}

	var req model.AnsibleConfigRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		result.Failed(ctx, 400, err.Error())
		return
	}

	config, err := c.service.Update(id, req)
	if err != nil {
		result.Failed(ctx, 500, err.Error())
		return
	}
	result.Success(ctx, config)
}

func (c *AnsibleConfigController) Delete(ctx *gin.Context) {
	id, ok := parseIDParam(ctx)
	if !ok {
		return
	}

	if err := c.service.Delete(id); err != nil {
		result.Failed(ctx, 500, err.Error())
		return
	}
	result.Success(ctx, nil)
}

func (c *AnsibleConfigController) GetByID(ctx *gin.Context) {
	id, ok := parseIDParam(ctx)
	if !ok {
		return
	}

	config, err := c.service.GetByID(id)
	if err != nil {
		result.Failed(ctx, 500, err.Error())
		return
	}
	result.Success(ctx, config)
}

func (c *AnsibleConfigController) List(ctx *gin.Context) {
	page := parsePositiveInt(ctx.DefaultQuery("page", "1"), 1)
	pageSize := parsePositiveInt(ctx.DefaultQuery("pageSize", ctx.DefaultQuery("size", "10")), 10)
	if pageSize > 100 {
		pageSize = 100
	}
	configType := parseNonNegativeInt(ctx.DefaultQuery("type", "0"), 0)
	name := ctx.DefaultQuery("name", "")

	list, total, err := c.service.List(page, pageSize, configType, name)
	if err != nil {
		result.Failed(ctx, 500, err.Error())
		return
	}
	result.SuccessWithPage(ctx, list, total, page, pageSize)
}

func parseIDParam(ctx *gin.Context) (uint, bool) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil || id == 0 {
		result.Failed(ctx, 400, "invalid id")
		return 0, false
	}
	return uint(id), true
}

func parsePositiveInt(value string, fallback int) int {
	parsed, err := strconv.Atoi(value)
	if err != nil || parsed <= 0 {
		return fallback
	}
	return parsed
}

func parseNonNegativeInt(value string, fallback int) int {
	parsed, err := strconv.Atoi(value)
	if err != nil || parsed < 0 {
		return fallback
	}
	return parsed
}
