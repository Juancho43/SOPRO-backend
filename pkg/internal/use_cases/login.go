package usecases

import (
	"fmt"

	"github.com/Juancho43/SOPRO-backend/pkg/internal/model/entities"
	"github.com/Juancho43/SOPRO-backend/pkg/internal/model/repositories"
)

type LoginUseCase struct {
	repo repositories.UserRepository
}

func NewLoginUseCase(repo repositories.UserRepository) *LoginUseCase {
	return &LoginUseCase{repo: repo}
}

func (uc *LoginUseCase) Execute(user *entities.User) error {
	//TODO: falla en obtener al user. cuando ya existe en la bd
	//HACER Que sea getByEmail
	fmt.Print("email", user.UID)
	_, err := uc.repo.GetUser(user.UID)
	if err != nil {
		errSave := uc.repo.SaveUser(user)
		if errSave != nil {
			return errSave
		}
	}

	return nil
}
