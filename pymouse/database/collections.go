package database

import (
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoCollection struct {
	Collection *mongo.Collection
}

// NewMongoCollection creates a new instance of MongoCollection
func NewMongoCollection(collName string) *MongoCollection {
	return &MongoCollection{
		Collection: Database.Collection(collName),
	}
}

// UpdateOne updates a document
func (mc *MongoCollection) UpdateOne(filter bson.M, data interface{}) (err error) {
	_, err = mc.Collection.UpdateOne(tdContext, filter, bson.M{"$set": data}, options.Update().SetUpsert(true))
	if err != nil {
		log.Printf("[Database][UpdateOne][Error]: %v", err)
	}
	return
}

// FindOne finds a document
func (mc *MongoCollection) FindOne(filter bson.M) (res *mongo.SingleResult) {
	res = mc.Collection.FindOne(tdContext, filter)
	return
}

// CountDocs counts documents
func (mc *MongoCollection) CountDocs(filter bson.M) (count int64, err error) {
	count, err = mc.Collection.CountDocuments(tdContext, filter)
	if err != nil {
		log.Printf("[Database][CountDocs][Error]: %v", err)
	}
	return
}

// FindAll finds all documents
func (mc *MongoCollection) FindAll(filter bson.M) (cur *mongo.Cursor, err error) {
	cur, err = mc.Collection.Find(tdContext, filter)
	if err != nil {
		log.Printf("[Database][FindAll][Error]: %v", err)
	}
	return
}

// DeleteOne deletes a document
func (mc *MongoCollection) DeleteOne(filter bson.M) (err error) {
	_, err = mc.Collection.DeleteOne(tdContext, filter)
	if err != nil {
		log.Printf("[Database][DeleteOne][Error]: %v", err)
	}
	return
}

// DeleteMany deletes multiple documents
func (mc *MongoCollection) DeleteMany(filter bson.M) (err error) {
	_, err = mc.Collection.DeleteMany(tdContext, filter)
	if err != nil {
		log.Printf("[Database][DeleteMany][Error]: %v", err)
	}
	return
}
