package weather

import (
	"context"
	"fmt"
	"io"
	"net/http"

	//"encoding/json"
	"log"
	//"reflect"

	"go.mongodb.org/mongo-driver/bson"
	// "go.mongodb.org/mongo-driver/mongo"
	// "go.mongodb.org/mongo-driver/mongo/options"

	"github.com/andriykutsevol/DDDCasbinExample/internal/domain/weather"
	"github.com/andriykutsevol/DDDCasbinExample/internal/infrastructure/mongo/storage"
)

type repository struct {
	storage *storage.MongoStorage
}

func NewRepository(storage *storage.MongoStorage) weather.Repository {
	return &repository{
		storage: storage,
	}
}

func (r *repository) Seed(ctx context.Context, params map[string]string) error {
	fmt.Println("weather Seed()")

	//https://open-meteo.com/en/docs
	//https://api.open-meteo.com/v1/forecast?latitude=52.52&longitude=13.41&hourly=temperature_2m
	resp, err := http.Get("https://api.open-meteo.com/v1/forecast?latitude=52.52&longitude=13.41&hourly=temperature_2m&past_days=92&forecast_days=1")
	if err != nil {
		log.Fatalln(err)
	}
	//We Read the response body on the line below.
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	doc := map[string]string{"city": "Dnipro", "data": string(body)}
	collection := r.storage.GetCollection()
	insertResult, err := collection.InsertOne(context.TODO(), doc)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Inserted a single document: ", insertResult.InsertedID)

	return nil
}

func (r *repository) Retrieve(ctx context.Context, city string) (string, error) {

	// filter := bson.D{{"city", city}}
	// _ = filter

	filter := bson.E{Key: "city", Value: city}

	collection := r.storage.GetCollection()

	var result bson.M
	err := collection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		log.Println("No city found")
		return "", nil
	}

	return result["data"].(string), nil
}
