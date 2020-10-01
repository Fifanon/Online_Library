package dbconfig


import (
   "database/sql"
   _ "github.com/lib/pq"
)

//GetMySQLDb **
func GetMySQLDb()(db *sql.DB, err error){
	HOST := "ec2-52-207-25-133.compute-1.amazonaws.com"
	dbDriver := "postgres"
	dbUser := "pxidleyobnqbqd"
	dbPass := "a456a7bacc7f3e0e314a841162445ae3ca624b37623bd59aa6a202e7712748cc"
	dbName := "d5r925fe5mtfbg"
	Port := "5432"
	db, err = sql.Open(dbDriver, dbUser+":"+dbPass+"@tcp("+HOST+":"+Port+")/"+dbName)
	return
}
