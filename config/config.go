package config

import (
	"log"
	"os"

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
	Streamers []*struct {
		ScreenId string `yaml:"screen-id" validate:"required"`
		Schedule string `yaml:"schedule" validate:"required"`
	} `yaml:"streamers" validate:"min=1,dive"`
}

func GetDefaultConfig() *cronConfig {
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

	configData, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	cronConfig := &cronConfig{}
	if err := yaml.Unmarshal(configData, cronConfig); err != nil {
		return nil, err
	}

	return cronConfig, validate.Struct(cronConfig)
}
