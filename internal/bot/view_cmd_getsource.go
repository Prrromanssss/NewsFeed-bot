package bot

import (
	"context"
	"strconv"

	"github.com/Prrromanssss/NewsFeed-bot/internal/botkit"
	"github.com/Prrromanssss/NewsFeed-bot/internal/model"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type SourceProvider interface {
	SourceByID(ctx context.Context, id int64) (*model.Source, error)
}

func ViewCmdGetSource(provider SourceProvider) botkit.ViewFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update *tgbotapi.Update) error {
		idStr := update.Message.CommandArguments()

		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			return err
		}

		source, err := provider.SourceByID(ctx, id)
		if err != nil {
			return err
		}

		reply := tgbotapi.NewMessage(update.Message.Chat.ID, formatSource(*source))
		reply.ParseMode = "MarkdownV2"

		if _, err := bot.Send(reply); err != nil {
			return err
		}

		return nil
	}
}
