package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"sync"
)

type Config struct {
	Listen struct {
		Port string `yaml:"port" env-default:":8000"`
	} `yaml:"listen"`
	Postrgres struct {
		Host     string `json:"host"`
		Port     string `json:"port"`
		Database string `json:"database"`
		Username string `json:"username"`
		Password string `json:"password"`
	} `yaml:"postgres"`
}

var instance *Config
var once sync.Once

func GetConfig() *Config {
	once.Do(func() {
		instance = &Config{}
		if err := cleanenv.ReadConfig("configs/config.yml", instance); err != nil {
			help, _ := cleanenv.GetDescription(instance, nil)
			log.Println(help)
		}
		log.Println(instance)
	})
	return instance
}
