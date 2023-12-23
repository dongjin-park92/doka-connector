package configs

import (
	"doka-connector/loggers"
	"fmt"

	"github.com/spf13/viper"
)

type ViperConfig struct {
	*viper.Viper
}

var configs *ViperConfig

func GetConfig() *ViperConfig {
	return configs
}

func init() {
	profile := initProfile()
	configs = readConfig(profile)
}

// 환경 인지
func initProfile() string {
	var profile string
	viper.AutomaticEnv()
	switch viper.GetString("ENV") {
	case "prod":
		profile = "prod"
	case "qa":
		profile = "qa"
	default:
		profile = "dev"
	}
	fmt.Println("doka-connector profile:", profile)
	return profile
}

// 각 환경별로 config.yaml 읽기
func readConfig(profile string) *ViperConfig {
	viperConfig := viper.New()
	viperConfig.AddConfigPath("./configs")
	viperConfig.SetConfigName(profile)
	err := viperConfig.ReadInConfig()
	if err != nil {
		loggers.GlobalLogger.Fatal(err)
	}
	return &ViperConfig{
		Viper: viperConfig,
	}
}
