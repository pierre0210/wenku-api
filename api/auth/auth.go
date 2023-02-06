package auth

import (
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/pierre0210/wenku-api/internal/database"
	userTable "github.com/pierre0210/wenku-api/internal/database/table/user"
)

type signUpBody struct {
	Username       string `json:"username"`
	Password       string `json:"password"`
	InvitationCode string `json:"code"`
}

type signUpResponse struct {
	Message string `json:"message"`
}

type loginBody struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type loginResponse struct {
	Message     string `json:"message"`
	AccessToken string `json:"accesstoken"`
}

func HandleSignup(c *gin.Context) {
	invitationCode := os.Getenv("INVITATION")
	body := signUpBody{}
	c.BindJSON(&body)

	isUser, _ := userTable.CheckIfUserExists(database.DB, body.Username)
	if isUser {
		c.JSON(401, signUpResponse{Message: "User already exist"})
		return
	} else if body.InvitationCode == invitationCode {
		_, err := userTable.SignUp(database.DB, body.Username, body.Password)
		if err != nil {
			c.JSON(401, signUpResponse{Message: err.Error()})
			return
		}
		c.JSON(200, signUpResponse{Message: "Signed up."})
	}
}

func HandleLogin(c *gin.Context) {
	body := loginBody{}
	c.BindJSON(&body)

	isUser, _ := userTable.CheckIfUserExists(database.DB, body.Username)
	if !isUser {
		c.JSON(401, loginResponse{Message: "Wrong username or password."})
		return
	}

	isMatch, _ := userTable.CheckPassword(database.DB, body.Username, body.Password)
	if isMatch {
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
			Subject:   body.Username,
			ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
		})

		tokeString, err := token.SignedString([]byte(os.Getenv("JWTSECRET")))
		log.Println(err)
		c.JSON(200, loginResponse{Message: "Logged in.", AccessToken: tokeString})
	} else {
		c.JSON(401, loginResponse{Message: "Wrong username or password."})
	}
}
