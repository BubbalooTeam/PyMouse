package main

import (
	"fmt"
	"pymouse/pymouse/database"
)

func main() {
	usersCollection := database.NewCollection("users")
	filter := map[string]interface{}{"name": "Alice"}
	for i := 0; i <= 1000; i++ {
		usersCollection.InsertOrUpdate(filter, map[string]interface{}{"name": "Alice", "parents": "John", "city": "Sorocaba", "age": i})
	}
	usersCollection.InsertOrUpdate(map[string]interface{}{"name": "Joana"}, map[string]interface{}{"name": "Joana", "parents": "Noah", "city": "SP", "age": 40})

	UserInfo := usersCollection.FindMatches(filter)
	fmt.Println(UserInfo)
}
