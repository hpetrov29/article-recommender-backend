package posts

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/hpetrov29/resttemplate/business/core/post"
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

	appNewPost.UserId = claims.Subject

	if err := web.Decode(r, &appNewPost); err != nil {
		return web.Respond(ctx, w, http.StatusBadRequest, err)
	}

	coreNewPost, err := toCoreNewPost(appNewPost)
	if err != nil {
		return web.Respond(ctx, w, http.StatusBadRequest, err)
	}

	post, err := h.post.Create(ctx, coreNewPost)
	if err != nil {
		return web.Respond(ctx, w, http.StatusInternalServerError, err)
	}

	return web.Respond(ctx, w, http.StatusOK, post)
}

func (h *Handlers) GetPost(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	postId := web.Param(r, "id")
	if strings.Trim(postId, " ") == "" {
		return web.Respond(ctx, w, http.StatusBadRequest, errors.New("must provide a post id"))
	}
	
	post, err := h.post.GetPostById(ctx, postId)
	if err != nil {
		if strings.Contains(err.Error(), "post not found") {	
			return web.Respond(ctx, w, http.StatusNotFound, err)
		}
		return web.Respond(ctx, w, http.StatusInternalServerError, err)
	}
	return web.Respond(ctx, w, http.StatusOK, post)
}

func (h *Handlers) GetPosts(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	h.post.GetPosts(ctx)
	return nil
}