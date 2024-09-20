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

type Collection struct {
	name string
}

func NewCollection(name string) *Collection {
	LoadDB()
	return &Collection{name: name}
}

func (c *Collection) FindMatches(filter map[string]interface{}) []map[string]interface{} {
	collectionData := db[c.name]
	var results []map[string]interface{}

	if len(filter) == 0 {
		return collectionData
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
			results = append(results, item)
		}
	}
	return results
}

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
				SaveDB()
				return true
			}
		}
	}
	db[c.name] = append(collectionData, info)
	SaveDB()
	return true
}

func (c *Collection) DeleteMatches(filter map[string]interface{}) bool {
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
				i--
			}
		}
	} else {
		db[c.name] = []map[string]interface{}{}
		deleted = true
	}
	SaveDB()
	return deleted
}
