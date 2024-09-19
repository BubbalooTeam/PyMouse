package main

import (
	"pymouse/pymouse/database"
)

func main() {
	usersCollection := database.NewCollection("users")
	filter := map[string]interface{}{"name": "Alice"}
	for i := 0; i <= 1000; i++ {
		usersCollection.InsertOrUpdate(filter, map[string]interface{}{"name": "Alice", "parents": "John", "city": "Sorocaba", "age": i})
	}
}
