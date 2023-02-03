package auth

import (
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/pierre0210/wenku-api/internal/database"
	userTable "github.com/pierre0210/wenku-api/internal/database/table/user"
)

type loginBody struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type loginResponse struct {
	Message     string `json:"message"`
	AccessToken string `json:"accesstoken"`
}

func HandleLogin(c *gin.Context) {
	body := loginBody{}
	c.BindJSON(&body)

	isUser, _ := userTable.CheckUser(database.DB, body.Username)
	if !isUser {
		c.JSON(401, loginResponse{Message: "Wrong username or password."})
	}

	isMatch, _ := userTable.CheckPassword(database.DB, body.Username, body.Password)
	if isMatch {
		token := jwt.NewWithClaims(jwt.SigningMethodES256, jwt.StandardClaims{
			Subject:   body.Username,
			ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
		})

		tokeString, _ := token.SignedString([]byte(os.Getenv("JWTSECRET")))
		c.JSON(200, loginResponse{Message: "Logged in.", AccessToken: tokeString})
	} else {
		c.JSON(401, loginResponse{Message: "Wrong username or password."})
	}
}
