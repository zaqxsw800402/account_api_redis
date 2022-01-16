package mongo

import (
	"fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
)

type Collection struct {
	collection *mgo.Collection
}

type Profile struct {
	UserID    int
	FirstName string
	LastName  string
	Email     string
	Password  string
}

func New(host, username, password string) Collection {
	url := fmt.Sprintf("mongodb://%s:27017", host)
	session, err := mgo.Dial(url)
	if err != nil {
		log.Println("failed to connect to mongodb")
	}
	cred := &mgo.Credential{
		Username: username,
		Password: password,
	}
	err = session.Login(cred)

	db := session.DB("user")
	c := db.C("user.profile")

	return Collection{c}
}

func (c *Collection) ReadProfile(userID int) (*Profile, error) {
	var p Profile
	err := c.collection.Find(bson.M{"userid": userID}).One(&p)
	if err != nil {
		log.Println("failed to read profile in mongodb")
		return nil, err
	}

	return &p, nil
}
