package controller

import (
	"net/http"

	"dodevops-api/api/domain/provider"
	"dodevops-api/api/domain/service"
	"dodevops-api/common/result"

	"github.com/gin-gonic/gin"
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

func (c *DomainController) ListDomains(ctx *gin.Context) {
	var query service.DomainListQuery
	if err := ctx.ShouldBindQuery(&query); err != nil {
		result.Failed(ctx, http.StatusBadRequest, err.Error())
		return
	}
	domains, err := c.service.ListDomains(ctx.Request.Context(), query)
	if err != nil {
		result.Failed(ctx, http.StatusOK, err.Error())
		return
	}
	result.SuccessWithPage(ctx, domains.List, domains.Total, normalizePage(query.Page), normalizePageSize(query.PageSize, query.Limit))
}

func (c *DomainController) DomainDetail(ctx *gin.Context) {
	domainID := ctx.Param("id")
	loginURL := ctx.Query("loginUrl") == "1" || ctx.Query("loginUrl") == "true"
	detail, err := c.service.DomainDetail(ctx.Request.Context(), domainID, loginURL)
	if err != nil {
		result.Failed(ctx, http.StatusOK, err.Error())
		return
	}
	result.Success(ctx, detail)
}

func (c *DomainController) ListRecords(ctx *gin.Context) {
	var query service.RecordListQuery
	if err := ctx.ShouldBindQuery(&query); err != nil {
		result.Failed(ctx, http.StatusBadRequest, err.Error())
		return
	}
	records, err := c.service.ListRecords(ctx.Request.Context(), query)
	if err != nil {
		result.Failed(ctx, http.StatusOK, err.Error())
		return
	}
	result.SuccessWithPage(ctx, records.List, records.Total, normalizePage(query.Page), normalizePageSize(query.PageSize, query.Limit))
}

func (c *DomainController) AddRecord(ctx *gin.Context) {
	domainID := ctx.Param("domainId")
	var record provider.DnsmgrRecordRequest
	if err := ctx.ShouldBindJSON(&record); err != nil {
		result.Failed(ctx, http.StatusBadRequest, err.Error())
		return
	}
	payload, err := c.service.AddRecord(ctx.Request.Context(), domainID, record)
	if err != nil {
		result.Failed(ctx, http.StatusOK, err.Error())
		return
	}
	result.Success(ctx, payload)
}

func (c *DomainController) UpdateRecord(ctx *gin.Context) {
	domainID := ctx.Param("domainId")
	recordID := ctx.Param("recordId")
	var record provider.DnsmgrRecordRequest
	if err := ctx.ShouldBindJSON(&record); err != nil {
		result.Failed(ctx, http.StatusBadRequest, err.Error())
		return
	}
	payload, err := c.service.UpdateRecord(ctx.Request.Context(), domainID, recordID, record)
	if err != nil {
		result.Failed(ctx, http.StatusOK, err.Error())
		return
	}
	result.Success(ctx, payload)
}

func (c *DomainController) DeleteRecord(ctx *gin.Context) {
	payload, err := c.service.DeleteRecord(ctx.Request.Context(), ctx.Param("domainId"), ctx.Param("recordId"))
	if err != nil {
		result.Failed(ctx, http.StatusOK, err.Error())
		return
	}
	result.Success(ctx, payload)
}

func (c *DomainController) SetRecordStatus(ctx *gin.Context) {
	var request struct {
		Status string `json:"status" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&request); err != nil {
		result.Failed(ctx, http.StatusBadRequest, err.Error())
		return
	}
	payload, err := c.service.SetRecordStatus(ctx.Request.Context(), ctx.Param("domainId"), ctx.Param("recordId"), request.Status)
	if err != nil {
		result.Failed(ctx, http.StatusOK, err.Error())
		return
	}
	result.Success(ctx, payload)
}

func (c *DomainController) SetRecordRemark(ctx *gin.Context) {
	var request struct {
		Remark string `json:"remark"`
	}
	if err := ctx.ShouldBindJSON(&request); err != nil {
		result.Failed(ctx, http.StatusBadRequest, err.Error())
		return
	}
	payload, err := c.service.SetRecordRemark(ctx.Request.Context(), ctx.Param("domainId"), ctx.Param("recordId"), request.Remark)
	if err != nil {
		result.Failed(ctx, http.StatusOK, err.Error())
		return
	}
	result.Success(ctx, payload)
}

func (c *DomainController) BatchRecords(ctx *gin.Context) {
	domainID := ctx.Param("domainId")
	var request provider.DnsmgrBatchRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		result.Failed(ctx, http.StatusBadRequest, err.Error())
		return
	}
	payload, err := c.service.BatchRecords(ctx.Request.Context(), domainID, request)
	if err != nil {
		result.Failed(ctx, http.StatusOK, err.Error())
		return
	}
	result.Success(ctx, payload)
}

func normalizePage(page int) int {
	if page <= 0 {
		return 1
	}
	return page
}

func normalizePageSize(pageSize, limit int) int {
	if pageSize > 0 {
		return pageSize
	}
	if limit > 0 {
		return limit
	}
	return 20
}
