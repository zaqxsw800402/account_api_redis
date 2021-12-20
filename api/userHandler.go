package main

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"red/dto"
	"red/encryption"
	"red/errs"
	"red/service"
	"strings"
)

type UserHandlers struct {
	service service.UserService
}

func (app *application) createAuthToken(c *gin.Context) {
	var user dto.UserRequest

	err := c.ShouldBindJSON(&user)
	if err != nil {
		badRequest(c, http.StatusBadRequest, err)
		return
	}

	token, err2 := app.uh.service.SaveToken(user)
	if err2 != nil {
		badRequest(c, http.StatusBadRequest, err2)
		return
	}

	var payload struct {
		Error   bool               `json:"error"`
		Message string             `json:"message"`
		Token   *dto.TokenResponse `json:"authentication_token"`
	}

	payload.Error = false
	payload.Message = "Success!"
	payload.Token = token

	c.JSON(http.StatusOK, payload)
}

func (app *application) authenticateToken(c *gin.Context) (*int, *errs.AppError) {
	request := c.Request
	authorizationHeader := request.Header.Get("Authorization")
	if authorizationHeader == "" {
		return nil, &errs.AppError{
			Code:    http.StatusBadRequest,
			Message: "no authorization header received",
		}
	}

	headersParts := strings.Split(authorizationHeader, " ")
	if len(headersParts) != 2 || headersParts[0] != "Bearer" {
		return nil, &errs.AppError{
			Code:    http.StatusBadRequest,
			Message: "no authorization header received",
		}
	}

	token := headersParts[1]
	if len(token) != 26 {
		return nil, &errs.AppError{
			Code:    http.StatusBadRequest,
			Message: "authentication token wrong size",
		}
	}

	// redis
	result := app.redis.GetUserID(token)
	i, err2 := result.Int()
	//if err2 != nil {
	//	log.Println(err2)
	//	return nil, &errs.AppError{
	//		Code:    http.StatusUnauthorized,
	//		Message: "user not in redis",
	//	}
	//}

	switch {
	case err2 == redis.Nil || i == 0:
		t, err := app.uh.service.GetUserWithToken(token)
		if err != nil {
			return nil, err
		}

		id := int(t.UserID)
		app.redis.SaveUserID(token, id)

		return &id, nil

	case err2 != nil:
		log.Println(err2)
		return nil, &errs.AppError{
			Code:    http.StatusBadRequest,
			Message: "user not in redis",
		}

	default:
		return &i, nil
	}

}

func (app *application) newUser(c *gin.Context) {
	var user dto.UserRequest

	err := c.ShouldBindJSON(&user)
	if err != nil {
		badRequest(c, http.StatusBadRequest, err)
		return
	}

	newHash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 12)
	if err != nil {
		badRequest(c, http.StatusBadRequest, err)
		return
	}

	_, appError := app.uh.service.SaveUser(user, string(newHash))
	if appError != nil {
		badRequest(c, appError.Code, appError)
		return
	}

	var nsqMessage struct {
		Op    string
		Email string
	}

	nsqMessage.Op = "newUser"
	nsqMessage.Email = user.Email

	marshal, err := json.Marshal(nsqMessage)
	if err != nil {
		log.Println(err)
	}
	err = app.mail.Publish("mailer", marshal)
	if err != nil {
		log.Println(err)
	}

	jsonResp(c, http.StatusOK)
}

func (app *application) SendPasswordResetEmail(c *gin.Context) {
	var payload struct {
		Email string `json:"email"`
	}

	err := c.ShouldBindJSON(&payload)
	if err != nil {
		badRequest(c, http.StatusBadRequest, err)
		return
	}

	_, appError := app.uh.service.GetUserByEmail(payload.Email)
	if appError != nil {
		badRequest(c, appError.Code, appError)
		return
	}

	var nsqMessage struct {
		Op    string
		Email string
	}

	nsqMessage.Op = "forgotPassword"
	nsqMessage.Email = payload.Email

	marshal, err := json.Marshal(nsqMessage)
	if err != nil {
		log.Println(err)
	}
	err = app.mail.Publish("mailer", marshal)
	if err != nil {
		log.Println(err)
	}

	jsonResp(c, http.StatusCreated)
}

func (app *application) ResetPassword(c *gin.Context) {
	var payload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := c.ShouldBindJSON(&payload)
	if err != nil {
		badRequest(c, http.StatusBadRequest, err)
	}

	encryptor := encryption.Encryption{Key: []byte(app.config.secretKey)}
	realEmail, err := encryptor.Decrypt(payload.Email)
	if err != nil {
		badRequest(c, http.StatusBadRequest, err)
		return
	}

	_, appError := app.uh.service.GetUserByEmail(realEmail)
	if appError != nil {
		badRequest(c, appError.Code, appError)
	}

	newHash, err := bcrypt.GenerateFromPassword([]byte(payload.Password), 12)
	if err != nil {
		badRequest(c, http.StatusBadRequest, err)
		return
	}

	appError = app.uh.service.UpdatePassword(realEmail, string(newHash))
	if err != nil {
		badRequest(c, http.StatusBadRequest, appError)
		return
	}

	var resp struct {
		Error   bool   `json:"error"`
		Message string `json:"message"`
	}

	resp.Error = false
	resp.Message = "password changed"

	c.JSON(http.StatusCreated, resp)
}
