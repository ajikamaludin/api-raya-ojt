package controllers

import (
	"github.com/ajikamaludin/api-raya-ojt/app/models"
	"github.com/ajikamaludin/api-raya-ojt/app/services"
	"github.com/ajikamaludin/api-raya-ojt/pkg/utils/constants"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

type AuthController struct {
	Service *services.Services
}

func (auth *AuthController) Login(c *fiber.Ctx) error {
	userRequest := new(models.UserLoginReq)

	_ = c.BodyParser(&userRequest)

	errors := auth.Service.Validator.ValidateRequest(userRequest)
	if errors != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"status":  constants.STATUS_FAIL,
			"message": "request body invalid",
			"error":   errors,
		})
	}

	user := &models.User{}
	auth.Service.Repository.GetUserByEmail(userRequest.Email, user)

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userRequest.Password))

	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  constants.STATUS_FAIL,
			"message": "credentials invalid",
			"error":   err.Error(),
		})
	}

	accessToken := auth.Service.JwtManager.CreateToken(user)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  constants.STATUS_SUCCESS,
		"message": "Success login",
		"data": map[string]interface{}{
			"user":         user.ToUserRes(),
			"accessToken":  accessToken,
			"refreshToken": "",
		},
	})
}

func (auth *AuthController) Register(c *fiber.Ctx) error {
	userRequest := new(models.UserRegisterReq)

	_ = c.BodyParser(&userRequest)

	errors := auth.Service.Validator.ValidateRequest(userRequest)
	if errors != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"status":  constants.STATUS_FAIL,
			"message": "request body invalid",
			"error":   errors,
		})
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(userRequest.Password), bcrypt.DefaultCost)
	hashedPin, _ := bcrypt.GenerateFromPassword([]byte(userRequest.Pin), bcrypt.DefaultCost)
	user := &models.User{
		Name:     userRequest.Name,
		Email:    userRequest.Email,
		Password: string(hashedPassword),
		Pin:      string(hashedPin),
	}

	err := auth.Service.Repository.GetUserByEmail(userRequest.Email, user)
	if err == nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"status":  constants.STATUS_FAIL,
			"message": "email is already exists",
			"error":   "",
		})
	}

	auth.Service.Repository.CreateUser(user)

	accessToken := auth.Service.JwtManager.CreateToken(user)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  constants.STATUS_SUCCESS,
		"message": "Success Register",
		"data": map[string]interface{}{
			"user":         user.ToUserRes(),
			"accessToken":  accessToken,
			"refreshToken": "",
		},
	})
}

func (auth *AuthController) ValidatePin(c *fiber.Ctx) error {
	pinRequest := new(models.PinReq)

	_ = c.BodyParser(&pinRequest)

	errors := auth.Service.Validator.ValidateRequest(pinRequest)
	if errors != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"status":  constants.STATUS_FAIL,
			"message": "request body invalid",
			"error":   errors,
		})
	}

	userId := auth.Service.JwtManager.GetUserId(c)
	err := auth.Service.Repository.ValidatePin(userId, pinRequest.Pin)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  constants.STATUS_FAIL,
			"message": "pin invalid",
			"data": map[string]interface{}{
				"status":      2,
				"status_text": "invalid",
			},
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  constants.STATUS_SUCCESS,
		"message": "pin valid",
		"data": map[string]interface{}{
			"status":      1,
			"status_text": "valid",
		},
	})
}

func (auth *AuthController) ErrorHandler(c *fiber.Ctx, err error) error {
	if err.Error() == "Missing or malformed JWT" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  constants.STATUS_FAIL,
			"message": "Unauthorized",
		})
	}

	// Other Error
	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
		"status":  constants.STATUS_FAIL,
		"message": "Unauthorized",
	})
}
