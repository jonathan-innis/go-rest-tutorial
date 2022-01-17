package database

import "context"

type Interface interface {
	Create(ctx context.Context, item interface{}) error
	List(ctx context.Context, model interface{}) ([]interface{}, error)
	// Get(ctx context.Context, model interface{}, id string) (interface{}, error)
}
