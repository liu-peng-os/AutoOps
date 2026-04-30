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
		domainGroup.GET("/domains", domainCtrl.ListDomains)
		domainGroup.GET("/domains/:id", domainCtrl.DomainDetail)
		domainGroup.GET("/records", domainCtrl.ListRecords)
		domainGroup.POST("/domains/:domainId/records", domainCtrl.AddRecord)
		domainGroup.PUT("/domains/:domainId/records/:recordId", domainCtrl.UpdateRecord)
		domainGroup.DELETE("/domains/:domainId/records/:recordId", domainCtrl.DeleteRecord)
		domainGroup.PUT("/domains/:domainId/records/:recordId/status", domainCtrl.SetRecordStatus)
		domainGroup.PUT("/domains/:domainId/records/:recordId/remark", domainCtrl.SetRecordRemark)
		domainGroup.POST("/domains/:domainId/records/batch", domainCtrl.BatchRecords)
	}
}
