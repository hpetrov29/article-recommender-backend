package postnosqldb

import (
	"context"
	"database/sql"

	"github.com/hpetrov29/resttemplate/business/core/post"
	"github.com/hpetrov29/resttemplate/business/data/dbnosql"
	"github.com/hpetrov29/resttemplate/business/data/order"
	"github.com/hpetrov29/resttemplate/internal/logger"
)

// Store manages the set of APIs for posts database access.
type Store struct {
	log    		*logger.Logger
	SQLstore 	post.Storer
	NOSQLstore 	dbnosql.NOSQLDBrepo
}

func NewStore (log *logger.Logger, sqlStore post.Storer, nosqlStore dbnosql.NOSQLDBrepo) *Store {
	return &Store{
		log:log, 
		SQLstore: sqlStore, 
		NOSQLstore: nosqlStore,
	}
}

func (s *Store) Create(ctx context.Context, post post.Post) (sql.Result, error) {
	res, err := s.SQLstore.Create(ctx, post); if err != nil {
		return nil, err
	}

	err = s.NOSQLstore.Insert(ctx, toDbContent(post)); if err != nil {
		return nil, err
	}

	return res, nil
}

func (s *Store) Delete(ctx context.Context, post post.Post) error {
	if err := s.NOSQLstore.Delete(ctx, 0); err != nil {
		return err
	}
	
	return s.SQLstore.Delete(ctx, post)
}

func (s *Store) QueryById(ctx context.Context, id uint64) (post.Post, error) {
	res, err := s.SQLstore.QueryById(ctx, id); if err != nil {
		return post.Post{}, nil
	}
	
	var content Content
	err = s.NOSQLstore.QueryById(ctx, res.ContentId, &content); if err != nil {
		return post.Post{}, err
	}
	
	res.Content = toCoreContent(content)

	return res, nil
}

func (s *Store) Query(ctx context.Context, filter post.QueryFilter, orderBy order.OrderBy, pageNumber int, rowsPerPage int) ([]post.Post, error) {
	return s.SQLstore.Query(ctx, filter, orderBy, pageNumber, rowsPerPage)
}
