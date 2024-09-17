package database 

import ( 
  "database/sql"
  _ "github.com/mattn/go-sqlite3"
)

func connectDB(db_file string) (*sql.DB,error) {
  return sql.Open("sqlite3",db_file)
}

