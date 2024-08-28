package main

import (
	"context"
	"crypto/tls"
	"database/sql"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/neo4j/neo4j-go-driver/neo4j"

	// "go.mongodb.org/mongo-driver/bson"
	// "go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDoc struct {
	MongoID      string `bson:"_id"`
	ID           string `bson:"id"`
	IdCourse     string `bson:"idCourse"`
	IdObjective  string `bson:"idObjective"`
	IdMaterial   string `bson:"idMaterial"`
	Transcript   string `bson:"transcript"`
	MaterialType string `bson:"materialType"`
	IsSuccessful bool   `bson:"isSuccessful"`
}

type Transcript struct {
	ID         string `bson:"_id"`
	Transcript string `bson:"transcript"`
}

func main() {
	db, err := sql.Open("postgres", "postgres://postgres:postgres@localhost:5430/homero?sslmode=disable")
	if err != nil {
		fmt.Println("Error connecting to POSTGRES:", err)
		return
	}
	defer db.Close()

	fmt.Println("Connected to POSTGRES!")

	rows, err := db.Query("SELECT learning_object_id, id, transcript_id, mimetype FROM media WHERE transcript_id IS NOT NULL LIMIT 10")
	if err != nil {
		fmt.Println("Error selecting from POSTGRES:", err)
		return
	}
	defer rows.Close()

	godotenv.Load()
	dbURL := os.Getenv("MONGOLOCAL")
	client, err := mongo.NewClient(options.Client().ApplyURI(dbURL))
	if err != nil {
		fmt.Println("Error creating client from MONGO:", err)
		return
	}

	err = client.Connect(context.TODO())
	if err != nil {
		fmt.Println("Error conecting to MONGO:", err)
		return
	}

	fmt.Println("Connected to MONGO!")

	defer client.Disconnect(context.TODO())


	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
	}
	driver, err := neo4j.NewDriver("bolt://localhost:7474", neo4j.BasicAuth("neo4j", "1234qwer", ""), func(c *neo4j.Config) { c.TLSConfig = tlsConfig })
	if err != nil {
		fmt.Println("Error creating neo4j driver:", err)
	}
	defer driver.Close()

	session, err := driver.Session(neo4j.AccessModeWrite)
	if err != nil {
		fmt.Println("Error creating neo4j session:", err)
	}
	defer session.Close()
	fmt.Println("Connected to NEO4J!")


	result, err := session.Run("MATCH (c:Course {id: $id})-[r]-() RETURN r, type(r)", map[string]interface{}{"id": "8713797a-c21f-4c5e-8b4a-86d2e0c2b6e6"})
	if err != nil {
		fmt.Println("Error fetching data from neo4j:", err)
		return
	}
	println(result)
	


	// dbName := "homero"
	// colName := "ClassMaterial"
	// collection := client.Database(dbName).Collection(colName)

	// _, err = collection.DeleteMany(context.TODO(), bson.M{})
	// if err != nil {
	// 	fmt.Println("Error deleting documents:", err)
	// 	return
	// }

	// for rows.Next() {
	// 	var idObjective string
	// 	var idMaterial string
	// 	var transcript_id string
	// 	var materialType string
	// 	err := rows.Scan(&idObjective, &idMaterial, &transcript_id, &materialType)
	// 	if err != nil {
	// 		fmt.Println("Error scanning MONGO:", err)
	// 		return
	// 	}

	// 	transcriptCol := client.Database(dbName).Collection("VideoLesson")
	// 	transcriptFilter := bson.M{"uuid": transcript_id}
	// 	var transcript Transcript
	// 	err = transcriptCol.FindOne(context.TODO(), transcriptFilter).Decode(&transcript)
	// 	if err != nil {
	// 		fmt.Println("Error fetching transcript:", err)
	// 		return
	// 	}

	// 	objId := primitive.NewObjectID()
	// 	mongoDoc := MongoDoc{
	// 		MongoID:      objId.Hex(),
	// 		ID:           transcript_id,
	// 		IdCourse:     "",
	// 		IdObjective:  idObjective,
	// 		IdMaterial:   idMaterial,
	// 		Transcript:   transcript.Transcript,
	// 		MaterialType: materialType,
	// 		IsSuccessful: true,
	// 	}

	// 	colName := "ClassMaterial"
	// 	collection := client.Database(dbName).Collection(colName)
	// 	_, err = collection.InsertOne(context.TODO(), mongoDoc)
	// 	if err != nil {
	// 		fmt.Println("Error inserting in MONGO:", err)
	// 		return
	// 	}
	// }

	// fmt.Println("Data inserted into MongoDB successfully!")
}
