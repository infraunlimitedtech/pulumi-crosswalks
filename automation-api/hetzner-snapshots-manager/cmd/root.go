package cmd

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile string
	rootCmd = &cobra.Command{
		Use:   "hetzner-snapshot-manager",
		Short: "manage snapshots based on pulumi preview events",
	}
)

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "/path/to/file for config. Required")
	rootCmd.MarkPersistentFlagRequired("config")
	rootCmd.PersistentFlags().BoolP("debug", "d", false, "debug")
	rootCmd.PersistentFlags().Bool("only-api-server", false, "Run only api server and do not stop it. For testing purposes.")
	rootCmd.PersistentFlags().Int("api-server-port", 0, "default is random")

	viper.BindPFlag("debug", rootCmd.PersistentFlags().Lookup("debug"))
	viper.BindPFlag("api-server-port", rootCmd.PersistentFlags().Lookup("api-server-port"))

	cobra.OnInitialize(initConfig)
}

func initConfig() {
	viper.SetConfigFile(cfgFile)

	if err := viper.ReadInConfig(); err != nil {
		log.Fatal(err)
	}
}
