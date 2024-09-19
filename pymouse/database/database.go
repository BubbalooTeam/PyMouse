package database

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

var dbFile = "pymouse/database/files/database.json"
var db = make(map[string][]map[string]interface{})

func LoadDB() {
	data, err := os.ReadFile(dbFile)
	if err != nil {
		if os.IsNotExist(err) {
			err = os.MkdirAll("pymouse/database/files", os.ModePerm)
			if err != nil {
				fmt.Printf("Error in Making a Directory: %v\n", err)
				return
			}
			return
		}
		for attempts := 0; attempts < 100; attempts++ {
			time.Sleep(1 * time.Second)
			data, err = os.ReadFile(dbFile)
			if err == nil {
				break
			}
			fmt.Printf("Error loading database... Attempt %d/100\n", attempts+1)
		}
		if err != nil {
			fmt.Printf("Failed to load database: %v\n", err)
			return
		}
	}
	json.Unmarshal(data, &db)
}

func SaveDB() {
	data, err := json.MarshalIndent(db, "", "    ")
	if err != nil {
		fmt.Printf("Error marshalling JSON database: %v\n", err)
		return
	}

	err = os.MkdirAll("pymouse/database/files", os.ModePerm)
	if err != nil {
		fmt.Printf("Error in making a database directory: %v\n", err)
		return
	}

	err = os.WriteFile(dbFile, data, 0644)
	if err != nil {
		fmt.Printf("Error saving database: %v\n", err)
	}
}

// Collection represents a database collection.
type Collection struct {
	name string
}

// NewCollection creates a new Collection and loads the database.
func NewCollection(name string) *Collection {
	LoadDB() // Carregar a base de dados ao criar a coleção
	return &Collection{name: name}
}

// FindOne searches for a single item in the collection that matches the filter.
func (c *Collection) FindOne(filter map[string]interface{}) map[string]interface{} {
	collectionData := db[c.name]
	if len(filter) == 0 {
		fmt.Println("[database/modules]: No filter provided for findOne.")
		return nil
	}
	for _, item := range collectionData {
		matches := true
		for k, v := range filter {
			if item[k] != v {
				matches = false
				break
			}
		}
		if matches {
			return item
		}
	}
	return nil
}

// InsertOrUpdate inserts a new item or updates an existing one based on the filter.
func (c *Collection) InsertOrUpdate(filter, info map[string]interface{}) bool {
	if len(info) == 0 {
		fmt.Println("[database/modules]: No information provided for insertOrUpdate.")
		return false
	}
	collectionData := db[c.name]
	if filter != nil {
		for i, item := range collectionData {
			matches := true
			for k, v := range filter {
				if item[k] != v {
					matches = false
					break
				}
			}
			if matches {
				for k, v := range info {
					item[k] = v
				}
				collectionData[i] = item
				SaveDB() // Save after update
				return true
			}
		}
	}
	db[c.name] = append(collectionData, info)
	SaveDB() // Save after insert
	return true
}

// Delete removes items from the collection based on the filter.
func (c *Collection) Delete(filter map[string]interface{}) bool {
	collectionData := db[c.name]
	deleted := false
	if filter != nil {
		for i := 0; i < len(collectionData); i++ {
			matches := true
			for k, v := range filter {
				if collectionData[i][k] != v {
					matches = false
					break
				}
			}
			if matches {
				collectionData = append(collectionData[:i], collectionData[i+1:]...)
				deleted = true
				i-- // adjust index after removal
			}
		}
	} else {
		db[c.name] = []map[string]interface{}{}
		deleted = true
	}
	SaveDB() // Save after delete
	return deleted
}
