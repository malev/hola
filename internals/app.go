package internals

import (
	"bytes"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"
)

type AppConfig struct {
	DryRun  bool
	Index   int
	Verbose bool
}

type App struct {
	AppConfig *AppConfig
	Config    string
	Compiled  string
	Requests  []*Request
}

func NewApp(dryRun bool, index int, verbose bool) *App {
	appConfig := &AppConfig{
		DryRun:  dryRun,
		Index:   index,
		Verbose: verbose,
	}

	return &App{
		AppConfig: appConfig,
	}
}

func (app *App) LoadConfiguration(configFile string) error {
	if !FileExists(configFile) {
		slog.Debug(fmt.Sprintf("%s doesn't exist.", configFile))
		return nil
	}

	stream, err := os.ReadFile(configFile)
	if err != nil {
		return fmt.Errorf("Error reading %s. %v", configFile, err)
	}

	app.Config = strings.TrimSpace(string(stream))
	return nil
}

func (app *App) LoadRequests(filename string) error {
	if filename == "-" {
		slog.Debug("Loading requests from STDIN is not supported yet")
		return fmt.Errorf("Can't load requests")
	}

	if !FileExists(filename) {
		return fmt.Errorf("%s doesn't exist", filename)
	}

	stream, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("Error reading %s. %v", filename, err)
	}

	raw := strings.TrimSpace(string(stream))

	compiler := NewCompiler(raw, app.Config)
	app.Compiled = compiler.Run()

	for _, warning := range compiler.Warnings {
		slog.Debug(warning)
	}

	parser := NewParser(app.Compiled)
	app.Requests, err = parser.Parse()
	if err != nil {
		return fmt.Errorf("Failed parsing requests %v", err)
	}

	return nil
}

func (app *App) printRequest(index int) {
	request := app.Requests[index]
	if app.AppConfig.Verbose {
		slog.Info(request.ToString())
	}
}

func (app App) printResponseInfo(resp *http.Response, elapsed time.Duration) {
	if app.AppConfig.Verbose {
		slog.Info(fmt.Sprintf("* Time to response: %s", elapsed))
		slog.Info(fmt.Sprintf("* %s %s", resp.Proto, resp.Status))
	}
}

func (app *App) printHeaders(headers http.Header) {
	if app.AppConfig.Verbose {
		for header, values := range headers {
			for _, value := range values {
				fmt.Printf("> %s: %s\n", header, value)
			}
		}
		slog.Info("")
	}
}

func (app *App) Send(index int) error {
	request := app.Requests[index]

	app.printRequest(index)

	req, err := http.NewRequest(
		request.Method,
		request.URL,
		bytes.NewBufferString(request.Body),
	)
	if err != nil {
		return fmt.Errorf("Error creating request %v", err)
	}

	for _, header := range request.Headers {
		req.Header.Set(header.Key, header.Value)
	}

	client := &http.Client{}
	start := time.Now()
	resp, err := client.Do(req)
	elapsed := time.Since(start)

	app.printResponseInfo(resp, elapsed)
	app.printHeaders(resp.Header)

	if err != nil {
		return fmt.Errorf("Error sending request %v", err)
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("Error reading body %v", err)
	}

	fmt.Println(string(body))
	return nil
}
