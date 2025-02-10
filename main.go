package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

func main() {
	LoadEnv()

	if os.Getenv("OPENAI_API_KEY") == "" {
		log.Fatal("OPENAI_API_KEY environment variable is required")
	}

	http.HandleFunc("/", HandleHome)
	http.HandleFunc("/generate", HandleGenerate)

	log.Println("Server starting on :9393...")
	log.Fatal(http.ListenAndServe(":9393", nil))
}

func generateImage(prompt string) (string, string, error) {
	apiKey := os.Getenv("OPENAI_API_KEY")
	url := "https://api.openai.com/v1/images/generations"

	reqBody := OpenAIRequest{
		Model:          "dall-e-3",
		Prompt:         prompt,
		N:              1,
		Size:           "1024x1024",
		ResponseFormat: "b64_json",
		Quality:        "standard",
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", "", fmt.Errorf("error marshaling request: %v", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", "", fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", "", fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", "", fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
	}

	var result OpenAIResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", "", fmt.Errorf("error decoding response: %v", err)
	}

	if len(result.Data) == 0 {
		return "", "", fmt.Errorf("no images returned from API")
	}

	return result.Data[0].B64, "image/jpeg", nil
}
