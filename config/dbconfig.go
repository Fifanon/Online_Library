package dbconfig


import (
   "database/sql"
)

//GetMySQLDb **
func GetMySQLDb()(db *sql.DB, err error){
	HOST := "freedb.tech"
	dbDriver := "mysql"
	dbUser := "freedbtech_lesley099"
	dbPass := "qwe123"
	dbName := "freedbtech_science_library"
	Port := "3306"
	db, err = sql.Open(dbDriver, dbUser+":"+dbPass+"@tcp("+HOST+":"+Port+")/"+dbName)
	return
}
