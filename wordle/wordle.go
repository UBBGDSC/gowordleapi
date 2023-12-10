package wordle

import (
	"encoding/base64"
	"math/rand"
	"os"
	"strings"
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
	wordlist        []string
}

func NewWordle() *Wordle {
	pref := WordlePreferences{
		Length:                 10,
		ContainsCapitalLetters: false,
		ContainsSpecialChars:   false,
		ContainsNumbers:        false,
	}
	wordlist, err := decodeBase64File("wordle/gdscubb.txt")
	if err != nil {
		// try within the package directory
		wordlist, err = decodeBase64File("gdscubb.txt")
		if err != nil {
			panic(err)
		}
	}
	worlde := &Wordle{pref, "", make(chan string), make(chan string), wordlist}
	word := worlde.Generate(pref)
	worlde.word = word
	return worlde
}

// Generate generates a random word based on the given WordlePreferences.
// If the length preference is 0, it returns a random word from the wordlist.
// Otherwise, it generates a word with the specified length by adding random characters
// with a 50% chance to be placed at the beginning or end.
// The generated word is returned as a string.
func (w *Wordle) Generate(pref WordlePreferences) string {
	randomWord := w.wordlist[rand.Intn(len(w.wordlist))]

	if pref.Length == 0 {
		return randomWord
	}

	// Calculate the number of characters to add
	numCharsToAdd := pref.Length - len(randomWord)

	// Create a slice for the new word
	word := make([]byte, pref.Length)
	beginning := 0
	end := pref.Length - 1

	// Add random characters with 50% chance to be placed at the beginning or end
	for i := 0; i < numCharsToAdd; i++ {
		if rand.Float32() < 0.5 {
			// Add random character at the beginning
			word[beginning] = getRandomChar(pref)
			beginning++
		} else {
			// Add random character at the end
			word[end] = getRandomChar(pref)
			end--
		}
	}

	// Add the random word
	copy(word[beginning:end+1], randomWord)

	return string(word)
}

// SetPreferences sets the preferences for the Wordle game.
// It takes a WordlePreferences struct as input and updates the Wordle's preferences accordingly.
// It also generates a new word based on the updated preferences and assigns it to the Wordle's word field.
func (w *Wordle) SetPreferences(pref WordlePreferences) {
	w.WordlePreferences = pref
	word := w.Generate(pref)
	w.word = word
}

func getRandomChar(pref WordlePreferences) byte {
	var charSet = letters
	if pref.ContainsCapitalLetters {
		charSet += "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	}
	if pref.ContainsSpecialChars {
		charSet += "!@#$%^&*()-_=+[]{}|;:'\",.<>/?"
	}
	if pref.ContainsNumbers {
		charSet += "0123456789"
	}

	return charSet[rand.Intn(len(charSet))]
}

func (w *Wordle) GetPreferences() WordlePreferences {
	return w.WordlePreferences
}

func decodeBase64File(filename string) ([]string, error) {
	// Read the base64-encoded file
	encodedData, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	// Decode the base64-encoded data
	decodedData, err := base64.StdEncoding.DecodeString(string(encodedData))
	if err != nil {
		return nil, err
	}

	// Convert the decoded data to a string and split it into lines
	decodedLines := strings.Split(string(decodedData), "\n")

	// Filter out empty lines
	var decodedWords []string
	for _, line := range decodedLines {
		if line != "" {
			decodedWords = append(decodedWords, line)
		}
	}

	return decodedWords, nil
}
