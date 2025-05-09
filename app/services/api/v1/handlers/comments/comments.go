package comments

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/hpetrov29/resttemplate/business/core/comment"
	"github.com/hpetrov29/resttemplate/business/web/v1/auth"
	"github.com/hpetrov29/resttemplate/internal/web"
)

// Handlers manages the set of user endpoints.
type Handlers struct {
	comment *comment.Core
}

// New constructs a new handlers struct for route access.
func New(cc *comment.Core, auth *auth.Auth) *Handlers {
	return &Handlers{
		comment: cc,
	}
}

func (h *Handlers) CreateComment(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	var newComment NewAppComment

	userId, err := strconv.ParseUint(auth.GetClaims(ctx).Subject, 10, 64) // Base 10, 64-bit unsigned integer
	if err != nil {
		return web.Respond(ctx, w, http.StatusUnauthorized, fmt.Errorf("authentication failed: %w", err))
	}

	if err := web.Decode(r, &newComment); err != nil {
		return web.Respond(ctx, w, http.StatusBadRequest, fmt.Errorf("error parsing the body of the request: %w", err))
	}

	postId, err := strconv.ParseUint(web.Param(r, "post_id"), 10, 64) // Base 10, 64-bit unsigned integer
	if err != nil {
		return web.Respond(ctx, w, http.StatusBadRequest, fmt.Errorf("error parsing post_id param: %w", err))
	}

	coreNewComment := toCoreNewComment(newComment, int64(userId), int64(postId))

	coreComement, err := h.comment.Create(ctx, coreNewComment)
	if err != nil {
		return web.Respond(ctx, w, http.StatusInternalServerError, err)
	}

	return web.Respond(ctx, w, http.StatusOK, toAppComment(coreComement))
}

func (h *Handlers) DeleteComment(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	var deleteComment DeleteComment

	userId, err := strconv.ParseUint(auth.GetClaims(ctx).Subject, 10, 64) // Base 10, 64-bit unsigned integer
	if err != nil {
		return web.Respond(ctx, w, http.StatusUnauthorized, fmt.Errorf("authentication failed: %w", err))
	}

	if err := web.Decode(r, &deleteComment); err != nil {
		return web.Respond(ctx, w, http.StatusBadRequest, fmt.Errorf("error parsing the body of the request: %w", err))
	}

	// prevents malicous comment deletion (a comment to be deleted by other than its author)
	if (userId != deleteComment.UserId) {
		return web.Respond(ctx, w, http.StatusUnauthorized, errors.New("unauthorized action"))
	}

	if err = h.comment.Delete(ctx, deleteComment.Id); err != nil{
		if errors.Is(err, comment.ErrNotFound) {
			return web.Respond(ctx, w, http.StatusBadRequest, err)
		}
		return web.Respond(ctx, w, http.StatusInternalServerError, err)
	}

	return web.Respond(ctx, w, http.StatusOK, fmt.Sprintf("Deletion of comment with id: %d successful.", deleteComment.Id))
}

func (h *Handlers) GetComments(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	id, err := strconv.ParseInt(web.Param(r, "post_id"), 10, 64); if err != nil {
		return web.Respond(ctx, w, http.StatusBadRequest, err)
	}

	comments, err := h.comment.QueryByPostId(ctx, id)
	if err != nil {
		return web.Respond(ctx, w, http.StatusInternalServerError, err)
	}

	return web.Respond(ctx, w, http.StatusOK, toAppComments(comments))
}