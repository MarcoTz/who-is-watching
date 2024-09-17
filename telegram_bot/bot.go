package telegram_bot

import ( 
  "context"
  "os/signal"
  "os"
  "bytes"
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

  RegisterHandlers(b)
  b.Start(ctx)
  return nil
}

func handler(ctx context.Context, b *bot.Bot, update *models.Update) {
  photo_data,err := os.ReadFile("./photo.jpg")
  if err != nil { return }
  b.SendPhoto(ctx,&bot.SendPhotoParams{ChatID: update.Message.Chat.ID, Photo:&models.InputFileUpload{Filename:"photo.jpg",Data:bytes.NewReader(photo_data)}})

}
