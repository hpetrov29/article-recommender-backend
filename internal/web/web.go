package web

import (
	"context"
	"errors"
	"net/http"
	"os"
	"syscall"
	"time"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"github.com/hpetrov29/resttemplate/internal/logger"
)

// Handler is a type definition that handles a http request within the mini framework.
type Handler func(ctx context.Context, w http.ResponseWriter, r *http.Request) error

// App is the entrypoint into the application and what configures our context
// object for each of our http handlers. Feel free to add any configuration
// data/logic on this App struct.
type App struct {
	mux *chi.Mux
	shutdown chan os.Signal
	middlewares []Middleware
	Log *logger.Logger
	Version string
}

// NewApp creates an App instance using the chi router.
func NewApp(shutdown chan os.Signal, log *logger.Logger, middlewares ...Middleware) *App {

	// TO DO: Create an OpenTelemetry HTTP Handler which wraps our router.
	mux := chi.NewMux();

	return &App{
		mux: mux,
		shutdown: shutdown,
		middlewares: middlewares,
		Log: log,
		Version: "v1",
	}
}

// SignalShutdown is used to gracefully shut down the app when an integrity
// issue is identified.
func (a *App) SignalShutdown() {
	a.shutdown <- syscall.SIGTERM
}

// ServeHTTP method implements the http.Handler interface for App. 
// It's the entry point for all http traffic.
func (a *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.mux.ServeHTTP(w, r)
}

// Handle wraps the `Handler` handler around with the specified middlewares and
// adds the route `path` that matches `method` http method to execute it.
// 
// Parameters:
// 	- method: http method.
// 	- path: the pattern to be matched against.
//	- handler: the handler function to be executed.
//	- middlewares: middlewares that run before the handler gets executed.
func (a *App) Handle(method string, path string, handler Handler, middlewares ...Middleware) {
	handler = wrapMiddleware(middlewares, handler)
	handler = wrapMiddleware(a.middlewares, handler)

	group := a.Version;
	a.handle(method, group, path, handler)
}

// =============================================================================

// handle sets a handler function for a given HTTP method and path pair
// to the application server mux.
func (a *App) handle(method string, group string, path string, handler Handler) {
	h := func(w http.ResponseWriter, r *http.Request) {
		vals := &Values{
			TraceId:    uuid.New().String(),
			Now:        time.Now(),
			StatusCode: 0,
		}
		ctx := context.WithValue(context.Background(), key, vals)

		if err := handler(ctx, w, r); err != nil {
			if validateShutdown(err) {
				a.SignalShutdown()
				return
			}
		}
	}

	finalPath := path
	if group != "" {
		finalPath = "/" + group + path
	}

	a.mux.MethodFunc(method, finalPath, h)
}

// validateShutdown validates the error for special conditions that do not
// warrant an actual shutdown by the system.
func validateShutdown(err error) bool {

	// Ignore syscall.EPIPE and syscall.ECONNRESET errors which occurs
	// when a write operation happens on the http.ResponseWriter that
	// has simultaneously been disconnected by the client (TCP
	// connections is broken). For instance, when large amounts of
	// data is being written or streamed to the client.
	// https://blog.cloudflare.com/the-complete-guide-to-golang-net-http-timeouts/
	// https://gosamples.dev/broken-pipe/
	// https://gosamples.dev/connection-reset-by-peer/

	switch {
	case errors.Is(err, syscall.EPIPE):

		// Usually, you get the broken pipe error when you write to the connection after the
		// RST (TCP RST Flag) is sent.
		// The broken pipe is a TCP/IP error occurring when you write to a stream where the
		// other end (the peer) has closed the underlying connection. The first write to the
		// closed connection causes the peer to reply with an RST packet indicating that the
		// connection should be terminated immediately. The second write to the socket that
		// has already received the RST causes the broken pipe error.
		return false

	case errors.Is(err, syscall.ECONNRESET):

		// Usually, you get connection reset by peer error when you read from the
		// connection after the RST (TCP RST Flag) is sent.
		// The connection reset by peer is a TCP/IP error that occurs when the other end (peer)
		// has unexpectedly closed the connection. It happens when you send a packet from your
		// end, but the other end crashes and forcibly closes the connection with the RST
		// packet instead of the TCP FIN, which is used to close a connection under normal
		// circumstances.
		return false
	}

	return true
}
