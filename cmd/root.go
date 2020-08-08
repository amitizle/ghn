package cmd

import (
	"fmt"
	"os"
	"path"

	"github.com/spf13/cobra"

	"github.com/amitizle/ghn/internal/config"
	"github.com/amitizle/ghn/pkg/logger"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var (
	cfgFile string

	cfg               *config.Config
	detaultConfigPath = path.Join(".config", "ghn")

	rootCmd = &cobra.Command{
		Use:   "ghn",
		Short: "Github notifications application",
	}
)

// Execute executes the root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		exitWithError(err)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	cobra.OnInitialize(initLog)
	rootCmd.PersistentFlags().StringVar(
		&cfgFile,
		"config",
		"",
		fmt.Sprintf("config file (default is %s", path.Join("$HOME", detaultConfigPath)),
	)
}

func initLog() {
	if err := logger.Init(cfg.Log.Level); err != nil {
		exitWithError(err)
	}
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	cfg = config.New()
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			exitWithError(err)
		}

		viper.AddConfigPath(path.Join(home, detaultConfigPath))
		viper.SetConfigName("config")
	}

	viper.SetEnvPrefix("ghn")
	viper.AutomaticEnv() // read in environment variables that match

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}

	if err := viper.Unmarshal(cfg); err != nil {
		exitWithError(err)
	}
}

func exitWithError(err error) {
	fmt.Println(err)
	os.Exit(1)
}
