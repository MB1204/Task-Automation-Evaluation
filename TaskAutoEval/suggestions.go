package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type ReplicateRequest struct {
	Input string `json:"input"`
}

type ReplicateResponse struct {
	Predictions []string `json:"predictions"`
}

func Handler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Parse form data
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Failed to parse form data", http.StatusInternalServerError)
		return
	}

	// Combine form data into a single string
	var inputText string
	for key, value := range r.Form {
		inputText += fmt.Sprintf("%s: %s\n", key, value[0])
	}

	// Prepare the payload for Replicate API
	replicateURL := "https://api.replicate.com/v1/predictions"
	replicateAPIKey := os.Getenv("r8_TXuSkiGfodml0LIiRBaSJCZ5pIYG4fH1cNZvG")

	requestBody := ReplicateRequest{
		Input: fmt.Sprintf("Give the best automation suggestions for this task:\n%s", inputText),
	}

	payload, err := json.Marshal(requestBody)
	if err != nil {
		http.Error(w, "Failed to create request payload", http.StatusInternalServerError)
		return
	}

	// Send the request to Replicate API
	req, err := http.NewRequest("POST", replicateURL, bytes.NewBuffer(payload))
	if err != nil {
		http.Error(w, "Failed to create request to Replicate API", http.StatusInternalServerError)
		return
	}
	req.Header.Set("Authorization", "Token "+replicateAPIKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Failed to send request to Replicate API", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Read and parse the response
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Failed to read response from Replicate API", http.StatusInternalServerError)
		return
	}

	var replicateResp ReplicateResponse
	err = json.Unmarshal(respBody, &replicateResp)
	if err != nil {
		http.Error(w, "Failed to parse Replicate API response", http.StatusInternalServerError)
		return
	}

	// Render the suggestions as HTML
	suggestionsHTML := "<h2>AI Suggestions</h2><ul>"
	for _, suggestion := range replicateResp.Predictions {
		suggestionsHTML += fmt.Sprintf("<li>%s</li>", suggestion)
	}
	suggestionsHTML += "</ul>"

	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(suggestionsHTML))
}
