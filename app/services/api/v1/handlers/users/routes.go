package users

import (
	"net/http"

	"github.com/hpetrov29/resttemplate/business/core/user"
	"github.com/hpetrov29/resttemplate/business/core/user/stores/usersqldb"
	"github.com/hpetrov29/resttemplate/business/web/v1/auth"
	"github.com/hpetrov29/resttemplate/business/web/v1/middleware"
	"github.com/hpetrov29/resttemplate/internal/idgenerator"
	"github.com/hpetrov29/resttemplate/internal/logger"
	"github.com/hpetrov29/resttemplate/internal/web"
	"github.com/jmoiron/sqlx"
)

// Config contains all the mandatory systems required by handlers.
type Config struct {
	Log   *logger.Logger
	Auth  *auth.Auth
	DB    *sqlx.DB
	IdGen *idgenerator.IdGenerator
}

// Routes initializes the required user specific repositories, service and handler,
// and sets up the API routes for the application with their respective handlers and middlewares.
//
// Parameters:
// 	- app: the web.App instance used to register the routes.
// 	- cfg: configuration including pointers to the logging, database, and authentication systems.
func Routes(app *web.App, cfg Config) {
	userRepository := usersqldb.NewStore(cfg.Log, cfg.DB)
	userService := user.NewCore(userRepository, cfg.Log, cfg.IdGen)
	handlers := New(userService, cfg.Auth)

	authenticated := middleware.Authenticate(cfg.Auth)
	_ = middleware.Authorize(cfg.Auth, auth.RuleAdminOnly)
	_ = middleware.Authorize(cfg.Auth, auth.RuleAdminOrSubject)

	// NATIVE AUTH
	app.Handle(http.MethodPost, "/users/token/{kid}", handlers.Signup)
	app.Handle(http.MethodGet, "/users/token/{kid}", handlers.Login)
	// PROTECTED ROUTES
	app.Handle(http.MethodGet, "/users", handlers.ProtectedRoute, authenticated)
}
