package database 

import ( 
  "rooxo/whoiswatching/types"
  "fmt"
  "strconv"
  "strings"
  "database/sql"
)

func GroupIdExists(group_id int, db *sql.DB) (bool,error){
  query := fmt.Sprintf("SELECT COUNT(*) FROM watchgroups WHERE rowid=%d",group_id)
  res,err := db.Query(query)
  if err != nil { return false,err }
  defer res.Close()

  res.Next()
  var num int
  err = res.Scan(&num)
  if err != nil { return false,err}
  
  return num>0,nil
}

func GetAllGroups(db *sql.DB) ([]types.WatchGroup,error){
  query := "SELECT rowid,show_id,current_ep FROM watchgroups"
  res, err := db.Query(query)
  if err != nil { return []types.WatchGroup{},err }
  defer res.Close()


  groups := make([]types.WatchGroup,0)
  for res.Next() {
    var group_id int
    var show_id int
    var current_ep int 
    err = res.Scan(&group_id,&show_id,&current_ep)
    if err != nil { return []types.WatchGroup{},err }
    watchers_query := fmt.Sprintf("SELECT watcher_id from watchers_groups WHERE group_id=%d",group_id);
    watchers_res,err := db.Query(watchers_query)
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

    group := types.WatchGroup{Id:group_id,Watchers:watchers,Current_ep:current_ep,Show:*show}
    groups = append(groups,group)

  }

  return groups,nil
}

func AddWatchGroup(show_id int, db *sql.DB) error{
  exists, err := ShowIdExists(show_id,db)
  if !exists{
    return &ShowIdDoesNotExist{show_id:show_id}
  }
  if err!=nil{
    return err
  }

  query := fmt.Sprintf("INSERT INTO watchgroups (show_id,current_ep) VALUES (%d,1)",show_id)
  _,err = db.Exec(query)
  if err !=nil{
    return err
  }
  return nil

}

func RemoveGroup(group_id int, db *sql.DB) error {
  exists,err := GroupIdExists(group_id,db)
  if err != nil { return err}
  if !exists { return &GroupIdDoesNotExistErr {group_id:group_id} }

  del_stmt := fmt.Sprintf("DELETE FROM watchgroups WHERE rowid=%d",group_id)
  _,err = db.Exec(del_stmt)
  if err != nil { return err }

  links_stmt := fmt.Sprintf("DELETE FROM watchers_groups WHERE group_id=%d",group_id)
  _,err = db.Exec(links_stmt)
  if err != nil { return err}

  return nil
}

func GetGroupsByShowName(show_name string, db *sql.DB) ([]types.WatchGroup,error) {
  show, err := GetShowByName(show_name,db)
  if err!= nil { return []types.WatchGroup{},err }

  query := fmt.Sprintf("SELECT rowid,current_ep from watchgroups where show_id=%d",show.Id)
  res,err := db.Query(query)
  if err!=nil {return []types.WatchGroup{},err}
  defer res.Close()

  groups := make([]types.WatchGroup,0)
  for res.Next(){
    var group_id int
    var current_ep int
    err = res.Scan(&group_id,&current_ep)
    if err != nil {return [] types.WatchGroup{},err}

    watcher_query := fmt.Sprintf("SELECT watcher_id FROM watchers_groups WHERE group_id=%d",group_id)
    watcher_res, err := db.Query(watcher_query)
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

    new_group := types.WatchGroup { Id: group_id, Show: *show, Current_ep:current_ep, Watchers:group_watchers} 
    groups = append(groups,new_group)
  }

  return groups,nil
}

func UpdateGroupEpisode(group_id int, new_ep int, db *sql.DB) error {
  exists, err := GroupIdExists(group_id,db)
  if err!=nil { return err}
  if !exists { return &GroupIdDoesNotExistErr{group_id:group_id} }

  stmt := fmt.Sprintf("UPDATE watchgroups SET current_ep=%d WHERE rowid=%d",new_ep,group_id)
  _,err = db.Exec(stmt)
  if err != nil { return err }

  return nil
}

func AddWatcherGroup(group_id int, watcher_name string, db *sql.DB) error {
  watcher,err := GetWatcherByName(watcher_name,db)
  if err != nil { return err } 

  query := fmt.Sprintf("INSERT INTO watchers_groups (group_id,watcher_id) VALUES (%d,%d)",group_id,watcher.Id)
  _,err = db.Exec(query)
  if err != nil { return err } 

  return nil
}

func RemoveWatcherGroup(group_id int, watcher_name string, db *sql.DB) error{
  watcher,err := GetWatcherByName(watcher_name,db)
  if err != nil { return err}

  is_watching_query := fmt.Sprintf("SELECT COUNT(*) FROM watchers_groups WHERE group_id=%d AND watcher_id=%d",group_id,watcher.Id)
  res,err := db.Query(is_watching_query)
  if err != nil { return err }

  res.Next()
  var count int 
  err = res.Scan(&count)
  if err != nil { return err }
  if count == 0 { return &NotAWatcherErr{watcher_name:watcher_name,group_id:group_id} }
  res.Close()

  
  query := fmt.Sprintf("DELETE FROM watchers_groups WHERE group_id=%d AND watcher_id=%d",group_id,watcher.Id)
  _,err = db.Exec(query)
  if err != nil { return err }

  return nil
}

func GetPossibleShows(watcher_names []string, db *sql.DB) ([]types.Show, error) {
  watcher_ids := make([]string,0)
  for _,watcher_name := range(watcher_names){
    watcher,err := GetWatcherByName(watcher_name,db)
    if err != nil { return []types.Show{},err }
    watcher_ids = append(watcher_ids,strconv.Itoa(watcher.Id))
  }

  query := fmt.Sprintf("SELECT g.show_id FROM watchgroups AS g JOIN watchers_groups AS wg ON g.rowid=wg.group_id WHERE watcher_id IN (%s)",strings.Join(watcher_ids,","))
  res, err := db.Query(query)
  if err != nil { return []types.Show{},err }
  defer res.Close()

  shows := make([]types.Show,0)
  for res.Next(){
    var show_id int 
    err = res.Scan(&show_id)
    if err != nil { return []types.Show{},err} 

    show,err := GetShowById(show_id,db)
    if err != nil { return []types.Show{},err}
    shows = append(shows,*show)
  }

  return shows,nil

}
