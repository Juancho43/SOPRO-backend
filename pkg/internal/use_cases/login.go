package usecases

import (
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
	// Verifica en PostgreSQL si el usuario (por uid de firebase) existe[cite: 4]
	_, err := uc.repo.GetUser(user.UID)
	if err != nil {
		// Si err no es nulo, asumimos que el usuario no existe y lo creamos automáticamente[cite: 5]
		// (Nota: Lo ideal es validar que sea un error del tipo 'Record Not Found' según tu ORM/Driver)
		errSave := uc.repo.SaveUser(user)
		if errSave != nil {
			return errSave
		}
	}

	// Si el usuario existe (err == nil), el flujo termina correctamente[cite: 5]
	return nil
}
