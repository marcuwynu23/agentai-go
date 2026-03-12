package core

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)


// geminiGenerateRequest is the JSON body for generateContent.
type geminiGenerateRequest struct {
	Contents []geminiContent `json:"contents"`
}

type geminiContent struct {
	Parts []geminiPart `json:"parts"`
}

type geminiPart struct {
	Text string `json:"text"`
}

// geminiGenerateResponse is the response from generateContent.
type geminiGenerateResponse struct {
	Candidates []struct {
		Content struct {
			Parts []struct {
				Text string `json:"text"`
			} `json:"parts"`
		} `json:"content"`
	} `json:"candidates"`
	Error *struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

// callGeminiGenerateContent sends one request to Gemini and returns the first text part.
// baseURL is the API base (e.g. https://generativelanguage.googleapis.com/v1beta/models).
func callGeminiGenerateContent(ctx context.Context, apiKey, model, baseURL, prompt string) (string, error) {
	if apiKey == "" {
		return "", fmt.Errorf("missing API key")
	}
	baseURL = strings.TrimSuffix(baseURL, "/")
	url := baseURL + "/" + model + ":generateContent?key=" + apiKey
	body := geminiGenerateRequest{
		Contents: []geminiContent{
			{Parts: []geminiPart{{Text: prompt}}},
		},
	}
	raw, err := json.Marshal(body)
	if err != nil {
		return "", err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(raw))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	var out geminiGenerateResponse
	if err := json.Unmarshal(data, &out); err != nil {
		return "", fmt.Errorf("decode response: %w", err)
	}
	if out.Error != nil {
		return "", fmt.Errorf("api error: %s", out.Error.Message)
	}
	if len(out.Candidates) == 0 {
		return "", fmt.Errorf("no candidates in response")
	}
	for _, p := range out.Candidates[0].Content.Parts {
		if p.Text != "" {
			return p.Text, nil
		}
	}
	return "", nil
}

// isRetryableError returns true for rate limit / quota errors.
func isRetryableError(err error) bool {
	if err == nil {
		return false
	}
	s := err.Error()
	return strings.Contains(s, "429") ||
		strings.Contains(s, "RESOURCE_EXHAUSTED") ||
		strings.Contains(s, "quota") ||
		strings.Contains(s, "rate")
}
