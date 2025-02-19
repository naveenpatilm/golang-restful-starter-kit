package apis

import (
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-ozzo/ozzo-routing"
	"github.com/go-ozzo/ozzo-routing/auth"
	"github.com/naveenpatilm/golang-restful-starter-kit/app"
	"github.com/naveenpatilm/golang-restful-starter-kit/errors"
	"github.com/naveenpatilm/golang-restful-starter-kit/models"
)

type Credential struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func Auth(signingKey string) routing.Handler {
	return func(c *routing.Context) error {
		var credential Credential
		if err := c.Read(&credential); err != nil {
			return errors.Unauthorized(err.Error())
		}

		identity := authenticate(credential)
		if identity == nil {
			return errors.Unauthorized("invalid credential")
		}

		token, err := auth.NewJWT(jwt.MapClaims{
			"id":   identity.GetID(),
			"name": identity.GetName(),
			"exp":  time.Now().Add(time.Hour * 72).Unix(),
		}, signingKey)
		if err != nil {
			return errors.Unauthorized(err.Error())
		}

		return c.Write(map[string]string{
			"token": token,
		})
	}
}

func authenticate(c Credential) models.Identity {
	if c.Username == "demo" && c.Password == "pass" {
		return &models.User{ID: "100", Name: "demo"}
	}
	return nil
}

func JWTHandler(c *routing.Context, j *jwt.Token) error {
	userID := j.Claims.(jwt.MapClaims)["id"].(string)
	app.GetRequestScope(c).SetUserID(userID)
	return nil
}
