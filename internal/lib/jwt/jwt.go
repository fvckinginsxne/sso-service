package jwt

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"time"

	"sso/internal/domain/models"
)

func NewToken(user models.User, app models.App, duration time.Duration) (string, error) {
	const op = "jwt.NewToken"

	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)

	claims["uid"] = user.ID
	claims["email"] = user.Email
	claims["exp"] = time.Now().Add(duration).Unix()
	claims["app_id"] = app.ID

	tokenString, err := token.SignedString([]byte(app.Secret))
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return tokenString, nil
}
