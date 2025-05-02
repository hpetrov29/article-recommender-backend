package postorchestrator

import (
	"context"

	"github.com/hpetrov29/resttemplate/business/core/post"
	"github.com/hpetrov29/resttemplate/business/data/order"
	"github.com/hpetrov29/resttemplate/internal/logger"
)

// implements post.Storer, represents the hybrid store that manages all other stores.
// Methods that have to be implemented:
/*
	Create(ctx context.Context, post Post) (error)
	Delete(ctx context.Context, post Post) error
	QueryById(ctx context.Context, id int64) (Post, error)
	Query(ctx context.Context, filter QueryFilter, orderBy order.OrderBy, pageNumber int, rowsPerPage int) ([]Post, error)
*/
type Store struct {
	log 	*logger.Logger
	Cache 	post.CacheStore
	SQL   	post.SQLstore
	NOSQL 	post.NOSQLStore
}

func NewStore(log *logger.Logger, cache post.CacheStore, sql post.SQLstore, nosql post.NOSQLStore) *Store {
	return &Store{
		log: log,
		Cache: cache,
		SQL: sql,
		NOSQL: nosql,
	}
}

func (o *Store) Create(ctx context.Context, post post.Post) error {
	// insert post metadata in sql repo
	// insert post's content in nosql repo
	// add the post in cache with some TTL
	if err := o.SQL.Create(ctx, post); err != nil {
		return err
	}
	if err := o.NOSQL.Create(ctx, post.Content, post.ContentId); err != nil {
		return err
	}
	return o.Cache.CreatePost(ctx, post)
}

func (o *Store) Delete(ctx context.Context, post post.Post) error {
	// first, delete post in the sql repo
	// second, delete the post's content in the nosql repo
	// third, if present in cache, delete
	if err := o.SQL.Delete(ctx, post.Id); err != nil {
		return err
	}
	if err := o.NOSQL.Delete(ctx, post.ContentId); err != nil {
		return err
	}
	return o.Cache.DeletePost(ctx, post.Id)
}

func (o *Store) QueryById(ctx context.Context, id int64) (post.Post, error) {
	p, ok, err := o.Cache.QueryPostById(ctx, id); if err != nil {
		return post.Post{}, err
	}

	if ok {
		return p, nil
	}
	
	p, err = o.SQL.QueryById(ctx, id)
	if err != nil {
		return post.Post{}, err
	}
	content, err := o.NOSQL.QueryById(ctx, p.ContentId)
	if err != nil {
		return post.Post{}, err
	}
	p.Content = content
	if err = o.Cache.CreatePost(ctx, p); err != nil {
		return post.Post{}, err
	}

	return p, nil
}

func (o *Store) Query(ctx context.Context, filter post.QueryFilter, orderBy order.OrderBy, pageNumber int, rowsPerPage int) ([]post.Post, error) {
	//only call the sql store since the nosql store doesnt care about post metadata
	return o.SQL.Query(ctx, filter, orderBy, pageNumber, rowsPerPage)
}
