package main

import (
	"log"
	"os"
	"path"
	"path/filepath"

	"github.com/joho/godotenv"
)

func LoadEnv() {
	p, err := os.Executable()

	if err != nil {
		log.Fatal(err)
	}

	p = filepath.Dir(p)
	err = godotenv.Load(path.Join(p, ".env"))

	if err != nil {
		log.Printf("Error loading .env file from %s", p)
	} else {
		return
	}

	err = godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}
