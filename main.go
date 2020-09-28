package main

import (
	"html/template"
	"database/sql"
	"net/http"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	_"github.com/gorilla/securecookie"
	_ "github.com/go-sql-driver/mysql"
	vars "github.com/Fifanon/online_library/varsAndFuncs"
	stct "github.com/Fifanon/online_library/structs"
	user "github.com/Fifanon/online_library/user"
	searchbk "github.com/Fifanon/online_library/searchBook"
	bk "github.com/Fifanon/online_library/operationsOnBook"
	mber "github.com/Fifanon/online_library/operationsOnMemb"
	brwbk "github.com/Fifanon/online_library/borrowBook"
	disp "github.com/Fifanon/online_library/displayBooks"
	s "github.com/Fifanon/online_library/session"

)

var db *sql.DB
var store *sessions.CookieStore

func main() {
	vars.GotFile = false
	r := mux.NewRouter()
	r.HandleFunc("/home", handler).Methods("GET")
	r.HandleFunc("/login-process", user.LoginProcessor).Methods("POST")
	r.HandleFunc("/sci-library/librarian/operations", librarianOp).Methods("GET")

	r.HandleFunc("/sci-library/welcome", user.MemberLoggedIn).Methods("GET")
	r.HandleFunc("/sci-library/book-searching", searchbk.ProcessBookSearch).Methods("POST")
	r.HandleFunc("/sci-library/op/member-searching", mber.ProcessMemberSearch).Methods("POST")

	r.HandleFunc("/sci-library/sign-up", signUp).Methods("GET")
	r.HandleFunc("/sci-library/sign-up/check-email", user.CheckEmail).Methods("POST")
	r.HandleFunc("/sign-up/processing", user.UploadPhotoFile).Methods("POST")
	r.HandleFunc("/signup/processing-continue", user.SignupProcessor)
	r.HandleFunc("/sci-library/log-in/forgot-password", restorePwd).Methods("GET")

	r.HandleFunc("/sci-library/librarian/operations/cancel-member", mber.CancelMember)
	r.HandleFunc("/sci-library/librarian/operations/cancel-this-member/{email}", mber.CancelThisMember)

	r.HandleFunc("/sci-library/librarian/operations/add-member", mber.AddMember)
	r.HandleFunc("/sci-library/librarian/op/add-member/{email}", mber.AddMembervalidate).Methods("GET")
	r.HandleFunc("/sci-library/librarian/op/delete-request-member/{email}", mber.DeleteRequest).Methods("GET")
	r.HandleFunc("/sci-library/librarian/operations/list-of-members", mber.MembersList)

	r.HandleFunc("/sci-library/librarian/op/issue-book/{isbn}/{email}", bk.SuccIssueBook)
	r.HandleFunc("/sci-library/librarian/operations/issue-book", bk.IssuedBook)
	r.HandleFunc("/sci-library/member/list-books-borrowed", bk.ListOfBooksBorrowed)

	r.HandleFunc("/sci-library/librarian/operations/add-book", bk.AddBook)
	r.HandleFunc("/sci-library/librarian/op/book-adding", bk.BookAdding).Methods("POST")
	r.HandleFunc("/sci-library/librarian/operations/update-book", bk.UpdateBook)
	r.HandleFunc("/sci-library/book-to-update-search", bk.UpdateBookSearch)
	r.HandleFunc("/sci-library/librarian/book-updating/{isbn}", bk.UpdateBookprocessing).Methods("POST")


	r.HandleFunc("/sci-library/librarian/operations/remove-book", bk.RemoveBook)
	r.HandleFunc("/sci-library/book-to-delete-search", bk.RemoveBookSearch)
	r.HandleFunc("/sci-library/librarian/book-removing/{isbn}", bk.RemoveBookprocessing)

	r.HandleFunc("/sci-library/librarian/operations/return-book", bk.ReturnBook)
	r.HandleFunc("/sci-library/book-return/member-search", bk.ReturnBookMemberSearch).Methods("POST")
	r.HandleFunc("/sci-library/librarian/book-return/{isbn}", bk.SuccReturnBook)

	r.HandleFunc("/sci-library/list-of-allbooks/ascending-order",  disp.DisplayListOfbook).Methods("GET")
	r.HandleFunc("/sci-library/physics-books", disp.DisplayPhysicsBooks).Methods("GET")
	r.HandleFunc("/sci-library/mathematics-books",  disp.DisplayMathBooks).Methods("GET")
	r.HandleFunc("/sci-library/electronic-eng-books",  disp.DisplayElecBooks).Methods("GET")
	r.HandleFunc("/sci-library/computer-eng-books",  disp.DisplayCompEngBooks).Methods("GET")
	r.HandleFunc("/sci-library/mechanical-eng-books",  disp.DisplayMechEngBooks).Methods("GET")
	r.HandleFunc("/sci-library/biology-books",  disp.DisplayBiologyBooks).Methods("GET")
	r.HandleFunc("/sci-library/librarian/operations/statistics", bk.ProcessStatistics).Methods("GET")
	r.HandleFunc("/sci-library/member/statistics", mber.ProcessMemberStatistics).Methods("GET")

	r.HandleFunc("/sci-library/member/borrowbook/{isbn}", brwbk.BorrowBook).Methods("GET")
	r.HandleFunc("/sci-library/librarian/operations/members-owe-fine", mber.MembersOwingFine).Methods("GET")
	r.HandleFunc("/sci-library/memberBooksborrowed", brwbk.SuccBorrow).Methods("GET")
	r.HandleFunc("/sci-library/log-out/", logOut)

	fs := http.FileServer(http.Dir("./project_files/"))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))

	http.ListenAndServe(":9080", r)
}

func init() {
	vars.Tpl = template.Must(template.ParseGlob("./project_files/templates/*.html"))
}
//rendering home page
func handler(w http.ResponseWriter, r *http.Request) {
	vars.Tpl.ExecuteTemplate(w, "home.html", stct.Msg)
	stct.Msg.LoginBefore = ""
	stct.Msg.EmailExistNot = ""
	stct.Msg.WrongPwd = ""
	return
}

func signUp(w http.ResponseWriter, r *http.Request) {
	//   var msg Msg
	// msg.fstmsg = fileuploadmsg
	//   msg.sndmsg = vars.Message
	vars.Tpl.ExecuteTemplate(w, "signup.html", vars.Fileuploadmsg)
	// vars.Message = ""
	vars.Fileuploadmsg = ""
	return
}

func librarianOp(w http.ResponseWriter, r *http.Request) {
	if	validated := s.GetSession(r);!validated{
		http.Redirect(w, r, "/home", http.StatusSeeOther)
		return
	}
	vars.Tpl.ExecuteTemplate(w, "operations.html", nil)
}

func restorePwd(w http.ResponseWriter, r *http.Request) {
	vars.Tpl.ExecuteTemplate(w, "pwdRestoration.html", nil)
}


func logOut(w http.ResponseWriter, r *http.Request) {
	s.ClearSessionHandler(w,r)
	return
}
