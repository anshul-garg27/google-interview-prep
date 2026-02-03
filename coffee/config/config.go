package config

import (
	"coffee/constants"
	"fmt"
	"os"
	"strconv"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func Setup() {
	viper.SetDefault(constants.ServerTimeoutConfig, 5)
	viper.SetDefault(constants.LogLevelConfig, constants.LogLevelDebug)
	viper.SetDefault(constants.PgMaxIdleConnConfig, 5)
	viper.SetDefault(constants.PgMaxOpenConnConfig, 5)

	viper.SetConfigFile(".env")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal(err)
	}
}

func getBoolEnv(key string) bool {
	val := getStrEnv(key)
	ret, err := strconv.ParseBool(val)
	if err != nil {
		panic(fmt.Sprintf("some error"))
	}
	return ret
}

func getStrEnv(key string) string {
	val := os.Getenv(key)
	if val == "" {
		panic(fmt.Sprintf("some error msg"))
	}
	return val
}

func AtoiWithPanic(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}
	return i
}
