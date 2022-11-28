package authorization

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
)

var SECRET = os.Getenv("AUTH_SECRET")

func CreateToken(customerId uint32, accountId uint32) (string, error) {
	claim := jwt.MapClaims{
		"customerId": customerId,
		"accountId":  accountId,
		"exp":        time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)

	t, err := token.SignedString([]byte(SECRET))

	if err != nil {
		return "", err
	}

	return t, nil
}

func VerifyToken() func(c *fiber.Ctx) error {

	f := jwtware.New(jwtware.Config{
		SigningKey: []byte(SECRET),
		AuthScheme: "Bearer",
	})
	return f
}

func GetTokenSecret(tokenStr string) (uint32, uint32) {
	tokenStr = strings.Replace(tokenStr, "Bearer ", "", 1)

	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return []byte(SECRET), nil
	})
	if err != nil {
		fmt.Println(err.Error())
	}

	subject := token.Claims.(jwt.MapClaims)
	return uint32(subject["accountId"].(float64)), uint32(subject["accountId"].(float64))
}
