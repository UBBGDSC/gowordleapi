package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strings"

	"github.com/UBBGDSC/gowordleapi/wordle"
)

const letters = "abcdefghijklmnopqrstuvwxyz"
const numbers = "0123456789"
const capitals = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
const specials = "!@#$%^&*()-_=+[]{}|;:'\",.<>/?"

func main() {
	// Create a Wordle instance
	wordleInstance := wordle.NewWordle()

	// Set up HTTP server
	server := &http.Server{
		Addr:    "127.0.0.1:8080", // Change the port as needed
		Handler: http.DefaultServeMux,
	}

	// ping handler
	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	})

	// Set up Wordle server endpoints
	wordle.SetupServer(server, wordleInstance)

	go func() {
		for easyWordUrl := range wordleInstance.EasyWordChannel {
			go easyWordRoutine(easyWordUrl)
		}
		fmt.Println("Easy word channel closed")
	}()
	go func() {
		for hardWordUrl := range wordleInstance.HardWordChannel {
			go easyWordRoutine(hardWordUrl)
		}
		fmt.Println("Hard word channel closed")
	}()

	// Alternative with select and 2 channels:
	// for {
	// 	select {
	// 	case easyWordUrl, ok := <-wordleInstance.EasyWordChannel:
	// 		if !ok {
	// 			log.Println("Easy word channel closed")
	// 			return
	// 		}
	// 		go easyWordRoutine(easyWordUrl)
	// 	case hardWordUrl, ok := <-wordleInstance.HardWordChannel:
	// 		if !ok {
	// 			log.Println("Hard word channel closed")
	// 			return
	// 		}
	// 		go easyWordRoutine(hardWordUrl)
	// 	}
	// }

	// Start the server
	log.Println("Wordle server is running on port 8080...")
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

func easyWordRoutine(wordUrl string) {
	// do a get requeest at wordURL to get the word length from WordlePreferences
	res, err := http.Get(wordUrl)
	if err != nil {
		log.Println("error getting word URL: ", err)
		return
	}

	var wordlePreferences wordle.WordlePreferences
	err = json.NewDecoder(res.Body).Decode(&wordlePreferences)
	if err != nil {
		log.Println("error decoding wordle preferences: ", err)
		return
	}

	easyWordLength := wordlePreferences.Length
	guessRequestBody := wordle.GuessRequest{
		Guess: strings.Repeat("a", easyWordLength),
	}

	for {
		jsonData, err := json.Marshal(guessRequestBody)
		if err != nil {
			log.Println("error marshalling guess request body: ", err)
			return
		}

		res, err := http.Post(wordUrl, "application/json", strings.NewReader(string(jsonData)))
		if err != nil {
			log.Println("error posting to word URL: ", err)
			return
		}

		var guessResponse wordle.GuessResponse
		err = json.NewDecoder(res.Body).Decode(&guessResponse)
		if err != nil {
			log.Println("error decoding guess response: ", err)
			return
		}

		// Check if the guess was successful
		if guessResponse.CorrectPositionCount == easyWordLength {
			log.Println("Successful guess with word", guessRequestBody.Guess, "at", wordUrl)
			break // Exit the loop if successful
		}

		guessRequestBody.Guess = generateNextGuess(guessRequestBody.Guess, guessResponse.Feedback, wordlePreferences)
	}
}

func generateNextGuess(guess string, feedback string, pref wordle.WordlePreferences) string {
	nextGuess := make([]byte, len(guess))
	charset := letters
	for i := 0; i < len(guess); i++ {
		if feedback[i] == '2' {
			nextGuess[i] = guess[i]
		} else {
			if pref.ContainsNumbers {
				charset += numbers
			}
			if pref.ContainsCapitalLetters {
				charset += capitals
			}
			if pref.ContainsSpecialChars {
				charset += specials
			}
			nextGuess[i] = charset[rand.Intn(len(charset))]
		}
	}
	return string(nextGuess)
}
