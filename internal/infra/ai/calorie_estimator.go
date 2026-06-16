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
		"model":      "claude-haiku-4-5-20251001",
		"max_tokens": 1024,
		"messages": []map[string]string{
			{"role": "user", "content": prompt},
		},
	})
	if err != nil {
		return EstimationResult{}, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://api.anthropic.com/v1/messages", bytes.NewReader(reqBody))
	if err != nil {
		return EstimationResult{}, err
	}
	req.Header.Set("x-api-key", e.apiKey)
	req.Header.Set("anthropic-version", "2023-06-01")
	req.Header.Set("content-type", "application/json")

	resp, err := e.httpClient.Do(req)
	if err != nil {
		return EstimationResult{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return EstimationResult{}, fmt.Errorf("anthropic API returned status %d", resp.StatusCode)
	}

	var apiResp struct {
		Content []struct {
			Type string `json:"type"`
			Text string `json:"text"`
		} `json:"content"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return EstimationResult{}, err
	}
	if len(apiResp.Content) == 0 {
		return EstimationResult{}, fmt.Errorf("empty response from Anthropic API")
	}

	var result EstimationResult
	if err := json.Unmarshal([]byte(apiResp.Content[0].Text), &result); err != nil {
		return EstimationResult{}, fmt.Errorf("failed to parse AI response as JSON: %w", err)
	}

	return result, nil
}
