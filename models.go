package main

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
