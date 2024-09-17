package database 

import (
  "rooxo/whoiswatching/types"
  "fmt" 
  "database/sql"
)

func WatcherExists(watcher_id int, db *sql.DB) (bool,error){

  query := fmt.Sprintf("SELECT count(*) from watchers where rowid=%d",watcher_id)
  res, err := db.Query(query)
  if err != nil { return false,err }

  res.Next()
  var num int 
  err = res.Scan(&num) 
  if err!=nil { return false,err }

  return num>0,nil
}

func GetWatcherById(watcher_id int, db *sql.DB) (*types.Watcher,error){
  query := fmt.Sprintf("SELECT name from watchers where rowid=%d",watcher_id)
  res, err := db.Query(query)
  if err !=nil {return nil,err}

  res.Next()
  var name string 
  err = res.Scan(&name)
  if err != nil { return nil,err } 
  
  watcher := types.Watcher{ Id:watcher_id, Name:name}
  return &watcher, nil

}
