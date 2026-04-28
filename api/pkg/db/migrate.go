package db

import (
	"database/sql"
	appmodel "dodevops-api/api/app/model"
	cmdbmodel "dodevops-api/api/cmdb/model"
	ccmodel "dodevops-api/api/configcenter/model"
	k8smodel "dodevops-api/api/k8s/model"
	monitormodel "dodevops-api/api/monitor/model"
	systemmodel "dodevops-api/api/system/model"
	taskmodel "dodevops-api/api/task/model"
	toolmodel "dodevops-api/api/tool/model"
	"fmt"
	"regexp"
	"strings"

	"github.com/pressly/goose/v3"
	"gorm.io/gorm"
)

var sqlManagedModels = []interface{}{
	&cmdbmodel.CmdbGroup{},
	&ccmodel.EcsAuth{},
	&ccmodel.KeyManage{},
	&ccmodel.SyncSchedule{},
	&ccmodel.AccountAuth{},
	&cmdbmodel.CmdbHost{},
	&cmdbmodel.CmdbSQLRecord{},
	&cmdbmodel.CmdbSQL{},
	&monitormodel.Agent{},
	&toolmodel.Tool{},
	&toolmodel.ServiceDeploy{},
	&systemmodel.SysAdmin{},
	&systemmodel.SysRole{},
	&systemmodel.SysMenu{},
	&systemmodel.SysDept{},
	&systemmodel.SysPost{},
	&systemmodel.SysAdminRole{},
	&systemmodel.SysRoleMenu{},
	&systemmodel.SysLoginInfo{},
	&systemmodel.SysOperationLog{},
	&k8smodel.KubeCluster{},
	&taskmodel.TaskTemplate{},
	&taskmodel.Task{},
	&taskmodel.TaskWork{},
	&taskmodel.TaskAnsible{},
	&taskmodel.TaskAnsibleWork{},
	&appmodel.Application{},
	&appmodel.JenkinsEnv{},
	&appmodel.QuickDeployment{},
	&appmodel.QuickDeploymentTask{},
}

var legacyAutoMigrateModels = []interface{}{}

func AutoMigrate(db *gorm.DB) error {
	if len(legacyAutoMigrateModels) == 0 {
		return nil
	}
	return db.AutoMigrate(legacyAutoMigrateModels...)
}

func SQLManagedTableNames() []string {
	return []string{
		"cmdb_group",
		"config_ecsauth",
		"config_keymanage",
		"config_sync_schedule",
		"config_account",
		"cmdb_host",
		"cmdb_sql_log",
		"cmdb_sql",
		"monitor_agent",
		"tool_link",
		"tool_service_deploy",
		"sys_admin",
		"sys_role",
		"sys_menu",
		"sys_dept",
		"sys_post",
		"sys_admin_role",
		"sys_role_menu",
		"sys_login_info",
		"sys_operation_log",
		"k8s_cluster",
		"task_template",
		"task_job",
		"task_work",
		"task_ansible",
		"task_ansiblework",
		"app_application",
		"app_jenkins_env",
		"quick_deployments",
		"quick_deployment_tasks",
	}
}

func RunMigrations(db *gorm.DB, migrationPath string) error {
	if migrationPath == "" {
		migrationPath = "./migrations"
	}

	gooseDialect, err := normalizeGooseDialectName("")
	if err != nil {
		return err
	}
	if db != nil && db.Config != nil && db.Config.Dialector != nil {
		gooseDialect, err = normalizeGooseDialectName(db.Config.Dialector.Name())
		if err != nil {
			return err
		}
	}
	if err := goose.SetDialect(gooseDialect); err != nil {
		return fmt.Errorf("configure goose dialect failed: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return err
	}

	if err := bootstrapGooseVersionTable(sqlDB, gooseDialect); err != nil {
		return err
	}

	if err := goose.Up(sqlDB, migrationPath); err != nil {
		return fmt.Errorf("goose up failed: %w", err)
	}

	return nil
}

func normalizeGooseDialectName(name string) (string, error) {
	switch strings.ToLower(name) {
	case "", "postgres", "postgresql":
		return "postgres", nil
	case "mysql":
		return "mysql", nil
	default:
		return "", fmt.Errorf("unsupported goose dialect: %s", name)
	}
}

func bootstrapGooseVersionTable(db *sql.DB, dialect string) error {
	hasGooseTable, err := tableExists(db, dialect, goose.DefaultTablename)
	if err != nil {
		return err
	}
	if hasGooseTable {
		return nil
	}

	hasLegacyTable, err := tableExists(db, dialect, "schema_migrations")
	if err != nil {
		return err
	}
	if !hasLegacyTable {
		return nil
	}

	versions, err := readLegacyMigrationVersions(db)
	if err != nil {
		return err
	}
	if len(versions) == 0 {
		return nil
	}

	if _, err := db.Exec(gooseVersionTableDDL(dialect, goose.DefaultTablename)); err != nil {
		return fmt.Errorf("create goose version table failed: %w", err)
	}

	for _, version := range versions {
		if _, err := db.Exec(gooseVersionInsertSQL(dialect, goose.DefaultTablename), version, true); err != nil {
			return fmt.Errorf("bootstrap goose version %d failed: %w", version, err)
		}
	}
	return nil
}

func tableExists(db *sql.DB, dialect, table string) (bool, error) {
	var query string
	switch dialect {
	case "postgres":
		query = `SELECT EXISTS (SELECT 1 FROM pg_tables WHERE (current_schema() IS NULL OR schemaname = current_schema()) AND tablename = $1)`
	case "mysql":
		query = `SELECT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_schema = DATABASE() AND table_name = ?)`
	default:
		return false, fmt.Errorf("unsupported dialect for tableExists: %s", dialect)
	}

	var exists bool
	if err := db.QueryRow(query, table).Scan(&exists); err != nil {
		return false, fmt.Errorf("check table %s failed: %w", table, err)
	}
	return exists, nil
}

var legacyMigrationVersionPattern = regexp.MustCompile(`^0*([0-9]+)_.+`)

func readLegacyMigrationVersions(db *sql.DB) ([]int64, error) {
	rows, err := db.Query(`SELECT name FROM schema_migrations ORDER BY id ASC`)
	if err != nil {
		return nil, fmt.Errorf("read legacy schema_migrations failed: %w", err)
	}
	defer rows.Close()

	var versions []int64
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, fmt.Errorf("scan legacy migration name failed: %w", err)
		}
		version, err := parseLegacyMigrationVersion(name)
		if err != nil {
			return nil, err
		}
		versions = append(versions, version)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate legacy schema_migrations failed: %w", err)
	}
	return versions, nil
}

func parseLegacyMigrationVersion(name string) (int64, error) {
	match := legacyMigrationVersionPattern.FindStringSubmatch(name)
	if len(match) != 2 {
		return 0, fmt.Errorf("parse legacy migration version failed for %s", name)
	}

	var version int64
	for _, ch := range match[1] {
		version = version*10 + int64(ch-'0')
	}
	return version, nil
}

func gooseVersionTableDDL(dialect, table string) string {
	switch dialect {
	case "postgres":
		return fmt.Sprintf(`CREATE TABLE %s (
		id integer PRIMARY KEY GENERATED BY DEFAULT AS IDENTITY,
		version_id bigint NOT NULL,
		is_applied boolean NOT NULL,
		tstamp timestamp NOT NULL DEFAULT now()
	)`, table)
	case "mysql":
		return fmt.Sprintf(`CREATE TABLE %s (
		id bigint unsigned NOT NULL AUTO_INCREMENT,
		version_id bigint NOT NULL,
		is_applied boolean NOT NULL,
		tstamp timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
		PRIMARY KEY (id)
	)`, table)
	default:
		return ""
	}
}

func gooseVersionInsertSQL(dialect, table string) string {
	switch dialect {
	case "postgres":
		return fmt.Sprintf(`INSERT INTO %s (version_id, is_applied) VALUES ($1, $2)`, table)
	case "mysql":
		return fmt.Sprintf(`INSERT INTO %s (version_id, is_applied) VALUES (?, ?)`, table)
	default:
		return ""
	}
}
