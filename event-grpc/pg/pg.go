package pg

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"init.bulbul.tv/bulbul-backend/event-grpc/config"
	"log"
	"sync"
	"time"
)

var (
	singletonPg *gorm.DB
	pgInit      sync.Once
	mutex *sync.Mutex = &sync.Mutex{}
)


func PG(config config.Config) *gorm.DB {
	pgInit.Do(func() {
		connectToPG(config)
		go pgConnectionCron(config)
	})
	return singletonPg
}

func pgConnectionCron(config config.Config)  {
	for {
		connectToPG(config)
		time.Sleep(1000*time.Millisecond)
	}
}

func connectToPG(config config.Config) {
	mutex.Lock()
	var err error

	if singletonPg == nil || singletonPg.Error != nil || singletonPg.DB().Ping() != nil {
		singletonPg, err = gorm.Open("postgres",
			"host="+config.POSTGRES_CONNECTION_CONFIG.POSTGRES_HOST+
				" port="+config.POSTGRES_CONNECTION_CONFIG.POSTGRES_PORT+
				" user="+config.POSTGRES_CONNECTION_CONFIG.POSTGRES_USER+
				" dbname="+config.POSTGRES_CONNECTION_CONFIG.POSTGRES_DBNAME+
				" password="+config.POSTGRES_CONNECTION_CONFIG.POSTGRES_PASSWORD)

		if err != nil {
			log.Printf("Error connecting to postgres: %v", err)
		} else {
			singletonPg.DB().SetMaxIdleConns(10)
			singletonPg.DB().SetMaxOpenConns(20)
			singletonPg.DB().SetConnMaxLifetime(0)

			gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
				return config.POSTGRES_CONNECTION_CONFIG.POSTGRES_EVENT_SCHEMA + "." + defaultTableName
			}
		}
	}

	mutex.Unlock()
}