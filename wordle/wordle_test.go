package wordle

import (
	"testing"
)

func TestGenerate(t *testing.T) {
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
		if !isCapitalLetter(char) {
			t.Errorf("Expected capital letter, but got %c", char)
		}
	}
}

func TestGenerateWithSpecialChars(t *testing.T) {
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
		if !isCapitalLetter(char) && !isSpecialChar(char) {
			t.Errorf("Expected capital letter, but got %c", char)
		}
	}
}

func TestGenerateWithNumbers(t *testing.T) {
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
		if !isCapitalLetter(char) && !isNumber(char) {
			t.Errorf("Expected capital letter, but got %c", char)
		}
	}
}
