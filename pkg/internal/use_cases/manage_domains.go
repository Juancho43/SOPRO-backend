package usecases

import (
	"context"
	"errors"
	"strings"

	"github.com/Juancho43/SOPRO-backend/pkg/internal/model/entities"
	"github.com/Juancho43/SOPRO-backend/pkg/internal/model/repositories"
)

type ListDomainsUseCase struct {
	domainRepo repositories.DomainRepository
}

func NewListDomainsUseCase(repo repositories.DomainRepository) *ListDomainsUseCase {
	return &ListDomainsUseCase{domainRepo: repo}
}

func (uc *ListDomainsUseCase) Execute(ctx context.Context, userID string) ([]*entities.DomainWithHabitCount, error) {
	if userID == "" {
		return nil, errors.New("el userID es obligatorio")
	}
	return uc.domainRepo.ListByUserID(ctx, userID)
}

type UpdateDomainUseCase struct {
	domainRepo repositories.DomainRepository
}

func NewUpdateDomainUseCase(repo repositories.DomainRepository) *UpdateDomainUseCase {
	return &UpdateDomainUseCase{domainRepo: repo}
}

type UpdateDomainCommand struct {
	ID     string `json:"id"`
	UserID string `json:"user_id"`
	Name   string `json:"name"`
	Color  string `json:"color"`
}

func (uc *UpdateDomainUseCase) Execute(ctx context.Context, cmd UpdateDomainCommand) error {
	name := strings.TrimSpace(cmd.Name)
	if name == "" {
		return errors.New("el nombre del dominio es obligatorio")
	}
	if len(name) > 100 {
		return errors.New("el nombre del dominio no puede superar los 100 caracteres")
	}

	color := strings.TrimSpace(cmd.Color)
	if color != "" && !hexColorRegex.MatchString(color) {
		return errors.New("el color debe tener un formato hexadecimal válido")
	}

	domain := &entities.Domain{
		ID:     cmd.ID,
		UserID: cmd.UserID,
		Name:   name,
		Color:  color,
	}

	return uc.domainRepo.Update(ctx, domain)
}

type DeleteDomainUseCase struct {
	domainRepo repositories.DomainRepository
}

func NewDeleteDomainUseCase(repo repositories.DomainRepository) *DeleteDomainUseCase {
	return &DeleteDomainUseCase{domainRepo: repo}
}

func (uc *DeleteDomainUseCase) Execute(ctx context.Context, id string, userID string) error {
	if id == "" || userID == "" {
		return errors.New("el ID del dominio y del usuario son requeridos")
	}
	return uc.domainRepo.Delete(ctx, id, userID)
}
