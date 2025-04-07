package main

import (
	"fmt"
	"os"
	"portal/parser"

	"github.com/spf13/cobra"
)

func main() {
	var rootCmd = &cobra.Command{Use: "portal", CompletionOptions: cobra.CompletionOptions{HiddenDefaultCmd: true}}

	var parseCmd = &cobra.Command{
		Use:   "parse <path>",
		Short: "Parse the given project",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			parser.ParseProject(args[0], parser.ParseOptions{})
		},
	}

	rootCmd.AddCommand(parseCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
