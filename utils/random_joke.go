package utils

import (
	"errors"
	"fmt"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
)

// FetchRandomJoke fetches a random joke from the API Ninjas endpoint
func FetchRandomJoke() (string, error) {
	// API endpoint
	url := "https://api.api-ninjas.com/v1/jokes"

	// Retrieve API key from environment variables
	apiKey := os.Getenv("API_NINJAS_KEY")
	if apiKey == "" {
		return "", errors.New("API_NINJAS_KEY environment variable not set")
	}

	// Create a new HTTP request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %v", err)
	}

	// Add API key to headers
	req.Header.Set("X-Api-Key", apiKey)

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to fetch random joke: %v", err)
	}
	defer resp.Body.Close()

	// Check if the response status is OK
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("received non-OK response: %d", resp.StatusCode)
	}

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %v", err)
	}

	// Parse the JSON response as an array
	var response []struct {
		Joke string `json:"joke"`
	}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return "", fmt.Errorf("failed to parse response: %v", err)
	}

	// Check if the array is empty
	if len(response) == 0 {
		return "", errors.New("no jokes found in the response")
	}

	// Return the first joke in the array
	return response[0].Joke, nil
}