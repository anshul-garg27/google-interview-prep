package config

import (
	b64 "encoding/base64"
	"github.com/joho/godotenv"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"
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
	case "PRE_PRODUCTION":
		err = godotenv.Load(".env.pre_production")
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

// Config for common runtime
type Config struct {
	Port    string
	GinPort string
	Env     string

	IDENTITY_URL      string
	SOCIAL_STREAM_URL string

	WEBENGAGE_API_KEY  string
	WEBENGAGE_URL      string
	WEBENGAGE_USER_URL string

	WEBENGAGE_BRAND_API_KEY  string
	WEBENGAGE_BRAND_URL      string
	WEBENGAGE_BRAND_USER_URL string

	CUSTOMER_APP_CLIENT_ID string
	HOST_APP_CLIENT_ID     string

	GCC_HOST_APP_CLIENT_ID string

	GCC_BRAND_CLIENT_ID string

	RedisClusterAddresses  []string
	REDIS_CLUSTER_PASSWORD string

	LOG_LEVEL int

	RABBIT_CONNECTION_CONFIG *RABBIT_CONNECTION_CONFIG
	EVENT_WORKER_POOL_CONFIG *EVENT_WORKER_POOL_CONFIG
	EVENT_SINK_CONFIG        *EVENT_SINK_CONFIG

	BRANCH_EVENT_WORKER_POOL_CONFIG *EVENT_WORKER_POOL_CONFIG
	BRANCH_EVENT_SINK_CONFIG        *EVENT_SINK_CONFIG

	VIDOOLY_EVENT_WORKER_POOL_CONFIG *EVENT_WORKER_POOL_CONFIG
	VIDOOLY_EVENT_SINK_CONFIG        *EVENT_SINK_CONFIG

	WEBENGAGE_EVENT_WORKER_POOL_CONFIG *EVENT_WORKER_POOL_CONFIG
	WEBENGAGE_EVENT_SINK_CONFIG        *EVENT_SINK_CONFIG

	SHOPIFY_EVENT_WORKER_POOL_CONFIG *EVENT_WORKER_POOL_CONFIG
	SHOPIFY_EVENT_SINK_CONFIG        *EVENT_SINK_CONFIG

	GRAPHY_EVENT_WORKER_POOL_CONFIG *EVENT_WORKER_POOL_CONFIG
	GRAPHY_EVENT_SINK_CONFIG        *EVENT_SINK_CONFIG

	POSTGRES_CONNECTION_CONFIG   *POSTGRES_CONNECTION_CONFIG
	CLICKHOUSE_CONNECTION_CONFIG *CLICKHOUSE_CONNECTION_CONFIG

	HMAC_SECRET []byte
}

type EVENT_SINK_CONFIG struct {
	EVENT_EXCHANGE string
}

type EVENT_WORKER_POOL_CONFIG struct {
	EVENT_WORKER_POOL_SIZE      int
	EVENT_BUFFERED_CHANNEL_SIZE int
}

type RABBIT_CONNECTION_CONFIG struct {
	RABBIT_USER      string
	RABBIT_PASSWORD  string
	RABBIT_HOST      string
	RABBIT_PORT      string
	RABBIT_HEARTBEAT int // in millisecods
	RABBIT_VHOST     string
}

type POSTGRES_CONNECTION_CONFIG struct {
	POSTGRES_HOST         string
	POSTGRES_PORT         string
	POSTGRES_USER         string
	POSTGRES_DBNAME       string
	POSTGRES_PASSWORD     string
	POSTGRES_EVENT_SCHEMA string
}

type CLICKHOUSE_CONNECTION_CONFIG struct {
	CLICKHOUSE_HOST     string
	CLICKHOUSE_PORT     string
	CLICKHOUSE_USER     string
	CLICKHOUSE_DB_NAMES []string
	CLICKHOUSE_PASSWORD string
	CLICKHOUSE_DB_NAME  string
}

var (
	config     Config
	configInit sync.Once
)

func New() Config {

	configInit.Do(func() {
		redisClusterAddressesStr := os.Getenv("REDIS_CLUSTER_ADDRESSES")

		EVENT_WORKER_POOL_SIZE, _ := strconv.Atoi(os.Getenv("EVENT_WORKER_POOL_SIZE"))
		EVENT_BUFFERED_CHANNEL_SIZE, _ := strconv.Atoi(os.Getenv("EVENT_BUFFERED_CHANNEL_SIZE"))

		BRANCH_EVENT_WORKER_POOL_SIZE, _ := strconv.Atoi(os.Getenv("BRANCH_EVENT_WORKER_POOL_SIZE"))
		BRANCH_EVENT_BUFFERED_CHANNEL_SIZE, _ := strconv.Atoi(os.Getenv("BRANCH_EVENT_BUFFERED_CHANNEL_SIZE"))

		VIDOOLY_EVENT_WORKER_POOL_SIZE, _ := strconv.Atoi(os.Getenv("VIDOOLY_EVENT_WORKER_POOL_SIZE"))
		VIDOOLY_EVENT_BUFFERED_CHANNEL_SIZE, _ := strconv.Atoi(os.Getenv("VIDOOLY_EVENT_BUFFERED_CHANNEL_SIZE"))

		WEBENGAGE_EVENT_WORKER_POOL_SIZE, _ := strconv.Atoi(os.Getenv("WEBENGAGE_EVENT_WORKER_POOL_SIZE"))
		WEBENGAGE_EVENT_BUFFERED_CHANNEL_SIZE, _ := strconv.Atoi(os.Getenv("WEBENGAGE_EVENT_BUFFERED_CHANNEL_SIZE"))

		SHOPIFY_EVENT_WORKER_POOL_SIZE, _ := strconv.Atoi(os.Getenv("SHOPIFY_EVENT_WORKER_POOL_SIZE"))
		SHOPIFY_EVENT_BUFFERED_CHANNEL_SIZE, _ := strconv.Atoi(os.Getenv("SHOPIFY_EVENT_BUFFERED_CHANNEL_SIZE"))

		GRAPHY_EVENT_WORKER_POOL_SIZE, _ := strconv.Atoi(os.Getenv("GRAPHY_EVENT_WORKER_POOL_SIZE"))
		GRAPHY_EVENT_BUFFERED_CHANNEL_SIZE, _ := strconv.Atoi(os.Getenv("GRAPHY_EVENT_BUFFERED_CHANNEL_SIZE"))

		RABBIT_HEARTBEAT, _ := strconv.Atoi(os.Getenv("RABBIT_HEARTBEAT"))
		rabbitConnectionConfig := RABBIT_CONNECTION_CONFIG{
			RABBIT_USER:      os.Getenv("RABBIT_USER"),
			RABBIT_PASSWORD:  os.Getenv("RABBIT_PASSWORD"),
			RABBIT_HOST:      os.Getenv("RABBIT_HOST"),
			RABBIT_PORT:      os.Getenv("RABBIT_PORT"),
			RABBIT_HEARTBEAT: RABBIT_HEARTBEAT,
			RABBIT_VHOST:     os.Getenv("RABBIT_VHOST"),
		}
		eventSinkConfig := EVENT_SINK_CONFIG{
			EVENT_EXCHANGE: os.Getenv("EVENT_EXCHANGE"),
		}
		eventWorkerPoolConfig := EVENT_WORKER_POOL_CONFIG{EVENT_WORKER_POOL_SIZE, EVENT_BUFFERED_CHANNEL_SIZE}

		branchEventSinkConfig := EVENT_SINK_CONFIG{
			EVENT_EXCHANGE: os.Getenv("BRANCH_EVENT_EXCHANGE"),
		}

		branchEventWorkerPoolConfig := EVENT_WORKER_POOL_CONFIG{BRANCH_EVENT_WORKER_POOL_SIZE, BRANCH_EVENT_BUFFERED_CHANNEL_SIZE}

		vidoolyEventWorkerPoolConfig := EVENT_WORKER_POOL_CONFIG{VIDOOLY_EVENT_WORKER_POOL_SIZE, VIDOOLY_EVENT_BUFFERED_CHANNEL_SIZE}
		vidoolyEventSinkConfig := EVENT_SINK_CONFIG{
			EVENT_EXCHANGE: os.Getenv("VIDOOLY_EVENT_EXCHANGE"),
		}
		webengageEventSinkConfig := EVENT_SINK_CONFIG{
			EVENT_EXCHANGE: os.Getenv("WEBENGAGE_EVENT_EXCHANGE"),
		}
		webengageEventWorkerPoolConfig := EVENT_WORKER_POOL_CONFIG{WEBENGAGE_EVENT_WORKER_POOL_SIZE, WEBENGAGE_EVENT_BUFFERED_CHANNEL_SIZE}

		shopifyEventSinkConfig := EVENT_SINK_CONFIG{
			EVENT_EXCHANGE: os.Getenv("SHOPIFY_EVENT_EXCHANGE"),
		}
		shopifyEventWorkerPoolConfig := EVENT_WORKER_POOL_CONFIG{SHOPIFY_EVENT_WORKER_POOL_SIZE, SHOPIFY_EVENT_BUFFERED_CHANNEL_SIZE}

		graphyEventSinkConfig := EVENT_SINK_CONFIG{
			EVENT_EXCHANGE: os.Getenv("GRAPHY_EVENT_EXCHANGE"),
		}
		graphyEventWorkerPoolConfig := EVENT_WORKER_POOL_CONFIG{GRAPHY_EVENT_WORKER_POOL_SIZE, GRAPHY_EVENT_BUFFERED_CHANNEL_SIZE}

		postgresConnectionConfig := POSTGRES_CONNECTION_CONFIG{
			POSTGRES_HOST:         os.Getenv("POSTGRES_HOST"),
			POSTGRES_PORT:         os.Getenv("POSTGRES_PORT"),
			POSTGRES_USER:         os.Getenv("POSTGRES_USER"),
			POSTGRES_DBNAME:       os.Getenv("POSTGRES_DBNAME"),
			POSTGRES_PASSWORD:     os.Getenv("POSTGRES_PASSWORD"),
			POSTGRES_EVENT_SCHEMA: os.Getenv("POSTGRES_EVENT_SCHEMA"),
		}

		clickhouseConnectionConfig := CLICKHOUSE_CONNECTION_CONFIG{
			CLICKHOUSE_HOST:     os.Getenv("CLICKHOUSE_HOST"),
			CLICKHOUSE_PORT:     os.Getenv("CLICKHOUSE_PORT"),
			CLICKHOUSE_USER:     os.Getenv("CLICKHOUSE_USER"),
			CLICKHOUSE_DB_NAMES: strings.Split(os.Getenv("CLICKHOUSE_DB_NAMES"), ","),
			CLICKHOUSE_DB_NAME:  os.Getenv("CLICKHOUSE_DBNAME"),
			CLICKHOUSE_PASSWORD: os.Getenv("CLICKHOUSE_PASSWORD"),
		}

		HMAC_SECRET, _ := b64.StdEncoding.DecodeString(os.Getenv("HMAC_SECRET"))

		config = Config{
			Port:    os.Getenv("PORT"),
			GinPort: os.Getenv("GIN_PORT"),
			Env:     os.Getenv("ENV"),

			IDENTITY_URL:      os.Getenv("IDENTITY_URL"),
			SOCIAL_STREAM_URL: os.Getenv("SOCIAL_STREAM_URL"),

			WEBENGAGE_API_KEY:  os.Getenv("WEBENGAGE_API_KEY"),
			WEBENGAGE_URL:      os.Getenv("WEBENGAGE_URL"),
			WEBENGAGE_USER_URL: os.Getenv("WEBENGAGE_USER_URL"),

			WEBENGAGE_BRAND_API_KEY:  os.Getenv("WEBENGAGE_BRAND_API_KEY"),
			WEBENGAGE_BRAND_URL:      os.Getenv("WEBENGAGE_BRAND_URL"),
			WEBENGAGE_BRAND_USER_URL: os.Getenv("WEBENGAGE_BRAND_USER_URL"),

			CUSTOMER_APP_CLIENT_ID: os.Getenv("CUSTOMER_APP_CLIENT_ID"),
			HOST_APP_CLIENT_ID:     os.Getenv("HOST_APP_CLIENT_ID"),
			GCC_HOST_APP_CLIENT_ID: os.Getenv("GCC_HOST_APP_CLIENT_ID"),
			GCC_BRAND_CLIENT_ID:    os.Getenv("GCC_BRAND_CLIENT_ID"),

			RedisClusterAddresses:  strings.Split(redisClusterAddressesStr, ","),
			REDIS_CLUSTER_PASSWORD: os.Getenv("REDIS_CLUSTER_PASSWORD"),

			LOG_LEVEL:                AtoiWithPanic(os.Getenv("LOG_LEVEL")),
			RABBIT_CONNECTION_CONFIG: &rabbitConnectionConfig,
			EVENT_SINK_CONFIG:        &eventSinkConfig,
			EVENT_WORKER_POOL_CONFIG: &eventWorkerPoolConfig,

			BRANCH_EVENT_SINK_CONFIG:        &branchEventSinkConfig,
			BRANCH_EVENT_WORKER_POOL_CONFIG: &branchEventWorkerPoolConfig,

			VIDOOLY_EVENT_WORKER_POOL_CONFIG: &vidoolyEventWorkerPoolConfig,
			VIDOOLY_EVENT_SINK_CONFIG:        &vidoolyEventSinkConfig,

			WEBENGAGE_EVENT_SINK_CONFIG:        &webengageEventSinkConfig,
			WEBENGAGE_EVENT_WORKER_POOL_CONFIG: &webengageEventWorkerPoolConfig,

			SHOPIFY_EVENT_SINK_CONFIG:        &shopifyEventSinkConfig,
			SHOPIFY_EVENT_WORKER_POOL_CONFIG: &shopifyEventWorkerPoolConfig,

			GRAPHY_EVENT_SINK_CONFIG:        &graphyEventSinkConfig,
			GRAPHY_EVENT_WORKER_POOL_CONFIG: &graphyEventWorkerPoolConfig,

			POSTGRES_CONNECTION_CONFIG:   &postgresConnectionConfig,
			CLICKHOUSE_CONNECTION_CONFIG: &clickhouseConnectionConfig,

			HMAC_SECRET: HMAC_SECRET,
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
