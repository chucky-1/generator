// Package configs has a configuration structure
package configs

// Config contains configuration data
type Config struct {
	Host string `env:"HOST" envDefault:"localhost"`
	Port string `env:"PORT" envDefault:"6379"`
}
