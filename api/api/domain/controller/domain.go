package controller

import (
	"net/http"

	"dodevops-api/api/domain/model"
	"dodevops-api/api/domain/service"
	"dodevops-api/common/result"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type DomainController struct {
	service *service.DomainService
}

func NewDomainController() *DomainController {
	return &DomainController{service: service.NewDomainService()}
}

func (c *DomainController) Health(ctx *gin.Context) {
	result.Success(ctx, c.service.HealthCheck(ctx.Request.Context()))
}

func (c *DomainController) Sync(ctx *gin.Context) {
	syncResult, err := c.service.SyncDnsmgr(ctx.Request.Context())
	if err != nil {
		if syncResult != nil {
			result.Failed(ctx, http.StatusOK, syncResult.Message)
			return
		}
		result.Failed(ctx, http.StatusOK, err.Error())
		return
	}
	result.Success(ctx, syncResult)
}

func (c *DomainController) ListZones(ctx *gin.Context) {
	var query model.DomainZoneListQuery
	if err := ctx.ShouldBindQuery(&query); err != nil {
		result.Failed(ctx, http.StatusBadRequest, err.Error())
		return
	}
	zones, total, err := c.service.ListZones(query)
	if err != nil {
		result.Failed(ctx, http.StatusOK, err.Error())
		return
	}
	if query.Page <= 0 {
		query.Page = 1
	}
	if query.PageSize <= 0 {
		query.PageSize = 20
	}
	result.SuccessWithPage(ctx, zones, total, query.Page, query.PageSize)
}

func (c *DomainController) ListRecords(ctx *gin.Context) {
	var query model.DomainRecordListQuery
	if err := ctx.ShouldBindQuery(&query); err != nil {
		result.Failed(ctx, http.StatusBadRequest, err.Error())
		return
	}
	records, total, err := c.service.ListRecords(query)
	if err != nil {
		result.Failed(ctx, http.StatusOK, err.Error())
		return
	}
	if query.Page <= 0 {
		query.Page = 1
	}
	if query.PageSize <= 0 {
		query.PageSize = 20
	}
	result.SuccessWithPage(ctx, records, total, query.Page, query.PageSize)
}

func (c *DomainController) LastSync(ctx *gin.Context) {
	run, err := c.service.LastSync()
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			result.Success(ctx, gin.H{})
			return
		}
		result.Failed(ctx, http.StatusOK, err.Error())
		return
	}
	result.Success(ctx, run)
}
