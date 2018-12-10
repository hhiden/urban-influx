package config

import (
	"os"
	"strings"
)

type config struct {
	InfluxURL string
	InfluxDB  string
}

func GetConfig() config {
	return config{
		InfluxURL: strings.ToLower(getEnv("INFLUX_DB_HOST", "http://influx-push-myproject.192.168.64.36.nip.io/")),
		InfluxDB:  strings.ToLower(getEnv("INFLUX_DB_NAME", "udp")),
	}
}

func getEnv(key string, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}
