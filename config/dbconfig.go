package dbconfig


import (
   "database/sql"
   _"github.com/lib/pq"
   "os"
)

//GetMySQLDb **
func GetMySQLDb()(db *sql.DB, err error){
	db, err = sql.Open("postgres", os.Getenv("DATABASE_URL"))
	return db,err
}
