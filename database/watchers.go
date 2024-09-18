package database

import (
  "database/sql"
  "fmt"
  "rooxo/whoiswatching/types"
)

func WatcherIdExists(watcher_id int, db *sql.DB) (bool, error) {

  query := fmt.Sprintf("SELECT count(*) FROM watchers where rowid=%d", watcher_id)
  res, err := db.Query(query)
  if err != nil {
    return false, err
  }
  defer res.Close()


  res.Next()
  var num int
  err = res.Scan(&num)
  if err != nil {
    return false, err
  }
  
  return num > 0, nil
}

func WatcherNameExists(watcher_name string, db *sql.DB) (bool, error) {
  query := fmt.Sprintf("SELECT count(*) FROM watchers where name='%s'", watcher_name)
  res, err := db.Query(query)
  if err != nil {
    return false, err
  }
  defer res.Close()

  res.Next()
  var num int
  err = res.Scan(&num)
  if err != nil {
    return false, err
  }

  return num > 0, nil
}

func GetWatcherById(watcher_id int, db *sql.DB) (*types.Watcher, error) {
  exists, err := WatcherIdExists(watcher_id,db)
  if err != nil { return nil,err}
  if !exists {return nil,&WatcherIdDoesNotExistErr{watcher_id:watcher_id}}

  query := fmt.Sprintf("SELECT name FROM watchers where rowid=%d", watcher_id)
  res, err := db.Query(query)
  defer res.Close()

  if err != nil {
    return nil, err
  }

  res.Next()
  var name string
  err = res.Scan(&name)
  if err != nil {
    return nil, err
  }

  watcher := types.Watcher{Id: watcher_id, Name: name}
  return &watcher, nil

}

func GetWatcherByName(watcher_name string, db *sql.DB) (*types.Watcher,error) {
  exists,err := WatcherNameExists(watcher_name,db)
  if err != nil { return nil,err}
  if !exists { return nil, &WatcherNameDoesNotExistErr{watcher_name:watcher_name} }

  query := fmt.Sprintf("SELECT rowid FROM watchers where name='%s'",watcher_name)
  res, err := db.Query(query)
  if err!=nil { return nil,err}
  defer res.Close()

  res.Next()
  var watcher_id int
  err = res.Scan(&watcher_id)
  if err!=nil {return nil,err}

  watcher := types.Watcher{Id:watcher_id,Name:watcher_name}
  return &watcher,nil
}

func GetAllWatchers(db *sql.DB) ([]types.Watcher, error) {
  query := "SELECT rowid,name FROM watchers"
  res, err := db.Query(query)
  if err != nil {
    return []types.Watcher{}, err
  }
  defer res.Close()

  watchers := make([]types.Watcher, 0)
  for res.Next() {
    var watcher_id int
    var watcher_name string
    err = res.Scan(&watcher_id, &watcher_name)
    if err != nil {
      return []types.Watcher{}, err
    }

    new_watcher := types.Watcher{Id: watcher_id, Name: watcher_name}
    watchers = append(watchers, new_watcher)
  }

  return watchers, nil
}

func AddWatcher(watcher_name string, db *sql.DB) error {
  exists,err := WatcherNameExists(watcher_name,db)
  if err != nil { return err } 
  if exists { return &WatcherExistsErr{watcher_name} }

  query := fmt.Sprintf("INSERT INTO watchers (name) VALUES ('%s');",watcher_name)
  _,err = db.Exec(query)
  if err != nil { return err }

  return nil
}

func RemoveWatcher(watcher_name string, db *sql.DB) error {
  watcher, err := GetWatcherByName(watcher_name,db)
  if err != nil { return err }

  query_del := fmt.Sprintf("DELETE FROM watchers WHERE name = '%s'",watcher_name)
  _,err = db.Exec(query_del)
  if err != nil { return err }

  query_groups := fmt.Sprintf("DELETE FROM watchers_groups WHERE watcher_id=%d",watcher.Id)
  _,err = db.Exec(query_groups)
  if err!=nil { return err }

  return nil
}
