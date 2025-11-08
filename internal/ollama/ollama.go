package ollama

import (
	"bytes"
	"encoding/json"
	"net/http"
)

// OllamaRequest represents the request to the Ollama API
type OllamaRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
}

// GenerateQA sends a prompt to the Ollama API and returns the response
func GenerateQA(model, url, prompt string) (string, error) {
	body, _ := json.Marshal(OllamaRequest{Model: model, Prompt: prompt})
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(body)) // #nosec G107 - URL is from config
	if err != nil {
		return "", err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	var fullResponse string
	dec := json.NewDecoder(resp.Body)
	for {
		var chunk map[string]interface{}
		if err := dec.Decode(&chunk); err != nil {
			break
		}
		if r, ok := chunk["response"].(string); ok {
			fullResponse += r
		}
		if done, ok := chunk["done"].(bool); ok && done {
			break
		}
	}
	return fullResponse, nil
}

// ParseFlashcards parses the Ollama response and returns a list of questions and answers
func ParseFlashcards(response string) ([]map[string]string, error) {
	var qas []map[string]string
	lines := splitLines(response)
	var q, a string
	for _, line := range lines {
		l := trimSpace(line)
		if len(l) < 2 {
			continue
		}
		switch l[:2] {
		case "Q:":
			q = trimSpace(l[2:])
		case "A:":
			a = trimSpace(l[2:])
			if q != "" && a != "" {
				qas = append(qas, map[string]string{"question": q, "answer": a})
				q = "" // Reset for next Q/A pair
			}
		}
	}
	return qas, nil
}

func splitLines(s string) []string {
	var lines []string
	start := 0
	for i, c := range s {
		if c == '\n' {
			lines = append(lines, s[start:i])
			start = i + 1
		}
	}
	if start < len(s) {
		lines = append(lines, s[start:])
	}
	return lines
}

func trimSpace(s string) string {
	for len(s) > 0 && (s[0] == '*' || s[0] == ' ' || s[0] == ':') {
		s = s[1:]
	}
	for len(s) > 0 && (s[len(s)-1] == '*' || s[len(s)-1] == ' ') {
		s = s[:len(s)-1]
	}
	return s
}
