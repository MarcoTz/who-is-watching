package database 

import ( 
  "rooxo/whoiswatching/types"
  "fmt"
  "database/sql"
)

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

func AddWatchGroup(show_id int, users []int, db *sql.DB) (int, error){
  exists, err := ShowIdExists(show_id,db)
  if !exists{
    return 0,&ShowIdDoesNotExist{show_id:show_id}
  }
  if err!=nil{
    return 0,err
  }

  query := fmt.Sprintf("INSERT INTO watchgroups show_id VALUES %d",show_id)
  _,err = db.Exec(query)
  if err !=nil{
    return 0,err
  }

  var group_id int
  id_query := "SELECT MAX(row_id) from watchgroups";
  res,err := db.Query(id_query)
  if err!=nil {return 0,err}
  defer res.Close() 

  res.Next()
  err = res.Scan(&group_id)
  if err!=nil { return 0,err}
  

  for _,watcher_id := range users{
    exists,err = WatcherIdExists(watcher_id,db)
    if err !=nil { return 0,err} 
    if !exists { return 0,&WatcherIdDoesNotExistErr{watcher_id:watcher_id} }
    insert_st := fmt.Sprintf("INSERT INTO watchers_groups (watcher_id,group_id) values (%d,%d)",group_id,watcher_id)
    _,err = db.Exec(insert_st)
    if err != nil { return 0,err}
  }

  return group_id,nil

}
