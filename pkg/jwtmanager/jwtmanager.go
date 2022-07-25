package jwtmanager

import (
	"time"

	"github.com/ajikamaludin/go-fiber-rest/app/configs"
	"github.com/ajikamaludin/go-fiber-rest/app/models"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"github.com/golang-jwt/jwt/v4"
)

func CreateToken(user *models.User) string {
	configs := configs.GetInstance()

	claims := jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Duration(configs.Jwtconfig.Expired) * time.Second).Unix(),
	}

	unsignToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	token, _ := unsignToken.SignedString([]byte(configs.Jwtconfig.Secret))

	return token
}

func GetUserId(c *fiber.Ctx) (UserId uuid.UUID) {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	UserId, _ = uuid.Parse(claims["user_id"].(string))
	return
}
