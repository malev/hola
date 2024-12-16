package internals

import (
	"fmt"
	"strings"
)

type Header struct {
	Key   string
	Value string
}

type Request struct {
	Title       string
	Method      string
	URL         string
	HttpVersion string
	Headers     []Header
	Body        string
}

func (r *Request) ToString() string {
	output := strings.TrimSpace(fmt.Sprintf("%s %s %s\n", r.Method, r.URL, r.HttpVersion))

	headers := "\n"
	if len(r.Headers) > 0 {
		for _, header := range r.Headers {
			headers += fmt.Sprintf("%s: %s\n", header.Key, header.Value)
		}
	}

	body := ""
	if len(r.Body) > 0 {
		body += "\n" + r.Body + "\n"
	}

	return output + headers + body
}

func NewRequest(title, rawRequest string) *Request {
	request := Request{Title: title}

	lines := strings.Split(rawRequest, "\n")
	firstLine := lines[0]

	request.Method, request.URL, request.HttpVersion = parseFirstLine(firstLine)
	lines = lines[1:]

	i := 0
	for {
		if i >= len(lines) {
			break
		}
		if strings.Contains(lines[i], ":") {
			request.Headers = append(request.Headers, NewHeader(lines[i]))
			i++
		} else {
			break
		}
	}

	if len(lines) > i+1 {
		request.Body = strings.TrimSpace(strings.Join(lines[i+1:], "\n"))
	}

	return &request
}

func NewHeader(rawHeader string) Header {
	parts := strings.Split(rawHeader, ":")

	return Header{
		Key:   strings.TrimSpace(parts[0]),
		Value: strings.TrimSpace(parts[1]),
	}
}

func parseFirstLine(firstLine string) (string, string, string) {
	elements := strings.Split(firstLine, " ")
	if len(elements) == 3 {
		return elements[0], elements[1], elements[2]
	} else {
		return elements[0], elements[1], ""
	}
}

func ExtractRawRequest(content []string, from int) string {
	var lines []string
	for i := from + 1; i < len(content); i++ {
		if strings.Contains(content[i], "###") {
			break
		}

		lines = append(lines, content[i])
	}

	return strings.Join(lines, "")
}
