package operator

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mogenius/punq/dtos"
	"github.com/mogenius/punq/services"
	"github.com/mogenius/punq/structs"
)

type LoginInput struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func InitAuthRoutes(router *gin.Engine) {
	router.POST("/auth/login", login)
	router.GET("/auth/authenticate", Auth(dtos.READER), authenticate)
}

// @Tags Login
// @Produce json
// @Success 200 {object} dtos.PunqToken
// // @Failure 400 {object} httputil.StatusBadRequest
// // @Failure 403 {object} httputil.StatusUnauthorized
// @Router /auth/login [post]
// @Param body body LoginInput true "LoginInput"
func login(c *gin.Context) {
	input := LoginInput{}

	structs.PrettyPrint(input)

	err := c.ShouldBindJSON(&input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}

	user := services.GetUserByEmail(input.Email)
	structs.PrettyPrint(user)

	_, err = user.PasswordCheck(input.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"err": "username or password is incorrect."})
		return
	}

	token, _ := services.GenerateToken(user)

	c.JSON(http.StatusOK, token)
}

func authenticate(c *gin.Context) {
	if temp, exists := c.Get("user"); exists {
		user, ok := temp.(dtos.PunqUser)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{})
			return
		}
		c.JSON(http.StatusOK, user)
		return
	}
	c.JSON(http.StatusUnauthorized, gin.H{})
}
