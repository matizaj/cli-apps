/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
	homedir ​"github.com/mitchellh/go-homedir"
)



// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "todoClient",
	Short: "A Todo API Client",
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.todoClient.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.PersistentFlags().String("api-root", "http://localhost:8080","Host todo api url")

	replacer:=strings.Replace("-", "_")
	viper.SetEnvKeyReplacer(replacer)
	viper.SetEnvPrefix("TODO")

	viper.BindPFlag("api-root", rootCmd.PersistentFlags().Lookup("api-root"))
}


