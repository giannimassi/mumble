package main

import (
	"crypto/rand"
	"flag"
	"fmt"
	"math"
	"math/big"
	"os/exec"
	"runtime"
	"strings"
)

var (
	length  = flag.Int("length", 12, "Length of the password")
	symbols = flag.Bool("symbols", false, "Include symbols (!@#)")
	copy    = flag.Bool("copy", false, "Copy password to clipboard")
)

var (
	vowels     = "aeiou"
	consonants = "bcdfghjkmnpqrstvwxyz" // removed l
	digits     = "23456789"             // removed 0 and 1
	symbolSet  = "!@#"
)

func main() {
	flag.Parse()

	for {
		pw := generatePronounceablePassword(*length, *symbols)
		if isValid(pw, *symbols) {
			entropy := estimateEntropy(pw, *symbols)
			if entropy >= 60 {
				fmt.Println(pw)
				fmt.Printf("Estimated entropy: %.2f bits\n", entropy)
				if *copy {
					if err := copyToClipboard(pw); err != nil {
						fmt.Println("Failed to copy to clipboard:", err)
					}
				}
				break
			}
		}
	}
}

func generatePronounceablePassword(length int, includeSymbols bool) string {
	var b strings.Builder
	positions := []int{}

	// Build pronounceable base (alternating consonant/vowel)
	for b.Len() < length {
		if b.Len() < length {
			b.WriteByte(randomChar(consonants))
			positions = append(positions, b.Len()-1)
		}
		if b.Len() < length {
			b.WriteByte(randomChar(vowels))
			positions = append(positions, b.Len()-1)
		}
	}

	pw := []byte(b.String()[:length])

	// Insert at least one digit at a random position
	digitPos := -1
	if len(digits) > 0 {
		digitPos = randomInt(len(pw))
		pw[digitPos] = randomChar(digits)
	}

	// Insert at least one symbol at a random position (if needed)
	if includeSymbols {
		symbolPos := randomInt(len(pw))
		// Ensure symbol doesn't overwrite the digit
		for symbolPos == digitPos && len(pw) > 1 {
			symbolPos = randomInt(len(pw))
		}
		pw[symbolPos] = randomChar(symbolSet)
	}

	return string(pw)
}

func randomChar(charset string) byte {
	nBig, _ := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
	return charset[nBig.Int64()]
}

func randomInt(n int) int {
	nBig, _ := rand.Int(rand.Reader, big.NewInt(int64(n)))
	return int(nBig.Int64())
}

func isValid(s string, requireSymbols bool) bool {
	hasVowel := false
	hasConsonant := false
	hasDigit := false
	hasSymbol := false

	for _, ch := range s {
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

	return hasVowel && hasConsonant && hasDigit && (!requireSymbols || hasSymbol)
}

func estimateEntropy(pw string, includeSymbols bool) float64 {
	charsetSize := 0
	if strings.ContainsAny(pw, consonants) {
		charsetSize += len(consonants)
	}
	if strings.ContainsAny(pw, vowels) {
		charsetSize += len(vowels)
	}
	if strings.ContainsAny(pw, digits) {
		charsetSize += len(digits)
	}
	if includeSymbols && strings.ContainsAny(pw, symbolSet) {
		charsetSize += len(symbolSet)
	}
	if charsetSize == 0 {
		return 0
	}
	return float64(len(pw)) * math.Log2(float64(charsetSize))
}

func copyToClipboard(s string) error {
	switch runtime.GOOS {
	case "darwin":
		cmd := exec.Command("pbcopy")
		in, _ := cmd.StdinPipe()
		if err := cmd.Start(); err != nil {
			return err
		}
		_, _ = in.Write([]byte(s))
		_ = in.Close()
		return cmd.Wait()
	case "linux":
		cmd := exec.Command("xclip", "-selection", "clipboard")
		in, _ := cmd.StdinPipe()
		if err := cmd.Start(); err != nil {
			return err
		}
		_, _ = in.Write([]byte(s))
		_ = in.Close()
		return cmd.Wait()
	default:
		return fmt.Errorf("clipboard not supported on this OS")
	}
}
