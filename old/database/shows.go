package database

import (
  "rooxo/whoiswatching/types"
  "strings"
  "strconv"
  "database/sql"
)

func GetAllShows(db *sql.DB) ([]types.Show,error){
  query := "SELECT rowid,name FROM shows"
  res, err := db.Query(query)
  if err!=nil { return []types.Show{},err }
  defer res.Close()

  shows := make([]types.Show,0)
  for res.Next() {
    var next_id int 
    var next_name string 
    err = res.Scan(&next_id,&next_name)
    if err != nil { return []types.Show{},err }
    new_show := types.Show{Id:next_id,Name:next_name}
    shows = append(shows,new_show)
  }

  return shows,nil

}

func GetShowById(show_id int, db *sql.DB) (*types.Show, error) {
    res, err := db.Query("SELECT name FROM shows WHERE rowid=?",show_id)
    if err!=nil { return nil,err} 
    defer res.Close()

    res.Next()
    var name string
    err = res.Scan(&name)
    if err != nil { return nil,err }

    show := types.Show{Id:show_id,Name:name}
    return &show,nil
}

func GetShowByName(name string, db *sql.DB) (*types.Show,error){
  exists, err := ShowNameExists(name,db)
  if err != nil { return nil,err}
  if !exists { return nil,&ShowNameDoesNotExistErr{show_name:name} }

  res, err := db.Query("SELECT rowid FROM shows WHERE name=?",name)
  if err != nil { return nil,err}
  defer res.Close()

  res.Next()
  var id int 
  err = res.Scan(&id) 
  if err!=nil { return nil,err}

  show := types.Show{Id:id,Name:name}
  return &show,nil  
}

func ShowIdExists(show_id int, db *sql.DB) (bool,error){
  res, err := db.Query("SELECT count(*) FROM shows where rowid=?",show_id)
  if err != nil { return false,err  }
  defer res.Close()

  res.Next()
  var num int
  err = res.Scan(&num)
  if err!=nil{
    return false,err
  }

  return num>0,nil
}

func ShowNameExists(show_name string, db *sql.DB) (bool,error){
  res, err := db.Query("SELECT count(*) FROM shows where name=?",show_name)
  if err != nil { return false,err }
  defer res.Close()

  res.Next()
  var num int
  err = res.Scan(&num)
  if err !=nil { return false,err} 
  return num>0,nil
}

func AddShow(show_name string, db *sql.DB) error {
  exists, err := ShowNameExists(show_name,db)
  if err != nil { return err }
  if exists { return &ShowExistsErr{show_name:show_name} }

  _,err = db.Exec("INSERT INTO shows (name) VALUES (?);",show_name)
  if err != nil { return err }

  return nil
}

func RemoveShow(show_name string, db *sql.DB) error {
  show, err := GetShowByName(show_name,db)
  if err != nil { return err }

  _,err = db.Exec("DELETE FROM shows WHERE name = ?",show_name)
  if err != nil { return err }

  _,err = db.Exec("DELETE FROM watchgroups WHERE show_id=?",show.Id)
  if err!=nil { return err }

  return nil
}

func GetUnwatchedShows(watcher_names []string, db *sql.DB) ([]types.Show,error){
  watcher_ids := make([]string,0)
  for _, watcher_name := range watcher_names {
    watcher, err := GetWatcherByName(watcher_name,db)
    if err != nil { return []types.Show{},err} 
    watcher_ids = append(watcher_ids,strconv.Itoa(watcher.Id))
  }

  res,err := db.Query("SELECT s.rowid FROM shows AS s JOIN watchgroups AS g on s.rowid=g.show_id JOIN watchers_groups AS wg on g.rowid=wg.group_id WHERE wg.watcher_id IN (?)",
    strings.Join(watcher_ids,",")) 
  if err != nil { return []types.Show{}, err}
  defer res.Close()

  exclude_ids := make([]string,0)
  for res.Next(){
    var show_id int
    err = res.Scan(&show_id)
    if err != nil { return []types.Show{}, err} 
    exclude_ids = append(exclude_ids,strconv.Itoa(show_id))
  }

  res, err = db.Query("SELECT rowid,name FROM shows WHERE rowid NOT IN (?)",strings.Join(exclude_ids,","))
  if err != nil { return []types.Show{}, err }
  defer res.Close()

  shows := make([]types.Show,0)
  for res.Next(){
    var show_id int 
    var show_name string 
    err = res.Scan(&show_id,&show_name)
    if err != nil { return []types.Show{},err }
    new_show := types.Show{Id:show_id,Name:show_name}
    shows = append(shows,new_show)
  }

  return shows,nil
}
