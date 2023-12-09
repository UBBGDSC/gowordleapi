package wordle

import (
	"math/rand"
	"time"
)

type WordlePreferences struct {
	Length                 int
	ContainsCapitalLetters bool
	ContainsSpecialChars   bool
	ContainsNumbers        bool
}

const letters = "abcdefghijklmnopqrstuvwxyz"

func init() {
	rand.Seed(time.Now().UnixNano())
}

type Wordle struct {
	WordlePreferences
	word            string
	EasyWordChannel chan string
	HardWordChannel chan string
}

func NewWordle() *Wordle {
	pref := WordlePreferences{
		Length:                 5,
		ContainsCapitalLetters: false,
		ContainsSpecialChars:   false,
		ContainsNumbers:        false,
	}
	worlde := &Wordle{pref, "", make(chan string), make(chan string)}
	word := worlde.Generate(pref)
	worlde.word = word
	return worlde
}

func (w *Wordle) Generate(pref WordlePreferences) string {
	word := make([]byte, pref.Length)

	for i := 0; i < pref.Length; i++ {
		word[i] = getRandomChar(pref)
	}

	return string(word)
}

func (w *Wordle) SetPreferences(pref WordlePreferences) {
	w.WordlePreferences = pref
	word := w.Generate(pref)
	w.word = word
}

func getRandomChar(pref WordlePreferences) byte {
	var charSet string
	if pref.ContainsCapitalLetters {
		charSet += "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	}
	if pref.ContainsSpecialChars {
		charSet += "!@#$%^&*()-_=+[]{}|;:'\",.<>/?"
	}
	if pref.ContainsNumbers {
		charSet += "0123456789"
	}

	if charSet == "" {
		charSet = letters
	}

	return charSet[rand.Intn(len(charSet))]
}

func (w *Wordle) GetPreferences() WordlePreferences {
	return w.WordlePreferences
}
