package pkg

import (
	"context"

	"github.com/msalemor/gorag/pkg/services"
)

type IStore interface {
	CreateTable(T any, ctx context.Context) (bool, error)
	CreateCollection(collection string, ctx context.Context) (bool, error)
	CollectionExists(collection string, ctx context.Context) bool
	AddMemory(memory Memory, ctx context.Context) (string, error)
	GetMemory(collection, key string, ctx context.Context) (Memory, error)
	GetAll(collection string, ctx context.Context) ([]Memory, error)
	DeleteMemory(collection, key string, ctx context.Context) (bool, error)
	Search(collection, query string, limit int, relevance float64, emb bool, ctx context.Context) ([]MemorySearchResult, error)
	DeleteCollection(collection string, ctx context.Context) (bool, error)
}

type IChatService interface {
	Chat(text string, temperature float64, maxTokens int, streaming bool) *services.OllamaChatResponse
}

type IEmbeddingService interface {
	Embed(text string) *[]float64
}
