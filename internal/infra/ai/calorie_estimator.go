package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type EstimationItem struct {
	Name     string `json:"name"`
	Calories int    `json:"calories"`
}

type EstimationResult struct {
	EstimatedCalories int              `json:"estimated_calories"`
	Confidence        string           `json:"confidence"`
	Items             []EstimationItem `json:"items"`
	Observation       string           `json:"observation"`
}

type CalorieEstimator struct {
	apiKey     string
	httpClient *http.Client
}

func NewCalorieEstimator(apiKey string) *CalorieEstimator {
	return &CalorieEstimator{
		apiKey:     apiKey,
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}
}

func (e *CalorieEstimator) Estimate(ctx context.Context, description string) (EstimationResult, error) {
	prompt := fmt.Sprintf(`Você é um nutricionista especializado em análise calórica de refeições.

Estime as calorias da seguinte refeição/alimento: "%s"

Responda APENAS com um JSON válido, sem texto adicional, no seguinte formato:
{
  "estimated_calories": <número inteiro total de calorias>,
  "confidence": "<low|medium|high>",
  "items": [
    {"name": "<nome do item>", "calories": <calorias do item>}
  ],
  "observation": "<observação sobre a estimativa ou dicas nutricionais>"
}

Critérios para confidence:
- high: descrição detalhada com quantidades precisas
- medium: descrição razoável mas com algumas incertezas
- low: descrição vaga ou muito genérica`, description)

	reqBody, err := json.Marshal(map[string]any{
		"contents": []map[string]any{
			{
				"parts": []map[string]string{
					{"text": prompt},
				},
			},
		},
	})
	if err != nil {
		return EstimationResult{}, err
	}

	url := "https://generativelanguage.googleapis.com/v1beta/models/gemini-2.0-flash:generateContent?key=" + e.apiKey
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(reqBody))
	if err != nil {
		return EstimationResult{}, err
	}
	req.Header.Set("content-type", "application/json")

	resp, err := e.httpClient.Do(req)
	if err != nil {
		return EstimationResult{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return EstimationResult{}, fmt.Errorf("Gemini API returned status %d", resp.StatusCode)
	}

	var apiResp struct {
		Candidates []struct {
			Content struct {
				Parts []struct {
					Text string `json:"text"`
				} `json:"parts"`
			} `json:"content"`
		} `json:"candidates"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return EstimationResult{}, err
	}
	if len(apiResp.Candidates) == 0 || len(apiResp.Candidates[0].Content.Parts) == 0 {
		return EstimationResult{}, fmt.Errorf("empty response from Gemini API")
	}

	text := apiResp.Candidates[0].Content.Parts[0].Text

	// Gemini às vezes envolve o JSON em ```json ... ```
	if start := indexOf(text, "{"); start >= 0 {
		if end := lastIndexOf(text, "}"); end >= start {
			text = text[start : end+1]
		}
	}

	var result EstimationResult
	if err := json.Unmarshal([]byte(text), &result); err != nil {
		return EstimationResult{}, fmt.Errorf("failed to parse AI response as JSON: %w", err)
	}

	return result, nil
}

func indexOf(s, substr string) int {
	for i := range s {
		if i+len(substr) <= len(s) && s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}

func lastIndexOf(s, substr string) int {
	result := -1
	for i := range s {
		if i+len(substr) <= len(s) && s[i:i+len(substr)] == substr {
			result = i
		}
	}
	return result
}
