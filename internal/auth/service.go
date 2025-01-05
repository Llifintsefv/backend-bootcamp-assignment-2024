package auth

import (
	"backend-bootcamp-assignment-2024/dto"
	"context"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type AuthService interface {
	generateJWToken(userType string) (string, error)
	registerUser(ctx context.Context, req dto.PostRegisterJSONRequestBody) (string, error)
	loginUser(ctx context.Context, req dto.PostLoginJSONRequestBody) (string, error)
}
type authService struct {
	secretKey      []byte
	authRepository AuthRepository
}

func NewAuthService(secretKey []byte, authRepository AuthRepository) AuthService {
	return &authService{
		secretKey:      secretKey,
		authRepository: authRepository,
	}
}

func (s *authService) generateJWToken(userType string) (string, error) {
	claims := jwt.MapClaims{
		"user_type": userType,
		"exp":       time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(s.secretKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, nil
}

func (s *authService) registerUser(ctx context.Context, req dto.PostRegisterJSONRequestBody) (string, error) {
	UserUUID, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}
	err = s.authRepository.RegisterUser(ctx, req, UserUUID.String())
	if err != nil {
		return "", err
	}
	return UserUUID.String(), nil
}

func (s *authService) loginUser(ctx context.Context, req dto.PostLoginJSONRequestBody) (string, error) {

	userType, err := s.authRepository.loginUser(ctx, req)
	if err != nil {
		return "", err
	}

	token, err := s.generateJWToken(userType)

	if err != nil {
		return "", err
	}

	return token, nil

}
