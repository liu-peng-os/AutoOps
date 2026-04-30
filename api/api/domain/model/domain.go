package model

import "time"

type DomainZone struct {
	ID             uint       `gorm:"primaryKey" json:"id"`
	Name           string     `gorm:"size:255;not null;index:idx_domain_zone_name" json:"name"`
	DisplayName    string     `gorm:"size:255" json:"displayName"`
	Provider       string     `gorm:"size:64" json:"provider"`
	Status         string     `gorm:"size:32;default:unknown" json:"status"`
	ExternalSource string     `gorm:"size:64;not null;index:idx_domain_zone_external,unique" json:"externalSource"`
	ExternalID     string     `gorm:"size:255;not null;index:idx_domain_zone_external,unique" json:"externalId"`
	RawData        string     `gorm:"type:text" json:"rawData"`
	LastSyncedAt   *time.Time `json:"lastSyncedAt"`
	CreatedAt      time.Time  `json:"createdAt"`
	UpdatedAt      time.Time  `json:"updatedAt"`
}

func (DomainZone) TableName() string {
	return "domain_zone"
}

type DomainRecord struct {
	ID             uint       `gorm:"primaryKey" json:"id"`
	ZoneID         uint       `gorm:"not null;index" json:"zoneId"`
	ZoneName       string     `gorm:"size:255;not null;index" json:"zoneName"`
	Name           string     `gorm:"size:255;not null;index" json:"name"`
	Type           string     `gorm:"size:32;not null;index" json:"type"`
	Value          string     `gorm:"type:text;not null" json:"value"`
	Line           string     `gorm:"size:128" json:"line"`
	TTL            int        `json:"ttl"`
	Priority       int        `json:"priority"`
	Status         string     `gorm:"size:32;default:unknown" json:"status"`
	ExternalSource string     `gorm:"size:64;not null;index:idx_domain_record_external,unique" json:"externalSource"`
	ExternalID     string     `gorm:"size:255;not null;index:idx_domain_record_external,unique" json:"externalId"`
	RawData        string     `gorm:"type:text" json:"rawData"`
	LastSyncedAt   *time.Time `json:"lastSyncedAt"`
	CreatedAt      time.Time  `json:"createdAt"`
	UpdatedAt      time.Time  `json:"updatedAt"`
}

func (DomainRecord) TableName() string {
	return "domain_record"
}

type DomainSyncRun struct {
	ID             uint       `gorm:"primaryKey" json:"id"`
	ExternalSource string     `gorm:"size:64;not null;index" json:"externalSource"`
	Status         string     `gorm:"size:32;not null;index" json:"status"`
	ZonesSynced    int        `json:"zonesSynced"`
	RecordsSynced  int        `json:"recordsSynced"`
	Message        string     `gorm:"type:text" json:"message"`
	StartedAt      time.Time  `json:"startedAt"`
	FinishedAt     *time.Time `json:"finishedAt"`
}

func (DomainSyncRun) TableName() string {
	return "domain_sync_run"
}

type DomainZoneListQuery struct {
	Page     int    `form:"page"`
	PageSize int    `form:"pageSize"`
	Keyword  string `form:"keyword"`
}

type DomainRecordListQuery struct {
	Page     int    `form:"page"`
	PageSize int    `form:"pageSize"`
	ZoneID   uint   `form:"zoneId"`
	Keyword  string `form:"keyword"`
	Type     string `form:"type"`
}
