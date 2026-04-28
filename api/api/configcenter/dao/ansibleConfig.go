package dao

import (
	"strings"
	"time"

	"dodevops-api/api/configcenter/model"
	"dodevops-api/common"
)

type AnsibleConfigDao struct{}

func NewAnsibleConfigDao() *AnsibleConfigDao {
	return &AnsibleConfigDao{}
}

func (d *AnsibleConfigDao) Create(config *model.SysConfig) error {
	now := time.Now()
	config.Status = 1
	config.CreateTime = now
	config.UpdateTime = now
	return common.GetDB().Create(config).Error
}

func (d *AnsibleConfigDao) Update(id uint, config *model.SysConfig) error {
	return common.GetDB().Model(&model.SysConfig{}).Where("id = ? AND config_key LIKE ?", id, "ansible:%").Updates(map[string]interface{}{
		"config_key":  config.ConfigKey,
		"config_type": config.ConfigType,
		"config_data": config.ConfigData,
		"remark":      config.Remark,
		"update_time": time.Now(),
	}).Error
}

func (d *AnsibleConfigDao) Delete(id uint) error {
	return common.GetDB().Where("config_key LIKE ?", "ansible:%").Delete(&model.SysConfig{}, id).Error
}

func (d *AnsibleConfigDao) GetByID(id uint) (*model.SysConfig, error) {
	var config model.SysConfig
	err := common.GetDB().Where("config_key LIKE ?", "ansible:%").First(&config, id).Error
	return &config, err
}

func (d *AnsibleConfigDao) List(page, pageSize, configType int, name string) ([]model.SysConfig, int64, error) {
	var configs []model.SysConfig
	var total int64

	db := common.GetDB().Model(&model.SysConfig{}).Where("config_key LIKE ?", "ansible:%")
	if configType > 0 {
		typeName, err := model.AnsibleConfigTypeName(configType)
		if err != nil {
			return nil, 0, err
		}
		db = db.Where("config_type = ?", typeName)
	}
	if strings.TrimSpace(name) != "" {
		db = db.Where("config_key ILIKE ?", "%"+strings.TrimSpace(name)+"%")
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err := db.Order("id DESC").Offset(offset).Limit(pageSize).Find(&configs).Error
	return configs, total, err
}
