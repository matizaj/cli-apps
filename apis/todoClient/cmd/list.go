/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"io"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all reported tasks",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("list called")
		hosturl, err := cmd.Flags().GetString("api-root")
		if err != nil {
			return err
		}
		return actionGetAll(os.Stdout, hosturl)
	},
}

func init() {
	rootCmd.AddCommand(listCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
func actionGetAll(out io.Writer, url string) error {
	items, err:= getAll(url)
	if err!= nil {
		return err
	} 

	return printAll(out, items)
}
func printAll(out io.Writer, items []item) error {
	w:=tabwriter.NewWriter(out,3,2,0, ' ',0)
	for k,v:=range items {
		done := "-"
		if v.Done {
			done = "X"
		}
		fmt.Fprintf(w, "%s\t%d\t%s\t\n", done, k+1, v.Task)
	}
	return w.Flush()
}