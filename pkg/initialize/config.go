package initialize

import (
	"admin/pkg/types"

	"github.com/spf13/viper"
)

func InitConfig(configFile string) *types.Config {
	v := viper.New()
	v.SetConfigFile(configFile)
	if err := v.ReadInConfig(); err != nil {
		panic(err)
	}

	var config types.Config
	if err := v.Unmarshal(&config); err != nil {
		panic(err)
	}

	return &config
}
