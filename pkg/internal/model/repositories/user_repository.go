package repositories

import "github.com/Juancho43/SOPRO-backend/app/internal/model/entities"

type UserRepository interface {
	SaveUser(user *entities.User) error
	GetUser(uid string) (*entities.User, error)
}
