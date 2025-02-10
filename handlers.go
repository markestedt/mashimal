package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

var animals = []string{
	"Lion", "Tiger", "Elephant", "Giraffe", "Panda", "Kangaroo", "Zebra", "Chimpanzee",
	"Dolphin", "Shark", "Whale", "Cheetah", "Leopard", "Wolf", "Fox", "Bear", "Hippopotamus",
	"Rhinoceros", "Crocodile", "Alligator", "Eagle", "Falcon", "Owl", "Penguin", "Peacock",
	"Parrot", "Snake", "Frog", "Turtle", "Rabbit", "Deer", "Horse", "Camel", "Goat", "Sheep",
	"Cow", "Pig", "Dog", "Cat", "Hamster", "Mouse", "Rat", "Octopus", "Lobster", "Crab",
	"Butterfly", "Bee", "Ant", "Duck", "Chicken", "Turkey",
}

func HandleHome(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/index.html"))

	data := PageData{
		Animals: animals,
	}

	tmpl.Execute(w, data)
}

func HandleGenerate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	animal1 := r.FormValue("animal1")
	animal2 := r.FormValue("animal2")

	prompt2 := fmt.Sprintf(`Create a realistic hybrid animal that combines %s and %s.
	The creature is a single animal, blending the features of the two species seamlessly.
	The create only has one head and one face.
	It looks cute and friendly, with big expressive eyes.
	The hybrid has the same texture all over the body, and that texture is a blend of the two species.
	The hybrid is design to live in a single ecological niche (land, water or air).
	The hybrid is the only subject in the image.
	The hybrid is set in a whimsicall, realistic environment fitting of its ecological niche.`, animal1, animal2)

	prompt := "My prompt contains all details dont add anything else. Here is the prompt: \n"

	prompt += fmt.Sprintf(`A highly detailed, ultra-realistic hybrid animal that combines a %s and a %s. The creature is a single, fully integrated and biologically plausible animal, blending the most iconic traits of both species into a seamless design. It has a cute, friendly appearance with expressive eyes and balanced body proportions. The textures—fur, feathers, scales, or skin—transition naturally between the two animals.

The hybrid fits a single ecological niche (land, water, or air), adapting traits from both animals to suit this niche. The hybrid is the only subject in the image, with no text or additional animals. It is set in a fitting natural environment, such as a forest, desert, ocean, or savanna, with soft, warm lighting that enhances its charm and realism. The animal should NEVER have more than one head.`, animal1, animal2)

	imageData, contentType, err := generateImage(prompt2)

	data := PageData{
		Animals: animals,
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
