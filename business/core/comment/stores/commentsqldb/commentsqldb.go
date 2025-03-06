package commentsqldb

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/hpetrov29/resttemplate/business/core/comment"
	"github.com/hpetrov29/resttemplate/business/data/dbsql/mysql"
	"github.com/hpetrov29/resttemplate/internal/logger"
	"github.com/jmoiron/sqlx"
)

// Store manages the set of APIs for comment database access.
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

// Create inserts a new comment record into the database.
//
// Parameters:
//   - ctx: the context for managing timeouts and cancellations.
//   - comment: the comment to be stored in the database.
//
// Returns:
//   - sql.Result: the result of the SQL insert operation, containing details such as rows affected.
//   - error: an error if the insertion fails.
func (s *Store) Create(ctx context.Context, comment comment.Comment) (sql.Result, error) {
	const q = `
	INSERT INTO comments
		(id, user_id, post_id, parent_id, content, created_at)
	VALUES
		(:id, :user_id, :post_id, :parent_id, :content, :created_at);`

	res, err := mysql.NamedExecContext(ctx, s.log, s.db, q, toDBComment(comment)); 
	if err != nil {
		return nil, fmt.Errorf("namedexeccontext: %w", err)
	}

	return res, nil
}

// Delete deletes a comment record from the database.
//
// Parameters:
//   - ctx: the context for managing timeouts and cancellations.
//   - id: the id of the comment to be deleted from the database.
//
// Returns:
//   - sql.Result: the result of the SQL insert operation, containing details such as rows affected.
//   - error: an error if the deletion fails.
func (s *Store) Delete(ctx context.Context, id uint64) (error) {
	data := struct {
		Id uint64 `db:"id"`
	}{
		Id: id,
	}

	const q = `DELETE FROM comments WHERE id = :id;`
	
	res, err := mysql.NamedExecContext(ctx, s.log, s.db, q, data)
	if err != nil {
		return fmt.Errorf("namedexeccontext: %w", err)
	}

	if ra, _ := res.RowsAffected(); ra == 0 {
		return comment.ErrNotFound
	}

	return nil
}