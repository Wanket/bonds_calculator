package connection

import (
	"bonds_calculator/internal/data/entgenerated"
	"bonds_calculator/internal/util"
	"context"
	"fmt"
	"time"
	// Register postgres driver.
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
)

type DBConnection struct {
	*entgenerated.Client
}

func NewDBConnection(config util.IGlobalConfig) *DBConnection {
	dbConfig := config.GetDataBaseConfig()

	dataBaseClient, err := tryOpenDBConnection(dbConfig, nil)
	for err != nil {
		log.WithError(err).Error("DBConnection: failed to open connection to database, retrying...")

		time.Sleep(time.Second)

		dataBaseClient, err = tryOpenDBConnection(dbConfig, nil)
	}

	log.Info("DBConnection: Database connection established")

	return &DBConnection{
		Client: dataBaseClient,
	}
}

func tryOpenDBConnection(config util.DBConfig, dataBaseClient *entgenerated.Client) (*entgenerated.Client, error) {
	var err error

	if dataBaseClient == nil {
		dataBaseClient, err = entgenerated.Open("postgres", fmt.Sprintf(
			"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
			config.Host, config.Port, config.User, config.Password, config.DBName,
		))

		if err != nil {
			return nil, fmt.Errorf("failed opening connection to database: %w", err)
		}
	}

	if err := dataBaseClient.Schema.Create(context.Background()); err != nil {
		return dataBaseClient, fmt.Errorf("failed creating database schema: %w", err)
	}

	return dataBaseClient, nil
}
