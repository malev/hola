/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/malev/hola/internals"
	"github.com/spf13/cobra"
)

type Match struct {
	ToBeReplaced string
	Key          string
	Value        string
}

func NewMatch(toBeReplaced, key string) *Match {
	return &Match{ToBeReplaced: toBeReplaced, Key: key}
}

func (m *Match) SetValue(value string) {
	m.Value = value
}

// compileCmd represents the compile command
var compileCmd = &cobra.Command{
	Use:   "compile",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			slog.Error("Missing .http file")
			os.Exit(1)
		}

		configFile, err := cmd.Flags().GetString("config")
		if err != nil {
			slog.Error("Error parsing --config")
			os.Exit(1)
		}

		app := internals.NewApp(false, 0, false, 0)
		err = app.LoadConfiguration(configFile)
		if err != nil {
			slog.Error(err.Error())
			os.Exit(1)
		}

		err = app.LoadRequests(args[0])
		if err != nil {
			slog.Error(err.Error())
			os.Exit(1)
		}

		for _, request := range app.Requests {
			fmt.Println(request.ToString())
		}
	},
}

func init() {
	rootCmd.AddCommand(compileCmd)
	compileCmd.Flags().StringP("config", "c", "config.json", "Configuration file")
}
