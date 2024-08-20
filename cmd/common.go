package cmd

const (
	DEFAULT_COLLECTION = "FAQ"
)

var (
	Collection              = "FAQ"
	OllamaChatEndpoint      = "http://localhost:11434/v1/chat/completions"
	ChatModel               = "llama3"
	OllamaEmbeddingEndpoint = "http://localhost:11434/v1/embeddings"
	EmbeddingModel          = "nomic-embed-text"
	Verbose                 = false
	Keep                    = false
)
