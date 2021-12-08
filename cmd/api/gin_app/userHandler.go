package gin_app

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"red/cmd/api/dto"
	"red/cmd/api/service"
)

type UserHandlers struct {
	service service.UserService
}

func (u UserHandlers) CreateAuthToken(c *gin.Context) {
	var user dto.UserRequest

	err := c.ShouldBindJSON(&user)
	if err != nil {
		badRequest(c, http.StatusBadRequest, err)
		return
	}

	token, err := u.service.SaveToken(user)
	if err != nil {
		badRequest(c, http.StatusBadRequest, err)
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

func (u *UserHandlers) getAllUsers(c *gin.Context) {

}

func (u *UserHandlers) getUser(c *gin.Context) {

}

func (u *UserHandlers) newUser(c *gin.Context) {

}

func (u *UserHandlers) editUser(c *gin.Context) {

}
