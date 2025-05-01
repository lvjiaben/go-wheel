package initialize

import (
	"flag"
	"fmt"

	"github.com/lvjiaben/go-wheel/pkg/global"
	"github.com/lvjiaben/go-wheel/pkg/utils/file"

	"github.com/fsnotify/fsnotify"

	"github.com/spf13/viper"
)

func ViperLoad() {
	path, err := file.GetRootDir()
	if err != nil {
		panic(err)
	}
	var c string
	flag.StringVar(&c, "c", "", "choose config file.")
	if c == "" {
		viper.SetConfigName("config")
	} else {
		viper.SetConfigName(c)
	}
	viper.SetConfigType("yaml")
	viper.AddConfigPath(path)
	err2 := viper.ReadInConfig()
	if err2 != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	viper.WatchConfig()
	viper.OnConfigChange(func(in fsnotify.Event) {
		fmt.Println("Config file changed:", in.Name)

		if err := viper.Unmarshal(&global.CONFIG); err != nil {
			fmt.Println(err)
		}
	})
	if err := viper.Unmarshal(&global.CONFIG); err != nil {
		panic(err)
	}
}
