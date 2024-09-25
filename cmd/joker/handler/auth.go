package handler

import (
	"jokegen/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type signInput struct {
	Login    string `json:"login" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=6"`
}
type Users interface {
	CreateUser(login string, password string) error
	LoginUser(login string, password string) (int, error)
	GetRandomJoke(userID int) (service.JokeResponse, error)
	GenerateToken(login, password string) (string, error)
	ParseToken(accessToken string) (int, error)
	AddAttempts(userID, count int) error
	AddAttemptsByLogin(login string, count int) error
}

// @Summary SignUp
// @Tags auth
// @Description create account
// @Accept  json
// @Produce  json
// @Param login query string true "your login"
// @Param password query string true "your password"
// @Router /user/sign-up [post]
func (h *Handler) signUp(c *gin.Context) {
	input := signInput{
		Login:    c.Query("login"),
		Password: c.Query("password"),
	}

	err := h.Users.CreateUser(input.Login, input.Password)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{
		Status: "successful",
	})
}

// @Summary SignIn
// @Tags auth
// @Description login
// @Accept  json
// @Produce  json
// @Param login query string true "your login"
// @Param password query string true "your password"
// @Router /user/sign-in [post]
func (h *Handler) signIn(c *gin.Context) {
	input := signInput{
		Login:    c.Query("login"),
		Password: c.Query("password"),
	}

	token, err := h.Users.GenerateToken(input.Login, input.Password)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"token": token,
	})
}

func (h *Handler) checkUser(c *gin.Context) {
	var input signInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	id, err := h.Users.LoginUser(input.Login, input.Password)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]int{
		"id": id,
	})
}
