package main

import (
	"context"
	"database/sql"
	"fmt"
	// "go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	// "github.com/lib/pq"
)

type MongoDoc struct {
	ID            string `bson:"_id"`
	IdCourse      int    `bson:"idCourse"`
	IdObjective   int    `bson:"idObjective"`
	IdMaterial    int    `bson:"idMaterial"`
	IdTranscript  int    `bson:"idTranscript"`
	MaterialType  string `bson:"materialType"`
	IsSuccessful  bool   `bson:"isSuccessful"`
	TranscriptTime string `bson:"transcriptTime"`
	KeyWord       string `bson:"keyWord"`
}

func main() {
	// Connect to Postgres
	db, err := sql.Open("postgres", "user:password@localhost/database")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer db.Close()

	// Fetch data from Postgres
	rows, err := db.Query("SELECT IdCourse, IdMaterial, IdTranscript, MaterialType, TranscriptTime FROM your_table_name")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer rows.Close()

	// Connect to MongoDB
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		fmt.Println(err)
		return
	}
	defer client.Disconnect(context.TODO())

	// Get a handle to the database and collection
	dbName := "your_database_name"
	colName := "your_collection_name"
	collection := client.Database(dbName).Collection(colName)

	// Iterate over the rows and send data directly to MongoDB
	for rows.Next() {
		var idCourse int
		var idMaterial int
		var idTranscript int
		var materialType string
		var transcriptTime string
		err := rows.Scan(&idCourse, &idMaterial, &idTranscript, &materialType, &transcriptTime)
		if err != nil {
			fmt.Println(err)
			return
		}
		mongoDoc := MongoDoc{
			ID:            fmt.Sprintf("%d", idCourse), // generate a unique ID
			IdCourse:      idCourse,
			IdObjective:   0, // set default value or fetch from another table
			IdMaterial:    idMaterial,
			IdTranscript:  idTranscript,
			MaterialType:  materialType,
			IsSuccessful:  true, // set default value or fetch from another table
			TranscriptTime: transcriptTime,
			KeyWord:       "", // set default value or fetch from another table
		}
		_, err = collection.InsertOne(context.TODO(), mongoDoc)
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	fmt.Println("Data inserted into MongoDB successfully!")
}