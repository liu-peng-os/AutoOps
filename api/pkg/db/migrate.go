package db

import (
	appmodel "dodevops-api/api/app/model"
	cmdbmodel "dodevops-api/api/cmdb/model"
	ccmodel "dodevops-api/api/configcenter/model"
	k8smodel "dodevops-api/api/k8s/model"
	monitormodel "dodevops-api/api/monitor/model"
	systemmodel "dodevops-api/api/system/model"
	taskmodel "dodevops-api/api/task/model"
	toolmodel "dodevops-api/api/tool/model"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"gorm.io/gorm"
)

var models = []interface{}{
	&cmdbmodel.CmdbGroup{},
	&ccmodel.EcsAuth{},
	&ccmodel.KeyManage{},
	&ccmodel.SyncSchedule{},
	&cmdbmodel.CmdbHost{},
	&cmdbmodel.CmdbSQLRecord{},
	&cmdbmodel.CmdbSQL{},
	&ccmodel.AccountAuth{},
	&taskmodel.TaskTemplate{},
	&taskmodel.Task{},
	&taskmodel.TaskWork{},
	&taskmodel.TaskAnsible{},
	&taskmodel.TaskAnsibleWork{},
	&monitormodel.Agent{},
	&k8smodel.KubeCluster{},
	&appmodel.Application{},
	&appmodel.JenkinsEnv{},
	&appmodel.QuickDeployment{},
	&appmodel.QuickDeploymentTask{},
	&systemmodel.SysOperationLog{},
	&toolmodel.Tool{},
	&toolmodel.ServiceDeploy{},
}

func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(models...)
}

type schemaMigration struct {
	ID        uint      `gorm:"primaryKey"`
	Name      string    `gorm:"size:255;uniqueIndex;not null"`
	AppliedAt time.Time `gorm:"not null"`
}

func (schemaMigration) TableName() string {
	return "schema_migrations"
}

func RunMigrations(db *gorm.DB, migrationPath string) error {
	if migrationPath == "" {
		migrationPath = "./migrations"
	}

	if err := db.AutoMigrate(&schemaMigration{}); err != nil {
		return fmt.Errorf("initialize schema_migrations failed: %w", err)
	}

	entries, err := os.ReadDir(migrationPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("read migration directory failed: %w", err)
	}

	var files []string
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		if strings.HasSuffix(entry.Name(), ".up.sql") {
			files = append(files, entry.Name())
		}
	}
	sort.Strings(files)

	for _, file := range files {
		var existing schemaMigration
		if err := db.Where("name = ?", file).First(&existing).Error; err == nil {
			continue
		} else if err != gorm.ErrRecordNotFound {
			return fmt.Errorf("query migration state failed for %s: %w", file, err)
		}

		sqlBytes, err := os.ReadFile(filepath.Join(migrationPath, file))
		if err != nil {
			return fmt.Errorf("read migration file failed for %s: %w", file, err)
		}
		sqlContent := strings.TrimSpace(string(sqlBytes))

		if err := db.Transaction(func(tx *gorm.DB) error {
			if sqlContent != "" {
				if execErr := tx.Exec(sqlContent).Error; execErr != nil {
					return fmt.Errorf("execute migration %s failed: %w", file, execErr)
				}
			}

			return tx.Create(&schemaMigration{
				Name:      file,
				AppliedAt: time.Now(),
			}).Error
		}); err != nil {
			return err
		}
	}

	return nil
}
