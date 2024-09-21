package bot 

import (
	"context"
  "strconv"
	"database/sql"
	"fmt"
  "os/exec"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"rooxo/whoiswatching/database"
	"rooxo/whoiswatching/types"
	"strings"
)

const SEP = ","

type Command string

const (
	Help          Command = "/help"
	PossibleShows Command = "/possible_shows" 
  RecommendShow Command = "/recommend"
  PushChanges   Command = "/push"

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
  MarkDone Command = "/finish_show"
  MarkNotDone Command = "/unfinish_show"

	//Watcher Commands
	ShowWatchers  Command = "/show_watchers"
	AddWatcher    Command = "/add_watcher"
	RemoveWatcher  Command = "/remove_watcher"
)

func handleHelp(ctx context.Context, b *bot.Bot, update *models.Update) {
	help_text := "Possible Commands\n"
  help_text += fmt.Sprintf("%s - %s\n", Help, "Get Help Message")
  help_text += fmt.Sprintf("%s - %s\n",ShowShows, "Show all shows")
  help_text += fmt.Sprintf("%s (%s) - %s\n",ShowGroups, "%show_name", "Show all groups (for show)")
  help_text += fmt.Sprintf("%s - %s\n",ShowWatchers, "Show all watchers")
  help_text += fmt.Sprintf("%s %s - %s\n",AddWatcher, "%watcher_name", "Add new watcher")
  help_text += fmt.Sprintf("%s %s - %s\n",RemoveWatcher, "%watcher_name", "Remove watcher")
  help_text += fmt.Sprintf("%s %s - %s\n",AddShow, "%show_name", "Add new show")
  help_text += fmt.Sprintf("%s %s - %s\n",RemoveShow, "%show_name", "Remove show")
  help_text += fmt.Sprintf("%s %s - %s\n",AddGroup, "%show_name", "Add watchgroup")
  help_text += fmt.Sprintf("%s %s - %s\n",    RemoveGroup, "%group_id", "Remove watchgrop")
  help_text += fmt.Sprintf("%s %s%s%s - %s\n",UpdateEp, "%group_id", SEP, "%episode_nr", "Update episode number for group")
  help_text += fmt.Sprintf("%s %s%s%s - %s\n",RemoveWatcherGroup, "%group_id", SEP, "%watcher_name", "Remove watcher from group")
  help_text += fmt.Sprintf("%s %s%s%s - %s\n",AddWatcherGroup, "%group_id", SEP, "%watcher_name", "Add watcher to group")
  help_text += fmt.Sprintf("%s (%s) - %s\n",RecommendShow, "%watcher_1 %watcher_2 ...", "Recommend show (for watchers)")
  help_text += fmt.Sprintf("%s %s - %s\n",PossibleShows, "%watcher_1 %watcher_2 ...", "Get possible show to watch with watchers")
  help_text += fmt.Sprintf("%s %s - %s\n",MarkDone, "%group_id", "Mark group as done")
  help_text += fmt.Sprintf("%s %s - %s\n",MarkNotDone, "%group_id", "Mark group as not done")
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

  new_id, err := database.AddWatchGroup(show.Id, db)
  if err!= nil {
    b.SendMessage(ctx, &bot.SendMessageParams{ChatID: update.Message.Chat.ID, Text: fmt.Sprintf("Could not create watchgroup: %s",err)})
    return
  }

  b.SendMessage(ctx,&bot.SendMessageParams{ChatID:update.Message.Chat.ID, Text: fmt.Sprintf("Successfully created watchgroup (ID %d) for %s",new_id,show_name)})

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
  input_sep := strings.Split(input,SEP)
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
  input_sep := strings.Split(input,SEP)
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
  input_sep := strings.Split(input,SEP)
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
    watchers := strings.Split(input,SEP)
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

func handlePossible(ctx context.Context, b *bot.Bot, update *models.Update){
  input := strings.TrimSpace(strings.Replace(update.Message.Text,string(PossibleShows),"",-1))
  if input == "" {
      b.SendMessage(ctx, &bot.SendMessageParams{ChatID: update.Message.Chat.ID, Text: "Could not parse input, please try again"})
      return
  }
  watchers := strings.Split(input,SEP)

	db, ok := ctx.Value("database").(*sql.DB)
	if !ok {
		b.SendMessage(ctx, &bot.SendMessageParams{ChatID: update.Message.Chat.ID, Text: "Database connection failed, bot probably needs to be restarted"})
		return
	}

  shows, err := database.GetPossibleShows(watchers,db)
	if err != nil {
    b.SendMessage(ctx, &bot.SendMessageParams{ChatID: update.Message.Chat.ID, Text: fmt.Sprintf("Could not get possible shows: %s",err)})
		return
	}
  if len(shows) == 0{
    b.SendMessage(ctx, &bot.SendMessageParams{ChatID:update.Message.Chat.ID, Text:fmt.Sprintf("%s do not have a show to watch",input)})
    return
  }
  rand_show := types.RandomShow(shows)
  b.SendMessage(ctx, &bot.SendMessageParams{ChatID:update.Message.Chat.ID, Text:fmt.Sprintf("%s can watch %s",input,rand_show.Name)})

}

func handleMarkDone(ctx context.Context, b *bot.Bot, update *models.Update){
  group_id,err := strconv.Atoi(strings.TrimSpace(strings.Replace(update.Message.Text,string(MarkDone),"",-1)))
  if err!= nil{
      b.SendMessage(ctx, &bot.SendMessageParams{ChatID: update.Message.Chat.ID, Text: "Please provide group id"})
      return
  }

	db, ok := ctx.Value("database").(*sql.DB)
	if !ok {
		b.SendMessage(ctx, &bot.SendMessageParams{ChatID: update.Message.Chat.ID, Text: "Database connection failed, bot probably needs to be restarted"})
		return
	}

  err = database.MarkDone(group_id,db)
  if err != nil { 
    b.SendMessage(ctx, &bot.SendMessageParams{ChatID: update.Message.Chat.ID, Text: fmt.Sprintf("Could not mark as done: %s",err)})
		return
  }

  b.SendMessage(ctx, &bot.SendMessageParams{ChatID: update.Message.Chat.ID, Text: fmt.Sprintf("Successfully marked group %d as done",group_id)})
}

func handleMarkNotDone(ctx context.Context, b *bot.Bot, update *models.Update){
  group_id,err := strconv.Atoi(strings.TrimSpace(strings.Replace(update.Message.Text,string(MarkNotDone),"",-1)))
  if err!= nil{
      b.SendMessage(ctx, &bot.SendMessageParams{ChatID: update.Message.Chat.ID, Text: "Please provide group id"})
      return
  }

	db, ok := ctx.Value("database").(*sql.DB)
	if !ok {
		b.SendMessage(ctx, &bot.SendMessageParams{ChatID: update.Message.Chat.ID, Text: "Database connection failed, bot probably needs to be restarted"})
		return
	}

  err = database.MarkNotDone(group_id,db)
  if err != nil { 
    b.SendMessage(ctx, &bot.SendMessageParams{ChatID: update.Message.Chat.ID, Text: fmt.Sprintf("Could not mark as not done: %s",err)})
		return
  }

  b.SendMessage(ctx, &bot.SendMessageParams{ChatID: update.Message.Chat.ID, Text: fmt.Sprintf("Successfully marked group %d as not done",group_id)})

}

func handlePush(ctx context.Context, b *bot.Bot, update *models.Update) {
  err := exec.Command("git","add", "-A").Run()
  if err != nil {
    b.SendMessage(ctx, &bot.SendMessageParams{ChatID: update.Message.Chat.ID, Text: fmt.Sprintf("Could not add changes: %s",err)})
    return
  }
  err = exec.Command("git","commit","-m", "autocommit").Run()
  if err != nil {
    b.SendMessage(ctx, &bot.SendMessageParams{ChatID: update.Message.Chat.ID, Text: fmt.Sprintf("Could not commit changes: %s",err)})
    return
  }
  err = exec.Command("git","push").Run()
  if err != nil {
    b.SendMessage(ctx, &bot.SendMessageParams{ChatID: update.Message.Chat.ID, Text: fmt.Sprintf("Could not push changes: %s",err)})
    return
  }
  b.SendMessage(ctx, &bot.SendMessageParams{ChatID: update.Message.Chat.ID, Text: "Successully pushed changes"})
}

func handleWhosBack(ctx context.Context, b *bot.Bot, update *models.Update){
  b.SendMessage(ctx, &bot.SendMessageParams{ChatID: update.Message.Chat.ID, Text: "back again"})
  b.SendMessage(ctx, &bot.SendMessageParams{ChatID: update.Message.Chat.ID, Text: "watchbot's back, tell a friend"})
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
  b.RegisterHandler(bot.HandlerTypeMessageText, string(PossibleShows), bot.MatchTypePrefix, handlePossible)
  b.RegisterHandler(bot.HandlerTypeMessageText, string(MarkDone), bot.MatchTypePrefix, handleMarkDone)
  b.RegisterHandler(bot.HandlerTypeMessageText, string(MarkNotDone), bot.MatchTypePrefix, handleMarkNotDone)
  b.RegisterHandler(bot.HandlerTypeMessageText, "/guess_whos_back",bot.MatchTypeExact, handleWhosBack)
  b.RegisterHandler(bot.HandlerTypeMessageText, string(PushChanges), bot.MatchTypeExact, handlePush)
}
