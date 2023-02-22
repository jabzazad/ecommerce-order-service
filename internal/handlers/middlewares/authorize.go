// Package middlewares is a internal handlers middlewares package
package middlewares

import (
	"ecommerce-order/internal/core/config"
	"ecommerce-order/internal/core/context"
	"ecommerce-order/internal/core/redis"
	"ecommerce-order/internal/models"
	"strings"

	"github.com/gofiber/fiber/v2"
)

// AuthAsAdmin authorize as admin
func AuthAsAdmin() fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := context.WithContext(c)
		currentRole := ctx.GetRole()
		if currentRole != models.RoleCustomer {
			return c.
				Status(config.RR.InvalidPermissionRole.HTTPStatusCode()).
				JSON(config.RR.InvalidPermissionRole.WithLocale(c))
		}

		return c.Next()
	}
}

// Authorize authorize
func Authorize() fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := context.WithContext(c)
		conn := redis.GetConnection()
		bearToken := ctx.Get("Authorization")
		strArr := strings.Split(bearToken, " ")
		var userID uint
		err := conn.Get(strArr[1], &userID)
		if err != nil {
			return c.
				Status(config.RR.InvalidToken.HTTPStatusCode()).
				JSON(config.RR.InvalidToken.WithLocale(c))
		}

		if userID != ctx.GetUserID() {
			return c.
				Status(config.RR.InvalidPermissionRole.HTTPStatusCode()).
				JSON(config.RR.InvalidPermissionRole.WithLocale(c))
		}

		return c.Next()
	}
}
