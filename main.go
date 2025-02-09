package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
)

type PageData struct {
	Animals   []string
	ImageData string
	Animal1   string
	Animal2   string
	Error     string
	Partial   bool
}

type OpenAIRequest struct {
	Model          string `json:"model"`
	Prompt         string `json:"prompt"`
	N              int    `json:"n"`
	Size           string `json:"size"`
	Quality        string `json:"quality"`
	ResponseFormat string `json:"response_format"`
}

type OpenAIResponse struct {
	Created int `json:"created"`
	Data    []struct {
		URL string `json:"url"`
		B64 string `json:"b64_json"`
	} `json:"data"`
}

func main() {
	LoadEnv()

	if os.Getenv("OPENAI_API_KEY") == "" {
		log.Fatal("OPENAI_API_KEY environment variable is required")
	}

	http.HandleFunc("/", handleHome)
	http.HandleFunc("/generate", handleGenerate)

	log.Println("Server starting on :8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleHome(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/index.html"))

	data := PageData{
		Animals: []string{
			"Lion", "Tiger", "Elephant", "Giraffe", "Panda", "Kangaroo", "Zebra", "Chimpanzee",
			"Dolphin", "Shark", "Whale", "Cheetah", "Leopard", "Wolf", "Fox", "Bear", "Hippopotamus",
			"Rhinoceros", "Crocodile", "Alligator", "Eagle", "Falcon", "Owl", "Penguin", "Peacock",
			"Parrot", "Snake", "Frog", "Turtle", "Rabbit", "Deer", "Horse", "Camel", "Goat", "Sheep",
			"Cow", "Pig", "Dog", "Cat", "Hamster", "Mouse", "Rat", "Octopus", "Lobster", "Crab",
			"Butterfly", "Bee", "Ant", "Duck", "Chicken", "Turkey",
		},
	}

	tmpl.Execute(w, data)
}

func handleGenerate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	animal1 := r.FormValue("animal1")
	animal2 := r.FormValue("animal2")

	prompt := "My prompt is fully detailed, dont add anything else. Here is the prompt: \n"

	prompt += fmt.Sprintf(`A highly detailed, ultra-realistic hybrid animal that combines a %s and a %s. The creature is a single, fully integrated and biologically plausible animal, blending the most iconic traits of both species into a seamless design. It has a cute, friendly appearance with expressive eyes and balanced body proportions. The textures—fur, feathers, scales, or skin—transition naturally between the two animals.

The hybrid fits a single ecological niche (land, water, or air), adapting traits from both animals to suit this niche. The hybrid is the only subject in the image, with no text or additional animals. It is set in a fitting natural environment, such as a forest, desert, ocean, or savanna, with soft, warm lighting that enhances its charm and realism. The animal should NEVER have more than one head.`, animal1, animal2)

	imageData, contentType, err := generateImage(prompt)

	data := PageData{
		Animals: []string{
			"Lion", "Tiger", "Elephant", "Giraffe", "Panda", "Kangaroo", "Zebra", "Chimpanzee",
			"Dolphin", "Shark", "Whale", "Cheetah", "Leopard", "Wolf", "Fox", "Bear", "Hippopotamus",
			"Rhinoceros", "Crocodile", "Alligator", "Eagle", "Falcon", "Owl", "Penguin", "Peacock",
			"Parrot", "Snake", "Frog", "Turtle", "Rabbit", "Deer", "Horse", "Camel", "Goat", "Sheep",
			"Cow", "Pig", "Dog", "Cat", "Hamster", "Mouse", "Rat", "Octopus", "Lobster", "Crab",
			"Butterfly", "Bee", "Ant", "Duck", "Chicken", "Turkey",
		},

		Animal1: animal1,
		Animal2: animal2,
		Partial: true,
	}

	if err != nil {
		data.Error = "Failed to generate image. Please try again."
		log.Printf("Error generating image: %v", err)
	} else {
		data.ImageData = fmt.Sprintf("data:%s;base64,%s", contentType, imageData)
	}

	// If this is an HTMX request, render only the result partial
	if r.Header.Get("HX-Request") == "true" {
		tmpl := template.Must(template.ParseFiles("templates/index.html"))
		tmpl.ExecuteTemplate(w, "result", data)
	} else {
		tmpl := template.Must(template.ParseFiles("templates/index.html"))
		tmpl.Execute(w, data)
	}
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
