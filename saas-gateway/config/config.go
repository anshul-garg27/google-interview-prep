package config

import (
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"

	log "github.com/sirupsen/logrus"

	"github.com/joho/godotenv"
)

func init() {
	env := os.Getenv("ENV")
	var err error

	_, b, _, _ := runtime.Caller(0)
	ProjectRootPath := filepath.Join(filepath.Dir(b), "../")

	switch env {
	case "PRODUCTION":
		err = godotenv.Load(".env.production")
		break
	case "STAGE":
		err = godotenv.Load(".env.stage")
		break
	case "LOCAL":
		err = godotenv.Load(".env.local")
		break
	default:
		err = godotenv.Load(ProjectRootPath + "/.env")
	}

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	New()
}

type Config struct {
	Port                   string
	HMAC_SECRET            string
	Env                    string
	SAAS_URL               string
	SAAS_DATA_URL          string
	IDENTITY_URL           string
	PARTNER_URL            string
	SENTRY_DSN             string
	ERP_CLIENT_ID          string
	LOG_LEVEL              int
	RedisClusterAddresses  []string
	REDIS_CLUSTER_PASSWORD string
}

var (
	config     Config
	configInit sync.Once
)

func New() Config {
	redisClusterAddressesStr := os.Getenv("REDIS_CLUSTER_ADDRESSES")
	configInit.Do(func() {
		config = Config{
			os.Getenv("PORT"),
			os.Getenv("HMAC_SECRET"),
			os.Getenv("ENV"),
			os.Getenv("SAAS_URL"),
			os.Getenv("SAAS_DATA_URL"),
			os.Getenv("IDENTITY_URL"),
			os.Getenv("PARTNER_URL"),
			os.Getenv("SENTRY_DSN"),
			os.Getenv("ERP_CLIENT_ID"),
			AtoiWithPanic(os.Getenv("LOG_LEVEL")),
			strings.Split(redisClusterAddressesStr, ","),
			os.Getenv("REDIS_CLUSTER_PASSWORD"),
		}
	})
	return config
}

func AtoiWithPanic(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}
	return i
}
