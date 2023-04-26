package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/go-yaml/yaml"
)

type Config struct {
	PostgresDb    PostgresDb `yaml:"db"`
	ServerAddress string     `yaml:"server_address"`
	JwtSecret     string     `yaml:"jwt_secret"`
}

type PostgresDb struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Dbname   string `yaml:"dbname"`
}

func ReadConfig(env string) (cfg Config, err error) {
	absPath, err := filepath.Abs(fmt.Sprintf("files/%s.yaml", env))
	if err != nil {
		return
	}
	file, err := os.Open(absPath)
	if err != nil {
		return
	}
	defer file.Close()

	d := yaml.NewDecoder(file)

	err = d.Decode(&cfg)

	return
}
