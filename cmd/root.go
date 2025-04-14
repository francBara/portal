package main

import (
	"log"

	"portal/internal/parser"
	"portal/internal/patcher/server"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{Use: "portal", CompletionOptions: cobra.CompletionOptions{HiddenDefaultCmd: true}}

func init() {
	var parseCmd = &cobra.Command{
		Use:   "parse <path>",
		Short: "Parse the given project",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			variables, err := parser.ParseProject(args[0], parser.ParseOptions{})
			if err != nil {
				log.Fatalln("Error parsing project: " + err.Error())
			}
			variables.DumpVariables()
		},
	}

	rootCmd.AddCommand(parseCmd)

	var port int
	var variablesPath string

	var patcherCmd = &cobra.Command{
		Use:   "patcher [options]",
		Short: "Start patcher Webserver, which serves the dashboard and patches the remote repo",
		Args:  cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			server.RunPatcher(port, variablesPath)
		},
	}
	patcherCmd.Flags().StringVarP(&variablesPath, "variables", "v", "./variables.json", "The path to the variables configuration file")
	patcherCmd.Flags().IntVarP(&port, "port", "p", 8080, "The port for the webserver")

	patcherCmd.Flags().String("repoOwner", "", "The owner of the Github repo")
	patcherCmd.Flags().String("repoName", "", "The name of the Github repo")
	patcherCmd.Flags().String("pac", "", "Your Github account personal access token")

	viper.BindPFlag("repoOwner", patcherCmd.Flags().Lookup("repoOwner"))
	viper.BindPFlag("repoName", patcherCmd.Flags().Lookup("repoName"))
	viper.BindPFlag("pac", patcherCmd.Flags().Lookup("pac"))

	rootCmd.AddCommand(patcherCmd)
}

func Execute() {
	rootCmd.Execute()
}
