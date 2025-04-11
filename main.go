package main

import (
	"fmt"
	"os"
	"portal/parser"
	"portal/patcher"

	"github.com/spf13/cobra"
)

var port int
var variablesPath string

func main() {
	var rootCmd = &cobra.Command{Use: "portal", CompletionOptions: cobra.CompletionOptions{HiddenDefaultCmd: true}}

	var parseCmd = &cobra.Command{
		Use:   "parse <path>",
		Short: "Parse the given project",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			variables := parser.ParseProject(args[0], parser.ParseOptions{})
			variables.DumpVariables()
		},
	}

	rootCmd.AddCommand(parseCmd)

	var patcherCmd = &cobra.Command{
		Use:   "patcher [options]",
		Short: "Start patcher Webserver, which serves the dashboard and patches the remote repo",
		Args:  cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			patcher.RunPatcher(port, variablesPath)
		},
	}
	patcherCmd.Flags().StringVarP(&variablesPath, "variables", "v", "./variables.json", "The path to the variables configuration file")
	patcherCmd.Flags().IntVarP(&port, "port", "p", 8080, "The port for the webserver")

	rootCmd.AddCommand(patcherCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
