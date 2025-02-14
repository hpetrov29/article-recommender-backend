// Package user provides an example of a core business API. Right now these
// calls are just wrapping the data/data layer.
package user

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/mail"
	"time"

	"github.com/google/uuid"
	"github.com/hpetrov29/resttemplate/internal/logger"
	"golang.org/x/crypto/bcrypt"
)

// Set of error variables for CRUD operations.
var (
	ErrNotFound              = errors.New("user not found")
	ErrUniqueEmail           = errors.New("email is not unique")
	ErrAuthenticationFailure = errors.New("authentication failed")
)

// =============================================================================

// Storer defines the methods required for storing and retrieving data from a user specific repository.
//
// The Storer interface includes methods for creating, deleting, and querying users.
// Implementation is found in business\core\user\stores\usersqldb\usersqldb.go
type Storer interface {
	Create(ctx context.Context, user User) (sql.Result, error)
	Delete(ctx context.Context, user User) error
	QueryByEmail(ctx context.Context, email mail.Address) (User, error)
}

// Core manages the set of APIs for user api access
type Core struct {
	storer Storer
	log *logger.Logger
}

// NewCore constructs and returns a new Core instance for user API access.
//
// Parameters:
//   - st: struct that implements the Storer interface for repository operations.
//   - log: pointer to the logger used for logging within the core.
func NewCore(st Storer, log *logger.Logger) *Core {
	return &Core{
		storer: st, 
		log: log,
	}
}

// Create adds a new user in the repository.
//
// Parameters:
//   - ctx: the context for the request, used for managing timeouts and cancellations.
//   - newUser: the details of the new user to be created.
func (c *Core) Create(ctx context.Context, newUser NewUser) (User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.DefaultCost)
	if err != nil {
		return User{}, fmt.Errorf("generatefrompassword: %w", err)
	}

	now := time.Now()

	usr := User{
		ID:           uuid.New(),
		Name:         newUser.Name,
		Email:        newUser.Email,
		PasswordHash: hash,
		Roles:        newUser.Roles,
		Department:   newUser.Department,
		Enabled:      true,
		DateCreated:  now,
		DateUpdated:  now,
	}

	if _, err := c.storer.Create(ctx, usr); err != nil {
		return User{}, fmt.Errorf("create: %w", err)
	}

	return usr, nil
}

// Delete removes a specified user from the repository.
//
// Parameters:
//   - ctx: the context for the request, used for managing timeouts and cancellations.
//   - usr: the details of the user to be deleted.
func (c *Core) Delete(ctx context.Context, usr User) error {
	if err := c.storer.Delete(ctx, usr); err != nil {
		return fmt.Errorf("delete: %w", err)
	}

	return nil
}

// QueryByEmail retrieves a user from the repository based on their email address.
//
// Parameters:
//   - ctx: the context for the request, used for managing timeouts and cancellations.
//   - email: the email address of the user to be retrieved.
func (c *Core) queryByEmail(ctx context.Context, email mail.Address) (User, error) {
	user, err := c.storer.QueryByEmail(ctx, email)
	if err != nil {
		return User{}, fmt.Errorf("query: email[%s]: %w", email, err)
	}

	return user, nil
}

// =============================================================================

// Authenticate verifies a user's email and password. On success, it returns the claims
// associated with that user which can be used for further authorization.
//
// Parameters:
//   - ctx: the context for managing timeouts and cancellations.
//   - email: the email address of the user to authenticate.
//   - password: the password provided by the user.
//
// Returns:
//   - User - the authenticated user if the email and password are correct.
//   - error - an error if authentication fails or if the user is not found.
func (c *Core) Authenticate(ctx context.Context, email mail.Address, password string) (User, error) {
	usr, err := c.queryByEmail(ctx, email)
	if err != nil {
		return User{}, fmt.Errorf("query: email[%s]: %w", email, err)
	}

	if err := bcrypt.CompareHashAndPassword(usr.PasswordHash, []byte(password)); err != nil {
		return User{}, fmt.Errorf("comparehashandpassword: %w", ErrAuthenticationFailure)
	}

	return usr, nil
}