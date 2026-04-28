package model

import (
	"fmt"
	"time"
)

const (
	AnsibleConfigTypeInventory  = 1
	AnsibleConfigTypeGlobalVars = 2
	AnsibleConfigTypeExtraVars  = 3
	AnsibleConfigTypeCLIArgs    = 4
)

var ansibleConfigTypeNames = map[int]string{
	AnsibleConfigTypeInventory:  "ansible_inventory",
	AnsibleConfigTypeGlobalVars: "ansible_global_vars",
	AnsibleConfigTypeExtraVars:  "ansible_extra_vars",
	AnsibleConfigTypeCLIArgs:    "ansible_cli_args",
}

type SysConfig struct {
	ID         uint      `gorm:"primaryKey;column:id"`
	ConfigKey  string    `gorm:"column:config_key"`
	ConfigType string    `gorm:"column:config_type"`
	ConfigData string    `gorm:"column:config_data"`
	Status     int       `gorm:"column:status"`
	Remark     string    `gorm:"column:remark"`
	CreateTime time.Time `gorm:"column:create_time"`
	UpdateTime time.Time `gorm:"column:update_time"`
}

func (SysConfig) TableName() string {
	return "sys_config"
}

type AnsibleConfigRequest struct {
	ID      uint   `json:"id"`
	Name    string `json:"name" binding:"required"`
	Content string `json:"content" binding:"required"`
	Remark  string `json:"remark"`
	Type    int    `json:"type" binding:"required"`
}

type AnsibleConfigVO struct {
	ID        uint      `json:"ID"`
	Name      string    `json:"Name"`
	Content   string    `json:"Content"`
	Remark    string    `json:"Remark"`
	Type      int       `json:"Type"`
	CreatedAt time.Time `json:"CreatedAt"`
	UpdatedAt time.Time `json:"UpdatedAt"`
}

func AnsibleConfigTypeName(configType int) (string, error) {
	name, ok := ansibleConfigTypeNames[configType]
	if !ok {
		return "", fmt.Errorf("unsupported ansible config type: %d", configType)
	}
	return name, nil
}

func AnsibleConfigKey(configType int, name string) (string, error) {
	typeName, err := AnsibleConfigTypeName(configType)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("ansible:%s:%s", typeName, name), nil
}

func AnsibleConfigTypeFromName(typeName string) int {
	for code, name := range ansibleConfigTypeNames {
		if name == typeName {
			return code
		}
	}
	return 0
}

func ToAnsibleConfigVO(config SysConfig) AnsibleConfigVO {
	return AnsibleConfigVO{
		ID:        config.ID,
		Name:      configNameFromKey(config.ConfigKey),
		Content:   config.ConfigData,
		Remark:    config.Remark,
		Type:      AnsibleConfigTypeFromName(config.ConfigType),
		CreatedAt: config.CreateTime,
		UpdatedAt: config.UpdateTime,
	}
}

func configNameFromKey(key string) string {
	parts := 0
	for i := 0; i < len(key); i++ {
		if key[i] == ':' {
			parts++
			if parts == 2 && i+1 < len(key) {
				return key[i+1:]
			}
		}
	}
	return key
}
