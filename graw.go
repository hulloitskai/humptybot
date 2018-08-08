package main

import (
	"github.com/turnage/graw"
	"log"
	"os"
)

func MakeGrawConfig() graw.Config {
	return graw.Config{
		Users:          []string{targetUser},
		CommentReplies: true,
		Mentions:       true,
		Messages:       true,
		Logger:         log.New(os.Stdout, "[GRAW] ", 0),
	}
}
