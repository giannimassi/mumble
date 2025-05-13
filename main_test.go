package main

import (
	"strings"
	"testing"
)

func TestGeneratePronounceablePassword(t *testing.T) {
	tests := []struct {
		name           string
		length         int
		includeSymbols bool
	}{
		{"default length with symbols", 12, true},
		{"longer password with symbols", 16, true},
		{"password without symbols", 12, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pw := generatePronounceablePassword(tt.length, tt.includeSymbols)

			// Check length
			if len(pw) != tt.length {
				t.Errorf("password length = %d, want %d", len(pw), tt.length)
			}

			// Check for required character types
			hasVowel := false
			hasConsonant := false
			hasDigit := false
			hasSymbol := false

			for _, ch := range pw {
				if strings.ContainsRune(vowels, ch) {
					hasVowel = true
				}
				if strings.ContainsRune(consonants, ch) {
					hasConsonant = true
				}
				if strings.ContainsRune(digits, ch) {
					hasDigit = true
				}
				if strings.ContainsRune(symbolSet, ch) {
					hasSymbol = true
				}
			}

			if !hasVowel {
				t.Error("password missing vowel")
			}
			if !hasConsonant {
				t.Error("password missing consonant")
			}
			if !hasDigit {
				t.Error("password missing digit")
			}
			if tt.includeSymbols && !hasSymbol {
				t.Error("password missing symbol when symbols are required")
			}
			if !tt.includeSymbols && hasSymbol {
				t.Error("password contains symbol when symbols are not required")
			}
		})
	}
}

func TestIsValid(t *testing.T) {
	tests := []struct {
		name           string
		password       string
		requireSymbols bool
		want           bool
	}{
		{"valid with symbols", "futiboda9@", true, true},
		{"valid without symbols", "futiboda9", false, true},
		{"missing vowel", "ftbdk9@", true, false},
		{"missing consonant", "aeiou9@", true, false},
		{"missing digit", "futiboda@", true, false},
		{"missing symbol when required", "futiboda9", true, false},
		{"has symbol when not required", "futiboda9@", false, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isValid(tt.password, tt.requireSymbols); got != tt.want {
				t.Errorf("isValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEstimateEntropy(t *testing.T) {
	tests := []struct {
		name           string
		password       string
		includeSymbols bool
		want           float64
	}{
		{"basic password", "futiboda9", false, 45.4},        // 6 chars * log2(20 consonants + 5 vowels + 8 digits)
		{"password with symbols", "futiboda9@", true, 51.7}, // 6 chars * log2(20 consonants + 5 vowels + 8 digits + 3 symbols)
		{"empty password", "", false, 0.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := estimateEntropy(tt.password, tt.includeSymbols)
			if (got < tt.want-0.01) || (got > tt.want+0.01) {
				t.Errorf("estimateEntropy() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRandomChar(t *testing.T) {
	// Test that randomChar returns characters from the given charset
	charset := "abc"
	iterations := 1000
	seen := make(map[byte]bool)

	for i := 0; i < iterations; i++ {
		ch := randomChar(charset)
		if !strings.ContainsRune(charset, rune(ch)) {
			t.Errorf("randomChar() returned %c, not in charset %s", ch, charset)
		}
		seen[ch] = true
	}

	// Check that we've seen all characters in the charset
	for _, ch := range charset {
		if !seen[byte(ch)] {
			t.Errorf("randomChar() never returned %c from charset %s", ch, charset)
		}
	}
}
