package likes

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/hpetrov29/resttemplate/business/core/like"
	"github.com/hpetrov29/resttemplate/business/web/v1/auth"
	"github.com/hpetrov29/resttemplate/internal/web"
)

// Handlers manages the set of like endpoints.
type Handlers struct {
	like *like.Core
	auth *auth.Auth
}

// New constructs a new handlers struct for route access.
func New(lc *like.Core, auth *auth.Auth) *Handlers {
	return &Handlers{
		like: lc,
		auth: auth,
	}
}

// Create adds a new like to the system.
func (h *Handlers) Like(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	claims := auth.GetClaims(ctx)
	if claims.Subject == "" {
		return web.Respond(ctx, w, http.StatusUnauthorized, errors.New("authentication failed"))
	}

	postIdStr := web.Param(r, "post_id")
	likeValueStr := web.Param(r, "is_like")
	
	userId, err := strconv.ParseUint(claims.Subject, 10, 64) // Base 10, 64-bit unsigned integer
	if err != nil {
		return web.Respond(ctx, w, http.StatusUnauthorized, fmt.Errorf("authentication failed: %w", err))
	}

	likeValue, err := strconv.ParseInt(likeValueStr, 10, 8)
	if err!= nil || likeValue < -1 || likeValue > 1 {
		return web.Respond(ctx, w, http.StatusBadRequest, fmt.Errorf("interaction type not recognized: %w", err))
	}

	postId, err := strconv.ParseUint(postIdStr, 10, 64)
	if err != nil {
		return web.Respond(ctx, w, http.StatusBadRequest, fmt.Errorf("action cannot be associated to a post: %w", err))
	}

	newLike := like.NewLike{
		Value: 		int8(likeValue),
		UserId: 	userId,
		PostId: 	postId,
	}

	if err = h.like.Publish(ctx, newLike); err != nil {
		return web.Respond(ctx, w, http.StatusInternalServerError, fmt.Errorf("error publishing the like to the queue: %w", err))
	}
	return web.Respond(ctx, w, http.StatusOK, toAppLike(newLike))
}