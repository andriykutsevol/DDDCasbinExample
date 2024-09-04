package storage

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// DatabaseService represents a service that interacts with MongoDB.
type DatabaseService struct {
	client *mongo.Client
}

type DatabaseDB struct {
	database *mongo.Database
}

type DatabaseCollection struct {
	collection *mongo.Collection
}

// connectToMongoDB establishes a connection to MongoDB and returns a client.
func ConnectToMongoDB() (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Set up MongoDB connection options
	//clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")

	// https://docs.docker.com/network/network-tutorial-standalone/

	// docker run --name template_go_react_golang --network template_go_react_network -p 8080:8080 template_go_react_golang
	// So if we use docker run with --nerwork, we have to use container name
	// clientOptions := options.Client().ApplyURI("mongodb://template_go_react_mongodb:27017")

	// But if we use docker-compose, we have to use service name.
	// So it is better to define environment variable for URL.
	// clientOptions := options.Client().ApplyURI("mongodb://mongodb:27017")

	uri, ok := os.LookupEnv("MONGOURI")
	if !ok {
		return nil, errors.New("MONGOURI environment variable is not set")
	}

	clientOptions := options.Client().ApplyURI(uri)

	// Create a MongoDB client
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}

	// Ping the MongoDB server to verify the connection
	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}

	fmt.Println("Connected to MongoDB!")

	return client, nil
}

func NewDatabaseService(client *mongo.Client) *DatabaseService {
	return &DatabaseService{client: client}
}

func (s *DatabaseService) ListDatabaseNames() (dblist []string) {
	// List databases
	databaseNames, err := s.client.ListDatabaseNames(context.TODO(), bson.M{})
	if err != nil {
		log.Fatal(err)
	}

	// Print the names of all databases
	// for _, dbName := range databaseNames {
	// 	dblist = append(dblist, dbName)
	// 	//fmt.Println("Database:", dbName)
	// }
	dblist = append(dblist, databaseNames...)

	return dblist
}

func (s *DatabaseService) GetDatabase(name string) (db *DatabaseDB) {
	dbt := s.client.Database(name)
	db = &DatabaseDB{database: dbt}
	return db

}

func (db *DatabaseDB) ListCollections() (clist []string) {

	clistCur, err := db.database.ListCollections(context.TODO(), bson.M{})
	if err != nil {
		log.Fatal(err)
	}

	// Iterate over the collections and print their names
	for clistCur.Next(context.TODO()) {
		var result bson.M
		err := clistCur.Decode(&result)
		if err != nil {
			log.Fatal(err)
		}
		clist = append(clist, result["name"].(string))
	}

	return clist
}

func (db *DatabaseDB) GetCollection(name string) (c *DatabaseCollection) {
	ct := db.database.Collection(name)
	return &DatabaseCollection{collection: ct}
}

// func (s *DatabaseCollection) SetCollection(c *DatabaseCollection) {
//     s.collection = c
// }

//===================================================================================

func NewCollectionService(collection *DatabaseCollection) *DatabaseCollection {
	return &DatabaseCollection{collection: collection.collection}
}

func (c *DatabaseCollection) ClearColection() (response string) {

	// Remove all documents from the collection
	result, err := c.collection.DeleteMany(context.TODO(), bson.M{})
	if err != nil {
		log.Fatal(err)
	}

	response = fmt.Sprintf("Deleted %v documents from the collection.\n", result.DeletedCount)
	return response
}

func (c *DatabaseCollection) GetDatabaseCollection() *mongo.Collection {
	return c.collection
}
