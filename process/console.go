package process

import (
	"fmt"
	"strings"

	"github.com/msalemor/gorag/pkg/services"
)

func ProcessConsole(chatendpoint, embeddingendpoint, collection, chatmodel, embeddingmodel string, keep bool, verbose bool) {

	ctx, chatService, store := initServices(chatendpoint, chatmodel, embeddingendpoint, embeddingmodel, verbose)

	ingestFAQ(store, collection, keep, verbose, ctx)

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
	fmt.Printf("assistant:\n%s\n", completion.Choices[0].Message.Content)
}
