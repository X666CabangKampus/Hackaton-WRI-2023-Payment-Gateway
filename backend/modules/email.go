package modules

import (
	"bytes"
	"fmt"
	"net/smtp"
	"os"
	"text/template"
)

func SendActivationMail(incEmailTo string, incNameTo string, incOTPCode string) {
	auth = smtp.PlainAuth("", "projectsuperapps@gmail.com", "xefolabosopgmcly", "smtp.gmail.com")
	templateData := struct {
		Name    string
		Email   string
		Message string
	}{
		Name:    incNameTo,
		Email:   incEmailTo,
		Message: incOTPCode,
	}

	fmt.Println(os.Getwd())
	var filepath = "./email.html"
	r := newRequest([]string{incEmailTo}, "INVOICE POLINEMA", "Hello, World!")
	err := r.parseTemplate(filepath, templateData)

	if err == nil {
		r.sendEmail()
		fmt.Println("Email Sent...")
	} else {
		fmt.Println(err)
	}

}

func newRequest(to []string, subject, body string) *Request {
	return &Request{
		to:      to,
		subject: subject,
		body:    body,
	}
}

type Request struct {
	from    string
	to      []string
	subject string
	body    string
}

var auth smtp.Auth

func (r *Request) parseTemplate(templateFileName string, data interface{}) error {
	t, err := template.ParseFiles(templateFileName)
	if err != nil {
		return err
	}
	buf := new(bytes.Buffer)
	if err = t.Execute(buf, data); err != nil {
		return err
	}
	r.body = buf.String()
	return nil
}

func (r *Request) sendEmail() (bool, error) {
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	subject := "Subject: " + r.subject + "\n"
	msg := []byte(subject + mime + "\n" + r.body)
	addr := "smtp.gmail.com:587"

	if err := smtp.SendMail(addr, auth, "projectsuperapps@gmail.com", r.to, msg); err != nil {
		return false, err
	}
	return true, nil
}
