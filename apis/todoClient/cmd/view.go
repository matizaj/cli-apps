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

// viewCmd represents the view command
var viewCmd = &cobra.Command{
	Use:   "view",
	Short: "Get single item",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("view called")
		hosturl, err:=cmd.Flags().GetString("api-root")
		if err!= nil {
			return err
		}
		return viewAction(os.Stdout, hosturl, args[0])
	},
	Args: cobra.ExactArgs(1),
}

func init() {
	rootCmd.AddCommand(viewCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// viewCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// viewCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
func viewAction(out io.Writer, hosturl string, id string) error {
	item, err := getOne(hosturl, id)
	if err != nil {
		return err
	}
	fmt.Fprintf(os.Stdout, item.Task)
	return nil
}