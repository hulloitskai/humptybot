package main

import (
	"github.com/turnage/graw/reddit"
	"log"
	"os"
)

// userAgent is an application identifier used internally by Reddit to keep
// track of the bot.
const userAgent = "humptybot:0.1.0"

// makeRedditBot configures and instantiates a Reddit bot.
func makeRedditBot() (b reddit.Bot) {
	// Predefine expected variables...
	var (
		user    = "humptybot"
		pass    string
		id      string
		secret  string
		varname string
		ok      bool
	)

	// missingArg notifies user about the missing argument, and exits.
	missingArg := func() {
		log.Fatalf("Missing environment variable: %s", varname)
	}

	// Read environment variables...
	varname = "HB_USER"
	if envuser, ok := os.LookupEnv(varname); ok {
		user = envuser
	}

	varname = "HB_PASS"
	if pass, ok = os.LookupEnv(varname); !ok {
		missingArg()
	}

	varname = "HB_ID"
	if id, ok = os.LookupEnv(varname); !ok {
		missingArg()
	}

	varname = "HB_SECRET"
	if secret, ok = os.LookupEnv(varname); !ok {
		missingArg()
	}

	// Create bot configuration...
	cfg := reddit.BotConfig{
		Agent: userAgent,
		App: reddit.App{
			ID:       id,
			Secret:   secret,
			Username: user,
			Password: pass,
		},
	}

	// Create the bot...
	bot, err := reddit.NewBot(cfg)
	if err != nil {
		log.Fatalf("Error while creating Reddit bot: %v", err)
	}

	return bot
}
