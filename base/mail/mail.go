package mail

import (
	"utils/error"

	"gopkg.in/gomail.v2"
)

//func SendToMail(user, password, host, to, subject, body, mailtype string) error {
// user := "yang**@yun*.com"
// 	password := "***"
// 	host := "smtp.exmail.qq.com:25"
// 	to := "397685131@qq.com"
//
// 	subject := "使用Golang发送邮件"
//
// 	body := `
// 		<html>
// 		<body>
// 		<h3>
// 		"Test send to email"
// 		</h3>
// 		</body>
// 		</html>
// 		`

func Send(user, userTitle, password, host string, port int, to, toTitle, subject, body string) {
	m := gomail.NewMessage()
	m.SetAddressHeader("From", user, userTitle)
	m.SetHeader("To", m.FormatAddress(to, toTitle))
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	d := gomail.NewPlainDialer(host, port, user, password)
	if err := d.DialAndSend(m); err != nil {
		error.Try(2000, 3, "utils/base/mail/Send/DialAndSend", err)
		return
	}
}
