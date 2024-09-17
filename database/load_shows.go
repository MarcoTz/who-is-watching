package database

import (
  "fmt"
  "database/sql"
  _ "github.com/mattn/go-sqlite3"
)

func connectDB(db_file string) (*sql.DB,error) {
  return sql.Open("sqlite3",db_file)
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

func userExists(user_id int, db *sql.DB) (bool,error){
  query := fmt.Sprintf("SELECT count(*) from users where rowid=%d",user_id)
  res, err := db.Query(query)
  if err != nil { return false,err }

  res.Next()
  var num int 
  err = res.Scan(&num) 
  if err!=nil { return false,err }

  return num>0,nil
}

func AddWatchGroup(show_id int, users []int, db *sql.DB) (int, error){
  exists, err := showExists(show_id,db)
  if !exists{
    return 0,&ShowIdNotFoundErr{show_id:show_id}
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
  res.Next()
  err = res.Scan(&group_id)
  if err!=nil { return 0,err}
  

  for _,user_id := range users{
    exists,err = userExists(user_id,db)
    if err !=nil { return 0,err} 
    if !exists { return 0,&UserIdNotFoundErr{user_id:user_id} }
    insert_st := fmt.Sprintf("INSERT INTO watchers_groups (watcher_id,group_id) values (%d,%d)",group_id,user_id)
    _,err = db.Exec(insert_st)
    if err != nil { return 0,err}
  }

  return group_id,nil

}
