package telegram_bot

import ( 
  "context"
  "github.com/go-telegram/bot"
  "github.com/go-telegram/bot/models"
)

type Command string

const(
  Help Command = "/help"
)

func handleHelp(ctx context.Context, b *bot.Bot, update *models.Update) {
  b.SendMessage(ctx, &bot.SendMessageParams{ChatID: update.Message.Chat.ID, Text:"help text"})
}

func RegisterHandlers(b *bot.Bot) {
  b.RegisterHandler(bot.HandlerTypeMessageText, string(Help), bot.MatchTypeExact, handleHelp)
}
