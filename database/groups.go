package database 

import ( 
  "rooxo/whoiswatching/types"
  "fmt"
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
