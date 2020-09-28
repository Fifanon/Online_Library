package structs

//BookStruct **
type BookStruct struct {
	ISBN          int
	Title         string
	Author        string
	Pages         int
	Subject       string
	Number        int
	BookImageName string
	ImageName     string
	Fine          int
	TimeLeft      string
	Availability  string
	BorrAllowed   bool
	MsgToClient   string
}

//Bk **
var Bk BookStruct
