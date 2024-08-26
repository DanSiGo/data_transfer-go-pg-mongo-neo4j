package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDoc struct {
	ID           string `bson:"_id"`
	IdCourse     string `bson:"idCourse"`
	IdObjective  string `bson:"idObjective"`
	IdMaterial   string `bson:"idMaterial"`
	Transcript string `bson:"transcript"`
	MaterialType string `bson:"materialType"`
	IsSuccessful bool   `bson:"isSuccessful"`
}

func main() {
	db, err := sql.Open("postgres", "postgres://postgres:postgres@localhost:5430/homero?sslmode=disable")
	if err != nil {
		fmt.Println("Error connecting to POSTGRES:", err)
		return
	}
	defer db.Close()

	fmt.Println("Connected to POSTGRES successfully!")

	rows, err := db.Query("SELECT learning_object_id, id, transcript_id, mimetype FROM media WHERE transcript_id IS NOT NULL LIMIT 10")
	if err != nil {
		fmt.Println("Error selecting from POSTGRES:", err)
		return
	}
	defer rows.Close()

	godotenv.Load()
	//dbURL := os.Getenv("MONGOURI")
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

	fmt.Println("Connected to MONGO succesfully!")

	defer client.Disconnect(context.TODO())

	dbName := "homero"
	colName := "ClassMaterial"
	collection := client.Database(dbName).Collection(colName)

	_, err = collection.DeleteMany(context.TODO(), bson.M{})
	if err != nil {
		fmt.Println("Error deleting documents:", err)
		return
	}

	for rows.Next() {
		var idObjective string
		var idMaterial string
		var transcript string
		var materialType string
		err := rows.Scan(&idObjective, &idMaterial, &transcript, &materialType)
		if err != nil {
			fmt.Println("Error scanning MONGO:", err)
			return
		}
		objID := primitive.NewObjectID()
		mongoDoc := MongoDoc{
			ID:           objID.Hex(), // generate a unique ID
			IdCourse:     "",
			IdObjective:  idObjective,
			IdMaterial:   idMaterial,
			Transcript: transcript,
			MaterialType: materialType,
			IsSuccessful: true,
		}
		_, err = collection.InsertOne(context.TODO(), mongoDoc)
		if err != nil {
			fmt.Println("Error inserting in MONGO:", err)
			return
		}
	}

	fmt.Println("Data inserted into MongoDB successfully!")
}
