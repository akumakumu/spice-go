// Package main is the entry point for the spice API application.
// It sets up a Gin-based HTTP server and defines routes for the API.
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	// Load the .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Get the MongoDB URI from the environment variable
	mongoURI := os.Getenv("MONGO_URI")

	// MongoDB Driver
	// Use the SetServerAPIOptions() method to set the version of the Stable API on the client
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(mongoURI).SetServerAPIOptions(serverAPI)

	// Create a new client and connect to the server
	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	// Send a ping to confirm a successful connection
	if err := client.Database("admin").RunCommand(context.TODO(), bson.D{{Key: "ping", Value: 1}}).Err(); err != nil {
		panic(err)
	}
	fmt.Println("Pinged your deployment. You successfully connected to MongoDB!")

	// Gin Init
	routes := gin.Default()

	// Get Fish Documents
	routes.GET("/fish", func(c *gin.Context) {
		coll := client.Database("resep").Collection("ikan")

		cursor, err := coll.Find(context.Background(), bson.D{})

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}

		defer cursor.Close(context.Background())

		var documents []bson.M

		if err = cursor.All(context.Background(), &documents); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		}

		c.JSON(http.StatusOK, gin.H{"message": documents})
	})

	routes.Run()
}
