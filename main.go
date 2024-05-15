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
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	// Load the .env file
	err := godotenv.Load()

	if err != nil {
		log.Fatalf("Error loading .env file : %v", err)
	}

	// Get Variables from the environment
	mongoURI := os.Getenv("MONGO_URI")
	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}

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

	// Get All Fish Documents
	routes.GET("/fish", func(c *gin.Context) {
		coll := client.Database("resep").Collection("ikan")

		cursor, err := coll.Find(context.TODO(), bson.D{})

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}

		defer cursor.Close(context.TODO())

		var docs []bson.M

		if err = cursor.All(context.TODO(), &docs); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		}

		c.JSON(http.StatusOK, gin.H{"message": docs})
	})

	// GET One Fish Documents
	routes.GET("/fish/:id", func(c *gin.Context) {
		coll := client.Database("resep").Collection("ikan")

		// Obtain ID from Parameter
		id := c.Param("id")
		objID, err := primitive.ObjectIDFromHex(id)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		}

		var doc bson.M

		if err = coll.FindOne(context.TODO(), bson.M{"_id": objID}).Decode(&doc); err != nil {
			if err == mongo.ErrNoDocuments {
				c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			}
			return
		}

		c.JSON(http.StatusOK, doc)
	})

	log.Printf("Running on Port %s", port)

	if err := routes.Run(":" + port); err != nil {
		log.Fatalf("Failed to Start : %v", err)
	}
}
