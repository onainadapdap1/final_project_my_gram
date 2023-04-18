package helpers

import (
	"errors"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var SECRET_KEY = []byte("@JevaniaDatubara")

func GenerateToken(id uint, email string) (string, error) {
	// menyimpan data user
	claims := jwt.MapClaims {
		"user_id": id,
		"email": email,
	}

	// enkripsi data user
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// parsing token into long string
	signedToken, err := token.SignedString(SECRET_KEY)
	if err != nil {
		return signedToken, err
	}

	return signedToken, nil
}

func VerifyToken(c *gin.Context) (*jwt.Token, error) {
	errResponse := errors.New("sign in to proceed")
	headerToken := c.GetHeader("Authorization")
	bearer := strings.HasPrefix(headerToken, "Bearer")

	if !bearer {
		return nil, errResponse
	}

	stringToken := strings.Split(headerToken, " ")[1]
	// parsing token menjadi struct pointer dari jwt.Token
	token, _ := jwt.Parse(stringToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errResponse
		}
		return []byte(SECRET_KEY), nil
	})

	if _, ok := token.Claims.(jwt.MapClaims); !ok && !token.Valid {
		return nil, errResponse
	}

	return token, nil
}