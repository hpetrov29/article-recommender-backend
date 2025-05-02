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

// QueryByPostId fetches the comment tree of the corresponding post from the database.
// It fetches up to 5 root-level comments and, for each comment, up to 5 children recursively, 
// to a maximum depth of 5 levels.
//
// Parameters:
//   - ctx: the context for managing timeouts and cancellations.
//   - id: the id of the post whose comment tree will be fetched.
//
// Returns:
//   - []Comment: the comment tree in a flat list
//   - error: an error if the fetch fails.
func (s *Store) QueryByPostId(ctx context.Context, id int64) ([]comment.Comment, error) {
	data := struct {
		Id int64 `db:"id"`
	}{
		Id: id,
	}

	const q = `WITH RECURSIVE ordComments AS
							(SELECT * ,
									row_number() OVER (PARTITION BY coalesce(parent_id, 0)
														ORDER BY created_at) rn
							FROM comments),
										r AS
							(SELECT 0 AS lvl,
									id AS root,
									t.*
							FROM ordComments t
							WHERE post_id = :id 
								AND parent_id IS NULL
								AND rn<6
							UNION ALL SELECT lvl+1 AS lvl,
												r.root,
												t.*
							FROM r
							INNER JOIN ordComments t ON t.parent_id=r.id
							AND t.rn<6
							)
							SELECT 
								r.id,
								r.user_id,
								r.parent_id,
								r.content,
								r.created_at,
								r.root,
								r.lvl
							FROM r
							ORDER BY root,
									lvl`
	
	var dbComments []Comment
	if err := mysql.NamedQuerySlice(ctx, s.log, s.db, q, data, &dbComments); err != nil {
		return []comment.Comment{}, fmt.Errorf("namedquerystruct: %w", err)
	}

	return ToCoreComments(dbComments), nil
}