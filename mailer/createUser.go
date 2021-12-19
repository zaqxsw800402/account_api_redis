package main

import (
	"encoding/json"
	"github.com/nsqio/go-nsq"
	"log"
)

func (app *application) newUser(m *nsq.Message) error {
	var message struct {
		Email string
	}

	err := json.Unmarshal(m.Body, &message)
	if err != nil {
		log.Println(err)
		return err
	}

	var data struct {
		Link string
	}

	data.Link = app.frontend

	// send mailer
	err = app.sendMail("info@ichih_bank.com", message.Email, "Create new account in ichih_bank", "new-user", data)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}
