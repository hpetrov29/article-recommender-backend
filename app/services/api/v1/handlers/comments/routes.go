package comments

import (
	"net/http"

	"github.com/hpetrov29/resttemplate/business/core/comment"
	"github.com/hpetrov29/resttemplate/business/core/comment/stores/commentsqldb"
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
	IdGen 		*idgenerator.IdGenerator
}

// Routes initializes the required comment specific repositories, services and handlers,
// and sets up the API routes for the application with their respective handlers and middlewares.
//
// Parameters:
// 	- app: the web.App instance used to register the routes.
// 	- cfg: configuration including pointers to the logging, database, and authentication systems.
func Routes(app *web.App, cfg Config) {
	sqlStore := commentsqldb.NewStore(cfg.Log, cfg.SQLDB)

	userService := comment.NewCore(sqlStore, cfg.Log, cfg.IdGen)

	handlers := New(userService, cfg.Auth)

	authenticated := middleware.Authenticate(cfg.Auth)

	//UNPROTECTED ROUTES
	app.Handle(http.MethodGet, "/comments/{post_id}", handlers.GetComments)

	// PROTECTED ROUTES
	app.Handle(http.MethodPost, "/comment/{post_id}", handlers.CreateComment, authenticated)
	app.Handle(http.MethodDelete, "/comment", handlers.DeleteComment, authenticated)
}
