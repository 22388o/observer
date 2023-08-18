package utils

import (
	"github.com/emersion/go-sasl"
	"github.com/emersion/go-smtp"
	"strings"
)

func SendVerifyCodeMail(email string) (id string, err error) {
	var answer string
	id, answer, _ = GenCaptcha(true)
	err = Mail("your obwallet verify code", answer, email)
	return
}

var mailServer string
var mailUser string
var mailPwd string

func InitMailAuth(server, user, pwd string) {
	mailServer = server
	mailUser = user
	mailPwd = pwd
}
func Mail(title, msgBody, toAddress string) error {
	auth := sasl.NewPlainClient("", mailUser, mailPwd)

	// Connect to the server, authenticate, set the sender and recipient,
	// and send the email all in one step.
	to := []string{toAddress}
	msg := strings.NewReader("To: " + toAddress + "\r\n" +
		"Subject: " + title + "\r\n" +
		"\r\n" +
		msgBody + "\r\n")
	err := smtp.SendMail("smtp.163.com:25", auth, mailUser, to, msg)
	return err
}
