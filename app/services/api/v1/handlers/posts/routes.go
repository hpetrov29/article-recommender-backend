package posts

import (
	"net/http"

	"github.com/hpetrov29/resttemplate/business/core/post"
	"github.com/hpetrov29/resttemplate/business/core/post/stores/postnosqldb"
	"github.com/hpetrov29/resttemplate/business/core/post/stores/postsqldb"
	"github.com/hpetrov29/resttemplate/business/data/dbnosql"
	"github.com/hpetrov29/resttemplate/business/web/v1/auth"
	"github.com/hpetrov29/resttemplate/business/web/v1/middleware"
	"github.com/hpetrov29/resttemplate/internal/idgenerator"
	"github.com/hpetrov29/resttemplate/internal/logger"
	"github.com/hpetrov29/resttemplate/internal/web"
	"github.com/jmoiron/sqlx"
)

// Config contains all the mandatory systems required by handlers.
type Config struct {
	Log  		*logger.Logger
	Auth 		*auth.Auth
	SQLDB   	*sqlx.DB
	NOSQLDB 	dbnosql.NOSQLDB
	IdGen 		*idgenerator.IdGenerator
}

// Routes initializes the required post specific repositories, services and handlers,
// and sets up the API routes for the application with their respective handlers and middlewares.
//
// Parameters:
// 	- app: the web.App instance used to register the routes.
// 	- cfg: configuration including pointers to the logging, database, and authentication systems.
func Routes(app *web.App, cfg Config) {
	nosqlRepo := cfg.NOSQLDB.GetRepository("posts")

	sqlStore := postsqldb.NewStore(cfg.Log, cfg.SQLDB)
	hybridStore := postnosqldb.NewStore(cfg.Log, sqlStore, nosqlRepo)

	userService := post.NewCore(hybridStore, cfg.Log, cfg.IdGen)

	handlers := New(userService, cfg.Auth)

	authenticated := middleware.Authenticate(cfg.Auth)
	_ = middleware.Authorize(cfg.Auth, auth.RuleAdminOnly)
	_ = middleware.Authorize(cfg.Auth, auth.RuleAdminOrSubject)

	// UNPROTECTED ROUTES
	app.Handle(http.MethodGet, "/post/{id}", handlers.QueryById)
	app.Handle(http.MethodGet, "/posts", handlers.Query)

	// PROTECTED ROUTES
	app.Handle(http.MethodPost, "/post", handlers.CreatePost, authenticated)
}
