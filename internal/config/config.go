package config

import "github.com/cristalhq/aconfig"

type config struct {
	Port     int `env:"Port" usage:"just a number"`
	LogLevel int `env:"LogLevel" usage:"set number"`
}

var Config config

func init() {
	loader := aconfig.LoaderFor(&Config, aconfig.Config{
		SkipFiles:        true,
		SkipFlags:        true,
		AllFieldRequired: true,
	})

	if err := loader.Load(); err != nil {
		panic(err)
	}
}
