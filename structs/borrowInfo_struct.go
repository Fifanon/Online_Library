package structs

//BorrowInfo **
type BorrowInfo struct {
	ISBN          int
	Title         string
	Author        string
	Pages         int
	Subject       string
	Number        int
	BookImageName string
	ImageName     string
	FirstName     string
	LastName      string
	Email         string
	Address         string
	Status         string
	Fine          int
	TimeLeft      int
	Deadline      string

}

//BooknBorr **
var BooknBorr BorrowInfo