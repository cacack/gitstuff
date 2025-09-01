package cmd

import (
	"os"

	"gitstuff/internal/verbosity"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var verboseCount int

var rootCmd = &cobra.Command{
	Use:   "gitstuff",
	Short: "A CLI tool for managing GitLab repositories",
	Long: `GitStuff is a command-line tool for managing your GitLab repositories.
It allows you to list repositories, clone them individually or all at once,
and check their status including current branch information.`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.gitstuff.yaml)")
	rootCmd.PersistentFlags().CountVarP(&verboseCount, "verbose", "v", "verbose output (use -v, -vv, -vvv for increasing levels)")

	cobra.OnInitialize(func() {
		verbosity.SetFromCount(verboseCount)
	})
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		viper.AddConfigPath(home)
		viper.AddConfigPath(".")
		viper.SetConfigType("yaml")
		viper.SetConfigName(".gitstuff")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		verbosity.Debug("Using config file: %s", viper.ConfigFileUsed())
	}
}
