package clickhouse

import (
	"coffee/constants"
	"fmt"

	"gorm.io/driver/clickhouse"

	// "log"

	"sync"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/spf13/viper"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	singletonDB *gorm.DB
	dbInit      sync.Once
)

func ClickhouseDB() *gorm.DB {
	dbInit.Do(connectToClickhouse)
	return singletonDB
}

func connectToClickhouse() {
	host := viper.GetString("CH_HOST")
	user := viper.GetString("CH_USER")
	pass := viper.GetString("CH_PASS")
	db := viper.GetString("CH_DB")
	port := viper.GetString("CH_PORT")
	dsn := fmt.Sprintf("tcp://%s:%s@%s:%s/%s?dial_timeout=10s&max_execution_time=60s",
		user, pass, host, port, db)
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
	log.Println("Connecting to clickhouse db")
	_db, err := gorm.Open(clickhouse.New(clickhouse.Config{
		DSN:                          dsn,
		DisableDatetimePrecision:     true,     // disable datetime64 precision, not supported before clickhouse 20.4
		DontSupportRenameColumn:      true,     // rename column not supported before clickhouse 20.4
		DontSupportEmptyDefaultValue: false,    // do not consider empty strings as valid default values
		SkipInitializeWithVersion:    false,    // smart configure based on used version
		DefaultGranularity:           3,        // 1 granule = 8192 rows
		DefaultCompression:           "LZ4",    // default compression algorithm. LZ4 is lossless
		DefaultIndexType:             "minmax", // index stores extremes of the expression
		DefaultTableEngineOpts:       "ENGINE=MergeTree() ORDER BY tuple()",
	}), &gorm.Config{Logger: gormLogger, DisableAutomaticPing: false})
	_db.Logger.LogMode(logger.Info)
	if err != nil {
		log.Println("Unable to connect to Clickhouse db", err)
	} else {
		log.Println("Clickhouse DB Connect Successful")
		singletonDB = _db
		_sqlDB, err := singletonDB.DB()
		if err != nil {
			log.Println("Unable to Init Gorm Db from clickhouse db", err)
		} else {
			log.Println("Gorm DB Init Successful")
			_sqlDB.SetMaxIdleConns(viper.GetInt(constants.ChMaxIdleConnConfig))
			_sqlDB.SetMaxOpenConns(viper.GetInt(constants.ChMaxOpenConnConfig))
		}
	}
}
