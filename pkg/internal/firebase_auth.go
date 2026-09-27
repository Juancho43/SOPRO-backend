package internal

import (
	"context"
	"fmt"
	"log/slog"

	"firebase.google.com/go/v4/auth"
	"github.com/Juancho43/SOPRO-backend/pkg/internal/model/entities"
)

type FirebaseVerifier struct {
	client *auth.Client
	logger *slog.Logger
}

func NewFirebaseVerifier(client *auth.Client, logger *slog.Logger) *FirebaseVerifier {
	return &FirebaseVerifier{
		client: client,
		logger: logger,
	}
}

func (f *FirebaseVerifier) Verify(idToken string) (*entities.User, error) {
	f.logger.Debug("Iniciando verificación de token en Firebase")

	ctx := context.Background()
	token, err := f.client.VerifyIDToken(ctx, idToken)
	if err != nil {
		f.logger.Warn("Error al verificar token en Firebase", "error", err)
		return nil, fmt.Errorf("error al verificar token en Firebase: %w", err)
	}

	if token == nil {
		f.logger.Warn("El token decodificado es nulo")
		return nil, fmt.Errorf("el token decodificado es nulo")
	}

	email := ""
	if emailClaim, ok := token.Claims["email"].(string); ok {
		email = emailClaim
	}

	name := ""
	if nameClaim, ok := token.Claims["name"].(string); ok {
		name = nameClaim
	}

	f.logger.Info("Token verificado exitosamente en Firebase", "uid", token.UID, "email", email)

	return &entities.User{
		UID:   token.UID,
		Email: email,
		Name:  name,
	}, nil
}
