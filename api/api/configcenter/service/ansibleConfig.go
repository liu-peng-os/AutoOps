package service

import (
	"errors"
	"strings"

	"dodevops-api/api/configcenter/dao"
	"dodevops-api/api/configcenter/model"
)

var ErrEmptyAnsibleConfigName = errors.New("ansible config name cannot be empty")

type AnsibleConfigService struct {
	dao *dao.AnsibleConfigDao
}

func NewAnsibleConfigService() *AnsibleConfigService {
	return &AnsibleConfigService{
		dao: dao.NewAnsibleConfigDao(),
	}
}

func (s *AnsibleConfigService) Create(req model.AnsibleConfigRequest) (*model.AnsibleConfigVO, error) {
	config, err := s.buildSysConfig(req)
	if err != nil {
		return nil, err
	}
	if err := s.dao.Create(config); err != nil {
		return nil, err
	}
	vo := model.ToAnsibleConfigVO(*config)
	return &vo, nil
}

func (s *AnsibleConfigService) Update(id uint, req model.AnsibleConfigRequest) (*model.AnsibleConfigVO, error) {
	config, err := s.buildSysConfig(req)
	if err != nil {
		return nil, err
	}
	if err := s.dao.Update(id, config); err != nil {
		return nil, err
	}
	updated, err := s.dao.GetByID(id)
	if err != nil {
		return nil, err
	}
	vo := model.ToAnsibleConfigVO(*updated)
	return &vo, nil
}

func (s *AnsibleConfigService) Delete(id uint) error {
	return s.dao.Delete(id)
}

func (s *AnsibleConfigService) GetByID(id uint) (*model.AnsibleConfigVO, error) {
	config, err := s.dao.GetByID(id)
	if err != nil {
		return nil, err
	}
	vo := model.ToAnsibleConfigVO(*config)
	return &vo, nil
}

func (s *AnsibleConfigService) List(page, pageSize, configType int, name string) ([]model.AnsibleConfigVO, int64, error) {
	configs, total, err := s.dao.List(page, pageSize, configType, name)
	if err != nil {
		return nil, 0, err
	}

	list := make([]model.AnsibleConfigVO, 0, len(configs))
	for _, config := range configs {
		list = append(list, model.ToAnsibleConfigVO(config))
	}
	return list, total, nil
}

func (s *AnsibleConfigService) buildSysConfig(req model.AnsibleConfigRequest) (*model.SysConfig, error) {
	name := strings.TrimSpace(req.Name)
	if name == "" {
		return nil, ErrEmptyAnsibleConfigName
	}
	key, err := model.AnsibleConfigKey(req.Type, name)
	if err != nil {
		return nil, err
	}
	typeName, err := model.AnsibleConfigTypeName(req.Type)
	if err != nil {
		return nil, err
	}

	return &model.SysConfig{
		ConfigKey:  key,
		ConfigType: typeName,
		ConfigData: req.Content,
		Remark:     req.Remark,
	}, nil
}
