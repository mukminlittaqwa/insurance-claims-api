package middleware

import (
    "insurance-claims-api/internal/config"
    "insurance-claims-api/internal/utils"
    "net/http"
    "strings"

    "github.com/gin-gonic/gin"
    "github.com/golang-jwt/jwt/v5"
    "go.mongodb.org/mongo-driver/bson/primitive"
)

type Claims struct {
    UserID primitive.ObjectID `json:"user_id"`
    Role   string            `json:"role"`
    jwt.RegisteredClaims
}

func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            utils.ErrorResponse(c, http.StatusUnauthorized, "Authorization header required")
            c.Abort()
            return
        }

        tokenString := strings.Split(authHeader, " ")[1]

        claims := &Claims{}
        token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
            return []byte(config.AppConfig.JWTSecret), nil
        })

        if err != nil || !token.Valid {
            utils.ErrorResponse(c, http.StatusUnauthorized, "Invalid token")
            c.Abort()
            return
        }

        c.Set("user_id", claims.UserID)
        c.Set("role", claims.Role)
        c.Next()
    }
}

func RoleRequired(roles ...string) gin.HandlerFunc {
    return func(c *gin.Context) {
        userRole, _ := c.Get("role")
        allowed := false
        for _, r := range roles {
            if userRole == r {
                allowed = true
                break
            }
        }
        if !allowed {
            utils.ErrorResponse(c, http.StatusForbidden, "Forbidden: insufficient role")
            c.Abort()
            return
        }
        c.Next()
    }
}