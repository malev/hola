package internals

import (
	"fmt"
	"strings"
)

type Parser struct {
	Content string
}

func NewParser(content string) *Parser {
	return &Parser{
		Content: content,
	}
}

func (p *Parser) Parse() ([]*Request, error) {
	var titles []string
	var requests []*Request

	if !strings.Contains(p.Content, "###") {
		lines := strings.SplitAfter(p.Content, "\n")
		rawRequest := ExtractRawRequest(lines, 0)
		request := NewRequest("0", rawRequest)
		requests = append(requests, request)

		return requests, nil
	}

	lines := strings.SplitAfter(p.Content, "\n")

	for i, line := range lines {
		if strings.Contains(line, "###") {
			if strings.TrimSpace(line) == "###" {
				titles = append(titles, fmt.Sprintf("%d", len(titles)+1))
			} else {
				title := strings.TrimSpace(strings.ReplaceAll(line, "###", ""))
				titles = append(titles, title)
			}
			rawRequest := ExtractRawRequest(lines, i)
			requests = append(requests, NewRequest(titles[len(titles)-1], rawRequest))
		}
	}

	if len(titles) != len(requests) {
		return requests, fmt.Errorf("Error parsing http requests")
	}

	return requests, nil
}
