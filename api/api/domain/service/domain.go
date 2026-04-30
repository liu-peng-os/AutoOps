package service

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"dodevops-api/api/domain/dao"
	"dodevops-api/api/domain/model"
	"dodevops-api/api/domain/provider"
	"dodevops-api/common/config"

	"gorm.io/gorm"
)

const sourceDnsmgr = "dnsmgr"

type DomainService struct {
	dao *dao.DomainDao
}

type HealthStatus struct {
	Source  string `json:"source"`
	Enabled bool   `json:"enabled"`
	Healthy bool   `json:"healthy"`
	Message string `json:"message"`
}

type SyncResult struct {
	Source        string    `json:"source"`
	Status        string    `json:"status"`
	ZonesSynced   int       `json:"zonesSynced"`
	RecordsSynced int       `json:"recordsSynced"`
	Message       string    `json:"message"`
	StartedAt     time.Time `json:"startedAt"`
	FinishedAt    time.Time `json:"finishedAt"`
}

func NewDomainService() *DomainService {
	return &DomainService{dao: dao.NewDomainDao()}
}

func (s *DomainService) HealthCheck(ctx context.Context) HealthStatus {
	system, ok := findDnsmgrSystem()
	if !ok || !system.Enabled {
		return HealthStatus{
			Source:  sourceDnsmgr,
			Enabled: false,
			Healthy: false,
			Message: "dnsmgr integration is not enabled",
		}
	}

	client := provider.NewDnsmgrClient(system)
	if err := client.HealthCheck(ctx); err != nil {
		return HealthStatus{
			Source:  sourceDnsmgr,
			Enabled: true,
			Healthy: false,
			Message: err.Error(),
		}
	}
	return HealthStatus{
		Source:  sourceDnsmgr,
		Enabled: true,
		Healthy: true,
		Message: "dnsmgr is reachable",
	}
}

func (s *DomainService) SyncDnsmgr(ctx context.Context) (*SyncResult, error) {
	system, ok := findDnsmgrSystem()
	if !ok || !system.Enabled {
		return nil, fmt.Errorf("dnsmgr integration is not enabled")
	}

	client := provider.NewDnsmgrClient(system)
	startedAt := time.Now()
	run := &model.DomainSyncRun{
		ExternalSource: sourceDnsmgr,
		Status:         "running",
		StartedAt:      startedAt,
	}
	_ = s.dao.CreateSyncRun(run)

	zones, err := client.ListZones(ctx)
	if err != nil {
		return s.failRun(run, startedAt, err)
	}

	now := time.Now()
	zonesSynced := 0
	recordsSynced := 0
	for _, zone := range zones {
		dbZone, err := s.upsertZone(zone, now)
		if err != nil {
			return s.failRun(run, startedAt, err)
		}
		zonesSynced++

		records, err := client.ListRecords(ctx, zone)
		if err != nil {
			return s.failRun(run, startedAt, fmt.Errorf("sync records for %s failed: %w", zone.Name, err))
		}
		for _, record := range records {
			if err := s.upsertRecord(dbZone.ID, record, now); err != nil {
				return s.failRun(run, startedAt, err)
			}
			recordsSynced++
		}
	}

	finishedAt := time.Now()
	run.Status = "success"
	run.ZonesSynced = zonesSynced
	run.RecordsSynced = recordsSynced
	run.Message = "sync completed"
	run.FinishedAt = &finishedAt
	_ = s.dao.UpdateSyncRun(run)

	return &SyncResult{
		Source:        sourceDnsmgr,
		Status:        "success",
		ZonesSynced:   zonesSynced,
		RecordsSynced: recordsSynced,
		Message:       "sync completed",
		StartedAt:     startedAt,
		FinishedAt:    finishedAt,
	}, nil
}

func (s *DomainService) ListZones(query model.DomainZoneListQuery) ([]model.DomainZone, int64, error) {
	return s.dao.ListZones(query)
}

func (s *DomainService) ListRecords(query model.DomainRecordListQuery) ([]model.DomainRecord, int64, error) {
	return s.dao.ListRecords(query)
}

func (s *DomainService) LastSync() (*model.DomainSyncRun, error) {
	return s.dao.LastSyncRun(sourceDnsmgr)
}

func (s *DomainService) upsertZone(zone provider.Zone, syncedAt time.Time) (*model.DomainZone, error) {
	existing, err := s.dao.FindZoneByExternal(sourceDnsmgr, zone.ExternalID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	if existing == nil {
		existing = &model.DomainZone{
			ExternalSource: sourceDnsmgr,
			ExternalID:     zone.ExternalID,
		}
	}

	existing.Name = zone.Name
	existing.DisplayName = zone.DisplayName
	existing.Provider = zone.Provider
	existing.Status = zone.Status
	existing.RawData = zone.RawData
	existing.LastSyncedAt = &syncedAt
	if err := s.dao.UpsertZone(existing); err != nil {
		return nil, err
	}
	return existing, nil
}

func (s *DomainService) upsertRecord(zoneID uint, record provider.Record, syncedAt time.Time) error {
	existing, err := s.dao.FindRecordByExternal(sourceDnsmgr, record.ExternalID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	if existing == nil {
		existing = &model.DomainRecord{
			ExternalSource: sourceDnsmgr,
			ExternalID:     record.ExternalID,
		}
	}

	existing.ZoneID = zoneID
	existing.ZoneName = record.ZoneName
	existing.Name = record.Name
	existing.Type = record.Type
	existing.Value = record.Value
	existing.Line = record.Line
	existing.TTL = record.TTL
	existing.Priority = record.Priority
	existing.Status = record.Status
	existing.RawData = record.RawData
	existing.LastSyncedAt = &syncedAt
	return s.dao.UpsertRecord(existing)
}

func (s *DomainService) failRun(run *model.DomainSyncRun, startedAt time.Time, err error) (*SyncResult, error) {
	finishedAt := time.Now()
	if run != nil {
		run.Status = "failed"
		run.Message = err.Error()
		run.FinishedAt = &finishedAt
		_ = s.dao.UpdateSyncRun(run)
	}
	return &SyncResult{
		Source:     sourceDnsmgr,
		Status:     "failed",
		Message:    err.Error(),
		StartedAt:  startedAt,
		FinishedAt: finishedAt,
	}, err
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
