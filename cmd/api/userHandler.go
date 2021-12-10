package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"red/cmd/api/dto"
	"red/cmd/api/errs"
	"red/cmd/api/service"
	"strings"
)

type UserHandlers struct {
	service service.UserService
}

func (u *UserHandlers) createAuthToken(c *gin.Context) {
	var user dto.UserRequest

	err := c.ShouldBindJSON(&user)
	if err != nil {
		badRequest(c, http.StatusBadRequest, err)
		return
	}

	token, err2 := u.service.SaveToken(user)
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

func (u *UserHandlers) authenticateToken(c *gin.Context) (*dto.UserResponse, *errs.AppError) {
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

	user, err := u.service.GetUserWithToken(token)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (u *UserHandlers) name() {

}

func (u *UserHandlers) getAllUsers(c *gin.Context) {

}

func (u *UserHandlers) getUser(c *gin.Context) {

}

func (u *UserHandlers) newUser(c *gin.Context) {

}

func (u *UserHandlers) editUser(c *gin.Context) {

}
