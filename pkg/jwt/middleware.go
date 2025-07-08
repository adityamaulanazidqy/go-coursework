package jwt

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"go-coursework/internal/dto"
	"golang.org/x/net/context"
	"os"
	"strings"
	"time"
)

type Claims struct {
	UserID int    `json:"user_id"`
	Email  string `json:"email"`
	Roles  string `json:"role"`
	jwt.RegisteredClaims
}

var (
	jwtSecret   = []byte(os.Getenv("JWT_SECRET"))
	redisClient *redis.Client
	logLogrus   logrus.Logger
)

func SetRedisClientMiddleware(rdb *redis.Client) {
	redisClient = rdb
}

func Middleware(allowedRoles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(dto.ErrorResponse{
				Message: "Missing token",
				Details: nil,
			})
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.Status(fiber.StatusUnauthorized).JSON(dto.ErrorResponse{
				Message: "Invalid token format",
				Details: nil,
			})
		}

		tokenStr := parts[1]
		claims := &Claims{}

		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (
			interface{},
			error,
		) {
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(dto.ErrorResponse{
				Message: "Invalid or expired token",
				Details: nil,
			})
		}

		if redisClient != nil {
			ctxRedis := context.Background()
			blacklisted, err := redisClient.Get(ctxRedis, "blacklist:"+tokenStr).Result()
			if err == nil && blacklisted == "true" {
				return c.Status(fiber.StatusUnauthorized).JSON(dto.ErrorResponse{
					Message: "Token has been logged out",
					Details: nil,
				})
			}
		}

		if len(allowedRoles) > 0 {
			roleMatch := false
			for _, role := range allowedRoles {
				if claims.Roles == role {
					roleMatch = true
					break
				}
			}
			if !roleMatch {
				return c.Status(fiber.StatusForbidden).JSON(dto.ErrorResponse{
					Message: "Forbidden",
					Details: nil,
				})
			}
		}

		c.Locals("user", claims)
		return c.Next()
	}
}

func MiddlewareSocket(allowedRoles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		tokenStr := c.Query("token")
		if tokenStr == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(dto.ErrorResponse{
				Message: "Missing token",
				Details: nil,
			})
		}

		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (
			interface{},
			error,
		) {
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(dto.ErrorResponse{
				Message: "Invalid or expired token",
				Details: nil,
			})
		}

		if redisClient != nil {
			ctxRedis := context.Background()
			blacklisted, err := redisClient.Get(ctxRedis, "blacklist:"+tokenStr).Result()
			if err == nil && blacklisted == "true" {
				return c.Status(fiber.StatusUnauthorized).JSON(dto.ErrorResponse{
					Message: "Token has been logged out",
					Details: nil,
				})
			}
		}

		if len(allowedRoles) > 0 {
			roleMatch := false
			for _, role := range allowedRoles {
				if claims.Roles == role {
					roleMatch = true
					break
				}
			}
			if !roleMatch {
				return c.Status(fiber.StatusForbidden).JSON(dto.ErrorResponse{
					Message: "Forbidden",
					Details: nil,
				})
			}
		}

		c.Locals("user", claims)
		return c.Next()
	}
}

func GenerateToken(userID int, email, role string) (
	string,
	error,
) {
	expirationTime := time.Now().Add(24 * time.Hour)

	claims := &Claims{
		UserID: userID,
		Email:  email,
		Roles:  role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func ExtractTokenFromHeader(c *fiber.Ctx) (
	string,
	error,
) {
	bearerToken := c.Get("Authorization")
	parts := strings.Split(bearerToken, " ")
	if len(parts) == 2 && parts[0] == "Bearer" {
		return parts[1], nil
	}

	return "", errors.New("invalid token format")
}

func VerifyToken(tokenStr string) (
	*Claims,
	error,
) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (
		interface{},
		error,
	) {
		return jwtSecret, nil
	})

	if err != nil || !token.Valid {
		logLogrus.WithFields(logrus.Fields{
			"error":   err,
			"message": "invalid token or expired token",
		}).Error("invalid token or expired token")

		return nil, errors.New("invalid token or expired token")
	}

	return claims, nil
}
