package postnosqldb

import (
	"context"

	"github.com/hpetrov29/resttemplate/business/core/post"
	"github.com/hpetrov29/resttemplate/business/data/dbnosql"
	"github.com/hpetrov29/resttemplate/internal/logger"
)

// Store manages the set of APIs for posts database access.
type Store struct {
	log    		*logger.Logger
	NOSQLstore 	dbnosql.NOSQLDBrepo
}

func NewStore (log *logger.Logger, nosqlStore dbnosql.NOSQLDBrepo) *Store {
	return &Store{
		log:log,
		NOSQLstore: nosqlStore,
	}
}

func (s *Store) Create(ctx context.Context, content post.Content, contentId int64) (error) {
	err := s.NOSQLstore.Insert(ctx, toDbContent(content, contentId)); if err != nil {
		return err
	}

	return nil
}

func (s *Store) Delete(ctx context.Context, id int64) error {
	return s.NOSQLstore.Delete(ctx, uint64(id))
}

func (s *Store) QueryById(ctx context.Context, id int64) (post.Content, error) {
	var content Content

	err := s.NOSQLstore.QueryById(ctx, id, &content); if err != nil {
		return post.Content{}, err
	}

	return toCoreContent(content), nil
}
