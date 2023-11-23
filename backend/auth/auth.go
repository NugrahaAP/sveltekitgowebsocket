package auth

import (
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

func GenerateRefreshToken(userId string) (string, error) {
	var secretKey = []byte(os.Getenv("SECRET_KEY"))

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId": userId,
		"StandarClaim": jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 6)),
		},
		"TokenType": "refresh",
	})
	refreshTokenString, err := refreshToken.SignedString(secretKey)
	if err != nil {
		log.Println("Error when signing JWT: ", err)
		return "", err
	}

	return refreshTokenString, nil
}

func GenerateToken(userId string) (string, error) {
	var secretKey = []byte(os.Getenv("SECRET_KEY"))

	// bikin JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"UserId": userId,
		"StandarClaim": jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 60)),
		},
		"TokenType": "access",
	})

	tokenString, err := token.SignedString(secretKey)

	if err != nil {
		log.Println("Error when signing JWT: ", err)
		return "", err
	}

	return tokenString, nil

}
