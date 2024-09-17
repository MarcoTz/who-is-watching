package database

import (
  "rooxo/whoiswatching/types"
  "fmt"
  "database/sql"
)

func GetAllShows(db *sql.DB) ([]types.Show,error){
  query := "SELECT rowid,name FROM shows"
  res, err := db.Query(query)
  if err!=nil { return []types.Show{},err }

  shows := make([]types.Show,0)
  for res.Next() {
    var next_id int 
    var next_name string 
    err = res.Scan(&next_id,&next_name)
    if err != nil { return []types.Show{},err }
    new_show := types.NewShow(next_id,next_name)
    shows = append(shows,new_show)
  }

  return shows,nil

}

func GetShowById(show_id int, db *sql.DB) (*types.Show, error) {
    query := fmt.Sprintf("SELECT name FROM shows WHERE rowid=%d",show_id);
    res, err := db.Query(query)
    if err!=nil { return nil,err} 

    res.Next()
    var name string
    err = res.Scan(&name)
    if err != nil { return nil,err }

    show := types.NewShow(show_id,name)
    return &show,nil
}

func GetShowByName(name string, db *sql.DB) (*types.Show,error){
  query := fmt.Sprintf("SELECT rowid FROM shows WHERE name=%s",name)
  res, err := db.Query(query)
  if err != nil { return nil,err}

  res.Next()
  var id int 
  err = res.Scan(&id) 
  if err!=nil { return nil,err}

  show := types.NewShow(id,name)
  return &show,nil

  
}

func showExists(show_id int, db *sql.DB) (bool,error){
  query := fmt.Sprintf("SELECT count(*) from shows where rowid=%d",show_id)
  res, err := db.Query(query)
  if err != nil {
    return false,err
  }

  res.Next()
  var num int
  err = res.Scan(&num)
  if err!=nil{
    return false,err
  }

  return num>0,nil
}
