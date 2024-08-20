package services

type ChatOpts struct {
	Messages    []Message
	Temperature float64
	MaxTokens   int
	Stream      bool
}

type EmbeddingOpts struct {
	Text string
}

type IChatService interface {
	Chat(opts *ChatOpts) *OllamaChatResponse
}

type IEmbeddingService interface {
	Embed(opts *EmbeddingOpts) *[]float64
}

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
	Collection string     `json:"collection"`
	Query      string     `json:"query"`
	Limit      int        `json:"limit"`
	Relevance  float64    `json:"relevance"`
	Messages   *[]Message `json:"messages"`
}
