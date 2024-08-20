package process

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/msalemor/gorag/pkg"
	"github.com/msalemor/gorag/pkg/services"
	"github.com/msalemor/gorag/pkg/stores"
)

func initServices(chatendpoint, chatmodel, embeddingendpoint, embeddingmodel string, verbose bool) (context.Context, *services.OllamaOpenAIChatService, *stores.SqliteStore) {
	client := &http.Client{}
	ctx := context.Background()

	chatService := &services.OllamaOpenAIChatService{
		Endpoint: chatendpoint,
		Model:    chatmodel,
		Client:   client,
	}

	embeddingService := &services.OllamaOpenAIEmbeddingService{
		Endpoint: embeddingendpoint,
		Model:    embeddingmodel,
		Client:   client,
	}

	store := &stores.SqliteStore{
		EmbeddingService: embeddingService,
		Verbose:          verbose,
	}
	return ctx, chatService, store
}

func ingestFAQ(store pkg.IStore, collection string, keep bool, verbose bool, ctx context.Context) {

	if !keep {
		store.DeleteCollection(collection, ctx)
	}

	if verbose {
		fmt.Printf("Ingesting FAQ into collection %s\n", collection)
	}

	// Ingest FAQ - Simple splitter based on paragraphs
	chunks := strings.Split(pkg.FAQ, "\n\n")
	for idx, chunk := range chunks {
		store.AddMemory(pkg.Memory{
			Collection: collection,
			Key:        fmt.Sprintf("faq-%d", idx),
			Text:       chunk,
		}, ctx)
	}
}
