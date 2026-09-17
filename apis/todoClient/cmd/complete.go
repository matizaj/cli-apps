/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"io"
	"os"
	"strconv"

	"github.com/spf13/cobra"
)

// completeCmd represents the complete command
var completeCmd = &cobra.Command{
	Use:   "complete",
	Short: "Mark item as completed",
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("complete called")

		hosturl, err:=cmd.Flags().GetString("api-root")
		if err != nil {
			return err
		}
		return completeAction(os.Stdout, hosturl, args[0])
	},
}

func init() {
	rootCmd.AddCommand(completeCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// completeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// completeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
func completeAction(out io.Writer, hosturl string, arg string) error {
	id, err := strconv.Atoi(arg)
	if err!= nil {
		return err
	}

	completeItem(hosturl, id)
	item, err:=getOne(hosturl, id)
	if err!= nil {
		return err
	}
	fmt.Fprintln(os.Stdout, item.Task)
	return nil
}