# gorag

## 1.0 - gorag

### 1.1 - Overview 

A simple Golang RAG package to store and recall vectors and text chunks from SQLite inspired grately on Semantic Kernel memories.

### 1.2 - Using gorag as a CLI

gorag is both a sample CLI and a package. The CLI is designed to showcase using gorag to store and recall vectors from a SQLite database. The CLI comes with two commands:

- `gorag console`
- `gorag ui`

### 1.3 - Using gorag as a package


#### Code

##### 1.3.1 - RAG


```go
package main

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/msalemor/gorag/pkg"
	"github.com/msalemor/gorag/pkg/services"
	"github.com/msalemor/gorag/pkg/stores"
)

var (
	chatendpoint      = "http://localhost:11434/api/chat"
	chatmodel         = "llama3"
	embeddingendpoint = "http://localhost:11434/api/embeddings"
	embeddingmodel    = "nomic-embed-text"
	collection        = "FAQ"
	verbose           = false
)

func main() {

	ctx := context.Background()
	client := &http.Client{}

	chatService := &services.OllamaChatService{
		Endpoint: chatendpoint,
		Model:    chatmodel,
		Client:   client,
	}

	embeddingService := &services.OllamaEmbeddingService{
		Endpoint: embeddingendpoint,
		Model:    embeddingmodel,
		Client:   client,
	}

	store := &stores.SqliteStore{
		EmbeddingService: embeddingService,
		Verbose:          verbose,
	}

	// Cleanup
	store.DeleteCollection(collection, ctx)

	// Ingest FAQ - Simple splitter based on paragraphs
	chunks := strings.Split(pkg.FAQ, "\n\n")
	for idx, chunk := range chunks {
		store.AddMemory(pkg.Memory{
			Collection: collection,
			Key:        fmt.Sprintf("faq-%d", idx),
			Text:       chunk,
		}, ctx)
	}

	// user question
	question := "What is the return policy?"

	// Search the vector database and find the nearest neighbors
	sb := strings.Builder{}
	results, _ := store.Search(collection, question, 3, 0.75, true, ctx)
	for _, result := range results {
		sb.WriteString(fmt.Sprintf("%s\n", result.Text))
	}

	// Augment the prompt
	augmentedPrompt := question + "\n" + sb.String()
	messages := []services.Message{
		{Role: "system", Content: "You are an AI Assistant for eChampShop an online shopping store for exercise equipment and provides maintenance and repair services. You can answer questions based on the context that is provided. If no context is provided, say I don't know."},
		{Role: "user", Content: augmentedPrompt},
	}

	// Process the completion
	completion := chatService.Chat(messages, 0.1, 4096, false)
	fmt.Printf("user:\n%s\n", question)
	fmt.Printf("assistant:\n%s\n", completion.Message.Content)
}
```