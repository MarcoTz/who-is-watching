package database 

import ( 
  "database/sql"
  _ "github.com/mattn/go-sqlite3"
)

func ConnectDB(db_file string) (*sql.DB,error) {
  return sql.Open("sqlite3",db_file)
}

