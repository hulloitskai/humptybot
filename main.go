package main

import (
	"github.com/turnage/graw"
	"log"
)

func main() {
	log.Println("Starting up HumptyBot...")

	bot := NewBot()
	cfg := MakeGrawConfig()
	_, wait, err := graw.Run(bot, bot.rb, cfg)

	if err != nil {
		log.Fatalf("Couldn't start graw, got an error: %v", err)
	}

	if err := wait(); err != nil {
		log.Printf("A graw handler errorred out: %v", err)
	}
}
