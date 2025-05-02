package postcache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/hpetrov29/resttemplate/business/core/post"
	"github.com/hpetrov29/resttemplate/business/data/cache"
	"github.com/hpetrov29/resttemplate/internal/logger"
)

// Store manages the set of APIs for posts database access.
type Store struct {
	log    			*logger.Logger
	CacheStore 		cache.Cache
}

func NewStore (log *logger.Logger, cache cache.Cache) *Store {
	return &Store{
		log:log, 
		CacheStore: cache,
	}
}

func (s *Store) CreatePost(ctx context.Context, post post.Post) (error) {
	data, _ := json.Marshal(toDBPost(post))
	return s.CacheStore.SetWithTTL(ctx, fmt.Sprintf("posts:%d", post.Id), data, 5*time.Second)
}

func (s *Store) DeletePost(ctx context.Context, id int64) error {
	return nil
}

func (s *Store) QueryPostById(ctx context.Context, id int64) (post.Post, bool, error) {
	var postData dbPost // change to db post

	data, ok, err := s.CacheStore.GetNonFatal(ctx, fmt.Sprintf("posts:%d", id))
	if err != nil {
		return post.Post{}, false, err
	}
	if !ok {
		return post.Post{}, false, nil
	}

	if err = json.Unmarshal(data, &postData); err != nil {
		return post.Post{}, false, err
	}
	return toCorePost(postData), true, nil
}