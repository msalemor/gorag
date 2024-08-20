package stores

import (
	"context"
	"encoding/json"
	"sort"

	"github.com/msalemor/gorag/pkg/services"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type SqliteStore struct {
	Db               *gorm.DB
	EmbeddingService services.IEmbeddingService
	Verbose          bool
}

var isTableFound bool

func (s *SqliteStore) checkTable(ctx context.Context) {
	if !isTableFound {
		s.CreateTable(services.Memory{}, ctx)
		isTableFound = true
	}
}

func (s *SqliteStore) checkDB(ctx context.Context) {
	if s.Db == nil {
		db, err := gorm.Open(sqlite.Open("memories.sqlite"), &gorm.Config{})
		if err != nil {
			logrus.Fatalf("failed to connect database: %v", err)
		}
		s.Db = db
		s.checkTable(ctx)
	}
}

func (s *SqliteStore) CreateTable(T any, ctx context.Context) (bool, error) {
	if s.Verbose {
		logrus.Info("Creating the memories table")
	}
	// Migrate the schema
	s.checkDB(ctx)
	s.Db.AutoMigrate(T)
	return true, nil
}

func (s *SqliteStore) CollectionExists(collection string, ctx context.Context) bool {
	if s.Verbose {
		logrus.Infof("Checking to see if the collection Exists: %s", collection)
	}
	s.checkDB(ctx)
	var table services.Memory
	s.Db.WithContext(ctx).First(&table, "collection = ?", collection)
	return table.Collection != ""
}

func (s *SqliteStore) CreateCollection(collection string, ctx context.Context) (bool, error) {
	if s.Verbose {
		logrus.Infof("Creating the collection: %s", collection)
	}

	s.checkDB(ctx)
	s.Db.WithContext(ctx).Create(&services.Memory{Collection: collection})
	return true, nil
}

func (s *SqliteStore) AddMemory(memory services.Memory, ctx context.Context) (string, error) {
	if s.Verbose {
		logrus.Infof("Adding Memory: (%s,%s)", memory.Collection, memory.Key)
	}

	s.checkDB(ctx)
	// if !s.Db.WithContext(ctx).Migrator().HasTable(pkg.Memory{}) {
	// 	s.CreateTable(pkg.Memory{}, ctx)
	// }
	// if !s.CollectionExists(memory.Collection, ctx) {
	// 	s.CreateCollection(memory.Collection, ctx)
	// }

	var result services.Memory
	// Search collection and key
	s.Db.WithContext(ctx).First(&result, "collection = ? AND key = ?", memory.Collection, memory.Key)

	v := s.EmbeddingService.Embed(&services.EmbeddingOpts{Text: memory.Text})
	if v == nil {
		return "", nil
	}

	// Serialize the embedding
	json, _ := json.Marshal(*v)

	memory.Embedding = string(json)

	// If the memory already exists, update it
	if result.Collection != "" {
		s.Db.WithContext(ctx).Model(&services.Memory{}).Where("collection = ? AND key = ?", memory.Collection, memory.Key).Updates(memory)
	} else {
		// If the memory does not exist, create it
		s.Db.WithContext(ctx).Create(&memory)
	}
	return memory.Key, nil
}

func (s *SqliteStore) GetMemory(collection, key string, ctx context.Context) (services.Memory, error) {
	if s.Verbose {
		logrus.Infof("Geting the memory: (%s,%s)", collection, key)
	}

	s.checkDB(ctx)
	var result services.Memory
	// Search collection and key
	s.Db.WithContext(ctx).First(&result, "collection = ? AND key = ?", collection, key)
	if result.Collection == "" {
		return services.Memory{}, nil
	}
	return result, nil
}

func (s *SqliteStore) GetAll(collection string, ctx context.Context) ([]services.Memory, error) {
	if s.Verbose {
		logrus.Infof("Getting all the memories from collection: %s", collection)
	}

	s.checkDB(ctx)
	var result []services.Memory
	s.Db.WithContext(ctx).Find(&result, "collection = ? AND key <> ''", collection)
	return result, nil
}

func (s *SqliteStore) DeleteCollection(collection string, ctx context.Context) (bool, error) {
	if s.Verbose {
		logrus.Infof("Delete collection: %s", collection)
	}

	s.checkDB(ctx)
	if s.CollectionExists(collection, ctx) {
		s.Db.WithContext(ctx).Where("collection = ?", collection).Delete(&services.Memory{})
		return true, nil
	}
	return false, nil
}

func (s *SqliteStore) DeleteMemory(collection, key string, ctx context.Context) (bool, error) {
	if s.Verbose {
		logrus.Infof("Delete memory: (%s,%s)", collection, key)
	}

	s.checkDB(ctx)
	test, _ := s.GetMemory(collection, key, ctx)
	if test.Collection != "" {
		s.Db.WithContext(ctx).Where("collection = ? AND key = ?", collection, key).Delete(&services.Memory{})
		return true, nil
	}
	return false, s.Db.Error
}

func (s *SqliteStore) Search(collection, query string, limit int, relevance float64, emb bool, ctx context.Context) ([]services.MemorySearchResult, error) {
	if s.Verbose {
		logrus.Infof("Searching memories in collection: %s", collection)
	}

	s.checkDB(ctx)
	records, _ := s.GetAll(collection, ctx)
	results := []services.MemorySearchResult{}
	for _, record := range records {
		// Deserialize the embedding
		var embedding []float64
		_ = json.Unmarshal([]byte(record.Embedding), &embedding)
		// Calculate the cosine similarity
		v := s.EmbeddingService.Embed(&services.EmbeddingOpts{Text: query})
		if v == nil {
			logrus.Error("Unable to get embedding")
		} else {
			similarity := services.CosineSimilarity(embedding, *v)
			// If the similarity is greater than the relevance, add it to the result
			if similarity > relevance {
				memoryResult := services.MemorySearchResult{
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
	}

	if len(results) > 0 {
		// Sort by relevance
		if s.Verbose {
			logrus.Infof("Sorting the memory results by relevance: %f", relevance)
		}
		sort.Slice(results, func(i, j int) bool {
			return results[i].Relevance > results[j].Relevance
		})

		// Limit the number of results
		if s.Verbose {
			logrus.Infof("Limiting the number of results from: %v to %v", len(results), limit)
		}
		if len(results) > limit {
			results = results[:limit]
		}
	}

	return results, nil
}
