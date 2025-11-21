package services

import (
    "errors"
    "insurance-claims-api/internal/config"
    "insurance-claims-api/internal/models"
    "insurance-claims-api/internal/repositories"
    "time"

    "github.com/golang-jwt/jwt/v5"
    // "go.mongodb.org/mongo-driver/bson/primitive"
    "golang.org/x/crypto/bcrypt"
)

type AuthService interface {
    Login(req models.LoginRequest) (*models.LoginResponse, error)
}

type authService struct {
    userRepo repositories.UserRepository
}

func NewAuthService(userRepo repositories.UserRepository) AuthService {
    return &authService{userRepo}
}

func (s *authService) Login(req models.LoginRequest) (*models.LoginResponse, error) {
    user, err := s.userRepo.FindByUsername(req.Username)
    if err != nil || bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)) != nil {
        return nil, errors.New("invalid credentials")
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "user_id": user.ID.Hex(),
        "role":    user.Role,
        "exp":     time.Now().Add(time.Hour * 24 * 7).Unix(),
    })

    tokenString, _ := token.SignedString([]byte(config.AppConfig.JWTSecret))

    resp := &models.LoginResponse{
        Token: tokenString,
    }
    resp.User.ID = user.ID
    resp.User.Username = user.Username
    resp.User.Role = user.Role

    return resp, nil
}