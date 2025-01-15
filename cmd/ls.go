package cmd

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/malev/hola/internals"
	"github.com/spf13/cobra"
)

var lsCmd = &cobra.Command{
	Use:   "ls <requests.http>",
	Short: "List requests available in a file",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			slog.Error("Missing .http file")
			os.Exit(1)
		}

		app := internals.NewApp(false, 0, 0, false, 0, "")
		err := app.LoadRequests(args[0])
		if err != nil {
			slog.Error("Error loading requests")
			os.Exit(1)
		}

		for _, request := range app.Requests {
			fmt.Printf("%s: %s %s\n", request.Title, request.Method, request.URL)
		}
	},
}

func init() {
	rootCmd.AddCommand(lsCmd)
}
