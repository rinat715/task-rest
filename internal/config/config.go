package config

import "github.com/cristalhq/aconfig"

type Config struct {
	Host     string `env:"Host" usage:"host url server"`
	Port     int    `env:"Port" usage:"just a number"`
	LogLevel int    `env:"LogLevel" usage:"set number"`
	PathDB   string `env:"PathDB" usage:"set conn path db"`
	Salt     string `env:"Salt" usage:"salt for hash pass"`
}

func NewConfig() *Config {
	var config Config
	loader := aconfig.LoaderFor(&config, aconfig.Config{
		SkipFiles:        true,
		SkipFlags:        true,
		AllFieldRequired: true,
	})

	if err := loader.Load(); err != nil {
		panic(err)
	}
	return &config
}
