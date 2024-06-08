package config

import (
	"github.com/spf13/viper"
	"log"
)

type Config struct {
	Viper *viper.Viper
}

func NewConfig() *Config {
	v := viper.New()
	v.SetConfigName("config")
	v.AddConfigPath("./config")
	v.AddConfigPath("$HOME")
	err := v.ReadInConfig()
	if err != nil {
		log.Fatal("Error while parsing configuration file", err)
	}

	return &Config{
		Viper: v,
	}
}
