package orm

import (
	"fmt"
	"sync"

	"github.com/tnngo/hammer/logger"
	"github.com/tnngo/lad"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	gl "gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"gorm.io/plugin/dbresolver"
)

type Orm struct {
	// Level 日志级别。
	Level logger.Level
	// MaxOpenConns 最大连接数，默认无限制。
	MaxOpenConns int `yaml:"max_open_conns"`
	// MaxIdleConns 最大空闲连接数据，默认无限制。
	MaxIdleConns int `yaml:"max_idle_conns"`
	// Options 配置选项。
	Options []*Option

	once sync.Once
}

// Options 数据库配置选项。
// 默认数据库类型mysql，
type Option struct {
	// DBType 数据库类型，0:mysql，1:pgsql。
	DBType string `yaml:"dbtype"`
	// User 数据库用户。
	User string
	// Password 数据库密码。
	Password string
	// IP 数据库IP地址。
	IP string
	// Port 数据库端口。
	Port int
	// Database 数据库名称。
	Database string
	// Charset连接编码。
	Charset string
}

var _db *gorm.DB

func gormLevel(level logger.Level) gl.LogLevel {
	switch level {
	case logger.DebugLevel, logger.InfoLevel:
		return gl.Info
	case logger.WarnLevel:
		return gl.Warn
	case logger.ErrorLevel:
		return gl.Error
	default:
		return gl.Silent
	}
}

func (o *Orm) buildOnce() {
	var (
		mysqlFormat      = "%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local"
		postgresqlFormat = "host=%s user=%s password=%s dbname=%s port=%d sslmode=disable"
	)

	var register *dbresolver.DBResolver
	for i, v := range o.Options {
		if v.IP == "" {
			v.IP = "localhost"
		}

		if v.Port < 1 {
			v.Port = 3306
		}

		if v.Charset == "" {
			v.Charset = "utf8mb4"
		}

		var dsn string
		switch v.DBType {
		case "mysql":
			dsn = fmt.Sprintf(mysqlFormat, v.User, v.Password, v.IP, v.Port, v.Database, v.Charset)
		case "pgsql":
			dsn = fmt.Sprintf(postgresqlFormat, v.IP, v.User, v.Password, v.Database, v.Port)
		default:
			dsn = fmt.Sprintf(mysqlFormat, v.User, v.Password, v.IP, v.Port, v.Database, v.Charset)
		}

		config := &gorm.Config{
			NamingStrategy:         schema.NamingStrategy{SingularTable: true},
			SkipDefaultTransaction: true,
		}
		config.Logger = gl.Default.LogMode(gormLevel(o.Level))

		if i == 0 {
			db, err := gorm.Open(mysql.Open(dsn), config)
			if err != nil {
				lad.L().Error(err.Error(), lad.Reflect("options", v))
				return
			}
			_db = db
			continue
		}

		if register == nil {
			register = dbresolver.Register(dbresolver.Config{
				Sources: []gorm.Dialector{mysql.Open(dsn)},
			}, v.Database)
		} else {
			register = register.Register(dbresolver.Config{
				Sources: []gorm.Dialector{mysql.Open(dsn)},
			}, v.Database)
		}
	}

	if register != nil {
		if err := _db.Use(register); err != nil {
			lad.L().Error(err.Error())
			return
		}
	}
	sqlDB, err := _db.DB()
	if err != nil {
		lad.L().Error(err.Error())
		return
	}

	if err := sqlDB.Ping(); err != nil {
		lad.L().Error(err.Error())
	} else {
		lad.L().Info("Database connection succeeded", lad.Reflect("options", o.Options))
	}
}

func (o *Orm) Build() {
	o.once.Do(o.buildOnce)
}

func DB(dbnames ...string) *gorm.DB {
	if len(dbnames) != 0 {
		ifs := make([]clause.Expression, 0)
		for _, v := range dbnames {
			ifs = append(ifs, dbresolver.Use(v))
		}
		return _db.Clauses(ifs...)
	}
	return _db
}
