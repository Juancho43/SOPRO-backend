package services

import (
	"github.com/Juancho43/SOPRO-backend/pkg/internal/model/entities"
	usecases "github.com/Juancho43/SOPRO-backend/pkg/internal/use_cases"
)

type AuthService struct {
	loginUseCase *usecases.LoginUseCase
}

func NewAuthService(loginUseCase *usecases.LoginUseCase) *AuthService {
	return &AuthService{
		loginUseCase: loginUseCase,
	}
}

func (s *AuthService) ExecuteGoogleLogin(user *entities.User) error {
	return s.loginUseCase.Execute(user)
}
