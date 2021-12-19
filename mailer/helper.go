package main

import (
	"bytes"
	"embed"
	"fmt"
	mail "github.com/xhit/go-simple-mail/v2"
	"html/template"
	"log"
	"time"
)

//go:embed templates
var emailTemplateFS embed.FS

func (app *application) sendMail(from, to, subject, tmpl string, data interface{}) error {

	templateToRender := fmt.Sprintf("templates/%s.html.tmpl", tmpl)

	t, err := template.New("email-html").ParseFS(emailTemplateFS, templateToRender)
	if err != nil {
		log.Println(err)
		return err
	}

	var tpl bytes.Buffer
	if err = t.ExecuteTemplate(&tpl, "body", data); err != nil {
		log.Println(err)
		return err
	}

	formattedMessage := tpl.String()
	//
	//templateToRender = fmt.Sprintf("templates/%s.plain.tmpl", tmpl)
	//
	//t, err = template.New("email-plain").ParseFS(emailTemplateFS, templateToRender)
	//if err != nil {
	//	log.Println(err)
	//	return err
	//}
	//
	//if err = t.ExecuteTemplate(&tpl, "body", data); err != nil {
	//	log.Println(err)
	//	return err
	//}
	//
	//plainMessage := tpl.String()

	// send the email
	server := mail.NewSMTPClient()

	// SMTP Server
	server.Host = app.smtp.host
	server.Port = app.smtp.port
	server.Username = app.smtp.username
	server.Password = app.smtp.password
	server.Encryption = mail.EncryptionSTARTTLS
	server.KeepAlive = false
	server.ConnectTimeout = 10 * time.Second
	server.SendTimeout = 10 * time.Second

	smtpClient, err := server.Connect()
	if err != nil {
		return err
	}

	email := mail.NewMSG()
	email.SetFrom(from).
		AddTo(to).
		SetSubject(subject)

	email.SetBody(mail.TextHTML, formattedMessage)
	//email.AddAlternative(mail.TextPlain, plainMessage)

	err = email.Send(smtpClient)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}
