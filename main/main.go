package main

import (
	"context"
	"database/sql"
	"fmt"
	// "go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	_ "github.com/lib/pq"
)

type MongoDoc struct {
	ID            string `bson:"_id"`
	IdCourse      string    `bson:"idCourse"`
	IdObjective   string    `bson:"idObjective"`
	IdMaterial    string    `bson:"idMaterial"`
	IdTranscript  string    `bson:"idTranscript"`
	MaterialType  string `bson:"materialType"`
	IsSuccessful  bool   `bson:"isSuccessful"`
	TranscriptTime string `bson:"transcriptTime"`
	KeyWord       string `bson:"keyWord"`
}

func main() {
	// Connect to Postgres
	
	db, err := sql.Open("postgres", "postgres://postgres:postgres@localhost:5430/homero?sslmode=disable")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer db.Close()

	// Test the connection
	err = db.Ping()
	if err != nil {
		fmt.Println("Error connecting to POSTGRES:", err)
		return
	}

	fmt.Println("Connected to POSTGRES successfully!")

	// Fetch data from Postgres
	rows, err := db.Query("SELECT learning_object_id, id, transcript_id, mimetype, duration FROM media WHERE transcript_id IS NOT NULL LIMIT 10")
	if err != nil {
		fmt.Println("Error selecting from POSTGRES:", err)
		return
	}
	defer rows.Close()

	// Connect to MongoDB
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb+srv://mongotreino:y77RSYBO1un6mGoD@cluster0.vrped.mongodb.net/?retryWrites=true&w=majority&appName=Cluster0"))
	if err != nil {
		fmt.Println(err)
		return
	}

	// connect to mongo
	err = client.Connect(context.TODO())
	if err != nil {
		fmt.Println(err)
		return
	}

	// Testar a conexão com o banco de dados
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		fmt.Println("Error conecting to MONGO:", err)
		return
	}
	fmt.Println("Connected to MONGO succesfully!")
	defer client.Disconnect(context.TODO())

	// Get a handle to the database and collection
	dbName := "homero"
	colName := "ClassMaterial"
	collection := client.Database(dbName).Collection(colName)

	// Iterate over the rows and send data directly to MongoDB
	for rows.Next() {
		var idCourse string
		var idMaterial string
		var idTranscript string
		var materialType string
		var transcriptTime string
		err := rows.Scan(&idCourse, &idMaterial, &idTranscript, &materialType, &transcriptTime)
		if err != nil {
			fmt.Println("Error scanning MONGO:", err)
			return
		}
		mongoDoc := MongoDoc{
			ID:            fmt.Sprintf("%d", idCourse), // generate a unique ID
			IdCourse:      idCourse,
			IdObjective:   "", // set default value or fetch from another table
			IdMaterial:    idMaterial,
			IdTranscript:  idTranscript,
			MaterialType:  materialType,
			IsSuccessful:  true, // set default value or fetch from another table
			TranscriptTime: transcriptTime,
			KeyWord:       "", // set default value or fetch from another table
		}
		_, err = collection.InsertOne(context.TODO(), mongoDoc)
		if err != nil {
			fmt.Println("Error inserting in MONGO:", err)
			return
		}
	}

	fmt.Println("Data inserted into MongoDB successfully!")
}