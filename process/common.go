package process

import (
	"context"
	"fmt"
	"strings"

	"github.com/msalemor/gorag/pkg"
)

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
