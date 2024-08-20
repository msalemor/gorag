package stores

import (
	"context"

	"github.com/msalemor/gorag/pkg/services"
)

type IStore interface {
	CreateTable(T any, ctx context.Context) (bool, error)
	CreateCollection(collection string, ctx context.Context) (bool, error)
	CollectionExists(collection string, ctx context.Context) bool
	AddMemory(memory services.Memory, ctx context.Context) (string, error)
	GetMemory(collection, key string, ctx context.Context) (services.Memory, error)
	GetAll(collection string, ctx context.Context) ([]services.Memory, error)
	DeleteMemory(collection, key string, ctx context.Context) (bool, error)
	Search(collection, query string, limit int, relevance float64, emb bool, ctx context.Context) ([]services.MemorySearchResult, error)
	DeleteCollection(collection string, ctx context.Context) (bool, error)
}
