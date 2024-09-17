package database 

import (
  "fmt" 
  "database/sql"
)

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
