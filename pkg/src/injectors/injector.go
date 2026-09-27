package injectors

import (
	"github.com/Juancho43/SOPRO-backend/pkg/core"
	"github.com/Juancho43/SOPRO-backend/pkg/external"
)

func InjectAll(app *core.Core) {
	external.GetFirebaseApp()
	HealthModule(app.Router)
	AuthModule(app)
}
