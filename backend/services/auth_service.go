package services

import (
	"errors"
	"fmt"
	"time"

	"cepm-backend/config"
	"cepm-backend/models"
	"cepm-backend/repositories"
	"cepm-backend/wechat"

	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

// Claims defines the JWT claims structure.
type Claims struct {
	UserID uint `json:"user_id"`
	jwt.RegisteredClaims
}

// AuthService defines the interface for authentication services.
type AuthService interface {
	WechatLogin(code string) (string, *models.User, error)
	GenerateJWT(userID uint) (string, error)
	ParseJWT(tokenString string) (*Claims, error)
}

type authService struct {
	userRepo     *repositories.UserRepository
	wechatClient *wechat.WechatClient
	jwtSecret    []byte
	jwtExpire    time.Duration
}

// NewAuthService creates a new instance of AuthService.
func NewAuthService(userRepo *repositories.UserRepository, wechatClient *wechat.WechatClient, jwtConfig *config.JWTConfig) AuthService {
	return &authService{
		userRepo:     userRepo,
		wechatClient: wechatClient,
		jwtSecret:    []byte(jwtConfig.SecretKey),
		jwtExpire:    time.Duration(jwtConfig.ExpireHours) * time.Hour,
	}
}

// WechatLogin handles the WeChat Work login process.
func (s *authService) WechatLogin(code string) (string, *models.User, error) {
	// 1. Get user info from WeChat Work using the code
	userInfo, err := s.wechatClient.GetUserInfoByCode(code)
	if err != nil {
		return "", nil, fmt.Errorf("failed to get user info from wechat work: %w", err)
	}

	// 2. Get user detail from WeChat Work using userid
	userDetail, err := s.wechatClient.GetUserDetail(userInfo.UserID)
	if err != nil {
		return "", nil, fmt.Errorf("failed to get user detail from wechat work: %w", err)
	}

	// 3. Find or create user in our database
	user, err := s.userRepo.FindUserByWechatUserid(userDetail.UserID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// User not found, create a new one
			// TODO: Assign default role and department if necessary
			newUser := &models.User{
				WechatUserid: userDetail.UserID,
				Name:         userDetail.Name,
				Email:        userDetail.Email, // WeChat Work might not always provide email
				Avatar:       userDetail.Avatar,
				IsActive:     true,
			}
			if err := s.userRepo.CreateUser(newUser); err != nil {
				return "", nil, fmt.Errorf("failed to create user in db: %w", err)
			}
			user = newUser
		} else {
			return "", nil, fmt.Errorf("failed to find user in db: %w", err)
		}
	}

	// 4. Generate JWT token
	token, err := s.GenerateJWT(user.ID)
	if err != nil {
		return "", nil, fmt.Errorf("failed to generate JWT: %w", err)
	}

	return token, user, nil
}

// GenerateJWT generates a JWT token for the given user ID.
func (s *authService) GenerateJWT(userID uint) (string, error) {
	expirationTime := time.Now().Add(s.jwtExpire)
	claims := &Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(s.jwtSecret)
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, nil
}

// ParseJWT parses and validates a JWT token string.
func (s *authService) ParseJWT(tokenString string) (*Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return s.jwtSecret, nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}
