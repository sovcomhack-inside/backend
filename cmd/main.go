package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/spf13/viper"

	"github.com/sovcomhack-inside/internal/api"
	"github.com/sovcomhack-inside/internal/pkg/store"
	"github.com/sovcomhack-inside/internal/pkg/store/xpgx"
)

const (
	CONFIG_PATH     = "config.yaml"
	DEFAULT_ADDRESS = "0.0.0.0"
	DEFAULT_PORT    = "8080"
)

func main() {
	// -------------------- Set up viper -------------------- //

	viper.SetConfigFile(CONFIG_PATH)
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("failed to read the config file: %s\n", err)
	}

	viper.SetDefault("service.bind.address", DEFAULT_ADDRESS)
	viper.SetDefault("service.bind.port", DEFAULT_PORT)

	// -------------------- Set up database -------------------- //

	dbPool, err := pgxpool.New(context.Background(), viper.GetString("postgres.connection_string"))
	if err != nil {
		log.Fatalf("failed to connect to the postgres database: %s", err)
	}
	defer dbPool.Close()

	store := store.NewStore(xpgx.NewPool(dbPool))

	// -------------------- Set up service -------------------- //

	svc, err := api.NewAPIService(store)
	if err != nil {
		log.Fatalf("error creating service instance: %s", err)
	}

	go svc.Serve(viper.GetString("service.bind.address") + ":" + viper.GetString("service.bind.port"))

	// -------------------- Listen for INT signal -------------------- //

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	if err := svc.Shutdown(); err != nil {
		log.Fatal(err)
	}
}
