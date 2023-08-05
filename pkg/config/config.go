package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Env        string   `yaml:"env"`
	DbConn     string   `yaml:"db_conn"`
	BotToken   string   `yaml:"bot_token"`
	Channels   []string `yaml:"channels"`
	CronScrape string   `yaml:"cron_scrape"`
	CronNotify string   `yaml:"cron_notify"`
}

func ReadFile(src string) (*Config, error) {
	b, err := os.ReadFile(src)
	if err != nil {
		return nil, err
	}

	conf := &Config{}
	err = yaml.Unmarshal(b, conf)
	if err != nil {
		return nil, err
	}

	return conf, nil
}

func (c *Config) IsProd() bool {
    return c.Env == "prod"
}
