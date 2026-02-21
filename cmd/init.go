package cmd

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/GoldenSheep402/Hermes/cmd/config"
	"github.com/GoldenSheep402/Hermes/cmd/create"
	"github.com/GoldenSheep402/Hermes/cmd/server"
)

var rootCmd = &cobra.Command{
	Use:          "hermes",
	SilenceUsage: true,
	Short:        "hermes is a Golang framework with unlimited creativity",
	Example:      "hermes server -c ./config.yaml",
}

func init() {
	rootCmd.AddCommand(server.StartCmd)
	rootCmd.AddCommand(config.StartCmd)
	rootCmd.AddCommand(create.StartCmd)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(-1)
	}
}
