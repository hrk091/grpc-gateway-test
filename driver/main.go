package driver

import "context"

type DatastoreDriver interface {
	Get(ctx context.Context, category, key string) (map[string]interface{}, error)
	GetAll(ctx context.Context, category string, queryName string) ([]map[string]interface{}, error)
	Create(ctx context.Context, category, key string, doc interface{}) error
	Update(ctx context.Context, category, key string, doc interface{}) (map[string]interface{}, error)
	Delete(ctx context.Context, category, key string) error
}
