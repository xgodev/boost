package main

import (
	"context"
	"fmt"
	"time"

	"github.com/xgodev/boost"
	"github.com/xgodev/boost/factory/contrib/go.mongodb.org/mongo-driver/v1"
	"github.com/xgodev/boost/wrapper/config"
	"github.com/xgodev/boost/wrapper/log"
)

func main() {
	// Initialize Boost framework
	boost.Start()

	// Create context with logger
	ctx := context.Background()
	logger := log.FromContext(ctx)

	// Test case 1: URI parameters should override explicit options
	logger.Info("Test Case 1: URI parameters should override explicit options")
	
	// Set explicit options
	config.Set("boost.factory.mongo.pool.max_size", 100)
	config.Set("boost.factory.mongo.read_preference.mode", "secondaryPreferred")
	
	// Set URI with conflicting options
	config.Set("boost.factory.mongo.uri", "mongodb://localhost:27017/mydb?maxPoolSize=50&readPreference=primary")
	
	// Create MongoDB connection
	conn1, err := mongo.NewConn(ctx)
	if err != nil {
		logger.Errorf("Failed to connect: %v", err)
		return
	}
	
	// Validate that URI options took precedence
	opts1 := conn1.ClientOptions
	logger.Infof("Max Pool Size: %d (should be 50 from URI, not 100 from explicit option)", *opts1.GetMaxPoolSize())
	logger.Infof("Read Preference: %s (should be primary from URI, not secondaryPreferred from explicit option)", 
		opts1.GetReadPreference().Mode())
	
	// Test case 2: Explicit options not in URI should still be applied
	logger.Info("\nTest Case 2: Explicit options not in URI should still be applied")
	
	// Reset configuration
	config.Set("boost.factory.mongo.pool.max_size", 100)
	config.Set("boost.factory.mongo.min_pool_size", 10)
	config.Set("boost.factory.mongo.write_concern.j", true)
	
	// Set URI without these options
	config.Set("boost.factory.mongo.uri", "mongodb://localhost:27017/mydb")
	
	// Create MongoDB connection
	conn2, err := mongo.NewConn(ctx)
	if err != nil {
		logger.Errorf("Failed to connect: %v", err)
		return
	}
	
	// Validate that explicit options were applied when not in URI
	opts2 := conn2.ClientOptions
	logger.Infof("Min Pool Size: %d (should be 10 from explicit option)", *opts2.GetMinPoolSize())
	j, _ := opts2.GetWriteConcern().GetJ()
	logger.Infof("Write Concern J: %v (should be true from explicit option)", j)
	
	// Close connections
	conn1.Client.Disconnect(ctx)
	conn2.Client.Disconnect(ctx)
	
	logger.Info("Validation complete")
}
