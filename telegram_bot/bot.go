package telegram_bot

import ( 
  "context"
  "os/signal"
  "os"
  "github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
  "strings"
)

func loadConfig() (string,error) {
  content, err := os.ReadFile("./api_key")
  if err != nil { return "",err }
  api_key := strings.TrimSpace(string(content))
  return api_key ,nil
}

func RunBot() error {
  token,err := loadConfig()
  if err != nil { 
    return err 
  }

  opts := []bot.Option{bot.WithDefaultHandler(handler),}
  ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
  defer cancel()

  b, err := bot.New(token,opts...)
  if err !=nil { return err }
  b.Start(ctx)
  return nil
}

func handler(ctx context.Context, b *bot.Bot, update *models.Update) {
  b.SendMessage(ctx,&bot.SendMessageParams{ChatID:update.Message.Chat.ID, Text:"Received message"})

}
