package telegram_bot

import (
	"context"
  "strconv"
	"database/sql"
	"fmt"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"rooxo/whoiswatching/database"
	"rooxo/whoiswatching/types"
	"strings"
)

type Command string

const (
	Help          Command = "/help"
	PossibleShows Command = "/possible_shows" //TODO
  RecommendShow Command = "/recommend" //TODO

	// Show Commands
	AddShow    Command = "/add_show"
	RemoveShow Command = "/remove_show"
	ShowShows  Command = "/show_shows"

	//Group Commands
	ShowGroups    Command = "/show_groups"
	AddGroup      Command = "/add_group" 
	UpdateEp      Command = "/update_ep" 
	AddWatcherGroup    Command = "/join_group"
  RemoveWatcherGroup Command = "/leave_group"
	RemoveGroup Command = "/remove_group" 

	//Watcher Commands
	ShowWatchers  Command = "/show_watchers"
	AddWatcher    Command = "/add_watcher"
	RemoveWatcher  Command = "/remove_watcher"
)

func handleHelp(ctx context.Context, b *bot.Bot, update *models.Update) {
	help_text := fmt.Sprintf(`Possible Commands
  %s - %s
  %s - %s
  %s (%s) - %s
  %s - %s
  %s %s - %s
  %s %s - %s
  %s %s - %s,
  %s %s - %s
  %s %s - %s
  %s %s - %s
  %s %s %s - %s
  %s %s %s - %s,
  %s %s %s - %s,
  %s (%s) - %s`,
		Help, "Get Help Message",
		ShowShows, "Show all shows",
		ShowGroups, "%show_name", "Show all groups (for show)",
		ShowWatchers, "Show all watchers",
    AddWatcher, "%name", "Add new watcher",
    RemoveWatcher, "%name", "Remove watcher",
    AddShow, "%name", "Add new show",
    RemoveShow, "%name", "Remove show", 
    AddGroup, "%name", "Add watchgroup",
    RemoveGroup, "%group_id", "Remove watchgrop",
    UpdateEp, "%group_id", "%episode_nr", "Update episode number for group",
    AddWatcherGroup, "%group_id","%watcher_name", "Add watcher to group",
    RemoveWatcherGroup, "%group_id", "%watcher_name", "Remove watcher from group",
    RecommendShow, "%watcher1,%watcher_2,...", "Recommend show (for watchers)")
	b.SendMessage(ctx, &bot.SendMessageParams{ChatID: update.Message.Chat.ID, Text: help_text})
}

func handleShowShows(ctx context.Context, b *bot.Bot, update *models.Update) {
	db, ok := ctx.Value("database").(*sql.DB)
	if !ok {
		b.SendMessage(ctx, &bot.SendMessageParams{ChatID: update.Message.Chat.ID, Text: "Database connection failed, bot probably needs to be restarted"})
		return
	}

	shows, err := database.GetAllShows(db)
	if err != nil {
		b.SendMessage(ctx, &bot.SendMessageParams{ChatID: update.Message.Chat.ID, Text: fmt.Sprintf("Could not get shows: %s", err)})
		return
	}

	var shows_str string
	for _, show := range shows {
		shows_str += "\n " + show.Name
	}
	b.SendMessage(ctx, &bot.SendMessageParams{ChatID: update.Message.Chat.ID, Text: fmt.Sprintf("All shows: %s", shows_str)})
}

func handleShowGroups(ctx context.Context, b *bot.Bot, update *models.Update) {
	db, ok := ctx.Value("database").(*sql.DB)
	if !ok {
		b.SendMessage(ctx, &bot.SendMessageParams{ChatID: update.Message.Chat.ID, Text: "Database connection failed, bot probably needs to be restarted"})
		return
	}

  show_name := strings.TrimSpace(strings.Replace(update.Message.Text,string(ShowGroups),"",-1))
  groups := make([]types.WatchGroup,0)
  if show_name == "" {
    loaded_groups, err := database.GetAllGroups(db)
    if err != nil {
      b.SendMessage(ctx, &bot.SendMessageParams{ChatID: update.Message.Chat.ID, Text: fmt.Sprintf("Could not load groups: %s", err)})
      return
    }
    groups =loaded_groups
  }else {
    loaded_groups, err := database.GetGroupsByShowName(show_name,db)
    if err!= nil { 
      b.SendMessage(ctx, &bot.SendMessageParams{ChatID: update.Message.Chat.ID, Text: fmt.Sprintf("Could not load groups: %s", err)})
      return
    }
    groups = loaded_groups

  }

	var groups_str string
	for _, group := range groups {
		groups_str += "\n " + types.DisplayGroup(group)
	}
	b.SendMessage(ctx, &bot.SendMessageParams{ChatID: update.Message.Chat.ID, Text: fmt.Sprintf("All Groups:%s", groups_str)})

}

func handleShowWatchers(ctx context.Context, b *bot.Bot, update *models.Update) {
	db, ok := ctx.Value("database").(*sql.DB)
	if !ok {
		b.SendMessage(ctx, &bot.SendMessageParams{ChatID: update.Message.Chat.ID, Text: "Database connection failed, bot probably needs to be restarted"})
		return
	}

	people, err := database.GetAllWatchers(db)
	if err != nil {
		b.SendMessage(ctx, &bot.SendMessageParams{ChatID: update.Message.Chat.ID, Text: fmt.Sprintf("Could not get watchers %s", err)})
	}

	var watcher_str string
	for _, watcher := range people {
		watcher_str += "\n " + watcher.Name
	}
	b.SendMessage(ctx, &bot.SendMessageParams{ChatID: update.Message.Chat.ID, Text: fmt.Sprintf("All watchers:\n %s", watcher_str)})

}

func handleAddWatcher(ctx context.Context, b *bot.Bot, update *models.Update) {
	watcher_name := strings.TrimSpace(strings.Replace(update.Message.Text, string(AddWatcher), "", -1))
  if watcher_name == ""{ 
		b.SendMessage(ctx, &bot.SendMessageParams{ChatID: update.Message.Chat.ID, Text: "Please provide watcher name"})
    return
  }
	db, ok := ctx.Value("database").(*sql.DB)
	if !ok {
		b.SendMessage(ctx, &bot.SendMessageParams{ChatID: update.Message.Chat.ID, Text: "Database connection failed, bot probably needs to be restarted"})
		return
	}

	err := database.AddWatcher(watcher_name, db)
	if err != nil {
		b.SendMessage(ctx, &bot.SendMessageParams{ChatID: update.Message.Chat.ID, Text: fmt.Sprintf("Could not create watcher: %s", err)})
		return
	}

	b.SendMessage(ctx, &bot.SendMessageParams{ChatID: update.Message.Chat.ID, Text: "Successfully added watcher"})

}

func handleRemoveWatcher(ctx context.Context, b *bot.Bot, update *models.Update) {
  watcher_name := strings.TrimSpace(strings.Replace(update.Message.Text,string(RemoveWatcher), "", -1))
  if watcher_name == "" { 
		b.SendMessage(ctx, &bot.SendMessageParams{ChatID: update.Message.Chat.ID, Text: "Please provide watcher name"})
    return
  }
  db, ok := ctx.Value("database").(*sql.DB)
	if !ok {
		b.SendMessage(ctx, &bot.SendMessageParams{ChatID: update.Message.Chat.ID, Text: "Database connection failed, bot probably needs to be restarted"})
		return
	}

  err := database.RemoveWatcher(watcher_name,db)
  if err!=nil{
    b.SendMessage(ctx,&bot.SendMessageParams{ChatID:update.Message.Chat.ID, Text: fmt.Sprintf("Could not remove watcher: %s",err)})
    return
  }

  b.SendMessage(ctx,&bot.SendMessageParams{ChatID:update.Message.Chat.ID, Text:fmt.Sprintf("Successfully removed watcher %s",watcher_name)})
}

func handleAddShow(ctx context.Context, b *bot.Bot, update *models.Update) {
  show_name := strings.TrimSpace(strings.Replace(update.Message.Text,string(AddShow),"",-1))
  if show_name == ""{
    b.SendMessage(ctx, &bot.SendMessageParams{ChatID:update.Message.Chat.ID, Text:"Please provide show name"})
    return
  }

  db, ok := ctx.Value("database").(*sql.DB)
	if !ok {
		b.SendMessage(ctx, &bot.SendMessageParams{ChatID: update.Message.Chat.ID, Text: "Database connection failed, bot probably needs to be restarted"})
		return
	}

  err := database.AddShow(show_name,db)
  if err!=nil { 
    b.SendMessage(ctx,&bot.SendMessageParams{ChatID:update.Message.Chat.ID, Text:fmt.Sprintf("Could not add new show: %s",err)})
    return 
  }
  b.SendMessage(ctx,&bot.SendMessageParams{ChatID:update.Message.Chat.ID, Text:fmt.Sprintf("Successfully added show: %s",show_name)})
  return 
}

func handleRemoveShow(ctx context.Context, b *bot.Bot, update *models.Update){
  show_name := strings.TrimSpace(strings.Replace(update.Message.Text,string(RemoveShow),"",-1))
  if show_name == ""{
    b.SendMessage(ctx, &bot.SendMessageParams{ChatID:update.Message.Chat.ID, Text:"Please provide show name"})
    return
  }

  db, ok := ctx.Value("database").(*sql.DB)
	if !ok {
		b.SendMessage(ctx, &bot.SendMessageParams{ChatID: update.Message.Chat.ID, Text: "Database connection failed, bot probably needs to be restarted"})
		return
	}

  err := database.RemoveShow(show_name,db)
  if err!=nil { 
    b.SendMessage(ctx,&bot.SendMessageParams{ChatID:update.Message.Chat.ID, Text:fmt.Sprintf("Could not remove show: %s",err)})
    return 
  }
  b.SendMessage(ctx,&bot.SendMessageParams{ChatID:update.Message.Chat.ID, Text:fmt.Sprintf("Successfully removed show: %s",show_name)})
  return 

}

func handleAddGroup(ctx context.Context, b *bot.Bot, update *models.Update){
  show_name := strings.TrimSpace(strings.Replace(update.Message.Text, string(AddGroup),"",-1))
  if show_name == ""{
    b.SendMessage(ctx,&bot.SendMessageParams{ChatID:update.Message.Chat.ID, Text:"Please provied show name"})
    return
  }
 
  db, ok := ctx.Value("database").(*sql.DB)
  if !ok {
    b.SendMessage(ctx, &bot.SendMessageParams{ChatID: update.Message.Chat.ID, Text: "Database connection failed, bot probaly needs to be restarted"})
    return
  }

  show, err := database.GetShowByName(show_name,db)
  if err != nil { 
    b.SendMessage(ctx, &bot.SendMessageParams{ChatID: update.Message.Chat.ID, Text: fmt.Sprintf("Could not get show %s: %s",show_name,err)})
    return
  }

  err = database.AddWatchGroup(show.Id, db)
  if err!= nil {
    b.SendMessage(ctx, &bot.SendMessageParams{ChatID: update.Message.Chat.ID, Text: fmt.Sprintf("Could not create watchgroup: %s",err)})
    return
  }

  b.SendMessage(ctx,&bot.SendMessageParams{ChatID:update.Message.Chat.ID, Text: fmt.Sprintf("Successfully created watchgroup for %s",show_name)})

}

func handleRemoveGroup(ctx context.Context, b *bot.Bot, update *models.Update){
  group_id_str := strings.TrimSpace(strings.Replace(update.Message.Text,string(RemoveGroup),"",-1))
  group_id,err := strconv.Atoi(group_id_str)
  if err != nil {
    b.SendMessage(ctx,&bot.SendMessageParams{ChatID:update.Message.Chat.ID, Text:fmt.Sprintf("Please provided group id (see %s",ShowGroups)})
    return
  }

  db, ok := ctx.Value("database").(*sql.DB)
  if !ok {
    b.SendMessage(ctx, &bot.SendMessageParams{ChatID: update.Message.Chat.ID, Text: "Database connection failed, bot probaly needs to be restarted"})
    return
  }

  err = database.RemoveGroup(group_id,db)
  if err != nil {
    b.SendMessage(ctx,&bot.SendMessageParams{ChatID: update.Message.Chat.ID, Text:fmt.Sprintf("Could not remove group: %s",err)})
    return 
  }
  b.SendMessage(ctx,&bot.SendMessageParams{ChatID: update.Message.Chat.ID, Text:"Successfully removed group"})

}

func handleUpdateEp(ctx context.Context, b *bot.Bot, update *models.Update){
  input := strings.TrimSpace(strings.Replace(update.Message.Text,string(UpdateEp),"",-1))
  if input == ""{
    b.SendMessage(ctx,&bot.SendMessageParams{ChatID: update.Message.Chat.ID, Text:"Please provide group id and episode number"})
    return 
  }
  input_sep := strings.Split(input," ")
  if len(input_sep) != 2 {
    b.SendMessage(ctx,&bot.SendMessageParams{ChatID: update.Message.Chat.ID, Text:"Could not parse inputs, please try again"})
    return
  }

  group_id,err := strconv.Atoi(input_sep[0])
  if err != nil { 
    b.SendMessage(ctx,&bot.SendMessageParams{ChatID: update.Message.Chat.ID, Text:"Could not parse group id, please try again"})
    return
  }
  ep_nr,err := strconv.Atoi(input_sep[1])
  if err != nil { 
    b.SendMessage(ctx,&bot.SendMessageParams{ChatID: update.Message.Chat.ID, Text:"Could not parse episode number, please try again"})
    return
  }

  db, ok := ctx.Value("database").(*sql.DB)
  if !ok {
    b.SendMessage(ctx, &bot.SendMessageParams{ChatID: update.Message.Chat.ID, Text: "Database connection failed, bot probaly needs to be restarted"})
    return
  }
  err = database.UpdateGroupEpisode(group_id,ep_nr,db)
  if err != nil {
    b.SendMessage(ctx, &bot.SendMessageParams{ChatID: update.Message.Chat.ID, Text: fmt.Sprintf("Could not update episode number: %s",err)})
    return
  }
  b.SendMessage(ctx, &bot.SendMessageParams{ChatID: update.Message.Chat.ID, Text: fmt.Sprintf("Successfully updated episode number %d for group %d ", group_id,ep_nr)})
}

func handleAddWatcherGroup(ctx context.Context, b *bot.Bot, update *models.Update){
  input := strings.TrimSpace(strings.Replace(update.Message.Text,string(AddWatcherGroup),"",-1))
  if input == ""{
    b.SendMessage(ctx,&bot.SendMessageParams{ChatID: update.Message.Chat.ID, Text:"Please provide group id and watcher name"})
    return 
  }
  input_sep := strings.Split(input," ")
  if len(input_sep) != 2 {
    b.SendMessage(ctx,&bot.SendMessageParams{ChatID: update.Message.Chat.ID, Text:"Could not parse inputs, please try again"})
    return
  }
  group_id,err := strconv.Atoi(input_sep[0])
  if err != nil {
    b.SendMessage(ctx,&bot.SendMessageParams{ChatID: update.Message.Chat.ID, Text:"Could not parse group id, please try again"})
    return
  }
  watcher_name := strings.TrimSpace(input_sep[1])

	db, ok := ctx.Value("database").(*sql.DB)
	if !ok {
		b.SendMessage(ctx, &bot.SendMessageParams{ChatID: update.Message.Chat.ID, Text: "Database connection failed, bot probably needs to be restarted"})
		return
	}

  err = database.AddWatcherGroup(group_id,watcher_name,db)
  if err != nil {
    b.SendMessage(ctx, &bot.SendMessageParams{ChatID: update.Message.Chat.ID, Text: fmt.Sprintf("Could not add watcher to group: %s",err)})
		return
  }
  b.SendMessage(ctx, &bot.SendMessageParams{ChatID: update.Message.Chat.ID, Text: fmt.Sprintf("Successfully added %s to group %d",watcher_name,group_id)})
}

func handleRemoveWatcherGroup(ctx context.Context, b *bot.Bot,update *models.Update){
  input := strings.TrimSpace(strings.Replace(update.Message.Text,string(RemoveWatcherGroup),"",-1))
  if input == ""{
    b.SendMessage(ctx,&bot.SendMessageParams{ChatID: update.Message.Chat.ID, Text:"Please provide group id and watcher name"})
    return 
  }
  input_sep := strings.Split(input," ")
  if len(input_sep) != 2 {
    b.SendMessage(ctx,&bot.SendMessageParams{ChatID: update.Message.Chat.ID, Text:"Could not parse inputs, please try again"})
    return
  }
  group_id,err := strconv.Atoi(input_sep[0])
  if err != nil {
    b.SendMessage(ctx,&bot.SendMessageParams{ChatID: update.Message.Chat.ID, Text:"Could not parse group id, please try again"})
    return
  }
  watcher_name := strings.TrimSpace(input_sep[1])

	db, ok := ctx.Value("database").(*sql.DB)
	if !ok {
		b.SendMessage(ctx, &bot.SendMessageParams{ChatID: update.Message.Chat.ID, Text: "Database connection failed, bot probably needs to be restarted"})
		return
	}

  err = database.RemoveWatcherGroup(group_id,watcher_name,db)
  if err != nil {
    b.SendMessage(ctx, &bot.SendMessageParams{ChatID: update.Message.Chat.ID, Text: fmt.Sprintf("Could not remove watcher from group: %s",err)})
		return
  }
  b.SendMessage(ctx, &bot.SendMessageParams{ChatID: update.Message.Chat.ID, Text: fmt.Sprintf("Successfully removed %s from group %d",watcher_name,group_id)})

}

func handleRecommendation(ctx context.Context, b * bot.Bot, update *models.Update) {
  db, ok := ctx.Value("database").(*sql.DB)
  if !ok {
    b.SendMessage(ctx, &bot.SendMessageParams{ChatID: update.Message.Chat.ID, Text: "Database connection failed, bot probably needs to be restarted"})
    return
  }

  input := strings.TrimSpace(strings.Replace(update.Message.Text,string(RecommendShow),"",-1))
  shows := make([]types.Show,0)
  if input == "" {
    loaded_shows,err := database.GetAllShows(db)
    if err != nil {
      b.SendMessage(ctx, &bot.SendMessageParams{ChatID: update.Message.Chat.ID, Text: fmt.Sprintf("Could not load shows: %s",err)})
      return
    }
    shows = loaded_shows
  }else {
    watchers := strings.Split(input," ")
    loaded_shows, err := database.GetUnwatchedShows(watchers,db)
    if err != nil {
      b.SendMessage(ctx, &bot.SendMessageParams{ChatID: update.Message.Chat.ID, Text:fmt.Sprintf("Could not get shows to watch: %s", err)})
      return
    }
    shows = loaded_shows
  }
    rand_show := types.RandomShow(shows)
    b.SendMessage(ctx, &bot.SendMessageParams{ChatID: update.Message.Chat.ID, Text:fmt.Sprintf("You should watch %s",rand_show.Name)})
    return
}

func RegisterHandlers(b *bot.Bot) {
	b.RegisterHandler(bot.HandlerTypeMessageText, string(Help), bot.MatchTypeExact, handleHelp)
	b.RegisterHandler(bot.HandlerTypeMessageText, string(ShowShows), bot.MatchTypeExact, handleShowShows)
	b.RegisterHandler(bot.HandlerTypeMessageText, string(ShowGroups), bot.MatchTypePrefix, handleShowGroups)
	b.RegisterHandler(bot.HandlerTypeMessageText, string(ShowWatchers), bot.MatchTypeExact, handleShowWatchers)
	b.RegisterHandler(bot.HandlerTypeMessageText, string(AddWatcher), bot.MatchTypePrefix, handleAddWatcher)
  b.RegisterHandler(bot.HandlerTypeMessageText, string(RemoveWatcher), bot.MatchTypePrefix, handleRemoveWatcher)
  b.RegisterHandler(bot.HandlerTypeMessageText, string(AddShow), bot.MatchTypePrefix, handleAddShow)
  b.RegisterHandler(bot.HandlerTypeMessageText, string(RemoveShow), bot.MatchTypePrefix, handleRemoveShow)
  b.RegisterHandler(bot.HandlerTypeMessageText, string(AddGroup), bot.MatchTypePrefix, handleAddGroup)
  b.RegisterHandler(bot.HandlerTypeMessageText, string(RemoveGroup), bot.MatchTypePrefix, handleRemoveGroup)
  b.RegisterHandler(bot.HandlerTypeMessageText, string(UpdateEp), bot.MatchTypePrefix, handleUpdateEp)
  b.RegisterHandler(bot.HandlerTypeMessageText, string(AddWatcherGroup), bot.MatchTypePrefix, handleAddWatcherGroup)
  b.RegisterHandler(bot.HandlerTypeMessageText, string(RemoveWatcherGroup), bot.MatchTypePrefix, handleRemoveWatcherGroup)
  b.RegisterHandler(bot.HandlerTypeMessageText, string(RecommendShow), bot.MatchTypePrefix, handleRecommendation)
}
