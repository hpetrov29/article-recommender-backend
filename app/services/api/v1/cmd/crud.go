package cmd

import (
	"github.com/hpetrov29/resttemplate/app/services/api/v1/handlers/comments"
	"github.com/hpetrov29/resttemplate/app/services/api/v1/handlers/likes"
	"github.com/hpetrov29/resttemplate/app/services/api/v1/handlers/posts"
	"github.com/hpetrov29/resttemplate/app/services/api/v1/handlers/users"
	v1 "github.com/hpetrov29/resttemplate/business/web/v1"
	"github.com/hpetrov29/resttemplate/internal/web"
)

// Routes constructs the add value which provides the implementation of
// of RouteAdder for specifying what routes to bind to this instance.
func Routes() add {
	return add{}
}

type add struct{}

// Add implements the RouterAdder interface.
func (add) Add(app *web.App, cfg v1.APIMuxConfig) {
	users.Routes(app, users.Config{
		Log:   		cfg.Log,
		Auth:  		cfg.Auth,
		DB:    		cfg.SQLDB,
		IdGen: 		cfg.IdGen,
	})
	posts.Routes(app, posts.Config{
		Log:   		cfg.Log,
		Auth:  		cfg.Auth,
		Cache: 		cfg.Cache,
		SQLDB:    	cfg.SQLDB,
		NOSQLDB: 	cfg.NOSQLDB,
		IdGen: 		cfg.IdGen,
	})
	likes.Routes(app, likes.Config{
		Log:   		cfg.Log,
		Auth:  		cfg.Auth,
		Messaging: 	cfg.Messaging,
	})
	comments.Routes(app, comments.Config{
		Log:   		cfg.Log,
		Auth:  		cfg.Auth,
		SQLDB:    	cfg.SQLDB,
		IdGen: 		cfg.IdGen,
	})
}
