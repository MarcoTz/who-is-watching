package bot 

import ( 
  "context"
  "os/signal"
  "os"
  "database/sql"
  "strings"
  "github.com/go-telegram/bot"
)

func loadConfig() (string,error) {
  content, err := os.ReadFile("./api_key")
  if err != nil { return "",err }
  api_key := strings.TrimSpace(string(content))
  return api_key ,nil
}

func RunBot(db *sql.DB) error {
  token,err := loadConfig()
  if err != nil { 
    return err 
  }

  opts := []bot.Option{}
  bg, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
  ctx := context.WithValue(bg,"database",db)
  defer cancel()

  b, err := bot.New(token,opts...)
  if err !=nil { return err }

  RegisterHandlers(b)
  b.Start(ctx)
  return nil
}
