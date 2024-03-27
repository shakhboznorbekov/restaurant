package config

import (
	"gopkg.in/yaml.v3"
	"log"
	"os"
)

type Config struct {
	BaseDestination string   `json:"base_destination" yaml:"base_destination"`
	FileDirectory   string   `json:"file_directory" yaml:"file_directory"`
	BOTToken        string   `json:"BOT_TOKEN" yaml:"BOT_TOKEN"`
	ChatID          string   `json:"ChatID" yaml:"ChatID"`
	ErrorBotToken   string   `json:"ERROR_BOT_TOKEN" yaml:"ERROR_BOT_TOKEN"`
	ErrorChatID     []string `json:"ERROR_CHAT_ID" yaml:"ERROR_CHAT_ID"`
}

func NewConfig() *Config {
	var c *Config

	yamlFile, err := os.ReadFile("config.yaml")
	if err != nil {
		log.Printf("yamlFile.Get err #%v", err)
	}
	err = yaml.Unmarshal(yamlFile, &c)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}

	return c
}
