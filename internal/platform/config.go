package platform

import "fmt"

type Config struct {
	Host  string
	Port  string
	Env   string
	Debug bool
}

func NewConfig() *Config {
	cfg := &Config{
		Host:  getEnvStr("HOST", "localhost"),
		Port:  getEnvStr("PORT", "8000"),
		Env:   getEnvStr("ENV", "dev"),
		Debug: getEnvBool("DEBUG", true),
	}
	return cfg
}

func (c *Config) ServerAddr() string {
	return fmt.Sprintf("%s:%s", c.Host, c.Port)
}

func (c *Config) StoreAddr() string {
	host := getEnvStr("POSTGRES_HOST", "localhost")
	port := getEnvStr("POSTGRES_PORT", "5432")
	user := getEnvStr("POSTGRES_USER", "postgres")
	password := getEnvStr("POSTGRES_PASSWORD", "postgres")
	db := getEnvStr("POSTGRES_DB", "postgres")
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", user, password, host, port, db)
}
