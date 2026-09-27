package injectors

import (
	"github.com/Juancho43/SOPRO-backend/app/core"
	"github.com/Juancho43/SOPRO-backend/app/external"
)

func InjectAll(app *core.Core) {
	external.GetFirebaseApp()
	HealthModule(app.Router)
	AuthModule(app)
}
