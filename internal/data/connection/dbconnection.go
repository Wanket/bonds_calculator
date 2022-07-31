package connection

import (
	"bonds_calculator/internal/data/entgenerated"
	"bonds_calculator/internal/util"
	"context"
	"fmt"
	// Register postgres driver.
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
)

type DBConnection struct {
	*entgenerated.Client
}

func NewDBConnection(config util.IGlobalConfig) *DBConnection {
	dbConfig := config.GetDataBaseConfig()

	dataBaseClient, err := entgenerated.Open("postgres", fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		dbConfig.Host, dbConfig.Port, dbConfig.User, dbConfig.Password, dbConfig.DBName,
	))
	if err != nil {
		log.Fatalf("failed opening connection to database: %v", err)
	}

	if err := dataBaseClient.Schema.Create(context.Background()); err != nil {
		log.Fatalf("failed creating database schema: %v", err)
	}

	return &DBConnection{
		Client: dataBaseClient,
	}
}
