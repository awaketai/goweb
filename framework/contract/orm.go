package contract

import (
	"net"
	"strconv"
	"time"

	"github.com/awaketai/goweb/framework"
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

const ORMKey = "web:orm"

// ORMService input parameter
type ORMService interface {
	GetDB(opt ...DBOption) (*gorm.DB, error)
}

// DBOption initialize option
type DBOption func(container framework.Container, config *DBConfig) error

// DBConfig
type DBConfig struct {
	// dsn config info
	WriteTimeout string `yaml:"write_timeout"` // write timeout
	Loc          string `yaml:"loc"`           // timezone
	Port         int    `yaml:"port"`          // port
	ReadTimeout  string `yaml:"read_timeout"`  // read timeout
	Charset      string `yaml:"charset"`       // character set
	ParseTime    bool   `yaml:"parse_time"`    // has parse time whether
	Protocol     string `yaml:"protocol"`      // protocol
	Dsn          string `yaml:"dsn"`           // if dsn not empty,other config about dsn inoperative
	Database     string `yaml:"database"`      // database
	Collation    string `yaml:"collation"`     // collation
	Timeout      string `yaml:"timeout"`       // connection timeout
	Username     string `yaml:"username"`
	Password     string `yaml:"password"`
	Driver       string `yaml:"driver"` // driver
	Host         string `yaml:"host"`   // database host

	// connection pool
	ConnMaxIdle     int    `yaml:"conn_max_idle"`     // max idle connection num
	ConnMaxOpen     int    `yaml:"conn_max_open"`     // max connection num
	ConnMaxLifetime string `yaml:"conn_max_lifetime"` // max connection lifetime
	ConnMaxIdletime string `yaml:"conn_max_idletime"` // max idle connection lifetime

	// gorm
	*gorm.Config
}

func (conf *DBConfig) FormatDsn() (string, error) {
	port := strconv.Itoa(conf.Port)
	timeout, err := time.ParseDuration(conf.Timeout)
	if err != nil {
		return "", err
	}
	readTimeout, err := time.ParseDuration(conf.ReadTimeout)
	if err != nil {
		return "", err
	}
	writeTimeout, err := time.ParseDuration(conf.WriteTimeout)
	if err != nil {
		return "", err
	}
	location, err := time.LoadLocation(conf.Loc)
	if err != nil {
		return "", err
	}
	driverConf := &mysql.Config{
		User:         conf.Username,
		Passwd:       conf.Password,
		Net:          conf.Protocol,
		Addr:         net.JoinHostPort(conf.Host, port),
		DBName:       conf.Database,
		Collation:    conf.Collation,
		Loc:          location,
		Timeout:      timeout,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
		ParseTime:    conf.ParseTime,
	}

	return driverConf.FormatDSN(), nil
}
