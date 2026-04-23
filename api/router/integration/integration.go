package integration

import (
	"dodevops-api/api/integration/controller"

	"github.com/gin-gonic/gin"
)

// RegisterIntegrationRoutes 注册外部系统接入骨架路由
func RegisterIntegrationRoutes(router *gin.RouterGroup) {
	router.GET("/integration/catalog", controller.GetIntegrationCatalog)
}
