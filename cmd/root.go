package main

import (
	"log"

	"portal/internal/parser"
	"portal/internal/server"

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
			variables, _, err := parser.ParseProject(args[0], parser.ParseOptions{})
			if err != nil {
				log.Fatalln("Error parsing project: " + err.Error())
			}
			variables.DumpVariables()
		},
	}

	rootCmd.AddCommand(parseCmd)

	var port int

	var patcherCmd = &cobra.Command{
		Use:   "patcher [options]",
		Short: "Start patcher Webserver, which serves the dashboard and patches the remote repo",
		Args:  cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			server.RunServer(port)
		},
	}
	patcherCmd.Flags().IntVarP(&port, "port", "p", 8080, "The port for the webserver")

	patcherCmd.Flags().String("repoOwner", "", "The owner of the Github repo")
	patcherCmd.Flags().String("repoName", "", "The name of the Github repo")
	patcherCmd.Flags().String("pac", "", "Your Github account personal access token")
	patcherCmd.Flags().Bool("preview", true, "Serve a live preview of your repo, requires Github integration")

	viper.BindPFlag("repoOwner", patcherCmd.Flags().Lookup("repoOwner"))
	viper.BindPFlag("repoName", patcherCmd.Flags().Lookup("repoName"))
	viper.BindPFlag("pac", patcherCmd.Flags().Lookup("pac"))
	viper.BindPFlag("servePreview", patcherCmd.Flags().Lookup("preview"))

	rootCmd.AddCommand(patcherCmd)
}

func Execute() {
	rootCmd.Execute()
}
