package cmd

import (
	"bufio"
	"fmt"
	"io"
	"log/slog"
	"os"
	"slices"

	"github.com/malev/hola/internals"
	"github.com/malev/hola/logger"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "hola <requests.http>",
	Short: "HTTP Client",
	Long: `hola is an HTTP Client that uses .http files and supports templating
to manage your secrets such as api-keys, api-secrets, etc.

  hola requests.http --index 0
  hola requests.http --index 1 --verbose
  hola requests.http --index 1 --dry-run
	`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			return
		}

		configFile, err := cmd.Flags().GetString("config")
		if err != nil {
			slog.Error("Error parsing --config")
			os.Exit(1)
		}

		dryRun, err := cmd.Flags().GetBool("dry-run")
		if err != nil {
			fmt.Println("Error parsing --dry-run")
			os.Exit(1)
		}

		index, err := cmd.Flags().GetInt("index")
		if err != nil {
			fmt.Println("Error parsing --index")
			os.Exit(1)
		}

		maxTimeout, err := cmd.Flags().GetInt("max-timeout")
		if err != nil {
			fmt.Println("Error parsing --max-timeout")
			os.Exit(1)
		}

		verbose, err := cmd.Flags().GetBool("verbose")
		if err != nil {
			fmt.Println("Error parsing --verbose")
			os.Exit(1)
		}

		output, err := cmd.Flags().GetString("output")
		if err != nil {
			fmt.Println("Error parsing --output")
			os.Exit(1)
		}

		validOutputs := []string{"json", "text"}
		if !slices.Contains(validOutputs, output) {
			fmt.Println("Only text and json are supported outputs. Defaulting to text.")
		}

		var input string
		if args[0] == "-" {
			reader := bufio.NewReader(os.Stdin)
			data, _ := io.ReadAll(reader)
			input = string(data)
		}

		app := internals.NewApp(dryRun, index, verbose, maxTimeout, output)
		err = app.LoadConfiguration(configFile)
		if err != nil {
			slog.Error(err.Error())
			os.Exit(1)
		}

		if args[0] == "-" {
			err = app.LoadRequest(input)
		} else {
			err = app.LoadRequests(args[0])
		}
		if err != nil {
			slog.Error(err.Error())
			os.Exit(1)
		}

		err = app.Send(index)
		if err != nil {
			slog.Error(err.Error())
			os.Exit(1)
		}
	},
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		debug, _ := cmd.Flags().GetBool("debug")
		level := slog.LevelInfo
		if debug {
			level = slog.LevelDebug
		}

		handler := logger.NewSimpleHanlder(level)
		logger := slog.New(handler)
		slog.SetDefault(logger)
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringP("config", "c", "config.json", "Configuration file")
	rootCmd.Flags().IntP("index", "", 0, "Index of the request to send")
	rootCmd.Flags().
		IntP("max-timeout", "", 0, "Maximum time in seconds that you allow the whole operation to take")
	rootCmd.Flags().BoolP("debug", "d", false, "Enable debug mode")
	rootCmd.Flags().BoolP("dry-run", "", false, "Prevent sending the request")
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCmd.Flags().BoolP("verbose", "v", false, "Set output to verbose")
	rootCmd.Flags().StringP("output", "o", "text", "Change the output")
}
