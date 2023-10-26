package config

import (
	"gopkg.in/yaml.v2"
	"log"
	"os"
	"time"
)

type Config struct {
	ServerUrl    string `yaml:"server_url"`
	PostInterval int    `yaml:"post_interval"`
}

func newConfig() *Config {
	cfg := &Config{}

	// config.yamlファイルを読み取る
	data, err := os.ReadFile("config.yaml")
	if err != nil {
		logAndExit("config.yamlの読み込みに失敗しました\nError: %v", err)
	}

	// YAMLを構造体にアンマーシャル
	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		logAndExit("config.yamlの解析に失敗しました\nError: %v", err)
	}

	return cfg
}

func logAndExit(format string, v ...interface{}) {
	log.Printf(format, v...)
	time.Sleep(5 * time.Second)
	os.Exit(1)
}

var AppConfig = newConfig()
