package internal

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Juancho43/SOPRO-backend/pkg/internal/model/entities"
)

type HTTPFocusmateAPI struct {
	client   *http.Client
	apiKey   string
	location *time.Location
}

func NewHTTPFocusmateAPI(apiKey string, timezoneStr string) (*HTTPFocusmateAPI, error) {
	loc, err := time.LoadLocation(timezoneStr)
	if err != nil {
		return nil, fmt.Errorf("error al cargar la zona horaria: %v", err)
	}

	return &HTTPFocusmateAPI{
		client:   &http.Client{Timeout: 10 * time.Second},
		apiKey:   apiKey,
		location: loc,
	}, nil
}

// fetchPartnerName es el método privado que consulta el endpoint del perfil
func (api *HTTPFocusmateAPI) fetchPartnerName(ctx context.Context, userID string) string {
	url := fmt.Sprintf("https://api.focusmate.com/v1/users/%s", userID)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return "Socio (Oculto)"
	}

	req.Header.Set("X-API-KEY", api.apiKey)

	resp, err := api.client.Do(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		return "Socio (Oculto)"
	}
	defer resp.Body.Close()

	var result struct {
		User struct {
			Name string `json:"name"`
		} `json:"user"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "Socio"
	}

	if result.User.Name != "" {
		return result.User.Name
	}

	return "Socio"
}

func (api *HTTPFocusmateAPI) FetchFutureSessions(ctx context.Context) ([]entities.FocusmateSession, error) {
	now := time.Now().In(api.location)

	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, api.location)
	endOfToday := startOfDay.AddDate(0, 0, 1).Add(-time.Nanosecond)

	startISO := startOfDay.Format(time.RFC3339)
	endISO := endOfToday.Format(time.RFC3339)

	url := fmt.Sprintf("https://api.focusmate.com/v1/sessions?start=%s&end=%s", startISO, endISO)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("X-API-KEY", api.apiKey)

	resp, err := api.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("fallo en la API de Focusmate: status %d", resp.StatusCode)
	}

	var result struct {
		Sessions []struct {
			SessionID   string    `json:"sessionId"`
			StartTime   time.Time `json:"startTime"`
			Duration    int       `json:"duration"`
			SessionType string    `json:"sessionType"`
			Title       *string   `json:"title"`
			Users       []struct {
				UserID    string `json:"userId"`
				Role      string `json:"role"`
				Completed *bool  `json:"completed"`
			} `json:"users"`
		} `json:"sessions"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	namesCache := make(map[string]string)
	var mappedSessions []entities.FocusmateSession

	for _, s := range result.Sessions {
		// ELIMINADO: Ya no usamos el 'continue' para saltar las sesiones canceladas.
		// Todo el registro del día se procesa.

		endTime := s.StartTime.Add(time.Duration(s.Duration) * time.Millisecond)

		// Evaluamos la realidad temporal de la sesión
		status := "upcoming"

		// Verificamos si la sesión fracasó (cancelada o inasistencia)
		isFailedSession := len(s.Users) > 0 && s.Users[0].Completed != nil && !*s.Users[0].Completed

		if isFailedSession {
			status = "canceled" // Transparencia absoluta sobre la falta
		} else if endTime.Before(now) {
			status = "completed"
		} else if s.StartTime.Before(now) && endTime.After(now) {
			status = "active"
		}

		// Resolución dinámica del nombre
		partnerName := "Socio"

		if s.SessionType == "group" {
			if s.Title != nil {
				partnerName = *s.Title
			} else if len(s.Users) > 1 {
				hostID := s.Users[1].UserID

				if cachedName, exists := namesCache[hostID]; exists {
					partnerName = cachedName + " (Host)"
				} else {
					fetchedName := api.fetchPartnerName(ctx, hostID)
					namesCache[hostID] = fetchedName
					partnerName = fetchedName + " (Host)"
				}
			}
		} else {
			// Sesión 1 a 1 (Paired)
			if len(s.Users) > 1 {
				partnerID := s.Users[1].UserID

				if cachedName, exists := namesCache[partnerID]; exists {
					partnerName = cachedName
				} else {
					fetchedName := api.fetchPartnerName(ctx, partnerID)
					namesCache[partnerID] = fetchedName
					partnerName = fetchedName
				}
			} else {
				partnerName = "Esperando socio..."
			}
		}

		mappedSessions = append(mappedSessions, entities.FocusmateSession{
			ID:          s.SessionID,
			StartTime:   s.StartTime.In(api.location),
			EndTime:     endTime.In(api.location),
			Status:      status,
			PartnerName: partnerName,
		})
	}

	return mappedSessions, nil
}
