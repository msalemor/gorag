package process

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/msalemor/gorag/pkg"
	"github.com/msalemor/gorag/pkg/services"
	"github.com/msalemor/gorag/pkg/stores"
)

func ConfigureRoutes(chatendpoint, embeddingendpoint, collection, model string, keep, verbose bool) *gin.Engine {

	client := &http.Client{}
	ctx := context.Background()

	chatService := &services.OllamaOpenAIChatService{
		Endpoint: chatendpoint,
		Model:    "llama3",
		Client:   client,
	}

	embeddingService := &services.OllamaOpenAIEmbeddingService{
		Endpoint: embeddingendpoint,
		Model:    "nomic-embed-text",
		Client:   client,
	}

	store := &stores.SqliteStore{
		EmbeddingService: embeddingService,
		Verbose:          verbose,
	}

	ingestFAQ(store, collection, keep, verbose, ctx)

	r := gin.Default()
	r.Use(gin.Logger())
	r.Use(cors.Default())
	r.Use(static.Serve("/", static.LocalFile("static/", false)))

	group := r.Group("/api")
	group.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	group.POST("/ingest", func(c *gin.Context) {
		var request struct {
			Collection string   `json:"collection"`
			URLs       []string `json:"urls"`
		}
		c.ShouldBindJSON(&request)

		if request.Collection == "" || len(request.URLs) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Provide a collection name and a list of URLs to ingest"})
			return
		}

		mem1 := pkg.Memory{
			Collection: "faq",
			Key:        "faq-1",
			Text:       "What is the return policy?",
		}
		store.AddMemory(mem1, ctx)

		c.JSON(http.StatusOK, gin.H{"message": "ingested", "urls": request.URLs})

	})

	group.POST("/chat", func(c *gin.Context) {
		var request struct {
			Collection  string             `json:"collection"`
			Messages    []services.Message `json:"messages"`
			MaxTokens   int                `json:"max_tokens"`
			Temperature float64            `json:"temperature"`
		}
		err := c.ShouldBindJSON(&request)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(request)
		if request.Collection == "" || len(request.Messages) == 0 {
			c.JSON(http.StatusNotFound, gin.H{"message": "Provide a collection name and a list messages"})
			return
		}
		result := chatService.Chat(request.Messages, request.Temperature, request.MaxTokens, false)
		c.JSON(http.StatusOK, gin.H{"content": result.Choices[0].Message.Content})
	})

	group.POST("/query", func(c *gin.Context) {
		var request pkg.QueryRequest
		err := c.ShouldBindJSON(&request)
		if err != nil {
			fmt.Println(err)
		}

		if request.Collection == "" || request.Query == "" {
			c.JSON(http.StatusNotFound, gin.H{"message": "Provide a collection name and a query"})
			return
		}
		// Search the vector database and find the nearest neighbors
		sb := strings.Builder{}
		results, _ := store.Search(request.Collection, request.Query, request.Limit, request.Relevance, true, ctx)
		for _, result := range results {
			sb.WriteString(fmt.Sprintf("%s\n", result.Text))
		}

		// Augment the prompt
		augmentedPrompt := request.Query + "\n" + sb.String()

		// Process the completion
		msg := services.Message{
			Role: "user", Content: augmentedPrompt,
		}
		var messages []services.Message
		if request.Messages != nil {
			messages = *request.Messages
		}
		messages = append(messages, msg)

		// Process the completion
		result := chatService.Chat(messages, 0.1, 4096, false)
		//result := chatService.Chat(request.Messages, request.Temperature, request.MaxTokens, false)
		c.JSON(http.StatusOK, gin.H{"content": result.Choices[0].Message.Content})
	})

	return r
}
