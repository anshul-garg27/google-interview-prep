package clickhouse

import (
	"gorm.io/driver/clickhouse"
	"gorm.io/gorm"
	"init.bulbul.tv/bulbul-backend/event-grpc/config"
	"log"
	"sync"
	"time"
)

var (
	singletonClickhouseMap map[string]*gorm.DB
	clickhouseInit         sync.Once
	mutex                  *sync.Mutex = &sync.Mutex{}
)

func Clickhouse(config config.Config, dbName *string) *gorm.DB {
	clickhouseInit.Do(func() {
		connectToClickhouse(config)
		go clickhouseConnectionCron(config)
	})

	if dbName != nil {
		if db, ok := singletonClickhouseMap[*dbName]; ok {
			return db
		}
	}

	return singletonClickhouseMap[config.CLICKHOUSE_CONNECTION_CONFIG.CLICKHOUSE_DB_NAME]
}

func clickhouseConnectionCron(config config.Config) {
	for {
		connectToClickhouse(config)
		time.Sleep(1000 * time.Millisecond)
	}
}

func connectToClickhouse(config config.Config) {
	mutex.Lock()
	//dsn := "tcp://clickhouse:9000?database=_e&username=default&password=Bul@Bul@Click@1!&read_timeout=10&write_timeout=20"

	for _, dbName := range config.CLICKHOUSE_CONNECTION_CONFIG.CLICKHOUSE_DB_NAMES {
		dsn := "tcp://" +
			config.CLICKHOUSE_CONNECTION_CONFIG.CLICKHOUSE_HOST + ":" + config.CLICKHOUSE_CONNECTION_CONFIG.CLICKHOUSE_PORT +
			"?database=" + dbName +
			"&username=" + config.CLICKHOUSE_CONNECTION_CONFIG.CLICKHOUSE_USER +
			"&password=" + config.CLICKHOUSE_CONNECTION_CONFIG.CLICKHOUSE_PASSWORD +
			"&read_timeout=10&write_timeout=20"

		connectionBroken := false

		if singletonClickhouseMap == nil {
			singletonClickhouseMap = map[string]*gorm.DB{}
		}

		if singletonClickhouseMap[dbName] == nil || singletonClickhouseMap[dbName].Error != nil {
			connectionBroken = true
		}

		if !connectionBroken {
			if _, err := singletonClickhouseMap[dbName].DB(); err != nil {
				connectionBroken = true
			}
		}

		if connectionBroken {
			db, err := gorm.Open(clickhouse.New(clickhouse.Config{
				DSN:                      dsn,
				DisableDatetimePrecision: true,
			}), &gorm.Config{})

			if err != nil {
				log.Printf("Error connecting to clickhouse: %v", err)
			} else {
				singletonClickhouseMap[dbName] = db
			}
		}
	}

	mutex.Unlock()
}
