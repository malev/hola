package internals

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"
)

type Printer interface {
	Meta(*http.Response, time.Duration)
	Headers(http.Header)
	Body([]byte)
}

type JSONPrinter struct{}

type Body struct {
	Type string `json:"_type"`
	Body string `json:"body"`
}

func (p JSONPrinter) Meta(resp *http.Response, elapsed time.Duration) {
	fmt.Printf(
		"{\"_type\": \"meta\", \"status\": \"%s\", \"elapsed\":%f}\n",
		resp.Status,
		elapsed.Seconds(),
	)
}

func (p JSONPrinter) Headers(headers http.Header) {
	fmt.Println("{\"_type\": \"headers\", \"headers\":[]}")
}

func (p JSONPrinter) Body(body []byte) {
	bodyObject := Body{Type: "body", Body: string(body)}
	jsonBody, _ := json.Marshal(bodyObject)

	fmt.Println(string(jsonBody))
}

type TextPrinter struct{ verbose bool }

func (p TextPrinter) Meta(resp *http.Response, elapsed time.Duration) {
	if p.verbose {
		slog.Info(fmt.Sprintf("* Time to response: %s", elapsed))
		slog.Info(fmt.Sprintf("* %s %s", resp.Proto, resp.Status))
	}
}

func (p TextPrinter) Headers(headers http.Header) {
	if p.verbose {
		for header, values := range headers {
			for _, value := range values {
				fmt.Printf("> %s: %s\n", header, value)
			}
		}
		slog.Info("")
	}
}

func (p TextPrinter) Body(body []byte) {
	fmt.Println(string(body))
}
