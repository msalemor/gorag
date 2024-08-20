package process

import (
	"fmt"
	"net/http"
	"path"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/msalemor/gorag/pkg/services"
	"github.com/msalemor/gorag/pkg/splitters"
)

func ConfigureRoutes(chatEndpoint, embeddingEndpoint, collection, chatModel, embeddingModel string, keep, verbose bool) *gin.Engine {

	ctx, chatService, store := initServices(chatEndpoint, chatModel, embeddingEndpoint, embeddingModel, verbose)
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
		err := c.ShouldBindJSON(&request)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}

		if request.Collection == "" || len(request.URLs) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Provide a collection name and a list of URLs to ingest"})
			return
		}

		// Get the text from the URLs
		for _, url := range request.URLs {
			content, err := urlService.GetURLText(url, chatService.Client)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
				return
			}
			if content != "" {
				chunks := splitter.Split(splitters.SplitterOpts{Content: content})
				filename := path.Base(url)
				filenameWithoutExt := path.Base(filename[:len(filename)-len(path.Ext(filename))])
				for idx, chunk := range chunks {
					store.AddMemory(services.Memory{
						Collection:  request.Collection,
						Key:         fmt.Sprintf("%s-%d", filenameWithoutExt, idx),
						Text:        chunk,
						Description: &url,
					}, ctx)
				}
			}
		}

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
		result := chatService.Chat(&services.ChatOpts{Messages: request.Messages, Temperature: request.Temperature, MaxTokens: request.MaxTokens})
		c.JSON(http.StatusOK, gin.H{"content": result.Choices[0].Message.Content})
	})

	group.POST("/query", func(c *gin.Context) {
		var request services.QueryRequest
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
		result := chatService.Chat(&services.ChatOpts{Messages: messages, Temperature: 0.1, MaxTokens: 4096})
		//result := chatService.Chat(request.Messages, request.Temperature, request.MaxTokens, false)
		c.JSON(http.StatusOK, gin.H{"content": result.Choices[0].Message.Content})
	})

	return r
}
