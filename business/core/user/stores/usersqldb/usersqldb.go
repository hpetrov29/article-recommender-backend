package usersqldb

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/mail"
	"strings"

	"github.com/hpetrov29/resttemplate/business/core/user"
	db "github.com/hpetrov29/resttemplate/business/data/dbsql/mysql"
	"github.com/hpetrov29/resttemplate/internal/logger"
	"github.com/jmoiron/sqlx"
)

// Store manages the set of APIs for user database access.
type Store struct {
	log *logger.Logger
	db  sqlx.ExtContext
}

// NewStore constructs the api for interacting with a relational database.
//
// Parameters:
//   - log: pointer to the logger used for logging within the store.
//   - db: pointer to the database connection used by the store.
//
// Returns:
//   - *Store: a pointer to the newly created Store instance.
func NewStore(log *logger.Logger, db *sqlx.DB) *Store {
	return &Store{
		log: log,
		db:  db,
	}
}

// Create inserts a new user record into the database.
//
// Parameters:
//   - ctx: the context for managing timeouts and cancellations.
//   - usr: the user data to be stored in the database.
//
// Returns:
//   - sql.Result: the result of the SQL insert operation, containing details such as rows affected.
//   - error: an error if the insertion fails, including a specific error if the email is not unique (Error 1062).
func (s *Store) Create(ctx context.Context, usr user.User) (sql.Result, error) {
	const q = `
	INSERT INTO users
		(user_id, name, email, password_hash, roles, enabled, department, date_created, date_updated)
	VALUES
		(:user_id, :name, :email, :password_hash, :roles, :enabled, :department, :date_created, :date_updated)`
	
	res, err := db.NamedExecContext(ctx, s.log, s.db, q, toDBUser(usr)); 
	
	if err != nil {
		if (strings.Split(err.Error(),":")[0] == "Error 1062 (23000)") {
			return nil, fmt.Errorf("namedexeccontext: %w", user.ErrUniqueEmail)
		}
		return nil, fmt.Errorf("namedexeccontext: %w", err)
	}

	return res, nil
}

// Delete removes a user record from the database based on the user's ID.
//
// Parameters:
//   - ctx: context for managing timeouts and cancellations.
//   - usr: user data containing the ID of the user to be deleted.
//
// Returns:
//   - error: an error if the deletion fails. If successful, returns nil.
func (s *Store) Delete(ctx context.Context, usr user.User) error {
	data := struct {
		UserID string `db:"user_id"`
	}{
		UserID: usr.ID.String(),
	}

	const q = `
	DELETE FROM
		users
	WHERE
		user_id = :user_id`

	if _, err := db.NamedExecContext(ctx, s.log, s.db, q, data); err != nil {
		return fmt.Errorf("namedexeccontext: %w", err)
	}

	return nil
}

// QueryByEmail retrieves a user from the database using their email address.
//
// Parameters:
//   - ctx: context for managing timeouts and cancellations.
//   - email: email address of the user to query.
//
// Returns:
//   - user.User: user associated with the given email address.
//   - error: an error if the query fails or the user is not found. Returns user.ErrNotFound if the email does not exist.
func (s *Store) QueryByEmail(ctx context.Context, email mail.Address) (user.User, error) {
	data := struct {
		Email string `db:"email"`
	}{
		Email: email.Address,
	}

	const q = `
	SELECT
        user_id, name, email, password_hash, roles, enabled, department, date_created, date_updated
	FROM
		users
	WHERE
		email = :email`

	var dbUsr dbUser
	if err := db.NamedQueryStruct(ctx, s.log, s.db, q, data, &dbUsr); err != nil {
		if errors.Is(err, db.ErrDBNotFound) {
			return user.User{}, fmt.Errorf("namedquerystruct: %w", user.ErrNotFound)
		}
		return user.User{}, fmt.Errorf("namedquerystruct: %w", err)
	}

	usr, err := toCoreUser(dbUsr)
	if err != nil {
		return user.User{}, err
	}

	return usr, nil
}