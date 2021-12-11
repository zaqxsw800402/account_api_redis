package main

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"red/cmd/api/dto"
	"red/cmd/api/errs"
	"red/cmd/api/service"
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

func (app *application) authenticateToken(c *gin.Context) (*dto.UserResponse, *errs.AppError) {
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

	user, err := app.uh.service.GetUserWithToken(token)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (app *application) name() {

}

func (app *application) getAllUsers(c *gin.Context) {

}

func (app *application) getUser(c *gin.Context) {

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

	jsonResp(c, http.StatusOK)
}

func (app *application) editUser(c *gin.Context) {

}
