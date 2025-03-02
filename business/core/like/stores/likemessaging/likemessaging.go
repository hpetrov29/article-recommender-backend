package likemessaging

import (
	"context"

	"github.com/hpetrov29/resttemplate/business/core/like"
	"github.com/hpetrov29/resttemplate/business/data/messaging"
	"github.com/hpetrov29/resttemplate/internal/logger"
)

// Store manages the set of APIs for posts database access.
type Store struct {
	log    			*logger.Logger
	MessagingQueue 	messaging.MessagingQueue
	Subject 		string
}

func NewStore (log *logger.Logger, mq messaging.MessagingQueue, subject string) *Store {
	return &Store{
		log:log, 
		MessagingQueue: mq,
		Subject: subject,
	}
}

func (s *Store) Publish(ctx context.Context, l like.Like) error {
	data, err := toBytes(toDbLike(l))
	if err != nil {
		return err
	}

	return s.MessagingQueue.Publish(s.Subject, data)
}