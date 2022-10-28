package config

import (
	"log"
	"os"

	yaml "gopkg.in/yaml.v3"
)

type Config struct {
	Environment EnvironmentConfig `yaml:"environment"`
	Logger      LoggerConfig      `yaml:"logger"`
	Server      ServerConfig      `yaml:"server"`
}

type EnvironmentConfig struct {
	Type string `yaml:"type"`
}

type LoggerConfig struct {
	Level string `yaml:"level"`
}

type ServerConfig struct {
	Host     string `yaml:"host"`
	HTTPPort string `yaml:"httpPort"`
	GrpcPort string `yaml:"grpcPort"`
}

func NewConfig() Config {
	return Config{}
}

func (config *Config) ReadConfig(configPath string) (err error) {
	configFile, err := os.Open(configPath)
	if err != nil {
		return err
	}
	defer configFile.Close()

	yamlDecoder := yaml.NewDecoder(configFile)
	err = yamlDecoder.Decode(&config)
	if err != nil {
		return err
	}
	log.Printf("config = [%+v]\n", config)
	return nil
}
