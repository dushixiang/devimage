package backend

import (
	"encoding/json"
	"github.com/mitchellh/go-homedir"
	"os"
)

type ConfigStore struct {
	ConfigPath string
	Error      string
}

type Config struct {
	Quality   int    `json:"quality"`
	OutputDir string `json:"outputDir"`
}

func NewConfigStore() (*ConfigStore, error) {
	cs := ConfigStore{}
	dir, err := homedir.Dir()
	if err != nil {
		return nil, err
	} else {
		configDir := dir + "/.dev-image"
		cs.ConfigPath = configDir + "/config.json"
		_, err := os.Stat(configDir)
		if err != nil {
			err := os.MkdirAll(configDir, os.ModePerm)
			if err != nil {
				return nil, err
			}

			var defaultConfig = &Config{
				Quality:   75,
				OutputDir: configDir + "/images",
			}
			if err := cs.Store(defaultConfig); err != nil {
				return nil, err
			}
		}
	}
	return &cs, nil
}

func (cs *ConfigStore) Store(config *Config) error {
	bytes, err := json.Marshal(config)
	if err != nil {
		return err
	}
	return os.WriteFile(cs.ConfigPath, bytes, os.ModePerm)
}

func (cs *ConfigStore) Load() (*Config, error) {
	bytes, err := os.ReadFile(cs.ConfigPath)
	if err != nil {
		return nil, err
	}
	var config Config
	if err := json.Unmarshal(bytes, &config); err != nil {
		return nil, err
	}
	return &config, nil
}
