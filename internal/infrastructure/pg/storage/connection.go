package storage

import (
	"context"
	"log"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
)

type DatabaseService struct {
	Pool *pgxpool.Pool
}

func NewDatabaseService(databaseUrl string) (DatabaseService, error) {

	// Define the connection configuration
	//connConfig, err := pgxpool.ParseConfig("postgres://postgres:okokokokd@localhost:5432/casbin")
	connConfig, err := pgxpool.ParseConfig(databaseUrl)
	if err != nil {
		log.Fatalf("Unable to parse connection string: %v\n", err)
	}
	// Explicitly set the password (optional, if the password might be problematic in URL form)
	//connConfig.ConnConfig.Password = "okokokokd"

	//dbPool, err := pgxpool.Connect(context.Background(), databaseUrl)

	var dbPool *pgxpool.Pool
	for i := 0; i < 20; i++ {
		dbPool, err = pgxpool.ConnectConfig(context.Background(), connConfig)
		if err != nil {
			log.Printf("Unable to connect to database: %v\n", err)
		} else {
			break
		}
		if i == 19 {
			log.Fatalf("Unable to connect to database: %v\n", err)
		}
		time.Sleep(1 * time.Second)
	}

	// Verify the connection
	err = dbPool.Ping(context.Background())
	if err != nil {
		log.Fatalf("Unable to ping database: %v\n", err)
	}

	return DatabaseService{Pool: dbPool}, err
}
