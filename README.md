# gorag

## 1.0 - gorag

### 1.1 - Overview 

A simple Golang RAG package to store and recall vectors and text chunks from SQLite inspired grately on Semantic Kernel memories.

### 1.2 - Using gorag as a CLI

gorag is both a sample CLI and a package. The CLI is designed to showcase using gorag to store and recall vectors from a SQLite database. The CLI comes with two commands:

- `gorag console`
- `gorag ui`

### 1.3 - Using gorag as a package

```go
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
```