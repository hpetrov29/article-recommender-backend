package postnosqldb

import (
	"context"
	"database/sql"

	"github.com/hpetrov29/resttemplate/business/core/post"
	"github.com/hpetrov29/resttemplate/business/data/order"
	"github.com/hpetrov29/resttemplate/internal/logger"
)

type NoSqlRepository interface {
	Insert() error
}

// Store manages the set of APIs for user database access.
type Store struct {
	log    		*logger.Logger
	SQLstore 	post.Storer
	NOSQLstore 	NoSqlRepository
}

func NewStore (log *logger.Logger, sqlStore post.Storer, nosqlStore NoSqlRepository) *Store {
	return &Store{
		log:log, 
		SQLstore: sqlStore, 
		NOSQLstore: nosqlStore,
	}
}

func (s *Store) Create(ctx context.Context, post post.Post) (sql.Result, error) {
	return s.SQLstore.Create(ctx, post)
}

func (s *Store) Delete(ctx context.Context, post post.Post) error {
	return s.SQLstore.Delete(ctx, post)
}

func (s *Store) QueryById(ctx context.Context, id string) (post.Post, error) {
	return s.SQLstore.QueryById(ctx, id)
}

func (s *Store) Query(ctx context.Context, filter post.QueryFilter, orderBy order.OrderBy, pageNumber int, rowsPerPage int) ([]post.Post, error) {
	return s.SQLstore.Query(ctx, filter, orderBy, pageNumber, rowsPerPage)
}
