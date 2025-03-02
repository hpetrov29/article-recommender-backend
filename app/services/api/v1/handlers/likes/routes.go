package likes

import (
	"net/http"

	"github.com/hpetrov29/resttemplate/business/core/like"
	"github.com/hpetrov29/resttemplate/business/core/like/stores/likemessaging"
	"github.com/hpetrov29/resttemplate/business/data/messaging"
	"github.com/hpetrov29/resttemplate/business/web/v1/auth"
	"github.com/hpetrov29/resttemplate/business/web/v1/middleware"
	"github.com/hpetrov29/resttemplate/internal/logger"
	"github.com/hpetrov29/resttemplate/internal/web"
)

// Config contains all the mandatory systems required by handlers.
type Config struct {
	Log   *logger.Logger
	Auth  *auth.Auth
	Messaging messaging.MessagingQueue
}

// Routes initializes the required user specific repositories, service and handler,
// and sets up the API routes for the application with their respective handlers and middlewares.
//
// Parameters:
// 	- app: the web.App instance used to register the routes.
// 	- cfg: configuration including pointers to the logging, database, and authentication systems.
func Routes(app *web.App, cfg Config) {
	
	LikesMessagingQueue := likemessaging.NewStore(cfg.Log, cfg.Messaging, "likes")
	
	likesCore := like.NewCore(LikesMessagingQueue, cfg.Log)

	handlers := New(likesCore, cfg.Auth)

	authenticated := middleware.Authenticate(cfg.Auth)

	// PROTECTED ROUTES
	app.Handle(http.MethodPost, "/like/{post_id}/{is_like}", handlers.Like, authenticated)
}