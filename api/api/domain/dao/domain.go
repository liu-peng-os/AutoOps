package dao

import (
	"dodevops-api/api/domain/model"
	"dodevops-api/common"

	"gorm.io/gorm"
)

type DomainDao struct {
	db *gorm.DB
}

func NewDomainDao() *DomainDao {
	return &DomainDao{db: common.GetDB()}
}

func (d *DomainDao) UpsertZone(zone *model.DomainZone) error {
	return d.db.Save(zone).Error
}

func (d *DomainDao) FindZoneByExternal(source, externalID string) (*model.DomainZone, error) {
	var zone model.DomainZone
	err := d.db.Where("external_source = ? AND external_id = ?", source, externalID).First(&zone).Error
	if err != nil {
		return nil, err
	}
	return &zone, nil
}

func (d *DomainDao) UpsertRecord(record *model.DomainRecord) error {
	return d.db.Save(record).Error
}

func (d *DomainDao) FindRecordByExternal(source, externalID string) (*model.DomainRecord, error) {
	var record model.DomainRecord
	err := d.db.Where("external_source = ? AND external_id = ?", source, externalID).First(&record).Error
	if err != nil {
		return nil, err
	}
	return &record, nil
}

func (d *DomainDao) ListZones(query model.DomainZoneListQuery) ([]model.DomainZone, int64, error) {
	page, pageSize := normalizePage(query.Page, query.PageSize)
	db := d.db.Model(&model.DomainZone{})
	if query.Keyword != "" {
		keyword := "%" + query.Keyword + "%"
		db = db.Where("name LIKE ? OR display_name LIKE ?", keyword, keyword)
	}

	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var zones []model.DomainZone
	err := db.Order("updated_at DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&zones).Error
	return zones, total, err
}

func (d *DomainDao) ListRecords(query model.DomainRecordListQuery) ([]model.DomainRecord, int64, error) {
	page, pageSize := normalizePage(query.Page, query.PageSize)
	db := d.db.Model(&model.DomainRecord{})
	if query.ZoneID > 0 {
		db = db.Where("zone_id = ?", query.ZoneID)
	}
	if query.Keyword != "" {
		keyword := "%" + query.Keyword + "%"
		db = db.Where("name LIKE ? OR value LIKE ? OR zone_name LIKE ?", keyword, keyword, keyword)
	}
	if query.Type != "" {
		db = db.Where("type = ?", query.Type)
	}

	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var records []model.DomainRecord
	err := db.Order("updated_at DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&records).Error
	return records, total, err
}

func (d *DomainDao) CreateSyncRun(run *model.DomainSyncRun) error {
	return d.db.Create(run).Error
}

func (d *DomainDao) UpdateSyncRun(run *model.DomainSyncRun) error {
	return d.db.Save(run).Error
}

func (d *DomainDao) LastSyncRun(source string) (*model.DomainSyncRun, error) {
	var run model.DomainSyncRun
	err := d.db.Where("external_source = ?", source).Order("started_at DESC").First(&run).Error
	if err != nil {
		return nil, err
	}
	return &run, nil
}

func normalizePage(page, pageSize int) (int, int) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}
	return page, pageSize
}
