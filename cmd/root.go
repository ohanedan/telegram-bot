package cmd

import (
	"telegram_bot/config"
	"telegram_bot/telegram_bot"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:  "subot",
	RunE: rootFunc,
	Args: cobra.ExactArgs(0),
}

func init() {
	rootCmd.Flags().StringP("config", "c", "config.json",
		"config path")
}

func Execute() {
	rootCmd.Execute()
}

func rootFunc(cmd *cobra.Command, args []string) error {

	configFlag, err := cmd.Flags().GetString("config")
	if err != nil {
		return err
	}

	conf, err := config.NewConfig(configFlag)
	if err != nil {
		return err
	}

	err = telegram_bot.RunBot(conf.Config)
	if err != nil {
		return err
	}

	return nil
}
