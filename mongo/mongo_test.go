package mongo

import (
	"context"
	"testing"

	"go.mongodb.org/mongo-driver/bson"
)

func TestMongoConnection(t *testing.T) {
	client := GetMongo()
	if client == nil {
		t.Fatal("Mongo client is nil")
	}

	// Create a test database and collection
	db := client.Database("testdb")
	collection := db.Collection("testcollection")

	// Insert a test document
	doc := MongoDoc{
		ID:           "test_id",
		IdCourse:     "test_course",
		IdObjective:  "test_objective",
		IdMaterial:   "test_material",
		Transcript:   "test_transcript",
		MaterialType: "test_material_type",
		IsSuccessful: true,
	}
	_, err := collection.InsertOne(context.TODO(), doc)
	if err != nil {
		t.Errorf("Error inserting document: %v", err)
	}

	// Verify that the document was inserted
	filter := bson.M{"_id": "test_id"}
	var result MongoDoc
	err = collection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		t.Errorf("Error finding document: %v", err)
	}

	// Clean up
	err = collection.Drop(context.TODO())
	if err != nil {
		t.Errorf("Error dropping collection: %v", err)
	}
}