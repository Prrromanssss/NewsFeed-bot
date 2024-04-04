package bot

import (
	"context"
	"fmt"
	"strings"

	"github.com/Prrromanssss/NewsFeed-bot/internal/botkit"
	"github.com/Prrromanssss/NewsFeed-bot/internal/botkit/markup"
	"github.com/Prrromanssss/NewsFeed-bot/internal/model"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/samber/lo"
)

type SourceLister interface {
	Sources(ctx context.Context) ([]model.Source, error)
}

func ViewCmdListSource(lister SourceLister) botkit.ViewFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update *tgbotapi.Update) error {
		sources, err := lister.Sources(ctx)
		if err != nil {
			return err
		}

		var (
			sourceInfos = lo.Map(sources, func(source model.Source, _ int) string {
				return formatSource(source)
			})
			msgText = fmt.Sprintf(
				"Список источников \\(всего %d\\):\n\n%s",
				len(sources),
				strings.Join(sourceInfos, "\n\n"),
			)
		)

		reply := tgbotapi.NewMessage(update.Message.Chat.ID, msgText)
		reply.ParseMode = "MarkdownV2"

		if _, err := bot.Send(reply); err != nil {
			return err
		}

		return nil
	}
}

func formatSource(source model.Source) string {
	return fmt.Sprintf(
		"🌐 *%s*\nID: `%d`\nURL фида: %s",
		markup.EscapeForMarkdown(source.Name),
		source.SourceID,
		markup.EscapeForMarkdown(source.FeedUrl),
	)
}
