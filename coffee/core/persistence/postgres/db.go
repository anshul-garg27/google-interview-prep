package postgres

import (
	"coffee/constants"
	"fmt"

	// "log"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	singletonDB *gorm.DB
	dbInit      sync.Once
)

func DB() *gorm.DB {
	host := viper.GetString("PG_HOST")
	user := viper.GetString("PG_USER")
	pass := viper.GetString("PG_PASS")
	db := viper.GetString("PG_DB")
	port := viper.GetString("PG_PORT")
	dbInit.Do(func() {
		dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Kolkata",
			host, user, pass, db, port)
		gormLogger := logger.New(
			// log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
			log.WithFields(log.Fields{}),
			logger.Config{
				SlowThreshold:             time.Nanosecond, // Slow SQL threshold
				LogLevel:                  logger.Info,     // Log level
				IgnoreRecordNotFoundError: true,            // Ignore ErrRecordNotFound error for logger
				Colorful:                  false,           // Disable color
			},
		)
		_db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{Logger: gormLogger, PrepareStmt: false})
		_db.Logger.LogMode(logger.Info)
		if err != nil {
			log.Fatal("Unable to connect to Postgres")
		} else {
			singletonDB = _db
		}
		sqlDB, err := singletonDB.DB()
		sqlDB.SetMaxIdleConns(viper.GetInt(constants.PgMaxIdleConnConfig))
		sqlDB.SetMaxOpenConns(viper.GetInt(constants.PgMaxOpenConnConfig))
	})
	return singletonDB
}
