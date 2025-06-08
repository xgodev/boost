package main

import (
	"context"
	"fmt"
	"time"

	"github.com/xgodev/boost"
	"github.com/xgodev/boost/factory/contrib/go.mongodb.org/mongo-driver/v1"
	"github.com/xgodev/boost/wrapper/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// TestDocument represents a simple document structure for testing
type TestDocument struct {
	ID        string    `bson:"_id"`
	Name      string    `bson:"name"`
	Value     int       `bson:"value"`
	CreatedAt time.Time `bson:"created_at"`
}

func main() {
	// Initialize Boost framework
	boost.Start()

	// Create context with logger
	ctx := context.Background()
	logger := log.FromContext(ctx)

	// Create MongoDB connection with default options
	logger.Info("Connecting to MongoDB...")
	conn, err := mongo.NewConn(ctx)
	if err != nil {
		logger.Errorf("Failed to connect to MongoDB: %v", err)
		return
	}

	// Get a collection reference
	collection := conn.Database.Collection("test_documents")

	// Insert a test document
	doc := TestDocument{
		ID:        "test-1",
		Name:      "Test Document",
		Value:     42,
		CreatedAt: time.Now(),
	}

	_, err = collection.InsertOne(ctx, doc)
	if err != nil {
		logger.Errorf("Failed to insert document: %v", err)
		return
	}
	logger.Infof("Inserted document with ID: %s", doc.ID)

	// Read the document back
	var result TestDocument
	err = collection.FindOne(ctx, bson.M{"_id": doc.ID}).Decode(&result)
	if err != nil {
		logger.Errorf("Failed to find document: %v", err)
		return
	}
	logger.Infof("Found document: %+v", result)

	// Validate connection pool settings
	stats := conn.Client.NumberSessionsInProgress()
	logger.Infof("Current sessions in progress: %d", stats)

	// Print client options to verify configuration
	printClientOptions(conn, logger)

	// Close the connection when done
	err = conn.Client.Disconnect(ctx)
	if err != nil {
		logger.Errorf("Failed to disconnect from MongoDB: %v", err)
	}
	logger.Info("Disconnected from MongoDB")
}

// printClientOptions prints the current client options configuration for validation
func printClientOptions(conn *mongo.Conn, logger log.Logger) {
	opts := conn.ClientOptions

	logger.Info("MongoDB Client Configuration:")
	logger.Infof("- URI: %s", opts.GetURI())
	
	if opts.GetMaxPoolSize() != nil {
		logger.Infof("- Max Pool Size: %d", *opts.GetMaxPoolSize())
	}
	if opts.GetMinPoolSize() != nil {
		logger.Infof("- Min Pool Size: %d", *opts.GetMinPoolSize())
	}
	if opts.GetMaxConnIdleTime() != nil {
		logger.Infof("- Max Conn Idle Time: %v", *opts.GetMaxConnIdleTime())
	}
	
	if opts.GetConnectTimeout() != nil {
		logger.Infof("- Connect Timeout: %v", *opts.GetConnectTimeout())
	}
	if opts.GetSocketTimeout() != nil {
		logger.Infof("- Socket Timeout: %v", *opts.GetSocketTimeout())
	}
	if opts.GetServerSelectionTimeout() != nil {
		logger.Infof("- Server Selection Timeout: %v", *opts.GetServerSelectionTimeout())
	}
	
	if opts.GetReplicaSet() != "" {
		logger.Infof("- Replica Set: %s", opts.GetReplicaSet())
	}
	
	if opts.GetReadConcern() != nil {
		logger.Infof("- Read Concern: %s", opts.GetReadConcern().GetLevel())
	}
	
	if opts.GetReadPreference() != nil {
		logger.Infof("- Read Preference: %s", opts.GetReadPreference().Mode())
	}
	
	if opts.GetWriteConcern() != nil {
		wc := opts.GetWriteConcern()
		w, _ := wc.GetW()
		j, _ := wc.GetJ()
		wt, _ := wc.GetWTimeout()
		logger.Infof("- Write Concern: w=%v, j=%v, wtimeout=%v", w, j, wt)
	}
	
	if opts.GetRetryReads() != nil {
		logger.Infof("- Retry Reads: %v", *opts.GetRetryReads())
	}
	if opts.GetRetryWrites() != nil {
		logger.Infof("- Retry Writes: %v", *opts.GetRetryWrites())
	}
	
	if opts.GetCompressors() != nil {
		logger.Infof("- Compressors: %v", opts.GetCompressors())
	}
	
	if opts.GetAppName() != "" {
		logger.Infof("- App Name: %s", opts.GetAppName())
	}
}
