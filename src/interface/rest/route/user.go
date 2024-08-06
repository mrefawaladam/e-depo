package route

import (
	"net/http"

	"e-depo/src/infra/helper"
	handlers "e-depo/src/interface/rest/handlers/user"

	"github.com/go-chi/chi/v5"
)

// HealthRouter a completely separate router for health check routes
func UserRouter(h handlers.UserHandlerInterface) http.Handler {
	r := chi.NewRouter()

	r.With(helper.RoleCheckMiddleware("admin")).Post("/create_user", h.CreateUser)
	r.Post("/login", h.Login)

	return r
}
