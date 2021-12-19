package main

import (
	"encoding/json"
	"fmt"
	"github.com/nsqio/go-nsq"
	"log"
	"red/urlsigner"
)

func (app *application) resetPassword(m *nsq.Message) error {
	var message struct {
		Email string
	}

	err := json.Unmarshal(m.Body, &message)
	if err != nil {
		log.Println(err)
		return err
	}

	link := fmt.Sprintf("%s/reset-password?email=%s", app.frontend, message.Email)

	sign := urlsigner.Signer{
		Secret: []byte(app.secretKey),
	}

	signedLink := sign.GenerateTokenFromString(link)

	var data struct {
		Link string
	}

	data.Link = signedLink

	// send mailer
	err = app.sendMail("info@ichih_bank.com", message.Email, "Password Reset Request", "password-reset", data)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}
