package telegram_bot

import ( 
  "rooxo/whoiswatching/types"
  "rooxo/whoiswatching/database"
  "database/sql"
  "context"
  "fmt"
  "github.com/go-telegram/bot"
  "github.com/go-telegram/bot/models"
)

type Command string

const(
  Help Command = "/help"
  PossibleShows Command = "/possible_shows" //TODO

  // Show Commands 
  UpdateShow    Command = "/update_show"    //TODO
  AddShow       Command = "/add_show"       //TODO
  RemoveShow    Command = "/remove_show"    //TODO
  ShowShows     Command = "/show_shows" 

  //Group Commands
  ShowGroups    Command = "/show_groups"    //TODO
  AddGroup      Command = "/add_group"      //TODO
  UpdateEp      Command = "/update_ep"      //TODO
  AddWatcher    Command = "/add_watcher"    //TODO
  RemoveWatcher Command = "/remove_watcher" //TODO

  //Watcher Commands 
  AddPerson     Command = "/add_person"     //TODO
  UpdatePerson  Command = "/update_person"  //TODO
  RemovePerson  Command = "/remove_person"  //TODO
)
func handleHelp(ctx context.Context, b *bot.Bot, update *models.Update) {
  help_text := fmt.Sprintf(`Possible Commands
  %s - %s
  %s - %s
  %s - %s`,
  Help, "Get Help Message", 
  ShowShows, "Show all shows",
  ShowGroups, "Show all groups");
  b.SendMessage(ctx, &bot.SendMessageParams{ChatID: update.Message.Chat.ID, Text:help_text})
}

func handleShowShows(ctx context.Context, b *bot.Bot, update *models.Update){
  db,ok := ctx.Value("database").(*sql.DB)
  if !ok {
    b.SendMessage(ctx, &bot.SendMessageParams{ChatID: update.Message.Chat.ID, Text:"Database connection failed, bot probably needs to be restarted"})
    return 
  }

  shows,err := database.GetAllShows(db)
  if err != nil {
    b.SendMessage(ctx,&bot.SendMessageParams{ChatID:update.Message.Chat.ID, Text: fmt.Sprintf("Could not get shows: %s",err)})
    return 
  }

  var shows_str string 
  for _,show := range shows {
    shows_str += "\n " + show.Name
  }
  b.SendMessage(ctx, &bot.SendMessageParams{ChatID:update.Message.Chat.ID, Text:fmt.Sprintf("All shows: %s", shows_str)})
}

func handleShowGroups(ctx context.Context, b *bot.Bot, update *models.Update){
  db,ok := ctx.Value("database").(*sql.DB)
  if !ok{
    b.SendMessage(ctx, &bot.SendMessageParams{ChatID: update.Message.Chat.ID, Text:"Database connection failed, bot probably needs to be restarted"})
    return 
  }
  groups,err := database.GetAllGroups(db)
  if err != nil { 
    b.SendMessage(ctx,&bot.SendMessageParams{ChatID:update.Message.Chat.ID, Text:fmt.Sprintf("Could not load groups %s", err)})
    return
  }

  var groups_str string 
  for _,group := range groups {
    groups_str += "\n " + types.DisplayGroup(group)
  }
  b.SendMessage(ctx,&bot.SendMessageParams{ChatID:update.Message.Chat.ID, Text:fmt.Sprintf("All Groups:\n %s", groups_str)})

}

func RegisterHandlers(b *bot.Bot) {
  b.RegisterHandler(bot.HandlerTypeMessageText, string(Help), bot.MatchTypeExact, handleHelp)
  b.RegisterHandler(bot.HandlerTypeMessageText, string(ShowShows), bot.MatchTypeExact, handleShowShows)
  b.RegisterHandler(bot.HandlerTypeMessageText, string(ShowGroups), bot.MatchTypeExact, handleShowGroups)

}
