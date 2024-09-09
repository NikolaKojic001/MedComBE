package Repository

import (
	"context"
	"fmt"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	uri           string
	clientOptions *options.ClientOptions
	client        *mongo.Client
	err           error
)

func InitConnection() {
	uri, _ = os.LookupEnv("MONGOCONNECTION")
	clientOptions = options.Client().ApplyURI(uri)
	client, err = mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to MongoDB")
}

func GetClient() *mongo.Client {
	return client
}
