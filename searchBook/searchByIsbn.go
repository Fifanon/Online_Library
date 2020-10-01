package searchbook

import (
	"database/sql"
	dbconfig "github.com/Fifanon/online_library/config"
	stct "github.com/Fifanon/online_library/structs"
)

//SearchByIsbn **
func SearchByIsbn(isbn int) (bks []stct.BookStruct, found bool, errEnc bool) {


	db, err := dbconfig.GetMySQLDb()
	if err != nil {
		panic(err)	
	}
	qr, err := db.Query(`select book_isbn,book_title,author_name,pages,subject_area,number,b_imagename from book_instances
                       where book_isbn = $1;`, isbn)
	if err != nil {
		panic(err)	
	}
	for qr.Next() {
		err = qr.Scan(&stct.Bk.ISBN, &stct.Bk.Title, &stct.Bk.Author, &stct.Bk.Pages, &stct.Bk.Subject, &stct.Bk.Number, &stct.Bk.BookImageName)
		if err != nil {
			if err == sql.ErrNoRows {
				found = false
			} else {
				errEnc = true
			}
		} else {
			found = true
		}
		if stct.Bk.Number != 0 {
			stct.Bk.Availability = "AVAILABLE"
		} else {
			stct.Bk.Availability = "NOT AVAILABLE"
		}
		bks = append(bks, stct.Bk)
	}
	return bks, found, errEnc
}
