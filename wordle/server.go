package wordle

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

const numWords = 100

var EasyWordlePreferences = WordlePreferences{
	Length:                 15,
	ContainsCapitalLetters: false,
	ContainsSpecialChars:   false,
	ContainsNumbers:        false,
}

type GuessRequest struct {
	Guess string `json:"guess"`
}

type GuessResponse struct {
	CorrectPositionCount int    `json:"correctPositionCount"`
	PartialMatchCount    int    `json:"partialMatchCount"`
	Feedback             string `json:"feedback"`
}

func SetupServer(server *http.Server, wordle *Wordle) {
	// setup channels
	http.HandleFunc("/wordle/info", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(wordle.WordlePreferences)
	})

	http.HandleFunc("/wordle/guess", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var guessReq GuessRequest
		err := json.NewDecoder(r.Body).Decode(&guessReq)
		if err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		response, err := handleGuess(guessReq, wordle.word, wordle.WordlePreferences)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	})

	// Create array to store word URLs
	wordURLs := make([]string, numWords)

	// Set up HTTP handlers synchronously
	for i := 0; i < numWords; i++ {
		wordEndpoint := fmt.Sprintf("/wordle/guess/word%d", i)
		wordURLs[i] = "http://" + server.Addr + wordEndpoint

		generatedEasyWord := wordle.Generate(EasyWordlePreferences)
		copyWord := generatedEasyWord

		// Create handler for each word endpoint
		http.HandleFunc(wordEndpoint, func(w http.ResponseWriter, r *http.Request) {
			var guessReq GuessRequest
			err := json.NewDecoder(r.Body).Decode(&guessReq)
			if err != nil {
				http.Error(w, "Invalid request body", http.StatusBadRequest)
				return
			}

			response, err := handleGuess(guessReq, copyWord, EasyWordlePreferences)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
		})
	}

	// Create and send endpoints to the easy channel in a goroutine
	go func() {
		for _, wordURL := range wordURLs {
			wordle.EasyWordChannel <- wordURL
		}
		close(wordle.EasyWordChannel)
	}()

	server.Handler = http.DefaultServeMux
}

func handleGuess(guessReq GuessRequest, wordToGuess string, preferences WordlePreferences) (GuessResponse, error) {
	if err := validateGuess(wordToGuess, guessReq.Guess, preferences); err != nil {
		return GuessResponse{}, err
	}

	correctPositionCount, partialMatchCount, feedback := calculateFeedback(wordToGuess, guessReq.Guess)

	return GuessResponse{
		CorrectPositionCount: correctPositionCount,
		PartialMatchCount:    partialMatchCount,
		Feedback:             feedback,
	}, nil
}

func calculateFeedback(secret, guess string) (int, int, string) {
	correctPositionCount := 0
	partialMatchCount := 0
	feedback := ""

	for i := range secret {
		if secret[i] == guess[i] {
			correctPositionCount++
			feedback += "2"
		} else if strings.Contains(secret, string(guess[i])) {
			partialMatchCount++
			feedback += "1"
		} else {
			feedback += "0"
		}
	}

	return correctPositionCount, partialMatchCount, feedback
}

func validateGuess(secret, guess string, pref WordlePreferences) error {
	if len(secret) != len(guess) {
		return errors.New("invalid guess length")
	}

	for _, char := range guess {
		if !pref.ContainsCapitalLetters && isCapitalLetter(char) {
			return errors.New("guess contains capital letters")
		}
		if !pref.ContainsSpecialChars && isSpecialChar(char) {
			return errors.New("guess contains special characters")
		}
		if !pref.ContainsNumbers && isNumber(char) {
			return errors.New("guess contains numbers")
		}
	}
	return nil
}
