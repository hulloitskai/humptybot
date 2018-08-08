package main

import (
	"github.com/joho/godotenv"
	"log"
	"math/rand"
)

func init() {
	readEnvFile()

	// Rig up the pseudorandom generator!
	rand.Seed(69420)
}

func readEnvFile() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf(
			"Failed to load environment variables from the .env file: %v",
			err,
		)
	}
}
