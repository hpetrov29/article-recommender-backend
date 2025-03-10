package posts

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/hpetrov29/resttemplate/business/core/post"
	"github.com/hpetrov29/resttemplate/business/data/page"
	"github.com/hpetrov29/resttemplate/business/web/v1/auth"
	"github.com/hpetrov29/resttemplate/internal/web"
)

// Handlers manages the set of user endpoints.
type Handlers struct {
	post *post.Core
	auth *auth.Auth
}

// New constructs a new handlers struct for route access.
func New(pc *post.Core, auth *auth.Auth) *Handlers {
	return &Handlers{
		post: pc,
		auth: auth,
	}
}

func (h *Handlers) CreatePost(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	var appNewPost AppNewPost
	claims := auth.GetClaims(ctx)
	if claims.Subject == "" {
		return web.Respond(ctx, w, http.StatusUnauthorized, errors.New("authentication failed"))
	}

	userId, err := strconv.ParseInt(claims.Subject, 10, 64) // Base 10, 64-bit unsigned integer
	if err != nil {
		return web.Respond(ctx, w, http.StatusUnauthorized, fmt.Errorf("authentication failed: %w", err))
	}

	if err := web.Decode(r, &appNewPost); err != nil {
		return web.Respond(ctx, w, http.StatusBadRequest, err)
	}

	coreNewPost := toCoreNewPost(appNewPost, userId)

	post, err := h.post.Create(ctx, coreNewPost)
	if err != nil {
		return web.Respond(ctx, w, http.StatusInternalServerError, err)
	}

	return web.Respond(ctx, w, http.StatusOK, toAppPost(post))
}

func (h *Handlers) QueryById(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	id, err := strconv.ParseInt(web.Param(r, "id"), 10, 64); if err != nil {
		return web.Respond(ctx, w, http.StatusBadRequest, err)
	}

	corePost, err := h.post.QueryById(ctx, id)
	if err != nil {
		if errors.Is(err, post.ErrNotFound) {	
			return web.Respond(ctx, w, http.StatusNotFound, err)
		}
		
		return web.Respond(ctx, w, http.StatusInternalServerError, err)
	}
	
	return web.Respond(ctx, w, http.StatusOK, toAppPost(corePost))
}

func (h *Handlers) Query(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	page, err := page.Parse(r)
	if err != nil {
		return web.Respond(ctx, w, http.StatusBadRequest, err)
	}

	filter, err := parseFilter(r)
	if err != nil {
		return web.Respond(ctx, w, http.StatusBadRequest, err)
	}

	orderBy, err := parseOrder(r)
	if err != nil {
		return web.Respond(ctx, w, http.StatusBadRequest, err)
	}

	posts, err := h.post.Query(ctx, filter, orderBy, page.Number, page.RowsPerPage)
	if err != nil {
		return web.Respond(ctx, w, http.StatusInternalServerError, err)
	}

	return web.Respond(ctx, w, http.StatusOK, toAppPosts(posts))
}