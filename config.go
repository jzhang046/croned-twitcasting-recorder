package main

import (
	"io/ioutil"
	"log"

	"github.com/go-playground/validator/v10"
	"gopkg.in/yaml.v3"
)

const (
	defaultConfigPath = "config.yaml"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

type cronConfig struct {
	Streamers []*streamerConfig `yaml:"streamers" validate:"min=1,dive"`
}

type streamerConfig struct {
	ScreenId string `yaml:"screen_id" validate:"required"`
	Schedule string `yaml:"schedule" validate:"required"`
}

func getDefaultConfig() *cronConfig {
	cronConfig, err := parseConfig(defaultConfigPath)
	if err != nil {
		log.Fatal("Error parsing config file: \n", err)
	}
	return cronConfig
}

func parseConfig(configPath string) (*cronConfig, error) {
	defer func() {
		if r := recover(); r != nil {
			log.Fatalln("Paniced parsing user config: ", r)
		}
	}()

	configData, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	cronConfig := &cronConfig{}
	if err := yaml.Unmarshal(configData, cronConfig); err != nil {
		return nil, err
	}

	return cronConfig, validate.Struct(cronConfig)
}
