package internals

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"
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

type Compiler struct {
	Input         string
	Configuration map[string]string
	Matches       []*Match
	Warnings      []string
}

func NewCompiler(input string, rawConfig string) *Compiler {
	var config map[string]string
	if rawConfig == "" {
		rawConfig = "{}"
	}
	jsonData := []byte(rawConfig)

	err := json.Unmarshal(jsonData, &config)
	if err != nil {
		panic("Error parsing json")
	}

	return &Compiler{
		Input:         input,
		Configuration: config,
	}
}

func (c *Compiler) Run() string {
	re := regexp.MustCompile(`{{(.*?)}}`)
	findings := re.FindAllStringSubmatch(c.Input, -1)

	for _, found := range findings {
		c.Matches = append(c.Matches, NewMatch(found[0], found[1]))
	}

	for _, match := range c.Matches {
		if strings.Contains(match.Key, "env|") {
			envVar := match.Key[4:]
			envValue := os.Getenv(envVar)
			if envValue == "" {
				c.Warnings = append(
					c.Warnings,
					fmt.Sprintf("%s missing from environment", envVar),
				)
			}
			match.SetValue(envValue)
			continue
		}

		value, ok := c.Configuration[match.Key]
		if !ok {
			c.Warnings = append(c.Warnings, fmt.Sprintf("%s missing from configuration", match.Key))
		}
		match.SetValue(value)
	}

	output := c.Input
	for _, match := range c.Matches {
		output = strings.Replace(output, match.ToBeReplaced, match.Value, 1)
	}

	return output
}
