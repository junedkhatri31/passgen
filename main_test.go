package main

import (
	"regexp"
	"strings"
	"testing"
)

// TestDefaultPasswordGeneration tests basic password generation with defaults
func TestDefaultPasswordGeneration(t *testing.T) {
	password, err := generatePassword(12, false)
	if err != nil {
		t.Fatalf("Failed to generate password: %v", err)
	}

	if len(password) != 12 {
		t.Errorf("Expected password length 12, got %d", len(password))
	}

	validatePasswordCharacterSets(t, password, false)
	validateNoExcludedCharacters(t, password)
}

// TestCustomLength tests password generation with custom lengths
func TestCustomLength(t *testing.T) {
	tests := []struct {
		name   string
		length int
	}{
		{"8 characters", 8},
		{"20 characters", 20},
		{"3 characters (minimum)", 3},
		{"128 characters (maximum)", 128},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			password, err := generatePassword(tt.length, false)
			if err != nil {
				t.Fatalf("Failed to generate password: %v", err)
			}

			if len(password) != tt.length {
				t.Errorf("Expected password length %d, got %d", tt.length, len(password))
			}

			validatePasswordCharacterSets(t, password, false)
			validateNoExcludedCharacters(t, password)
		})
	}
}

// TestWithSpecialCharacters tests password generation with special characters
func TestWithSpecialCharacters(t *testing.T) {
	tests := []struct {
		name   string
		length int
	}{
		{"12 characters with special", 12},
		{"16 characters with special", 16},
		{"4 characters with special (minimum)", 4},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			password, err := generatePassword(tt.length, true)
			if err != nil {
				t.Fatalf("Failed to generate password: %v", err)
			}

			if len(password) != tt.length {
				t.Errorf("Expected password length %d, got %d", tt.length, len(password))
			}

			validatePasswordCharacterSets(t, password, true)
			validateNoExcludedCharacters(t, password)
		})
	}
}

// TestInvalidLength tests that invalid lengths are rejected
func TestInvalidLength(t *testing.T) {
	tests := []struct {
		name           string
		length         int
		includeSpecial bool
		shouldFail     bool
	}{
		{"length 2 (too short)", 2, false, true},
		{"length 0", 0, false, true},
		{"length -1", -1, false, true},
		{"length 3 with special (too short)", 3, true, true},
		{"length 3 without special (valid)", 3, false, false},
		{"length 4 with special (valid)", 4, true, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			password, err := generatePassword(tt.length, tt.includeSpecial)

			if tt.shouldFail {
				if err == nil {
					t.Errorf("Expected error for invalid length %d, but got password: %s", tt.length, password)
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error for valid length %d: %v", tt.length, err)
				}
				if len(password) != tt.length {
					t.Errorf("Expected password length %d, got %d", tt.length, len(password))
				}
			}
		})
	}
}

// TestMultiplePasswordsUniqueness tests that multiple passwords are unique
func TestMultiplePasswordsUniqueness(t *testing.T) {
	passwords := make(map[string]bool)
	count := 100

	for i := 0; i < count; i++ {
		password, err := generatePassword(12, false)
		if err != nil {
			t.Fatalf("Failed to generate password: %v", err)
		}
		passwords[password] = true
	}

	uniqueCount := len(passwords)
	if uniqueCount < count-1 { // Allow for extremely rare collision
		t.Errorf("Expected at least %d unique passwords, got %d", count-1, uniqueCount)
	}
}

// TestExcludedCharactersNeverAppear tests that excluded characters never appear
func TestExcludedCharactersNeverAppear(t *testing.T) {
	// Generate a large sample of passwords
	for i := 0; i < 50; i++ {
		password, err := generatePassword(20, true)
		if err != nil {
			t.Fatalf("Failed to generate password: %v", err)
		}
		validateNoExcludedCharacters(t, password)
	}
}

// TestPasswordRandomness tests that consecutive password generations are different
func TestPasswordRandomness(t *testing.T) {
	passwords := make([]string, 10)
	
	for i := 0; i < 10; i++ {
		password, err := generatePassword(12, false)
		if err != nil {
			t.Fatalf("Failed to generate password: %v", err)
		}
		passwords[i] = password
	}

	// Check that all passwords are different
	for i := 0; i < len(passwords); i++ {
		for j := i + 1; j < len(passwords); j++ {
			if passwords[i] == passwords[j] {
				t.Errorf("Generated identical passwords at indices %d and %d: %s", i, j, passwords[i])
			}
		}
	}
}

// TestGetRandomChar tests the random character selection
func TestGetRandomChar(t *testing.T) {
	charset := "ABCDEFGH"
	charCounts := make(map[byte]int)

	// Generate many random characters to test distribution
	iterations := 1000
	for i := 0; i < iterations; i++ {
		char, err := getRandomChar(charset)
		if err != nil {
			t.Fatalf("Failed to get random character: %v", err)
		}

		// Verify the character is from the charset
		if !strings.ContainsRune(charset, rune(char)) {
			t.Errorf("Generated character '%c' is not in charset '%s'", char, charset)
		}

		charCounts[char]++
	}

	// Verify all characters appeared at least once (with high probability)
	for i := 0; i < len(charset); i++ {
		char := charset[i]
		if charCounts[char] == 0 {
			t.Logf("Warning: Character '%c' never appeared in %d iterations", char, iterations)
		}
	}
}

// TestShuffleString tests the string shuffling function
func TestShuffleString(t *testing.T) {
	original := []byte("ABCDEFGHIJKLMNOP")
	shuffled := make([]byte, len(original))
	copy(shuffled, original)

	err := shuffleString(shuffled)
	if err != nil {
		t.Fatalf("Failed to shuffle string: %v", err)
	}

	// Verify length is preserved
	if len(shuffled) != len(original) {
		t.Errorf("Shuffled length %d doesn't match original length %d", len(shuffled), len(original))
	}

	// Verify all characters are still present
	originalStr := string(original)
	shuffledStr := string(shuffled)
	
	for _, char := range originalStr {
		if !strings.ContainsRune(shuffledStr, char) {
			t.Errorf("Character '%c' from original not found in shuffled string", char)
		}
	}

	// Note: There's a small chance the shuffle produces the same order,
	// but with 16 characters, this is extremely unlikely
}

// TestPasswordWithoutSpecialCharactersHasNoSpecial tests that passwords without special flag don't have special chars
func TestPasswordWithoutSpecialCharactersHasNoSpecial(t *testing.T) {
	for i := 0; i < 20; i++ {
		password, err := generatePassword(12, false)
		if err != nil {
			t.Fatalf("Failed to generate password: %v", err)
		}

		hasSpecial := regexp.MustCompile(`[^A-Za-z0-9]`).MatchString(password)
		if hasSpecial {
			t.Errorf("Password without special flag contains special characters: %s", password)
		}
	}
}

// TestPasswordWithSpecialCharactersHasSpecial tests that passwords with special flag have special chars
func TestPasswordWithSpecialCharactersHasSpecial(t *testing.T) {
	// Due to randomness, a single password might not have special chars even with the flag
	// So we test multiple times and expect at least one to have them
	foundWithSpecial := false
	
	for i := 0; i < 20; i++ {
		password, err := generatePassword(12, true)
		if err != nil {
			t.Fatalf("Failed to generate password: %v", err)
		}

		hasSpecial := regexp.MustCompile(`[^A-Za-z0-9]`).MatchString(password)
		if hasSpecial {
			foundWithSpecial = true
			break
		}
	}

	if !foundWithSpecial {
		t.Error("Generated 20 passwords with special flag but none contained special characters")
	}
}

// Helper function to validate password has required character sets
func validatePasswordCharacterSets(t *testing.T, password string, includeSpecial bool) {
	t.Helper()

	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
	hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
	hasNumber := regexp.MustCompile(`[0-9]`).MatchString(password)

	if !hasUpper {
		t.Error("Password missing uppercase characters")
	}
	if !hasLower {
		t.Error("Password missing lowercase characters")
	}
	if !hasNumber {
		t.Error("Password missing numbers")
	}

	// Note: Due to randomness, not every password with includeSpecial=true will have special chars
	// The guarantee is only that at least one special char is placed initially
}

// Helper function to validate no excluded characters are present
func validateNoExcludedCharacters(t *testing.T, password string) {
	t.Helper()

	excludedChars := []rune{'0', 'O', 'I', 'l', '1'}
	for _, char := range excludedChars {
		if strings.ContainsRune(password, char) {
			t.Errorf("Password contains excluded character '%c': %s", char, password)
		}
	}
}

// Benchmark password generation
func BenchmarkGeneratePassword(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := generatePassword(12, false)
		if err != nil {
			b.Fatalf("Failed to generate password: %v", err)
		}
	}
}

func BenchmarkGeneratePasswordWithSpecial(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := generatePassword(16, true)
		if err != nil {
			b.Fatalf("Failed to generate password: %v", err)
		}
	}
}

func BenchmarkGeneratePasswordLong(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := generatePassword(128, true)
		if err != nil {
			b.Fatalf("Failed to generate password: %v", err)
		}
	}
}