package main

import (
	"fmt"
	"github.com/turnage/graw/reddit"
	"log"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

const (
	targetUser          = "xhumptyDumptyx"
	maxActivityTimesLen = 250
)

type Bot struct {
	rb            reddit.Bot
	activityTimes []time.Time
}

func NewBot() *Bot {
	return &Bot{
		rb: makeRedditBot(),
	}
}
func parseVariableValue(s string) (val int, ok bool) {
	equalsIndex := strings.IndexByte(s, '=')
	if equalsIndex == -1 {
		return val, false
	}

	trimmed := s[equalsIndex+1:]
	if spaceIndex := strings.IndexByte(s, ' '); spaceIndex != -1 {
		trimmed = s[:spaceIndex]
	}

	var err error
	val, err = strconv.Atoi(trimmed)
	if err != nil {
		log.Printf("String conversion error while parsing variable value: %v", err)
		ok = false
	}
	return
}

func (b *Bot) Replyf(parentName, format string, a ...interface{}) {
	reply := fmt.Sprintf(format, a...)
	log.Printf("Sending reply: "+format, a...)
	b.rb.Reply(parentName, reply)
}

func (b *Bot) Message(m *reddit.Message) (err error) {
	log.Printf("Got message from /u/%s: %s", m.Author, m.Body)

	if m.Author == targetUser {
		settableVarNames := []string{"maxHourlyActivity, maxDailyActivity"}

		for _, varname := range settableVarNames {
			if i := strings.Index(m.Body, varname+"="); i != -1 {
				val, ok := parseVariableValue(m.Body[i:])
				if !ok {
					b.Replyf(m.Name, "Unable to set variable: %s", varname)
					return
				}

				switch varname {
				case "maxHourlyActivity":
					maxHourlyActivity = val
				case "maxDailyActivity":
					maxDailyActivity = val
				default:
					log.Printf(
						"Tried (and failed) to set an unknown variable: %s",
						varname,
					)
					return
				}

				b.Replyf(m.Name, "Successfully set variable '%s' to '%d'", varname, val)
			}
		}
	}

	return
}

func (b *Bot) CommentReply(reply *reddit.Message) (err error) {
	log.Printf("Got comment reply from /u/%s: %s", reply.Author, reply.Body)

	if reply.Author == targetUser {
		b.Replyf(reply.Name, "y'know this reply is still considered reddit "+
			"activity")
		return
	}

	b.Replyf(reply.Name, "no u")
	return
}

func (b *Bot) Mention(mention *reddit.Message) (err error) {
	if mention.Author == targetUser {
		b.Replyf(mention.Name, "are you talking about me behind my back again")
		return
	}

	b.Replyf(mention.Name, "👀")
	return
}

// User-activity-handling...

// markActivity increments the bot's activity count if and only if
func (b *Bot) markActivity() {
	b.activityTimes = append(b.activityTimes, time.Now())

	// Cap the number of recorded activities...
	if len(b.activityTimes) > maxActivityTimesLen {
		b.activityTimes = b.activityTimes[1:]
	}
}

func (b *Bot) countActivitySince(point time.Time) (count int) {
	for i, t := range b.activityTimes {
		if t.Before(point) {
			return i + 1
		}
	}
	return len(b.activityTimes)
}

func randomMessage() string {
	messages := []string{
		"stop hitting that yeet",
		"😤",
		"😒",
		"boi.",
		"stop it, get some help",
		"how will you ever get that anime waifu of your dreams if you keep " +
			"procrastinating like this?",
		"why are you like this",
		"hey 👀",
	}
	return messages[rand.Intn(len(messages))]
}

func (b *Bot) UserPost(post *reddit.Post) (err error) {
	log.Printf("Got post from /u/%s with title: %s", post.Author, post.Title)
	b.markActivity()
	b.Replyf(post.Name, randomMessage())
	return
}

var (
	maxHourlyActivity = 2
	maxDailyActivity  = 7
)

func (b *Bot) UserComment(c *reddit.Comment) (err error) {
	log.Printf("Got comment from /u/%s: %s", c.Author, c.Body)
	b.markActivity()

	lastHourActivity := b.countActivitySince(time.Now().Add(-time.Hour))
	tooManyInLastHour := lastHourActivity > maxHourlyActivity

	now := time.Now()
	lastDayActivity := b.countActivitySince(time.Date(
		now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location(),
	))
	tooManyToday := lastDayActivity > maxDailyActivity

	if tooManyInLastHour {
		b.Replyf(
			c.Name,
			"yo you've been active like %d times in the last hour btw",
			lastHourActivity,
		)

		if tooManyToday {
			b.Replyf(
				c.Name,
				"and also like %d times today so far... maybe slow down a lil' bit?",
				lastDayActivity,
			)
		}

		b.Replyf(c.Name, "sorry but i'm literally programmed to stalk you 👀")
		return
	}

	if tooManyToday {
		b.Replyf(
			c.Name,
			"yo you've been active %d times today btw 👀",
			lastDayActivity,
		)
		return
	}

	b.Replyf(c.Name, randomMessage())
	return
}
