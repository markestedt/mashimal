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

	prompt := fmt.Sprintf(`Create a realistic image of a hybrid animal that combines %[1]s and %[2]s. The animal is referred to as The Hybrid from here on.
	The Hybrid is a single animal, blending the defining anatomy of %[1]s with the texture and surface of %[2]s.
	The Hybrid only has one head and one face.
	The Hybrid looks cute and friendly, with big expressive eyes.
	The Hybrid is designed to live in a single ecological niche (land, water or air).
	The Hybrid is the one and only subject in the image.
	The background is a whimsicall, realistic environment fitting of its ecological niche.`, animal1, animal2)
	imageData, contentType, err := generateImage(prompt)

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
