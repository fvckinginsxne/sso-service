package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v5"

	"auth/internal/domain/models"
)

func NewToken(
	user *models.User,
	jwtSecret string,
	duration time.Duration,
) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["uid"] = user.ID
	claims["exp"] = time.Now().Add(duration).Unix()

	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", err
	}

	return tokenString, err
}

func ParseUserID(tokenStr, jwtSecret string) (int64, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}

		return []byte(jwtSecret), nil
	})
	if err != nil {
		return 0, err
	}

	claims := token.Claims.(jwt.MapClaims)

	uid := claims["uid"].(float64)

	return int64(uid), nil
}
