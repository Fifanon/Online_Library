package member

import(
	vars "online_library/varsAndFuncs"
	stct "online_library/structs"
	"net/http"
	dbconfig "online_library/config"
	s "online_library/session"
)

//MembersOwingFine **
func MembersOwingFine(w http.ResponseWriter, r *http.Request) {
	if	validated := s.GetSession(r);!validated{
		http.Redirect(w, r, "/home", http.StatusSeeOther)
		return
	}
		var bookBorrow stct.BorrowInfo
		booksBorrow := []stct.BorrowInfo{}
		var mbemail string
		mbemails := []string{}
		var fine int
		db, err := dbconfig.GetMySQLDb()
		if err != nil {
			panic(err)
		}
		qr, err := db.Query(`select distinct member_email from books_borrowed where fine > 0;`)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		for qr.Next() {
			err = qr.Scan(&mbemail)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			mbemails = append(mbemails, mbemail)
		}
		bookBorrow.Fine = 0
		for _, mbemail := range mbemails {
			qr2, err := db.Query(`select fine from books_borrowed where member_email = ? and fine > 0;`, mbemail)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			for qr2.Next() {
				err = qr2.Scan(&fine)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
				}
				bookBorrow.Fine = bookBorrow.Fine + fine
			}
			qr3, err := db.Query(`select m_firstname,m_lastname from members where m_email = ?;`, mbemail)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			for qr3.Next() {
				err = qr3.Scan(&bookBorrow.FirstName, &bookBorrow.LastName)

				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
				}
			}
			booksBorrow = append(booksBorrow, bookBorrow)
		}
		vars.Tpl.ExecuteTemplate(w, "membersOwingFine.html", booksBorrow)
		return
}