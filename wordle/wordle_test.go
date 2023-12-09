package wordle

import (
	"testing"
)

func TestGenerateWithNoSpecialNoNumbers(t *testing.T) {
	pref := WordlePreferences{
		Length:                 5,
		ContainsCapitalLetters: true,
		ContainsSpecialChars:   false,
		ContainsNumbers:        false,
	}

	w := NewWordle()
	word := w.Generate(pref)

	if len(word) != pref.Length {
		t.Errorf("Expected word length %d, but got %d", pref.Length, len(word))
	}

	for _, char := range word {
		if isNumber(char) || isSpecialChar(char) {
			t.Errorf("Did not expect number or special character, but got %c", char)
		}
	}
}

func TestGenerateWithNoNumbers(t *testing.T) {
	pref := WordlePreferences{
		Length:                 5,
		ContainsCapitalLetters: true,
		ContainsSpecialChars:   true,
		ContainsNumbers:        false,
	}

	w := NewWordle()
	word := w.Generate(pref)

	if len(word) != pref.Length {
		t.Errorf("Expected word length %d, but got %d", pref.Length, len(word))
	}

	for _, char := range word {
		if isNumber(char) {
			t.Errorf("Did not expect numbers in word, but got %c", char)
		}
	}
}

func TestGenerateWithNoSpecials(t *testing.T) {
	pref := WordlePreferences{
		Length:                 5,
		ContainsCapitalLetters: true,
		ContainsSpecialChars:   false,
		ContainsNumbers:        true,
	}

	w := NewWordle()
	word := w.Generate(pref)

	if len(word) != pref.Length {
		t.Errorf("Expected word length %d, but got %d", pref.Length, len(word))
	}

	for _, char := range word {
		if isSpecialChar(char) {
			t.Errorf("Did not expect special characters in word, but got %c", char)
		}
	}
}
func TestGenerateMultipleRandomWords(t *testing.T) {
	pref := WordlePreferences{
		Length:                 5,
		ContainsCapitalLetters: true,
		ContainsSpecialChars:   false,
		ContainsNumbers:        false,
	}

	w := NewWordle()

	for i := 0; i < 100; i++ {
		word := w.Generate(pref)

		if len(word) != pref.Length {
			t.Errorf("Expected word length %d, but got %d", pref.Length, len(word))
		}

		for _, char := range word {
			if isNumber(char) || isSpecialChar(char) {
				t.Errorf("Expected only capital and lowercase letters, but got %c", char)
			}
		}
	}
}
