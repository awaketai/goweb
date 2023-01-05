package orm

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/awaketai/goweb/framework"
	"github.com/awaketai/goweb/framework/contract"
	"gorm.io/driver/clickhouse"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

type WebOrmService struct {
	container  framework.Container
	configPath string
	gormConfig *gorm.Config
	dbs        map[string]*gorm.DB // key:dbs val:gorm.DB
	lock       *sync.RWMutex
}

func NewWebOrm(params ...any) (any, error) {
	container := params[0].(framework.Container)
	dbs := make(map[string]*gorm.DB)
	lock := &sync.RWMutex{}
	return &WebOrmService{
		container: container,
		dbs:       dbs,
		lock:      lock,
	}, nil
}

func WithConfigPath(path string) contract.DBOption {
	return func(orm contract.ORM) error {
		webGorm := orm.(*WebOrmService)
		webGorm.configPath = path
		return nil
	}
}

func WithGormConfig(config *gorm.Config) contract.DBOption {
	return func(orm contract.ORM) error {
		webGorm := orm.(*WebOrmService)
		webGorm.gormConfig = config
		return nil
	}
}

// GetDB implement contract.ORMService
func (service *WebOrmService) GetDB(options ...contract.DBOption) (*gorm.DB, error) {
	for _, option := range options {
		option(service)
	}
	if service.configPath == "" {
		WithConfigPath("database.default")
	}
	configService := service.container.MustMake(contract.ConfigKey).(contract.Config)
	loggerService := service.container.MustMake(contract.LogKey).(contract.Log)
	config, err := GetBaseConfig(service.container)
	if err != nil {
		return nil, fmt.Errorf("WebOrmService.GetDB err:%v", err)
	}
	if err := configService.Load(service.configPath, config); err != nil {
		return nil, fmt.Errorf("WebOrmService.Load config err:%v", err)
	}
	if config.Dsn == "" {
		dsn, err := config.FormatDsn()
		if err != nil {
			return nil, fmt.Errorf("WebOrmService.FormatDsn err:%v", err)
		}
		config.Dsn = dsn
	}

	service.lock.RLock()
	if db, ok := service.dbs[config.Dsn]; ok {
		service.lock.RUnlock()
		return db, nil
	}
	service.lock.RUnlock()
	ormLogger := NewOrmLogger(loggerService)
	gormConfig := service.gormConfig
	if gormConfig == nil {
		gormConfig = &gorm.Config{}
	}
	if gormConfig.Logger == nil {
		gormConfig.Logger = ormLogger
	}
	if config.Dsn == "" {
		return nil, fmt.Errorf("db dsn can not blank")
	}
	var (
		db *gorm.DB
	)
	switch config.Driver {
	case Mysql:
		db, err = gorm.Open(mysql.Open(config.Dsn), gormConfig)
	case Postgres:
		db, err = gorm.Open(postgres.Open(config.Dsn), gormConfig)
	case Sqlite:
		db, err = gorm.Open(sqlite.Open(config.Dsn), gormConfig)
	case Sqlserver:
		db, err = gorm.Open(sqlserver.Open(config.Dsn), gormConfig)
	case Clickhouse:
		db, err = gorm.Open(clickhouse.Open(config.Dsn), gormConfig)
	}

	if err != nil {
		return nil, fmt.Errorf("open db err:%v", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("db.DB err:%v", err)
	}
	if config.ConnMaxIdle > 0 {
		sqlDB.SetMaxIdleConns(config.ConnMaxIdle)
	}
	if config.ConnMaxOpen > 0 {
		sqlDB.SetMaxOpenConns(config.ConnMaxOpen)
	}
	if config.ConnMaxLifetime != "" {
		lifeTime, err := time.ParseDuration(config.ConnMaxLifetime)
		if err != nil {
			loggerService.Error(context.Background(), "conn max life time error", map[string]any{
				"err": err,
			})
		} else {
			sqlDB.SetConnMaxLifetime(lifeTime)
		}
	}

	if err != nil {
		service.lock.Lock()
		service.dbs[config.Dsn] = db
		service.lock.Unlock()
	}
	return db, nil
}
