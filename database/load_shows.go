package database

import (
  "fmt"
  "database/sql"
  "os"
  _ "github.com/mattn/go-sqlite3"
)

func Load() {
  entries, err := os.ReadDir(".")
  if err != nil {
    fmt.Printf("Could not read dir %s",err)
    return 
  }
  fmt.Printf("Dir entries %s\n",entries)

  db, err := sql.Open("sqlite3","./watchers.db")
  if err != nil {
    fmt.Println("Failed to create database")
    return 
  }

  res, err:= db.Query("select id,name from shows");
  if err != nil {
    fmt.Printf("Failed to execute statement %s",err);
    return
  }

  fmt.Println("Id | Name")
  for res.Next(){
    var id int
    var name string
    err = res.Scan(&id,&name)
    if err != nil{
      fmt.Printf("Failed to get row: %s",err);
      return
    }
    fmt.Printf("%d | %s",id,name)
    fmt.Println("")
  }
}
