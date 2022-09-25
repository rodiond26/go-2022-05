package main

import (
	"fmt"

	"github.com/BurntSushi/toml"
)

// При желании конфигурацию можно вынести в internal/config.
// Организация конфига в main принуждает нас сужать API компонентов, использовать
// при их конструировании только необходимые параметры, а также уменьшает вероятность циклической зависимости.
type Config struct {
	Logger LoggerConf
	DB     DBConf
}

type LoggerConf struct {
	Level string `yaml:"c"`
	Path  string
}

type DBConf struct {
	DSN string
}

func NewConfig(configPath string) (config Config, err error) {
	_, err = toml.DecodeFile(configPath, &config)
	if err != nil {
		err = fmt.Errorf("when decode config [%s] then error [%w]", configPath, err)
		return Config{}, err
	}
	return config, nil
}
