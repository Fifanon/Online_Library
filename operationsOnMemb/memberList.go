package member

import(
	s "github.com/Fifanon/online_library/session"
	vars "github.com/Fifanon/online_library/varsAndFuncs"
	stct "github.com/Fifanon/online_library/structs"
	"net/http"
	dbconfig "github.com/Fifanon/online_library/config"
)

//MembersList **
func MembersList(w http.ResponseWriter, r *http.Request) {
	if	validated := s.GetSession(r);!validated{
		http.Redirect(w, r, "/home", http.StatusSeeOther)
		return
	}
		db, err := dbconfig.GetMySQLDb()
		if err != nil {
			panic(err)
		}
		qr, err := db.Query(`select * from members;`)
		if err != nil {
			panic(err)
		}
		var member stct.Users
		members := []stct.Users{}
		for qr.Next() {
			err = qr.Scan(&member.FirstName, &member.LastName, &member.Email, &member.Address, &member.PhoneNum, &member.Password, &member.Status, &member.ImageName, &member.DateTime)
			if err != nil {
				panic(err)
			}
			members = append(members, member)
		}
		vars.Tpl.ExecuteTemplate(w, "ListOfAllMembers.html", members)
		return
}