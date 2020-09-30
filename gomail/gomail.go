package gomail

import (
   "fmt"
   "net/smtp"
)

// smtpServer data to smtp server.
type smtpServer struct {
	host string
	port string
}

// Address URI to smtp server.
func (s *smtpServer) Address() string {
	return s.host + ":" + s.port
}

//SendEmail **
func SendEmail(toEmail string, msg string, subj string) (done bool, err error){
	// Sender data.
	from := "scilibrary6@gmail.com"
	password := "zxcasd,.23"

	// Receiver email address.
	to := []string{toEmail}

	// smtp server configuration.
	smtpServer := smtpServer{host: "smtp.gmail.com", port: "587"}

	// Message.
	mime := "MIME-version: 1.0;\nContent-Type: text/plain; charset=\"UTF-8\";\n\n"
	subject := "Subject: " + subj + "!\n"
	message := []byte(subject + mime +"\n" + msg)

	// Authentication.
	auth := smtp.PlainAuth("", from, password, smtpServer.host)

	// Sending email.
	err = smtp.SendMail(smtpServer.Address(), auth, from, to, message)
	if err != nil {
		fmt.Println(err)
		return false, err
	}
		
	fmt.Println("Email Sent!")
	return true, nil
}

