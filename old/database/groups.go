package database 

import ( 
  "slices"
  "strings"
  "rooxo/whoiswatching/types"
  "database/sql"
)

func GroupIdExists(group_id int, db *sql.DB) (bool,error){
  res,err := db.Query("SELECT COUNT(*) FROM watchgroups WHERE rowid=?",group_id)
  if err != nil { return false,err }
  defer res.Close()

  res.Next()
  var num int
  err = res.Scan(&num)
  if err != nil { return false,err}
  
  return num>0,nil
}

func GetGroupById(group_id int, db *sql.DB) (*types.WatchGroup, error){
  res,err := db.Query("SELECT show_id,current_ep,done FROM watchgroups WHERE rowid=?",group_id)
  if err!=nil { return nil,err }

  var show_id int
  var current_ep int
  var done bool
  res.Next()
  err = res.Scan(&show_id,&current_ep,&done)
  if err != nil { return nil,err} 
  res.Close()
  
  show,err := GetShowById(show_id,db)
  if err != nil { return nil,err }

  res,err = db.Query("SELECT watcher_id FROM watchers_groups WHERE group_id=?",group_id)
  if err != nil { return nil,err }
  defer res.Close()

  watchers := make([]types.Watcher,0)
  for res.Next() {
    var watcher_id int 
    err = res.Scan(&watcher_id)
    if err != nil { return nil,err }

    watcher,err := GetWatcherById(watcher_id,db)
    if err != nil { return nil,err} 

    watchers = append(watchers,*watcher)

  }

  group := types.WatchGroup{Id:group_id,Show:*show,Current_ep:current_ep, Done:done,Watchers:watchers}
  return &group,nil
}

func GetAllGroups(db *sql.DB) ([]types.WatchGroup,error){
  query := "SELECT rowid,show_id,current_ep,done FROM watchgroups"
  res, err := db.Query(query)
  if err != nil { return []types.WatchGroup{},err }
  defer res.Close()


  groups := make([]types.WatchGroup,0)
  for res.Next() {
    var group_id int 
    var show_id int
    var current_ep int 
    var done bool
    err = res.Scan(&group_id,&show_id,&current_ep,&done)
    if err != nil { return []types.WatchGroup{},err }
    watchers_res,err := db.Query("SELECT watcher_id from watchers_groups WHERE group_id=?",group_id)
    if err!=nil { return []types.WatchGroup{},err }

    show,err := GetShowById(show_id,db)
    if err!=nil{ return []types.WatchGroup{},err}

    watchers := make([]types.Watcher,0)
    for watchers_res.Next(){
      var watcher_id int 
      err = watchers_res.Scan(&watcher_id)
      if err!=nil {return []types.WatchGroup{},err}

      watcher,err := GetWatcherById(watcher_id,db)
      if err != nil { return []types.WatchGroup{},err}

      watchers = append(watchers,*watcher)
    }

    group := types.WatchGroup{Id:group_id,Watchers:watchers,Current_ep:current_ep,Show:*show,Done:done}
    groups = append(groups,group)

  }

  return groups,nil
}

func AddWatchGroup(show_id int, db *sql.DB) (int,error){
  exists, err := ShowIdExists(show_id,db)
  if !exists{
    return 0,&ShowIdDoesNotExistErr{show_id:show_id}
  }
  if err!=nil{
    return 0,err
  }

  _,err = db.Exec("INSERT INTO watchgroups (show_id,current_ep,done) VALUES (?,1,false)",show_id)
  if err !=nil{
    return 0,err
  }

  res, err := db.Query("SELECT MAX(rowid) FROM watchgroups")
  if err != nil { return 0,err} 
  defer res.Close()
  res.Next()
  var group_id int 
  err = res.Scan(&group_id)
  if err != nil { return 0,err} 

  return group_id,nil

}

func RemoveGroup(group_id int, db *sql.DB) error {
  exists,err := GroupIdExists(group_id,db)
  if err != nil { return err}
  if !exists { return &GroupIdDoesNotExistErr {group_id:group_id} }

  _,err = db.Exec("DELETE FROM watchgroups WHERE rowid=?",group_id)
  if err != nil { return err }

  _,err = db.Exec("DELETE FROM watchers_groups WHERE group_id=?",group_id)
  if err != nil { return err}

  return nil
}

func GetGroupsByShowName(show_name string, db *sql.DB) ([]types.WatchGroup,error) {
  show, err := GetShowByName(show_name,db)
  if err!= nil { return []types.WatchGroup{},err }

  res,err := db.Query("SELECT rowid,current_ep,done from watchgroups where show_id=?",show.Id)
  if err!=nil {return []types.WatchGroup{},err}
  defer res.Close()

  groups := make([]types.WatchGroup,0)
  for res.Next(){
    var group_id int
    var current_ep int
    var done bool
    err = res.Scan(&group_id,&current_ep,&done)
    if err != nil {return [] types.WatchGroup{},err}

    watcher_res, err := db.Query("SELECT watcher_id FROM watchers_groups WHERE group_id=?",group_id)
    if err != nil { return [] types.WatchGroup{},err} 
    defer watcher_res.Close()

    group_watchers := make([]types.Watcher,0)
    for watcher_res.Next() {
      var watcher_id int
      err = watcher_res.Scan(&watcher_id)
      if err != nil { return []types.WatchGroup{},err}

      watcher,err := GetWatcherById(watcher_id,db)
      if err != nil { return []types.WatchGroup{},err}
      group_watchers = append(group_watchers,*watcher)
      
    }

    new_group := types.WatchGroup { Id: group_id, Show: *show, Current_ep:current_ep, Watchers:group_watchers, Done:done} 
    groups = append(groups,new_group)
  }

  return groups,nil
}

func UpdateGroupEpisode(group_id int, new_ep int, db *sql.DB) error {
  exists, err := GroupIdExists(group_id,db)
  if err!=nil { return err}
  if !exists { return &GroupIdDoesNotExistErr{group_id:group_id} }

  _,err = db.Exec("UPDATE watchgroups SET current_ep=? WHERE rowid=?",new_ep,group_id)
  if err != nil { return err }

  return nil
}

func AddWatcherGroup(group_id int, watcher_name string, db *sql.DB) error {
  watcher,err := GetWatcherByName(watcher_name,db)
  if err != nil { return err } 

  _,err = db.Exec("INSERT INTO watchers_groups (group_id,watcher_id) VALUES (?,?)",group_id,watcher.Id)
  if err != nil { return err } 

  return nil
}

func RemoveWatcherGroup(group_id int, watcher_name string, db *sql.DB) error{
  watcher,err := GetWatcherByName(watcher_name,db)
  if err != nil { return err}

  res,err := db.Query("SELECT COUNT(*) FROM watchers_groups WHERE group_id=? AND watcher_id=?",group_id,watcher.Id)
  if err != nil { return err }

  res.Next()
  var count int 
  err = res.Scan(&count)
  if err != nil { return err }
  if count == 0 { return &NotAWatcherErr{watcher_name:watcher_name,group_id:group_id} }
  res.Close()

  
  _,err = db.Exec("DELETE FROM watchers_groups WHERE group_id=? AND watcher_id=?",group_id,watcher.Id)
  if err != nil { return err }

  return nil
}

func GetGroupsByWatcher(watcher_name string, db *sql.DB) ([]types.WatchGroup, error){
  res,err := db.Query(`
  SELECT wg.group_id 
  FROM watchers_groups AS wg 
  JOIN watchers AS w ON w.rowid=wg.watcher_id 
  JOIN watchgroups AS g ON wg.group_id=g.rowid
  WHERE w.name=? AND g.done=false
  `,watcher_name)
  if err!=nil { return []types.WatchGroup{},err}
  defer res.Close()

  groups := make([]types.WatchGroup,0)
  for res.Next(){
    var group_id int 
    err = res.Scan(&group_id)
    if err !=nil { return []types.WatchGroup{},err}
    group,err := GetGroupById(group_id,db)
    if err != nil { return [] types.WatchGroup{},err}
    groups = append(groups,*group)
  }

  return groups,nil
}

func GetPossibleShows(watcher_names []string, db *sql.DB) ([]types.Show, error){
  if len(watcher_names) == 0 { return []types.Show{}, &NoInput{} }
  groups,err := GetGroupsByWatcher(watcher_names[0],db)
  if err != nil { return []types.Show{}, err}
  
  slices.SortFunc(watcher_names,strings.Compare)
  groups_filtered := make([]types.WatchGroup,0)
  for _,group := range groups{
    group_watchers := group.GetWatchers()
    slices.SortFunc(group_watchers,strings.Compare)
    keep := true
    for i := range watcher_names{ 
      if watcher_names[i] != group_watchers[i]{
        keep = false
        break
      }
    }
    if keep { groups_filtered = append(groups_filtered,group) }
  }

  shows := make([]types.Show,0)
  for _,group := range groups_filtered{
    shows = append(shows,group.Show)
  }

  return shows,nil
}

func MarkDone(group_id int, db *sql.DB) error {
  exists,err := GroupIdExists(group_id,db)
  if err != nil { return err}
  if !exists { return &GroupIdDoesNotExistErr{group_id:group_id} }

  _,err = db.Exec("UPDATE watchgroups SET done=true WHERE rowid=?",group_id)
  if err != nil { return err }

  return nil
}

func MarkNotDone(group_id int, db *sql.DB) error {
  exists, err := GroupIdExists(group_id,db)
  if err != nil { return err }
  if !exists{ return &GroupIdDoesNotExistErr{group_id:group_id} }

  _, err = db.Exec("UPDATE watchgroups SET done=false WHERE rowid=?",group_id)
  if err != nil { return err }
  return nil
}
