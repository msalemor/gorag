package pkg

import "github.com/msalemor/gorag/pkg/services"

type Memory struct {
	Collection         string `gorm:"primaryKey"`
	Key                string `gorm:"primaryKey"`
	Text               string
	Description        *string
	AdditionalMetadata *string
	Embedding          string
}

type MemorySearchResult struct {
	Collection         string `gorm:"primaryKey"`
	Key                string `gorm:"primaryKey"`
	Text               string
	Description        *string
	AdditionalMetadata *string
	Embedding          *[]float64
	Relevance          float64
}

type QueryRequest struct {
	Collection string              `json:"collection"`
	Query      string              `json:"query"`
	Limit      int                 `json:"limit"`
	Relevance  float64             `json:"relevance"`
	Messages   *[]services.Message `json:"messages"`
}
