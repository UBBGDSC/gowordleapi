package wordle

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"strings"
)

const numWords = 50

type GuessRequest struct {
	Guess string `json:"guess"`
}

type GuessResponse struct {
	CorrectPositionCount int    `json:"correctPositionCount"`
	PartialMatchCount    int    `json:"partialMatchCount"`
	Feedback             string `json:"feedback"`
}

func SetupServer(server *http.Server, wordle *Wordle) {
	http.HandleFunc("/wordle/guess", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			handlePostGuess(w, r, wordle.word, wordle.WordlePreferences)
		} else if r.Method == http.MethodGet {
			handleGetGuess(w, r, wordle.WordlePreferences)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
	})

	// Create array to store word URLs
	easyWordURLs := make([]string, numWords)
	hardWordURLs := make([]string, numWords)

	// Set up HTTP handlers synchronously
	for i := 0; i < numWords; i++ {
		wordEndpoint := fmt.Sprintf("/wordle/guess/word%d", i)
		easyWordURLs[i] = "http://" + server.Addr + wordEndpoint

		prefs := WordlePreferences{
			Length:                 0,
			ContainsCapitalLetters: false,
			ContainsSpecialChars:   false,
			ContainsNumbers:        false,
		}
		generatedEasyWord := wordle.Generate(prefs)
		prefs.Length = len(generatedEasyWord)
		copyWord := generatedEasyWord

		// Create handler for each word endpoint
		http.HandleFunc(wordEndpoint, func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodPost {
				handlePostGuess(w, r, copyWord, prefs)
			} else if r.Method == http.MethodGet {
				handleGetGuess(w, r, prefs)
			} else {
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
				return
			}
		})
	}

	// Set up HTTP handlers synchronously for hard words
	for i := 0; i < numWords; i++ {
		hardWordEndpoint := fmt.Sprintf("/wordle/guess/hardWord%d", i)
		hardWordURLs[i] = "http://" + server.Addr + hardWordEndpoint

		// Generate hard word with a random length
		prefs := WordlePreferences{
			Length:                 rand.Intn(10) + 10, // Random length between 10 and 20
			ContainsCapitalLetters: rand.Intn(2) == 0,
			ContainsSpecialChars:   rand.Intn(2) == 0,
			ContainsNumbers:        rand.Intn(2) == 0,
		}
		generatedHardWord := wordle.Generate(prefs)
		prefs.Length = len(generatedHardWord)
		copyWord := generatedHardWord

		// Create handler for each hard word endpoint
		http.HandleFunc(hardWordEndpoint, func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodPost {
				handlePostGuess(w, r, copyWord, prefs)
			} else if r.Method == http.MethodGet {
				handleGetGuess(w, r, prefs)
			} else {
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
				return
			}
		})
	}

	// Create and send endpoints to the easy channel in a goroutine
	go func() {
		for _, wordURL := range easyWordURLs {
			wordle.EasyWordChannel <- wordURL
		}
		close(wordle.EasyWordChannel)
	}()
	go func() {
		for _, wordURL := range hardWordURLs {
			wordle.HardWordChannel <- wordURL
		}
		close(wordle.HardWordChannel)
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

func handlePostGuess(w http.ResponseWriter, r *http.Request, word string, prefs WordlePreferences) {
	var guessReq GuessRequest
	err := json.NewDecoder(r.Body).Decode(&guessReq)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	response, err := handleGuess(guessReq, word, prefs)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func handleGetGuess(w http.ResponseWriter, r *http.Request, prefs WordlePreferences) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(prefs)
}
