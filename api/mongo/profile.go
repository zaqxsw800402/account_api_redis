package mongo

import (
	"fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
	"red/dto"
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
		log.Println(err.Error())
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

func (c *Collection) SaveProfile(userID int, user dto.UserRequest) error {
	p := &Profile{
		UserID:    userID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Password:  user.Password,
	}
	err := c.collection.Insert(p)
	if err != nil {
		log.Println("failed to save profile in mongodb")
		return err
	}
	return nil
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

func (c *Collection) UpdatedProfile(userID int, user dto.UserRequest) error {
	original, err := c.ReadProfile(userID)
	if err != nil {
		log.Println("failed to read profile in mongodb")
		return err
	}
	if len(user.FirstName) > 0 {
		original.FirstName = user.FirstName
	}

	if len(user.LastName) > 0 {
		original.LastName = user.LastName
	}

	selector := bson.M{"userid": userID}

	err = c.collection.Update(selector, original)
	if err != nil {
		log.Println("failed to read profile in mongodb")
		return err
	}

	return nil
}
