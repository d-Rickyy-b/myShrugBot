package bot

import (
	"fmt"
	"log"
	"strconv"
	"time"
	"unicode/utf8"

	"github.com/d-Rickyy-b/myShrugBot/internal/config"
	tb "gopkg.in/tucnak/telebot.v2"
)

// define bot in global context to be accessible by all handlers
var bot *tb.Bot

// ellipsis takes a text and an integer to generate a new string that's <= 'max' characters (not bytes) long.
// Excess characters will be removed from the middle and replaced with "..."
func ellipsis(text string, max int) string {
	// subtract length of "..."
	max -= 5

	if max <= 5 {
		max = 5
	}

	if len(text) <= max {
		return text
	}

	charCount := 0

	endFirst := max / 2
	startLast := utf8.RuneCountInString(text) - (max / 2)
	result := ""

	for _, r := range text {
		charCount++

		if charCount > endFirst && charCount <= startLast {
			continue
		} else if charCount == endFirst {
			result += string(r)
			result += " ... "
			continue
		}

		result += string(r)
	}

	return result
}

// shrugHandler handles all the Telegram inline queries.
// It responds to incoming inline queries with a set of different options.
// 1) The classic shrug
// 2) The (male) shrug emoji
// 3) The text the user entered + classic shrug
// 4) The text the user entered + shrug emoji
func shrugHandler(q *tb.Query) {
	shrugs := []string{
		"¯\\_(ツ)_/¯",
		"🤷‍♂️",
	}

	if q.Text != "" {
		shrugs = append(shrugs, q.Text+" ¯\\_(ツ)_/¯")
		shrugs = append(shrugs, q.Text+" 🤷‍♂️")
	}

	results := make(tb.Results, len(shrugs)) // []tb.Result
	for i, shrug := range shrugs {
		result := &tb.ArticleResult{
			ResultBase: tb.ResultBase{},
			Title:      ellipsis(shrug, 45),
			Text:       shrug,
			HideURL:    true,
		}

		results[i] = result
		// needed to set a unique string ID for each result
		results[i].SetResultID(strconv.Itoa(i))
	}

	err := bot.Answer(q, &tb.QueryResponse{
		Results:   results,
		CacheTime: 60, // a minute
	})

	if err != nil {
		log.Println(err)
	}
}

// StartBot makes all the initializations and configurations for running the Telegram bot
func StartBot(c config.Config) {
	var poller tb.Poller

	if c.Webhook.Enabled {
		// Currently, we don't support all the fields of Webhook.
		// We could add stuff like TLS/Cert config later still
		poller = &tb.Webhook{
			Listen: c.Webhook.Listen,
			Endpoint: &tb.WebhookEndpoint{
				PublicURL: c.Webhook.Url,
			},
		}
		log.Println("Using webhook from config!")
	} else {
		poller = &tb.LongPoller{Timeout: 10 * time.Second}
		log.Println("Using long polling as fallback!")
	}

	b, err := tb.NewBot(tb.Settings{
		Token:  c.Token,
		Poller: poller,
	})

	if err != nil {
		log.Fatal(err)
		return
	}
	bot = b

	b.Handle("/start", func(m *tb.Message) {
		welcomeText := fmt.Sprintf("Hi! Thanks for using this bot. You can just type '@%s' in any chat and I'll provide you with some shruggies ¯\\_(ツ)_/¯", b.Me.Username)
		b.Send(m.Sender, welcomeText)
	})

	b.Handle(tb.OnQuery, shrugHandler)

	// We must remove the webhook if we want to use Polling. Otherwise we'll run into an error and can't poll new updates.
	if !c.Webhook.Enabled {
		b.RemoveWebhook()
	}
	log.Printf("Starting @%s", b.Me.Username)
	b.Start()
}
