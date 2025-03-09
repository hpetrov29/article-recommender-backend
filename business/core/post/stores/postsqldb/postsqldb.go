package postsqldb

import (
	"bytes"
	"context"
	"errors"
	"fmt"

	"github.com/hpetrov29/resttemplate/business/core/post"

	mysql "github.com/hpetrov29/resttemplate/business/data/dbsql/mysql"
	"github.com/hpetrov29/resttemplate/business/data/order"
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

/*
type SQLstore interface {
	Create(context.Context, Post) (error)
	Delete(context.Context, Post) error
	QueryById(context.Context, int64) (Post, error)
	Query(ctx context.Context, filter QueryFilter, orderBy order.OrderBy, pageNumber int, rowsPerPage int) ([]Post, error)
}
*/

// Create inserts a new post record into the database.
//
// Parameters:
//   - ctx: the context for managing timeouts and cancellations.
//   - post: the post contents to be stored in the database.
//
// Returns:
//   - sql.Result: the result of the SQL insert operation, containing details such as rows affected.
//   - error: an error if the insertion fails.
func (s *Store) Create(ctx context.Context, post post.Post) (error) {
	const q = `
	INSERT INTO posts
		(id, user_id, title, description, content_id, created_at, updated_at)
	VALUES
		(:id, :user_id, :title, :description, :content_id, :created_at, :updated_at);`
	
	_, err := mysql.NamedExecContext(ctx, s.log, s.db, q, toDBPost(post)); 
	
	if err != nil {
		return fmt.Errorf("namedexeccontext: %w", err)
	}

	return nil
}

// Delete removes a post from the database based on the post's Id.
//
// Parameters:
//   - ctx: context for managing timeouts and cancellations.
//   - post: the post containing the Id of the post to be deleted.
//
// Returns:
//   - error: an error if the deletion fails. If successful, returns nil.
func (s *Store) Delete(ctx context.Context, id int64) error {
	data := struct {
		PostId int64 `db:"id"`
	}{
		PostId: id,
	}

	const q = `
	DELETE FROM
		posts
	WHERE
		id = :id`

	if _, err := mysql.NamedExecContext(ctx, s.log, s.db, q, data); err != nil {
		return fmt.Errorf("namedexeccontext: %w", err)
	}

	return nil
}

func (s *Store) QueryById(ctx context.Context, id int64) (post.Post, error) {
	data := struct {
		Id int64 `db:"id"`
	}{
		Id: id,
	}

	const postQuery = `SELECT id, user_id, title, description, front_image, content_id, created_at, updated_at FROM posts WHERE id = :id;`
	const commentsQuery = `WITH RECURSIVE comment_tree AS (SELECT id, user_id, parent_id, content, created_at, 1 AS level FROM comments WHERE post_id = :id AND parent_id IS NULL UNION ALL SELECT c.id, c.user_id, c.parent_id, c.content, c.created_at, ct.level + 1 FROM comments c INNER JOIN comment_tree ct ON c.parent_id = ct.id) SELECT id, user_id, parent_id, content, created_at FROM comment_tree ORDER BY level ASC, created_at ASC;`
	
	var dbPost dbPost
	var dbComments []dbComment

	if err := mysql.NamedQueryStruct(ctx, s.log, s.db, postQuery, data, &dbPost); err != nil {
		if errors.Is(err, mysql.ErrDBNotFound) {
			return post.Post{}, fmt.Errorf("namedquerystruct: %w", post.ErrNotFound)
		}
		return post.Post{}, fmt.Errorf("namedquerystruct: %w", err)
	}

	if err := mysql.NamedQuerySlice(ctx, s.log, s.db, commentsQuery, data, &dbComments); err != nil {
		return post.Post{}, fmt.Errorf("namedquerystruct: %w", err)
	}

	post := toCorePost(dbPost)
	post.Comments = toCoreComments(dbComments)

	return post, nil
}

func (s *Store) Query(ctx context.Context, filter post.QueryFilter, orderBy order.OrderBy, pageNumber int, rowsPerPage int) ([]post.Post, error) {
	data := map[string]interface{}{
		"offset":        (pageNumber - 1) * rowsPerPage,
		"rows_per_page": rowsPerPage,
	}
	const q = `
	SELECT
	    *
	FROM
		posts`

	buf := bytes.NewBufferString(q)

	s.applyFilter(filter, data, buf)
	s.orderByClause(orderBy, buf)
	buf.WriteString(" LIMIT :rows_per_page OFFSET :offset;")

	var dbPosts []dbPost
	if err := mysql.NamedQuerySlice(ctx, s.log, s.db, buf.String(), data, &dbPosts); err != nil {
		return nil, fmt.Errorf("namedqueryslice: %w", err)
	}

	return toCorePostSlice(dbPosts), nil
}