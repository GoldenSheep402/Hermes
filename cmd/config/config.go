package config

import (
	"errors"
	"fmt"
	"github.com/GoldenSheep402/Hermes/cmd/server/modList"
	"github.com/GoldenSheep402/Hermes/conf"
	"github.com/GoldenSheep402/Hermes/pkg/fsx"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
	"os"
)

var (
	configPath string
	forceGen   bool
	StartCmd   = &cobra.Command{
		Use:     "config",
		Short:   "Generate config file",
		Example: "jframe config -p ./config.yaml -f",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Generating config...")
			err := GenYamlConfig(configPath, forceGen)
			if err != nil {
				fmt.Println(err.Error())
				os.Exit(1)
			}
		},
	}
)

func init() {
	StartCmd.PersistentFlags().StringVarP(&configPath, "path", "p", "./config.yaml", "Generate config in provided path")
	StartCmd.PersistentFlags().BoolVarP(&forceGen, "force", "f", false, "Force generate config in provided path")
}

func GenYamlConfig(path string, force bool) error {
	if fsx.FileExist(path) && !force {
		return errors.New(path + " already exist, use -f to Force coverage")
	}

	globalConfig := conf.GlobalConfig{MODE: "debug"}
	data, err := yaml.Marshal(&globalConfig)
	if err != nil {
		return errors.New("Error marshaling global config: " + err.Error())
	}

	err = os.WriteFile(path, data, 0644)
	if err != nil {
		return errors.New("Generate file with error: " + err.Error())
	}
	fmt.Printf("Config file config to successfully at: [%v]\n", path)

	file, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return errors.New("Error opening file: " + err.Error())
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			zap.S().Error("Error closing file: " + err.Error())
		}
	}(file)

	// append mod config to globalconfig file
	_conf := make(map[string]interface{})

	needToAppend := false
	if len(modList.ModList) > 0 {
		for _, mod := range modList.ModList {
			if mod.Config() != nil {
				needToAppend = true
				_conf[mod.Name()] = mod.Config()
			}
		}

		if needToAppend {
			data, err = yaml.Marshal(&_conf)
			if err != nil {
				return errors.New("Error marshaling combined config: " + err.Error())
			}

			_, err = file.Write(data)
			if err != nil {
				return errors.New("Error writing to file: " + err.Error())
			}
			fmt.Printf("Append mod config to successfully at: [%v]\n", path)
			return nil
		}
	}

	return nil
}
