package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"sync"

)

const(
	config = "config.yml"
)

type Config struct {
	Listen struct{
		Port string `yaml:"port"`
		BindIP string `yaml:"bind_ip"`
	} `yaml:"listen"`
	Nats struct{
		ClusterID string `yaml:"cluster_id"`
		ClientID string `yaml:"client_id"`
		Channel string `yaml:"channel"`
	} `yaml:"nats"`
	DB struct{
		Username string `yaml:"username"`
		Password string `yaml:"password"`
		Host string `yaml:"host"`
		Port string `yaml:"port"`
		Name string `yaml:"db"`
		MaxAttempts int `yaml:"max_attempts"`
	} `yaml:"db"`
}

var instance *Config
var once sync.Once

func GetConfig() *Config {
	once.Do(func() {
		instance = &Config{}
		err := cleanenv.ReadConfig(config, instance)
		if err != nil {
			log.Printf("err: %s\n", err.Error())
			return
		}
	})

	return instance
}