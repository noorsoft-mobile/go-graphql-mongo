package main

import (
	"context"
	"github.com/graphql-go/graphql"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

func usersCollection() *mongo.Collection {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://mongo:27017"))
	if err != nil {
		log.Panic("Error when creating mongodb connection client", err)
	}
	collection := client.Database("testing").Collection("users")
	err = client.Connect(ctx)
	if err != nil {
		log.Panic("Error when connecting to mongodb", err)
	}

	return collection
}

func usersResolver(_ graphql.ResolveParams) (interface{}, error) {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	collection := usersCollection()
	result, err := collection.Find(ctx, bson.D{})
	if err != nil {
		log.Print("Error when finding user", err)
		return nil, err
	}

	defer result.Close(ctx)

	var r []bson.M
	err = result.All(ctx, &r)
	if err != nil {
		log.Print("Error when reading users from cursor", err)
	}

	return r, nil
}

func addUserResolver(p graphql.ResolveParams) (interface{}, error) {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	collection := usersCollection()
	id, err := collection.InsertOne(ctx, p.Args["input"])
	if err != nil {
		log.Print("Error when inserting user", err)
		return nil, err
	}

	var result bson.M
	err = collection.FindOne(ctx, bson.M{"_id": id.InsertedID}).Decode(&result)
	if err != nil {
		log.Print("Error when finding the inserted user by its id", err)
		return nil, err
	}

	return result, nil
}
