package v1

import (
	"net/http"
	"os"

	"github.com/hpetrov29/resttemplate/business/data/dbnosql"
	"github.com/hpetrov29/resttemplate/business/data/messaging"
	"github.com/hpetrov29/resttemplate/business/web/v1/auth"
	"github.com/hpetrov29/resttemplate/internal/idgenerator"
	"github.com/hpetrov29/resttemplate/internal/logger"
	"github.com/hpetrov29/resttemplate/internal/web"
	"github.com/jmoiron/sqlx"
)

// APIMuxConfig contains all mandatory systems required by handlers.
type APIMuxConfig struct {
	Build    	string
	Shutdown 	chan os.Signal
	Log      	*logger.Logger
	Auth	 	*auth.Auth
	SQLDB       *sqlx.DB
	NOSQLDB 	dbnosql.NOSQLDB
	Messaging 	messaging.MessagingQueue
	IdGen 	 	*idgenerator.IdGenerator
}

// RouteAdder defines behavior that sets the routes to bind for an instance
// of the service.
type RouteAdder interface {
	Add(app *web.App, cfg APIMuxConfig)
}

func NewAPIMux(config APIMuxConfig, routeAdder RouteAdder) http.Handler {
	app := web.NewApp(config.Shutdown,config.Log ,nil)
	
	// constructs the handlers and binds them to the API endpoints
	routeAdder.Add(app, config)
	return app
}