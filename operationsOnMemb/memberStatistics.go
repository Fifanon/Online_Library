package member

import (
	"net/http"

	dbconfig "github.com/Fifanon/online_library/config"
	stct "github.com/Fifanon/online_library/structs"
	vars "github.com/Fifanon/online_library/varsAndFuncs"
	s "github.com/Fifanon/online_library/session"

)

//ProcessMemberStatistics **
func ProcessMemberStatistics(w http.ResponseWriter, r *http.Request) {
	if	validated := s.GetSession(r);!validated{
		http.Redirect(w, r, "/home", http.StatusSeeOther)
		return
	}
		var bookBorrow stct.BorrowInfo
		booksBorrow := []stct.BorrowInfo{}
		mbemail := stct.User.Email
		var fine int
		db, err := dbconfig.GetMySQLDb()
		if err != nil {
			panic(err)
		}
		bookBorrow.Fine = 0
		bookBorrow.Number = 0
		qr2, err := db.Query(`select fine from books_borrowed where member_email = ?;`, mbemail)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		for qr2.Next() {
			err = qr2.Scan(&fine)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			bookBorrow.Number = bookBorrow.Number + 1
			bookBorrow.Fine = bookBorrow.Fine + fine
			booksBorrow = append(booksBorrow, bookBorrow)
		}
		vars.Tpl.ExecuteTemplate(w, "memberStatistics.html", booksBorrow)
}
