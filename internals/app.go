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
	DryRun     bool
	Number     int
	Line       int
	Verbose    bool
	MaxTimeout int
}

type App struct {
	AppConfig *AppConfig
	Config    string
	Compiled  string
	Requests  []*Request
	Printer   Printer
}

func NewApp(dryRun bool, number int, line int, verbose bool, maxTimeout int, output string) *App {
	appConfig := &AppConfig{
		DryRun:     dryRun,
		Number:     number,
		Line:       number,
		Verbose:    verbose,
		MaxTimeout: maxTimeout,
	}

	var printer Printer
	if output == "json" {
		printer = JSONPrinter{}
	} else {
		printer = TextPrinter{verbose: verbose}
	}

	return &App{
		AppConfig: appConfig,
		Printer:   printer,
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

func (app *App) LoadRequest(raw string) error {
	var err error
	raw = strings.TrimSpace(raw)

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

func (app *App) LoadRequests(filename string) error {
	if filename == "-" {
		slog.Debug("Loading requests from STDIN is not supported yet")
		return fmt.Errorf("Can't load requests")
	}

	if !FileExists(filename) {
		return fmt.Errorf("File %s doesn't exist", filename)
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

func (app *App) Send(number int) error {
	index := number - 1
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

	req.Header.Set("User-Agent", fmt.Sprintf("hola/v0.0.4"))
	for _, header := range request.Headers {
		req.Header.Set(header.Key, header.Value)
	}

	client := &http.Client{}
	if app.AppConfig.MaxTimeout > 0 {
		client = &http.Client{
			Timeout: time.Second * time.Duration(app.AppConfig.MaxTimeout),
		}
	}

	start := time.Now()
	resp, err := client.Do(req)
	elapsed := time.Since(start)

	app.Printer.Meta(resp, elapsed)
	app.Printer.Headers(resp.Header)

	if err != nil {
		return fmt.Errorf("Error sending request %v", err)
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("Error reading body %v", err)
	}

	app.Printer.Body(body)
	return nil
}
