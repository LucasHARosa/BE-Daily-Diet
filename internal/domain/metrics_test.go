package domain

import "testing"

func TestCalculateBestStreak(t *testing.T) {
	tests := []struct {
		name     string
		input    []bool
		expected int
	}{
		{
			name:     "sequência básica do roadmap",
			input:    []bool{true, true, false, true, true, true},
			expected: 3,
		},
		{
			name:     "todas dentro da dieta",
			input:    []bool{true, true, true, true},
			expected: 4,
		},
		{
			name:     "todas fora da dieta",
			input:    []bool{false, false, false},
			expected: 0,
		},
		{
			name:     "nenhuma refeição",
			input:    []bool{},
			expected: 0,
		},
		{
			name:     "uma refeição dentro",
			input:    []bool{true},
			expected: 1,
		},
		{
			name:     "uma refeição fora",
			input:    []bool{false},
			expected: 0,
		},
		{
			name:     "começa fora e termina dentro",
			input:    []bool{false, false, true, true, true, true, true},
			expected: 5,
		},
		{
			name:     "múltiplas sequências iguais — retorna a maior",
			input:    []bool{true, true, false, true, true, false, true, true},
			expected: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CalculateBestStreak(tt.input)
			if result != tt.expected {
				t.Errorf("CalculateBestStreak(%v) = %d, esperado %d", tt.input, result, tt.expected)
			}
		})
	}
}
