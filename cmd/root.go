package cmd

import (
	"bufio"
	"fmt"
	"io"
	"log/slog"
	"os"
	"slices"
	"strings"

	"github.com/malev/hola/internals"
	"github.com/malev/hola/logger"
	"github.com/spf13/cobra"
)

const VERSION = "v0.0.4"

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "hola <requests.http>",
	Short: "HTTP Client",
	Long: `hola is an HTTP Client that uses .http files and supports templating
to manage your secrets such as api-keys, api-secrets, etc.

  hola requests.http --number 0
  hola requests.http --number 1 --verbose
  hola requests.http --number 1 --dry-run
	`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return fmt.Errorf("Missing .http files")
		}
		if len(args) > 1 {
			return fmt.Errorf("hola takes only one .http file")
		}
		return nil
	},
	Version: VERSION,
	Run: func(cmd *cobra.Command, args []string) {
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

		line, err := cmd.Flags().GetInt("line")
		if err != nil {
			fmt.Println("Error parsing --line")
			os.Exit(1)
		}

		if line != 0 {
			panic("line not implemented")
		}

		number, err := cmd.Flags().GetInt("number")
		if err != nil {
			fmt.Println("Error parsing --number")
			os.Exit(1)
		}

		if number < 1 {
			fmt.Println("Number should be bigger than 0")
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
			input = strings.TrimSpace(string(data))
		}

		app := internals.NewApp(dryRun, number, line, verbose, maxTimeout, output)
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

		err = app.Send(number)
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
	rootCmd.Flags().IntP("line", "l", 0, "Line number of the request to send")
	rootCmd.Flags().IntP("number", "n", 1, "Number of request to send")
	rootCmd.Flags().
		IntP("max-timeout", "", 0, "Maximum time in seconds that you allow the whole operation to take")
	rootCmd.Flags().BoolP("debug", "d", false, "Enable debug mode")
	rootCmd.Flags().BoolP("dry-run", "", false, "Prevent sending the request")
	rootCmd.Flags().BoolP("verbose", "v", false, "Set output to verbose")
	rootCmd.Flags().StringP("output", "o", "text", "Change the output")
	rootCmd.Flags().BoolP("version", "", false, "Print hola's version")
}
