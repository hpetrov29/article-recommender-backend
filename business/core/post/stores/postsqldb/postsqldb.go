package postsqldb

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/hpetrov29/resttemplate/business/core/post"

	db "github.com/hpetrov29/resttemplate/business/data/dbsql/mysql"
	"github.com/hpetrov29/resttemplate/internal/logger"
	"github.com/jmoiron/sqlx"
)

// Store manages the set of APIs for user database access.
type Store struct {
	log *logger.Logger
	db  sqlx.ExtContext
}

// NewStore constructs the api required for interacting with a relational database.
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

// Create inserts a new post record into the database.
//
// Parameters:
//   - ctx: the context for managing timeouts and cancellations.
//   - post: the post contents to be stored in the database.
//
// Returns:
//   - sql.Result: the result of the SQL insert operation, containing details such as rows affected.
//   - error: an error if the insertion fails.
func (s *Store) Create(ctx context.Context, post post.Post) (sql.Result, error) {
	const q = `
	INSERT INTO posts
		(id, title, content, user_id, date_created, date_updated)
	VALUES
		(:id, :title, :content, :user_id, :date_created, :date_updated);`
	
	res, err := db.NamedExecContext(ctx, s.log, s.db, q, toDBPost(post)); 
	
	if err != nil {
		return nil, fmt.Errorf("namedexeccontext: %w", err)
	}

	return res, nil
}

// Delete removes a post from the database based on the post's Id.
//
// Parameters:
//   - ctx: context for managing timeouts and cancellations.
//   - post: the post containing the Id of the post to be deleted.
//
// Returns:
//   - error: an error if the deletion fails. If successful, returns nil.
func (s *Store) Delete(ctx context.Context, post post.Post) error {
	data := struct {
		PostId string `db:"id"`
	}{
		PostId: post.Id.String(),
	}

	const q = `
	DELETE FROM
		posts
	WHERE
		id = :id`

	if _, err := db.NamedExecContext(ctx, s.log, s.db, q, data); err != nil {
		return fmt.Errorf("namedexeccontext: %w", err)
	}

	return nil
}

func (s *Store) QueryById(ctx context.Context, id string) (post.Post, error) {
	data := struct {
		id string `db:"id"`
	}{
		id: id,
	}

	const q = `
	SELECT
        *
	FROM
		posts
	WHERE
		id = :id;`

	var dbPost dbPost
	if err := db.NamedQueryStruct(ctx, s.log, s.db, q, data, &dbPost); err != nil {
		if errors.Is(err, db.ErrDBNotFound) {
			return post.Post{}, fmt.Errorf("namedquerystruct: %w", post.ErrNotFound)
		}
		return post.Post{}, fmt.Errorf("namedquerystruct: %w", err)
	}

	post := toCorePost(dbPost)

	return post, nil
}

func (s *Store) GetPosts(ctx context.Context) ([]post.Post, error) {
	return []post.Post{}, nil
}