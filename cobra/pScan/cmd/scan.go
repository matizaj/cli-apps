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

// scanCmd represents the scan command
var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Run a port scan on the host",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("scan called")

		hostsFile, err := cmd.Flags().GetString("hosts-file")
		if err != nil {
			return err
		}

		ports, err := cmd.Flags().GetIntSlice("ports")
		if err != nil {
			return err
		}

		return scanAction(os.Stdout, hostsFile, ports)
	},
}

func init() {
	rootCmd.AddCommand(scanCmd)
	scanCmd.Flags().IntSliceP("ports", "p", []int{80, 22, 443}, "ports to scan")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// scanCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// scanCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}


func scanAction(out io.Writer, hostsFile string, ports []int) error{
	hl := &scan.HostsList{}
	err := hl.Load(hostsFile)
	if err != nil {
			return err
	}

	result := scan.Run(hl, ports)
	return printResult(out, result)
}

func printResult(out io.Writer, res []scan.Result) error {
	message:=""
	for  _, r := range res {
		message += fmt.Sprintf("%s", r.Host)

		if r.NotFound {
			message+=fmt.Sprintf(" Host not found\n\n")
			continue
		}

		message+= fmt.Sprintln()
		for _, p := range r.PortStates {
			message+= fmt.Sprintf("\t%d: %s\n", p.Port, p.Open)
		}
		message += fmt.Sprintln()
	}
	_, err := fmt.Fprintln(out, message)
	
	return err
}