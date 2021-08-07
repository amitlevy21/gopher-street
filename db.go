package main

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func WriteExpensesToDB(collection *mongo.Collection, ctx context.Context, expenses *Expenses) error {
	converted := make([]interface{}, len(*expenses))
	for i, e := range *expenses {
		converted[i] = e
	}
	log.Println("Writing expenses to DB")
	_, err := collection.InsertMany(ctx, converted)

	return err
}

func openDB() (context.Context, context.CancelFunc, *mongo.Client) {
	uri := "mongodb://localhost"
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	return ctx, cancel, openDBWithURI(ctx, uri)
}

func openDBWithURI(ctx context.Context, uri string) *mongo.Client {
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		panic(err)
	}
	return client
}

func closeDB(ctx context.Context, client *mongo.Client) {
	if err := client.Disconnect(ctx); err != nil {
		panic(err)
	}
}
