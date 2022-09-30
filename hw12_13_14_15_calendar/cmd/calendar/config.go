package main

import (
	"bytes"
	"os"

	"github.com/spf13/viper"
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
	DSN      string
	Password string
}

func NewConfig(configPath string) (config Config, err error) {
	err = initConfig(configPath)
	if err != nil {
		return Config{}, err
	}
	return readConfig(), err
}

func initConfig(file string) (err error) {
	yaml, err := os.ReadFile(file)
	if err != nil {
		return err
	}
	viper.SetConfigType("yaml")
	err = viper.ReadConfig(bytes.NewBuffer(yaml))
	if err != nil {
		return err
	}
	return nil
}

func readConfig() (config Config) {
	return Config{
		Logger: LoggerConf{
			Level: viper.GetString("level"),
		},
		DB: DBConf{
			DSN:      "",
			Password: os.Getenv("DB_PASSWORD"),
		},
	}
}
