package dbconn

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client = DBinstance()

func DBinstance() *mongo.Client {
	// calorie123@mongodb:27017/ the moongodb here refers to the service created by docker-compose.
	// and we are connecting to 27017 port of the mongodb service.
	// we cant use localhost here because all 3 services inside docker-compose are running in their
	// own serperate containers and each of them recognizes itself as the localhost.
	MongoDb := "mongodb://calorie-user:calorie123@mongodb:27017/?authSource=calorie-tracker-db"

	client, err := mongo.NewClient(options.Client().ApplyURI(MongoDb))
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to MongoDB")
	return client
}

func OpenCollection(client *mongo.Client, collectionName string) *mongo.Collection {
	var collection *mongo.Collection = client.Database("calorie-tracker-db").Collection(collectionName)
	return collection
}
