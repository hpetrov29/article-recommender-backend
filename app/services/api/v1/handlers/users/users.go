package users

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/mail"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/hpetrov29/resttemplate/business/core/user"
	"github.com/hpetrov29/resttemplate/business/web/v1/auth"
	"github.com/hpetrov29/resttemplate/internal/web"
)

// Handlers manages the set of user endpoints.
type Handlers struct {
	user *user.Core
	auth *auth.Auth
}

// New constructs a new handlers struct for route access.
func New(uc *user.Core, auth *auth.Auth) *Handlers {
	return &Handlers{
		user: uc,
		auth: auth,
	}
}

// Create adds a new user to the system.
func (h *Handlers) Signup(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	var app AppNewUser
	if err := web.Decode(r, &app); err != nil {
		return web.Respond(ctx, w, http.StatusBadRequest, err)
	}

	nc, err := toCoreNewUser(app)
	if err != nil {
		return web.Respond(ctx, w, http.StatusBadRequest, err)
	}

	usr, err := h.user.Create(ctx, nc)
	if err != nil {
		if errors.Is(err, user.ErrUniqueEmail) {
			return web.Respond(ctx, w, http.StatusConflict, errors.New("email is already in use"))
		}
		return web.Respond(ctx, w, http.StatusInternalServerError, err)
	}

	kid := web.Param(r, "kid")
	claims := auth.Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   strconv.FormatInt(usr.Id, 10),
			Issuer:    "service",
			ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		},
	}

	token, err := h.auth.GenerateToken(kid, claims)
	if err != nil {
		fmt.Println(err)
		return web.Respond(ctx, w, http.StatusInternalServerError, errors.New("failed to generate a token"))
	}

	return web.Respond(ctx, w, http.StatusCreated, toAppUser(usr,token))
}

func (h *Handlers) ProtectedRoute(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	w.WriteHeader(200)
	return nil
}

func (h *Handlers) Login(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	kid := web.Param(r, "kid")
	if kid == "" {
		return errors.New("key id not provided")
	}

	email, pass, ok := r.BasicAuth()
	if !ok || strings.Trim(pass, " ") == "" {
		return web.Respond(ctx, w, http.StatusBadRequest, errors.New("must provide email and password"))
	}

	addr, err := mail.ParseAddress(email)
	if err != nil {
		return web.Respond(ctx, w, http.StatusBadRequest, errors.New("invalid email format"))
	}

	usr, err := h.user.Authenticate(ctx, *addr, pass)
	if err != nil {
		return web.Respond(ctx, w, http.StatusUnauthorized,  errors.New("invalid email or password"))
	}

	claims := auth.Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   strconv.FormatInt(usr.Id, 10),
			Issuer:    "service",
			ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		},
		Email: email,
	}
	
	token, err := h.auth.GenerateToken(kid, claims)
	if err != nil {
		return web.Respond(ctx, w,  http.StatusInternalServerError, errors.New("failed to generate a token"))
	}

	http.SetCookie(w, &http.Cookie{
        Name:    "authCookie",
        Value:   token,
        Expires: time.Now().Add(24 * time.Hour),
        Path:    "/",
		Secure: false,
		HttpOnly: true,
		SameSite: http.SameSiteNoneMode,
    })
	
	return web.Respond(ctx, w, http.StatusOK, toToken(token))
}

func (h *Handlers) Me(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	claims, err := h.auth.Authenticate(ctx, r.Header.Get("authorization"))
	if err != nil {
		return web.Respond(ctx, w, http.StatusUnauthorized, err)
	}

	return web.Respond(ctx, w, http.StatusOK, claims);
}