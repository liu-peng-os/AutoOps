package service

import "dodevops-api/common/config"

type ModuleCatalogItem struct {
	Key               string   `json:"key"`
	MenuName          string   `json:"menuName"`
	EntryPath         string   `json:"entryPath"`
	Source            string   `json:"source"`
	Status            string   `json:"status"`
	SupportsWriteback bool     `json:"supportsWriteback"`
	Capabilities      []string `json:"capabilities"`
	Notes             string   `json:"notes"`
}

type ExternalSystemItem struct {
	Key          string   `json:"key"`
	DisplayName  string   `json:"displayName"`
	Category     string   `json:"category"`
	Provider     string   `json:"provider"`
	Mode         string   `json:"mode"`
	BaseURL      string   `json:"baseUrl"`
	Enabled      bool     `json:"enabled"`
	Capabilities []string `json:"capabilities"`
}

type Catalog struct {
	Stage      string               `json:"stage"`
	Principles []string             `json:"principles"`
	Modules    []ModuleCatalogItem  `json:"modules"`
	Systems    []ExternalSystemItem `json:"systems"`
}

func BuildCatalog() Catalog {
	catalog := Catalog{
		Stage: "stage-1",
		Principles: []string{
			"所有外部系统都作为可替换依赖处理，而不是平台核心模型本身。",
			"平台内部只保留通用业务字段、外部映射关系、同步状态和审计信息。",
			"新增模块优先按统一入口、适配器隔离、同步任务复用的方式演进。",
		},
		Modules: []ModuleCatalogItem{
			{
				Key:               "domain-management",
				MenuName:          "域名管理",
				EntryPath:         "/integration/domain",
				Source:            "external-system",
				Status:            "planned",
				SupportsWriteback: false,
				Capabilities:      []string{"统一入口", "后续可增强汇总视图", "保留替换空间"},
				Notes:             "短期先作为 dnsmgr 统一接入口，后续再按需要增强平台内汇总能力。",
			},
			{
				Key:               "operations-workorder",
				MenuName:          "运营工单",
				EntryPath:         "/integration/workorder",
				Source:            "external-master-data",
				Status:            "planned",
				SupportsWriteback: true,
				Capabilities:      []string{"同步", "展示", "基础编辑", "回写", "审计"},
				Notes:             "长期以外部表格为主数据源，平台提供统一查询与回写入口。",
			},
			{
				Key:               "site-management",
				MenuName:          "站点管理",
				EntryPath:         "/cmdb/site",
				Source:            "external-master-data",
				Status:            "planned",
				SupportsWriteback: true,
				Capabilities:      []string{"同步", "展示", "基础编辑", "回写", "审计"},
				Notes:             "挂在资产管理下，长期以外部表格为主数据源。",
			},
		},
	}

	if config.Config == nil {
		return catalog
	}

	for _, system := range config.Config.Integrations.Systems {
		catalog.Systems = append(catalog.Systems, ExternalSystemItem{
			Key:          system.Key,
			DisplayName:  system.DisplayName,
			Category:     system.Category,
			Provider:     system.Provider,
			Mode:         system.Mode,
			BaseURL:      system.BaseURL,
			Enabled:      system.Enabled,
			Capabilities: append([]string{}, system.Capabilities...),
		})
	}

	return catalog
}
