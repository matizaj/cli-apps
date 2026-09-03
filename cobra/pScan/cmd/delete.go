/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"io"
	"matizaj/cli-apps/cobra/pScan/scan"
	"os"

	"github.com/spf13/cobra"
)

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Remove host from list",
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error{
		fmt.Println("delete called")
		filename, err := cmd.Flags().GetString("hosts-file")
		if err != nil {
			return err
		}
		return deleteAction(os.Stdout, filename, args)
	},
}

func init() {
	hostsCmd.AddCommand(deleteCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// deleteCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// deleteCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
func deleteAction(out io.Writer, filename string, args []string)error{
	hl := &scan.HostsList{}

	if err := hl.Load(filename); err != nil {
		return err
	}
	for _, host := range args {
		if err := hl.Remove(host);err!= nil {
			return err
		}
		fmt.Fprintln(out, "Remove host %s", host)
	}
	return hl.Save(filename)
}