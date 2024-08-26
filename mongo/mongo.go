package mongo

import (
	"context"
	"fmt"
	// "os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDoc struct {
	ID           string `bson:"_id"`
	IdCourse     string `bson:"idCourse"`
	IdObjective  string `bson:"idObjective"`
	IdMaterial   string `bson:"idMaterial"`
	Transcript   string `bson:"transcript"`
	MaterialType string `bson:"materialType"`
	IsSuccessful bool   `bson:"isSuccessful"`
}

var (
	mongoClient *mongo.Client
	mongoCtx    context.Context
)

func init()  {
	godotenv.Load()
	//dbURL := os.Getenv("MONGOURI")
	// mongodb://localhost:27017/
	// dbURL := os.Getenv("MONGOLOCAL")
	mongoCtx, _ = context.WithTimeout(context.Background(), 10*time.Second)
	mongoClient, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017/"))
	if err != nil {
		fmt.Println("Error creating client from MONGO:", err)
		return
	}

	err = mongoClient.Connect(context.TODO())
	if err != nil {
		fmt.Println("Error conecting to MONGO:", err)
		return
	}
	fmt.Println("Connected to MONGO succesfully!")
}

func GetMongo() *mongo.Client{
	return mongoClient
}

func CloseMongo() {
	mongoClient.Disconnect(mongoCtx)
}