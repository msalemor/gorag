package stores

import (
	"context"
	"encoding/json"
	"log"
	"sort"

	"github.com/msalemor/gorag/pkg"
	"github.com/msalemor/gorag/pkg/services"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type SqliteStore struct {
	Db               *gorm.DB
	EmbeddingService pkg.IEmbeddingService
	Verbose          bool
}

var isTableFound bool

func (s *SqliteStore) checkTable(ctx context.Context) {
	if !isTableFound {
		s.CreateTable(pkg.Memory{}, ctx)
		isTableFound = true
	}
}

func (s *SqliteStore) checkDB(ctx context.Context) {
	if s.Db == nil {
		db, err := gorm.Open(sqlite.Open("memories.sqlite"), &gorm.Config{})
		if err != nil {
			log.Fatalf("failed to connect database: %v", err)
		}
		s.Db = db
		s.checkTable(ctx)
	}
}

func (s *SqliteStore) CreateTable(T any, ctx context.Context) (bool, error) {
	if s.Verbose {
		log.Println("Create Table")
	}
	// Migrate the schema
	s.checkDB(ctx)
	s.Db.AutoMigrate(T)
	return true, nil
}

func (s *SqliteStore) CollectionExists(collection string, ctx context.Context) bool {
	if s.Verbose {
		log.Println("Check Collection Exists")
	}
	s.checkDB(ctx)
	var table pkg.Memory
	s.Db.WithContext(ctx).First(&table, "collection = ?", collection)
	return table.Collection != ""
}

func (s *SqliteStore) CreateCollection(collection string, ctx context.Context) (bool, error) {
	if s.Verbose {
		log.Println("Create Collection")
	}

	s.checkDB(ctx)
	s.Db.WithContext(ctx).Create(&pkg.Memory{Collection: collection})
	return true, nil
}

func (s *SqliteStore) AddMemory(memory pkg.Memory, ctx context.Context) (string, error) {
	if s.Verbose {
		log.Printf("Adding Memory: %s", memory.Key)
	}

	s.checkDB(ctx)
	// if !s.Db.WithContext(ctx).Migrator().HasTable(pkg.Memory{}) {
	// 	s.CreateTable(pkg.Memory{}, ctx)
	// }
	// if !s.CollectionExists(memory.Collection, ctx) {
	// 	s.CreateCollection(memory.Collection, ctx)
	// }

	var result pkg.Memory
	// Search collection and key
	s.Db.WithContext(ctx).First(&result, "collection = ? AND key = ?", memory.Collection, memory.Key)

	v := s.EmbeddingService.Embed(memory.Text)
	if v == nil {
		return "", nil
	}

	// Serialize the embedding
	json, _ := json.Marshal(*v)

	memory.Embedding = string(json)

	// If the memory already exists, update it
	if result.Collection != "" {
		s.Db.WithContext(ctx).Model(&pkg.Memory{}).Where("collection = ? AND key = ?", memory.Collection, memory.Key).Updates(memory)
	} else {
		// If the memory does not exist, create it
		s.Db.WithContext(ctx).Create(&memory)
	}
	return memory.Key, nil
}

func (s *SqliteStore) GetMemory(collection, key string, ctx context.Context) (pkg.Memory, error) {
	if s.Verbose {
		log.Println("Get Memory")
	}

	s.checkDB(ctx)
	var result pkg.Memory
	// Search collection and key
	s.Db.WithContext(ctx).First(&result, "collection = ? AND key = ?", collection, key)
	if result.Collection == "" {
		return pkg.Memory{}, nil
	}
	return result, nil
}

func (s *SqliteStore) GetAll(collection string, ctx context.Context) ([]pkg.Memory, error) {
	if s.Verbose {
		log.Println("Get All")
	}

	s.checkDB(ctx)
	var result []pkg.Memory
	s.Db.WithContext(ctx).Find(&result, "collection = ? AND key <> ''", collection)
	return result, nil
}

func (s *SqliteStore) DeleteCollection(collection string, ctx context.Context) (bool, error) {
	if s.Verbose {
		log.Println("Delete Collection")
	}

	s.checkDB(ctx)
	if s.CollectionExists(collection, ctx) {
		s.Db.WithContext(ctx).Where("collection = ?", collection).Delete(&pkg.Memory{})
		return true, nil
	}
	return false, nil
}

func (s *SqliteStore) DeleteMemory(collection, key string, ctx context.Context) (bool, error) {
	if s.Verbose {
		log.Println("Delete Memory")
	}

	s.checkDB(ctx)
	test, _ := s.GetMemory(collection, key, ctx)
	if test.Collection != "" {
		s.Db.WithContext(ctx).Where("collection = ? AND key = ?", collection, key).Delete(&pkg.Memory{})
		return true, nil
	}
	return false, s.Db.Error
}

func (s *SqliteStore) Search(collection, query string, limit int, relevance float64, emb bool, ctx context.Context) ([]pkg.MemorySearchResult, error) {
	if s.Verbose {
		log.Println("Search")
	}

	s.checkDB(ctx)
	records, _ := s.GetAll(collection, ctx)
	results := []pkg.MemorySearchResult{}
	for _, record := range records {
		// Deserialize the embedding
		var embedding []float64
		_ = json.Unmarshal([]byte(record.Embedding), &embedding)
		// Calculate the cosine similarity
		v := s.EmbeddingService.Embed(query)
		similarity := services.CosineSimilarity(embedding, *v)
		// If the similarity is greater than the relevance, add it to the result
		if similarity > relevance {
			memoryResult := pkg.MemorySearchResult{
				Collection:  record.Collection,
				Key:         record.Key,
				Text:        record.Text,
				Description: record.Description,
				Relevance:   similarity}
			if emb {
				var floatList *[]float64
				json.Unmarshal([]byte(record.Embedding), floatList)
				memoryResult.Embedding = floatList
			}
			results = append(results, memoryResult)
		}
	}
	// Sort by relevance
	sort.Slice(results, func(i, j int) bool {
		return results[i].Relevance > results[j].Relevance
	})

	// Limit the number of results
	if len(results) > limit {
		results = results[:limit]
	}

	return results, nil
}
