package main

import (
	"context"
	"errors"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Prrromanssss/NewsFeed-bot/internal/config"
	"github.com/Prrromanssss/NewsFeed-bot/internal/fetcher"
	"github.com/Prrromanssss/NewsFeed-bot/internal/notifier"
	"github.com/Prrromanssss/NewsFeed-bot/internal/storage"
	"github.com/Prrromanssss/NewsFeed-bot/internal/summary"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/jmoiron/sqlx"
)

func main() {
	botAPI, err := tgbotapi.NewBotAPI(config.Get().TelegramBotToken)
	if err != nil {
		log.Printf("[ERROR] Failed to create bot: %v", err)
		return
	}

	db, err := sqlx.Connect("postgres", config.Get().DatabaseDSN)
	if err != nil {
		log.Printf("[ERROR] Failed to connect to database: %v", err)
		return
	}
	defer db.Close()

	var (
		articleStorage = storage.NewArticleStorage(db)
		sourceStorage  = storage.NewSourceStorage(db)
		fetcher        = fetcher.NewFetcher(
			articleStorage,
			sourceStorage,
			config.Get().FetchInterval,
			config.Get().FilterKeywords,
		)
		notifier = notifier.NewNotifier(
			articleStorage,
			summary.NewOpenAISummarizer(config.Get().OpenAIKey, config.Get().OpenAIPrompt),
			botAPI,
			config.Get().NotificationInterval,
			2*config.Get().FetchInterval,
			config.Get().TelegramChannelID,
		)
	)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	go func(ctx context.Context) {
		if err := fetcher.Start(ctx); err != nil {
			if !errors.Is(err, context.Canceled) {
				log.Printf("[ERROR] Failed to start fetcher: %v", err)
				return
			}

			log.Printf("[INFO] Fetcher stopped")
		}
	}(ctx)

	if err := notifier.Start(ctx); err != nil {
		if !errors.Is(err, context.Canceled) {
			log.Printf("[ERROR] Failed to start notifier: %v", err)
			return
		}

		log.Printf("[INFO] Notifier stopped")
	}
}
