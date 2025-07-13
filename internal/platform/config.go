package platform

import "fmt"

type Config struct {
	HttpHost string
	HttpPort string

	MonitorHost string
	MonitorPort string

	Env           string
	Debug         bool
	MonitorEnable bool
}

func NewConfig() *Config {
	cfg := &Config{
		HttpHost:      getEnvStr("HTTP_HOST", "localhost"),
		HttpPort:      getEnvStr("HTTP_PORT", "8000"),
		MonitorHost:   getEnvStr("HTTP_HOST", "localhost"),
		MonitorPort:   getEnvStr("HTTP_PORT", "9000"),
		Env:           getEnvStr("ENV", "dev"),
		Debug:         getEnvBool("DEBUG", true),
		MonitorEnable: getEnvBool("MONITOR_ENABLE", true),
	}
	return cfg
}

func (c *Config) HttpServerAddr() string {
	return fmt.Sprintf("%s:%s", c.HttpHost, c.HttpPort)
}

func (c *Config) MonitorServerAddr() string {
	return fmt.Sprintf("%s:%s", c.MonitorHost, c.MonitorPort)
}

func (c *Config) StoreAddr() string {
	host := getEnvStr("POSTGRES_HOST", "localhost")
	port := getEnvStr("POSTGRES_PORT", "5432")
	user := getEnvStr("POSTGRES_USER", "postgres")
	password := getEnvStr("POSTGRES_PASSWORD", "postgres")
	db := getEnvStr("POSTGRES_DB", "postgres")
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", user, password, host, port, db)
}
