package bot

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"unicode/utf8"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/filters/inlinequery"
	"github.com/d-Rickyy-b/myShrugBot/internal/config"
)

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
func shrugHandler(b *gotgbot.Bot, ctx *ext.Context) error {
	query := ctx.InlineQuery.Query

	shrugs := []string{
		"¯\\_(ツ)_/¯",
		"🤷‍♂️",
	}

	log.Printf("Received query: %s\n", query)

	if query != "" {
		shrugs = append(shrugs, query+" ¯\\_(ツ)_/¯")
		shrugs = append(shrugs, query+" 🤷‍♂️")
	}

	results := make([]gotgbot.InlineQueryResult, len(shrugs))
	for i, shrug := range shrugs {
		result := gotgbot.InlineQueryResultArticle{
			Id:    strconv.Itoa(i), // needed to set a unique string ID for each result
			Title: ellipsis(shrug, 45),
			InputMessageContent: gotgbot.InputTextMessageContent{
				MessageText: shrug,
			},
			HideUrl: true,
		}

		results[i] = result
	}

	_, err := b.AnswerInlineQuery(ctx.InlineQuery.Id, results, nil)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

// StartBot makes all the initializations and configurations for running the Telegram bot
func StartBot(botConfig config.Config) {
	bot, createBotErr := gotgbot.NewBot(botConfig.BotToken, &gotgbot.BotOpts{
		Client: http.Client{},
	})
	if createBotErr != nil {
		log.Println(botConfig)
		log.Fatalln("Something went wrong:", createBotErr)
	}

	dispatcher := ext.NewDispatcher(&ext.DispatcherOpts{
		Error: func(b *gotgbot.Bot, ctx *ext.Context, err error) ext.DispatcherAction {
			log.Println("an error occurred while handling update:", err.Error())
			return ext.DispatcherActionNoop
		},
	})

	updater := ext.NewUpdater(&ext.UpdaterOpts{
		ErrorLog:   nil,
		Dispatcher: dispatcher,
	})

	if botConfig.Webhook.Enabled {
		parsedURL, parseErr := url.Parse(botConfig.Webhook.URL)
		if parseErr != nil {
			log.Fatalln("Can't parse webhook url:", parseErr)
		}
		log.Printf("Starting webhook on '%s:%d%s'...\n", botConfig.Webhook.ListenIP, botConfig.Webhook.ListenPort, botConfig.Webhook.ListenPath)
		// TODO add support for custom certificates
		startErr := updater.StartWebhook(bot, parsedURL.Path, ext.WebhookOpts{
			ListenAddr: fmt.Sprintf("%s:%d", botConfig.Webhook.ListenIP, botConfig.Webhook.ListenPort),
		})
		if startErr != nil {
			panic("failed to start webhook: " + startErr.Error())
		}
		_, setWebhookErr := bot.SetWebhook(botConfig.Webhook.URL, &gotgbot.SetWebhookOpts{})
		if setWebhookErr != nil {
			panic("failed to set webhook: " + setWebhookErr.Error())
		}
	} else {
		log.Println("Start polling...")
		_, _ = bot.DeleteWebhook(nil)
		err := updater.StartPolling(bot, &ext.PollingOpts{DropPendingUpdates: false})
		if err != nil {
			panic("failed to start polling: " + err.Error())
		}
	}

	dispatcher.AddHandler(handlers.NewCommand("start", func(b *gotgbot.Bot, ctx *ext.Context) error {
		welcomeText := fmt.Sprintf("Hi! Thanks for using this bot. You can just type '@%s' in any chat and I'll provide you with some shruggies ¯\\_(ツ)_/¯", b.Username)
		b.SendMessage(ctx.EffectiveUser.Id, welcomeText, nil)
		return nil
	}))

	dispatcher.AddHandler(handlers.NewInlineQuery(inlinequery.All, shrugHandler))

	log.Printf("Bot has been started as @%s...\n", bot.User.Username)
	updater.Idle()
}
