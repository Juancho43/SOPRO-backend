package usecases

import (
	"context"
	"errors"
	"regexp"
	"strings"

	"github.com/Juancho43/SOPRO-backend/pkg/internal/model/entities"
	"github.com/Juancho43/SOPRO-backend/pkg/internal/model/repositories"
)

var hexColorRegex = regexp.MustCompile(`^#([A-Fa-f0-9]{6}|[A-Fa-f0-9]{3})$`)

type CreateDomainUseCase struct {
	domainRepo repositories.DomainRepository
}

func NewCreateDomainUseCase(repo repositories.DomainRepository) *CreateDomainUseCase {
	return &CreateDomainUseCase{domainRepo: repo}
}

type CreateDomainCommand struct {
	UserID string `json:"user_id"`
	Name   string `json:"name"`
	Color  string `json:"color"`
}

func (uc *CreateDomainUseCase) Execute(ctx context.Context, cmd CreateDomainCommand) (*entities.Domain, error) {
	name := strings.TrimSpace(cmd.Name)
	if name == "" {
		return nil, errors.New("el nombre del dominio es obligatorio")
	}
	if len(name) > 100 {
		return nil, errors.New("el nombre del dominio no puede superar los 100 caracteres")
	}

	color := strings.TrimSpace(cmd.Color)
	if color == "" {
		color = "#3f51b5"
	} else if !hexColorRegex.MatchString(color) {
		return nil, errors.New("el color debe tener un formato hexadecimal válido (ej. #3f51b5)")
	}

	domain := &entities.Domain{
		UserID: cmd.UserID,
		Name:   name,
		Color:  color,
	}

	if err := uc.domainRepo.Create(ctx, domain); err != nil {
		return nil, err
	}

	return domain, nil
}
