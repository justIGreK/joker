package handler

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
)

// @Summary generate a joke
// @Security BearerAuth
// @Tags joke
// @Description generate joke
// @Accept  json
// @Produce  json
// @Router /jokes/ [get]
func (h *Handler) GetJoke(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, err.Error())
	}
	joke, err := h.services.Users.GetRandomJoke(userId)
	if err != nil {
		if err == sql.ErrNoRows {
			newErrorResponse(c, http.StatusUnauthorized, "login or password are incorrect")
			return
		}
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, joke)
}
