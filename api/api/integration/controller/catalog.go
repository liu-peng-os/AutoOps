package controller

import (
	"dodevops-api/api/integration/service"
	"dodevops-api/common/result"

	"github.com/gin-gonic/gin"
)

// GetIntegrationCatalog 返回当前阶段的外部接入目录和首批模块信息
func GetIntegrationCatalog(c *gin.Context) {
	result.Success(c, service.BuildCatalog())
}
