package database

import (
  "database/sql"
  "fmt"
  "rooxo/whoiswatching/types"
)

func WatcherIdExists(watcher_id int, db *sql.DB) (bool, error) {

  query := fmt.Sprintf("SELECT count(*) from watchers where rowid=%d", watcher_id)
  res, err := db.Query(query)
  if err != nil {
    return false, err
  }

  res.Next()
  var num int
  err = res.Scan(&num)
  if err != nil {
    return false, err
  }
  res.Close()

  return num > 0, nil
}

func WatcherNameExists(watcher_name string, db *sql.DB) (bool, error) {
  query := fmt.Sprintf("SELECT count(*) from watchers where name='%s'", watcher_name)
  res, err := db.Query(query)
  if err != nil {
    return false, err
  }

  res.Next()
  var num int
  err = res.Scan(&num)
  if err != nil {
    return false, err
  }
  res.Close()

  return num > 0, nil
}

func GetWatcherById(watcher_id int, db *sql.DB) (*types.Watcher, error) {
  query := fmt.Sprintf("SELECT name from watchers where rowid=%d", watcher_id)
  res, err := db.Query(query)
  if err != nil {
    return nil, err
  }

  res.Next()
  var name string
  err = res.Scan(&name)
  if err != nil {
    return nil, err
  }
  res.Close()

  watcher := types.Watcher{Id: watcher_id, Name: name}
  return &watcher, nil

}

func GetAllWatchers(db *sql.DB) ([]types.Watcher, error) {
  query := "SELECT rowid,name from watchers"
  res, err := db.Query(query)
  if err != nil {
    return []types.Watcher{}, err
  }

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
  res.Close()

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
