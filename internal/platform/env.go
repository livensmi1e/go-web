package platform

import (
	"os"
	"strconv"
)

func getEnvStr(key string, fallback string) string {
	if val, exist := os.LookupEnv(key); exist {
		return val
	} else {
		return fallback
	}
}

func getEnvBool(key string, fallback bool) bool {
	if val, exist := os.LookupEnv(key); exist {
		valBool, _ := strconv.ParseBool(val)
		return valBool
	} else {
		return fallback
	}
}
