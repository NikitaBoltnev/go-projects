package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"PR_service/api"
	"PR_service/handler"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	connString := os.Getenv("DATABASE_URL")
	if connString == "" {
		connString = "postgres://postgres:1@localhost:5432/postgres"
	}

	dbpool, err := pgxpool.New(context.Background(), connString)
	if err != nil {
		log.Fatalf("error create pool: %v\n", err)
	}
	defer dbpool.Close()

	h := handler.NewHandler(dbpool)

	httpHandler := api.Handler(h)

	port := "8080"
	log.Fatal(http.ListenAndServe(":"+port, httpHandler))
}
