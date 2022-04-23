package main

import (
	"context"
	"log"
	"sync"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var lock = &sync.Mutex{}

type connection struct {
	client *mongo.Client
	ctx    context.Context
}

type DB struct {
	connection
	database *mongo.Database
}

type cursor interface {
	Next(context.Context) bool
	Decode(interface{}) error
	Err() error
	Close(context.Context) error
}

var instance *DB

func Instance(ctx context.Context) *DB {
	if instance == nil {
		lock.Lock()
		defer lock.Unlock()
		if instance == nil {
			client := openDB(ctx)
			instance = &DB{connection{client, ctx}, client.Database("user")}
		}
	}
	return instance
}

func openDB(ctx context.Context) *mongo.Client {
	uri := "mongodb://localhost"
	return openDBWithURI(ctx, uri)
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

func (db *DB) closeDB(ctx context.Context) {
	lock.Lock()
	defer lock.Unlock()
	if err := db.client.Disconnect(ctx); err != nil {
		panic(err)
	}
	instance = nil
}
func (db *DB) dropDB(ctx context.Context) error {
	return db.database.Drop(ctx)
}

func (db *DB) WriteExpenses(ctx context.Context, expenses *Expenses) error {
	col := db.database.Collection("expenses")
	converted := make([]interface{}, len(expenses.ToSlice()))
	for i, e := range expenses.ToSlice() {
		converted[i] = e
	}
	log.Println("Writing expenses to DB")
	_, err := col.InsertMany(ctx, converted)

	return err
}

func (db *DB) GetExpenses(ctx context.Context) (*Expenses, error) {
	col := db.database.Collection("expenses")
	cur, err := col.Find(ctx, bson.D{})
	if err != nil {
		return &Expenses{}, err
	}
	defer db.closeCursor(ctx, cur)
	return db.getExpensesFromCur(ctx, cur)
}

func (db *DB) getExpensesFromCur(ctx context.Context, cur cursor) (*Expenses, error) {
	exps := Expenses{}
	for cur.Next(ctx) {
		var exp Expense
		if err := cur.Decode(&exp); err != nil {
			return &exps, err
		}
		exps.Classified = append(exps.Classified, &exp)
	}

	if err := cur.Err(); err != nil {
		return &exps, err
	}

	return &exps, nil
}

func (db *DB) closeCursor(ctx context.Context, cursor cursor) {
	if err := cursor.Close(ctx); err != nil {
		log.Printf("error while closing cursor: %s", err)
	}
}
