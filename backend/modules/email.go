package modules

import (
	"bytes"
	"fmt"
	"net/smtp"
	"os"
	"text/template"
)

type EmailData struct {
	ClientName string `json:"name"`
	Email string `json:"email"`
	InvoiceNumber string `json:"invoice_number"`
	Amount string `json:"amount"`
	Status string `json:"status"`
	Semester string `json:"semester"`
	DateTime string `json:"date_time"`
}

func SendActivationMail(data EmailData) {
	auth = smtp.PlainAuth("", "projectsuperapps@gmail.com", "xefolabosopgmcly", "smtp.gmail.com")

	fmt.Println(os.Getwd())
	var filepath = "./conf/email.html"
	r := newRequest([]string{data.Email}, fmt.Sprintf("Invoice [%s] for [%s] on [%s]", data.InvoiceNumber, data.Semester, data.DateTime), "Hello, World!")
	err := r.parseTemplate(filepath, data)

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
