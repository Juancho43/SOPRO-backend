package usecases_test

import (
	"testing"
	"time"

	"github.com/Juancho43/SOPRO-backend/pkg/internal/model/entities"
	"github.com/Juancho43/SOPRO-backend/pkg/internal/model/repositories"
	usecases "github.com/Juancho43/SOPRO-backend/pkg/internal/use_cases"
	// Ajusta este import a tu ruta real
)

// Mock del repositorio para controlar las condiciones de éxito
type mockHabitStreakRepo struct {
	mockGet  func(userID, habitName string) (*entities.HabitStreak, error)
	mockSave func(streak *entities.HabitStreak) error
}

func (m *mockHabitStreakRepo) GetByUserIDAndHabit(userID, habitName string) (*entities.HabitStreak, error) {
	return m.mockGet(userID, habitName)
}

func (m *mockHabitStreakRepo) Save(streak *entities.HabitStreak) error {
	return m.mockSave(streak)
}
func (m *mockHabitStreakRepo) GetAllByUserID(userID string) ([]*entities.HabitStreak, error) {
	// Un retorno vacío cumple la ley del contrato sin afectar la lógica de esta prueba específica
	return nil, nil
}
func TestMarkHabitAsDone_MomentumManagement(t *testing.T) {
	// Anclamos el test a una zona horaria controlada para resultados predecibles
	loc, _ := time.LoadLocation("America/Buenos_Aires")
	now := time.Date(2026, 10, 13, 10, 0, 0, 0, loc)

	yesterday := now.AddDate(0, 0, -1)
	twoDaysAgo := now.AddDate(0, 0, -2)

	tests := []struct {
		name           string
		setupMock      func() repositories.HabitStreakRepository
		expectedErr    string
		expectedStreak int
	}{
		{
			name: "1. Nuevo Comienzo (Deseo Ardiente): Racha inicial en 1",
			setupMock: func() repositories.HabitStreakRepository {
				return &mockHabitStreakRepo{
					mockGet: func(userID, habitName string) (*entities.HabitStreak, error) {
						return nil, nil // No existe el hábito aún
					},
					mockSave: func(streak *entities.HabitStreak) error {
						if streak.CurrentStreak != 1 {
							t.Errorf("Se esperaba racha de 1, se obtuvo %d", streak.CurrentStreak)
						}
						return nil
					},
				}
			},
			expectedErr:    "",
			expectedStreak: 1,
		},
		{
			name: "2. Cero Días (Enfoque y Prioridades): Rechaza múltiples registros el mismo día",
			setupMock: func() repositories.HabitStreakRepository {
				return &mockHabitStreakRepo{
					mockGet: func(userID, habitName string) (*entities.HabitStreak, error) {
						return &entities.HabitStreak{
							CurrentStreak: 5,
							UpdatedAt:     now, // Actualizado hoy mismo
						}, nil
					},
					mockSave: func(streak *entities.HabitStreak) error {
						return nil // No debería llegar aquí
					},
				}
			},
			expectedErr:    "ya has forjado tu disciplina hoy; ahora enfócate en ejecutar",
			expectedStreak: 5,
		},
		{
			name: "3. Un Día (El Poder del Momentum): Incrementa la racha",
			setupMock: func() repositories.HabitStreakRepository {
				return &mockHabitStreakRepo{
					mockGet: func(userID, habitName string) (*entities.HabitStreak, error) {
						return &entities.HabitStreak{
							CurrentStreak: 5,
							MaxStreak:     5,
							UpdatedAt:     yesterday, // Actualizado ayer
						}, nil
					},
					mockSave: func(streak *entities.HabitStreak) error {
						if streak.CurrentStreak != 6 {
							t.Errorf("Se esperaba racha de 6, se obtuvo %d", streak.CurrentStreak)
						}
						return nil
					},
				}
			},
			expectedErr:    "",
			expectedStreak: 6,
		},
		{
			name: "4. Más de Un Día (Superando el Fracaso Temporal): Reinicia la racha a 1",
			setupMock: func() repositories.HabitStreakRepository {
				return &mockHabitStreakRepo{
					mockGet: func(userID, habitName string) (*entities.HabitStreak, error) {
						return &entities.HabitStreak{
							CurrentStreak: 5,
							MaxStreak:     5,
							UpdatedAt:     twoDaysAgo, // Dejó pasar más de un día
						}, nil
					},
					mockSave: func(streak *entities.HabitStreak) error {
						if streak.CurrentStreak != 1 {
							t.Errorf("La racha debía reiniciarse a 1, se obtuvo %d", streak.CurrentStreak)
						}
						if streak.MaxStreak != 5 {
							t.Errorf("La racha máxima nunca debe perderse, se esperaba 5, se obtuvo %d", streak.MaxStreak)
						}
						return nil
					},
				}
			},
			expectedErr:    "",
			expectedStreak: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := tt.setupMock()
			uc, err := usecases.NewMarkHabitAsDone(repo, "America/Buenos_Aires")
			if err != nil {
				t.Fatalf("Error inesperado al crear caso de uso: %v", err)
			}

			err = uc.Execute("user123", "Lectura Diaria", now)

			// Validación de errores
			if tt.expectedErr != "" {
				if err == nil || err.Error() != tt.expectedErr {
					t.Errorf("Se esperaba el error %q, pero se obtuvo %v", tt.expectedErr, err)
				}
			} else if err != nil {
				t.Errorf("No se esperaba error, se obtuvo %v", err)
			}
		})
	}
}
