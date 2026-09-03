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

// listCmd represents the list command
var listCmd = &cobra.Command{
	Aliases: []string{"l"},
	Use:   "list",
	Short: "List hosts in hosts list",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("list called")
		hostsFile, err := cmd.Flags().GetString("hosts-file")
		if err != nil{
			return err
		}
		return listAction(os.Stdout, hostsFile, args)
	},
}

func init() {
	hostsCmd.AddCommand(listCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func listAction(out io.Writer, file string, args []string) error{
	hl := &scan.HostsList{}
	if err := hl.Load(file); err!= nil {
		return err
	}

	for _, host := range hl.Hosts {
		if _, err := fmt.Fprintln(out, host); err != nil {
			return err
		}
	}
	return nil
}
