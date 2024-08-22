package main

import (
	"reflect"
	"testing"
)

// TestFindAnagramGroups проверяет функцию FindAnagramGroups.
func TestFindAnagramGroups(t *testing.T) {
	tests := []struct {
		name   string
		input  []string
		output map[string][]string
	}{
		{
			name:  "simple test",
			input: []string{"пятак", "пятка", "тяпка", "листок", "слиток", "столик"},
			output: map[string][]string{
				"пятак":  {"пятак", "пятка", "тяпка"},
				"листок": {"листок", "слиток", "столик"},
			},
		},
		{
			name:   "no anagrams",
			input:  []string{"apple", "banana", "cherry"},
			output: map[string][]string{},
		},
		{
			name:   "all anagrams",
			input:  []string{"abc", "bca", "cab"},
			output: map[string][]string{"abc": {"abc", "bca", "cab"}},
		},
		{
			name:   "some duplicates",
			input:  []string{"dog", "god", "dog", "god", "cat"},
			output: map[string][]string{"dog": {"dog", "dog", "god", "god"}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FindAnagramGroups(tt.input)
			if !reflect.DeepEqual(result, tt.output) {
				t.Errorf("FindAnagramGroups() = %v, want %v", result, tt.output)
			}
		})
	}
}
