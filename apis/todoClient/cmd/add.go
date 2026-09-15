/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add new task",
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error{
		fmt.Println("add called")
		hosturl, err := cmd.Flags().GetString("api-root")
		if err != nil {
			return err
		}

		return addAction(os.Stdout, hosturl, args[0])
	},
}

func init() {
	rootCmd.AddCommand(addCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// addCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// addCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func addAction(out io.Writer, url string, taskName string) error {
	err:=addItem(url, taskName)
	if err != nil {
		return err
	}
	if _, err := fmt.Fprintf(out, "Added task name: %s", taskName); err != nil {
		return err
	}
	return nil
}
