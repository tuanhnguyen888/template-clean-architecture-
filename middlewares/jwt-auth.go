package middlewares

import (
	"log"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/tuanhnguyen888/server/service"
)

// AuthorizeJWT validates the token from the http request, returning a 401 if it's not valid
func AuthorizeJWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		const BEARER_SCHEMA = "Bearer "
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(401, gin.H{
				"message": "Not log in yet",
			})
		}

		tokenString := authHeader[len(BEARER_SCHEMA):]

		token, _ := service.NewJWTService().ValidateToken(tokenString)

		if token.Valid {
			claims := token.Claims.(jwt.MapClaims)
			log.Println("Claims[Name]: ", claims["name"])
			log.Println("Claims[Admin]: ", claims["admin"])
			log.Println("Claims[Issuer]: ", claims["iss"])
			log.Println("Claims[IssuedAt]: ", claims["iat"])
			log.Println("Claims[ExpiresAt]: ", claims["exp"])

			c.Next()
		} else {
			// log.Println(err)
			c.JSON(401, gin.H{
				"message": "NonAuthoritative",
			})
			return
		}
	}
}
