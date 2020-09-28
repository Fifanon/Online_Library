package bkop

import(
	"time"
	"online_library/modules/github.com/gorilla/mux"
	vars "online_library/varsAndFuncs"
	stct "online_library/structs"
	"net/http"
	dbconfig "online_library/config"
	s "online_library/session"

)

//IssuedBook **
func IssuedBook(w http.ResponseWriter, r *http.Request) {

	if	validated := s.GetSession(r);!validated{
		http.Redirect(w, r, "/home", http.StatusSeeOther)
		return
	}
		db, err := dbconfig.GetMySQLDb()
		if err != nil {
			panic(err)
		}
		qr, err := db.Query(`select book_isbn,book_title,author_name,pages,subject_area,number,b_imagename,m_imagename,m_firstname,m_lastname, m_email,m_status,m_address from tmp_borrow,members,book_instances
                       where book_isbn = bk_isbn and m_email = mb_email;`)
		if err != nil {
			panic(err)
		}
		var bookndMem stct.BorrowInfo
		bookndMembs := []stct.BorrowInfo{}
		for qr.Next() {
			err = qr.Scan(&bookndMem.ISBN, &bookndMem.Title, &bookndMem.Author, &bookndMem.Pages, &bookndMem.Subject,
				&bookndMem.Number, &bookndMem.BookImageName, &bookndMem.ImageName, &bookndMem.FirstName, &bookndMem.LastName, &bookndMem.Email,&bookndMem.Status,&bookndMem.Address)
			if err != nil {
				panic(err)
			}
			bookndMembs = append(bookndMembs, bookndMem)
		}
		if len(bookndMembs) == 0 {
			vars.Message = "NO BOOK REQUESTED"
			vars.Tpl.ExecuteTemplate(w, "NoData.html", vars.Message)
			vars.Message = ""
			return
		}
		vars.Tpl.ExecuteTemplate(w, "bookIssuing.html", bookndMembs)
		return
}

//SuccIssueBook **
func SuccIssueBook(w http.ResponseWriter, r *http.Request) {
	if	validated := s.GetSession(r);!validated{
		http.Redirect(w, r, "/home", http.StatusSeeOther)
		return
	}
		params := mux.Vars(r)
		email := params["email"]
		isbn := params["isbn"]
		fine := 0
		now := time.Now()
		deadline := now.AddDate(0, 0, 15)
		db, err := dbconfig.GetMySQLDb()
		if err != nil {
			panic(err)
		}
		temp, err := db.Prepare(`insert into books_borrowed (isbn,member_email,fine,bowd_time,deadline)
              values(?,?,?,NOW(),?);`)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		_, err = temp.Exec(&isbn, &email, &fine, &deadline)
		if err != nil {
			panic(err)
		}
		_, err = db.Query(`delete from tmp_borrow where mb_email = ? and bk_isbn = ?;`, email, isbn)
		if err != nil {
			panic(err)
		}
		http.Redirect(w, r, "/sci-library/librarian/operations/issue-book", http.StatusSeeOther)
}