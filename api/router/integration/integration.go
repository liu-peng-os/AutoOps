package integration

import (
	domainController "dodevops-api/api/domain/controller"
	"dodevops-api/api/integration/controller"

	"github.com/gin-gonic/gin"
)

// RegisterIntegrationRoutes 注册外部系统接入骨架路由
func RegisterIntegrationRoutes(router *gin.RouterGroup) {
	router.GET("/integration/catalog", controller.GetIntegrationCatalog)

	domainCtrl := domainController.NewDomainController()
	domainGroup := router.Group("/domain")
	{
		domainGroup.GET("/health", domainCtrl.Health)
		domainGroup.POST("/sync", domainCtrl.Sync)
		domainGroup.GET("/sync/last", domainCtrl.LastSync)
		domainGroup.GET("/zones", domainCtrl.ListZones)
		domainGroup.GET("/records", domainCtrl.ListRecords)
	}
}
