package service

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"dodevops-api/api/domain/provider"
	"dodevops-api/common/config"
)

const sourceDnsmgr = "dnsmgr"

type DomainService struct{}

type HealthStatus struct {
	Source  string `json:"source"`
	Enabled bool   `json:"enabled"`
	Healthy bool   `json:"healthy"`
	Message string `json:"message"`
}

type DomainListQuery struct {
	Page     int    `form:"page"`
	PageSize int    `form:"pageSize"`
	Offset   int    `form:"offset"`
	Limit    int    `form:"limit"`
	Keyword  string `form:"keyword"`
}

type RecordListQuery struct {
	DomainID string `form:"domainId" binding:"required"`
	Page     int    `form:"page"`
	PageSize int    `form:"pageSize"`
	Offset   int    `form:"offset"`
	Limit    int    `form:"limit"`
	Keyword  string `form:"keyword"`
	Type     string `form:"type"`
}

func NewDomainService() *DomainService {
	return &DomainService{}
}

func (s *DomainService) HealthCheck(ctx context.Context) HealthStatus {
	client, enabled, err := s.client()
	if err != nil {
		return HealthStatus{Source: sourceDnsmgr, Enabled: enabled, Healthy: false, Message: err.Error()}
	}
	if err := client.HealthCheck(ctx); err != nil {
		return HealthStatus{Source: sourceDnsmgr, Enabled: true, Healthy: false, Message: err.Error()}
	}
	return HealthStatus{Source: sourceDnsmgr, Enabled: true, Healthy: true, Message: "dnsmgr is reachable"}
}

func (s *DomainService) ListDomains(ctx context.Context, query DomainListQuery) (*provider.DnsmgrPageResult, error) {
	client, _, err := s.client()
	if err != nil {
		return nil, err
	}
	values := buildPageValues(query.Page, query.PageSize, query.Offset, query.Limit, query.Keyword, "")
	if query.Keyword != "" {
		values.Set("kw", query.Keyword)
	}
	return client.ListDomains(ctx, values)
}

func (s *DomainService) DomainDetail(ctx context.Context, domainID string, loginURL bool) (map[string]interface{}, error) {
	client, _, err := s.client()
	if err != nil {
		return nil, err
	}
	return client.DomainDetail(ctx, domainID, loginURL)
}

func (s *DomainService) ListRecords(ctx context.Context, query RecordListQuery) (*provider.DnsmgrPageResult, error) {
	client, _, err := s.client()
	if err != nil {
		return nil, err
	}
	values := buildPageValues(query.Page, query.PageSize, query.Offset, query.Limit, query.Keyword, query.Type)
	return client.ListRecords(ctx, query.DomainID, values)
}

func (s *DomainService) AddRecord(ctx context.Context, domainID string, record provider.DnsmgrRecordRequest) (interface{}, error) {
	client, _, err := s.client()
	if err != nil {
		return nil, err
	}
	return client.AddRecord(ctx, domainID, record)
}

func (s *DomainService) UpdateRecord(ctx context.Context, domainID, recordID string, record provider.DnsmgrRecordRequest) (interface{}, error) {
	client, _, err := s.client()
	if err != nil {
		return nil, err
	}
	return client.UpdateRecord(ctx, domainID, recordID, record)
}

func (s *DomainService) DeleteRecord(ctx context.Context, domainID, recordID string) (interface{}, error) {
	client, _, err := s.client()
	if err != nil {
		return nil, err
	}
	return client.DeleteRecord(ctx, domainID, recordID)
}

func (s *DomainService) SetRecordStatus(ctx context.Context, domainID, recordID, status string) (interface{}, error) {
	client, _, err := s.client()
	if err != nil {
		return nil, err
	}
	return client.SetRecordStatus(ctx, domainID, recordID, status)
}

func (s *DomainService) SetRecordRemark(ctx context.Context, domainID, recordID, remark string) (interface{}, error) {
	client, _, err := s.client()
	if err != nil {
		return nil, err
	}
	return client.SetRecordRemark(ctx, domainID, recordID, remark)
}

func (s *DomainService) BatchRecords(ctx context.Context, domainID string, batch provider.DnsmgrBatchRequest) (interface{}, error) {
	client, _, err := s.client()
	if err != nil {
		return nil, err
	}
	return client.BatchRecords(ctx, domainID, batch)
}

func (s *DomainService) client() (*provider.DnsmgrClient, bool, error) {
	system, ok := findDnsmgrSystem()
	if !ok {
		return nil, false, fmt.Errorf("dnsmgr integration is not configured")
	}
	if !system.Enabled {
		return nil, false, fmt.Errorf("dnsmgr integration is not enabled")
	}
	return provider.NewDnsmgrClient(system), true, nil
}

func buildPageValues(page, pageSize, offset, limit int, keyword, recordType string) url.Values {
	if limit <= 0 {
		limit = pageSize
	}
	if limit <= 0 {
		limit = 20
	}
	if offset <= 0 && page > 1 {
		offset = (page - 1) * limit
	}
	values := url.Values{
		"offset": []string{fmt.Sprintf("%d", offset)},
		"limit":  []string{fmt.Sprintf("%d", limit)},
	}
	if keyword != "" {
		values.Set("keyword", keyword)
		values.Set("search", keyword)
	}
	if recordType != "" {
		values.Set("type", strings.ToUpper(recordType))
	}
	return values
}

func findDnsmgrSystem() (config.ExternalSystem, bool) {
	if config.Config == nil {
		return config.ExternalSystem{}, false
	}
	for _, system := range config.Config.Integrations.Systems {
		if strings.EqualFold(system.Provider, sourceDnsmgr) || strings.EqualFold(system.Key, sourceDnsmgr) || strings.EqualFold(system.Key, "domain-management") {
			return system, true
		}
	}
	return config.ExternalSystem{}, false
}
