package member

import(
	"database/sql"
	vars "github.com/Fifanon/online_library/varsAndFuncs"
	stct "github.com/Fifanon/online_library/structs"
	"net/http"
	dbconfig "github.com/Fifanon/online_library/config"
)

//ProcessMemberSearch **
func ProcessMemberSearch(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	var myuser stct.Users
	myuser.Email = r.Form.Get("email")

	db, err := dbconfig.GetMySQLDb()
	if err != nil {
		panic(err)
	}
	rows := db.QueryRow(`select m_firstname,m_lastname,m_email,m_address,m_telephone,m_status,m_imagename FROM members where m_email = ?;`, myuser.Email)
	err = rows.Scan(&myuser.FirstName, &myuser.LastName, &myuser.Email, &myuser.Address, &myuser.PhoneNum, &myuser.Status, &myuser.ImageName)
	db.Close()
	if err != nil {
		if err == sql.ErrNoRows {
			stct.Msg.EmailExistNot = "user does not exist"
			vars.Tpl.ExecuteTemplate(w, "memberCancellation.html", stct.Msg)
			stct.Msg.EmailExistNot = ""
			return
		}
	}
	vars.Tpl.ExecuteTemplate(w, "MemberSearched.html", myuser)
	stct.Msg.EmailExistNot = ""
}