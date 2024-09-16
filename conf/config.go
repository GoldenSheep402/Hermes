package conf

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"os"
)

var SysVersion = "dev"

var serveConfig *GlobalConfig

func LoadConfig(configPath ...string) {
	if len(configPath) == 0 || configPath[0] == "" {
		viper.SetConfigName("config")
		viper.AddConfigPath(".")
		viper.AddConfigPath("./config")
	} else {
		viper.SetConfigFile(configPath[0])
	}

	loadConfig := func() {
		newConf := new(GlobalConfig)
		err := viper.ReadInConfig()
		if err != nil {
			fmt.Println("Config Read failed: " + err.Error())
			os.Exit(1)
		}
		err = viper.Unmarshal(newConf)
		if err != nil {
			fmt.Println("Config Unmarshal failed: " + err.Error())
			os.Exit(1)
		}
		serveConfig = newConf
	}

	loadConfig()

	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("Config fileHandle changed: ", e.Name)
		loadConfig()
	})
	viper.WatchConfig()
}

func Get() *GlobalConfig {
	return serveConfig
}
