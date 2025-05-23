package main

import (
	"context"
	"time"

	"github.com/xgodev/boost"
	"github.com/xgodev/boost/factory/contrib/go.mongodb.org/mongo-driver/v1"
	"github.com/xgodev/boost/wrapper/config"
	"github.com/xgodev/boost/wrapper/log"
)

func main() {
	// Initialize Boost framework
	boost.Start()

	// Set advanced configuration options
	setupAdvancedConfig()

	// Create context with logger
	ctx := context.Background()
	logger := log.FromContext(ctx)

	// Create MongoDB connection with advanced options
	logger.Info("Connecting to MongoDB with advanced configuration...")
	conn, err := mongo.NewConn(ctx)
	if err != nil {
		logger.Errorf("Failed to connect to MongoDB: %v", err)
		return
	}

	// Ping the database to verify connection
	err = conn.Client.Ping(ctx, nil)
	if err != nil {
		logger.Errorf("Failed to ping MongoDB: %v", err)
		return
	}

	logger.Info("Successfully connected to MongoDB with advanced configuration")

	// Perform database operations here...

	// Close the connection when done
	err = conn.Client.Disconnect(ctx)
	if err != nil {
		logger.Errorf("Failed to disconnect from MongoDB: %v", err)
	}
}

// setupAdvancedConfig demonstrates how to configure advanced MongoDB options
func setupAdvancedConfig() {
	// MongoDB connection settings
	config.Set("boost.factory.mongo.uri", "mongodb://localhost:27017/mydb")
	
	// Connection pool settings
	config.Set("boost.factory.mongo.pool.max_size", 100)
	config.Set("boost.factory.mongo.pool.min_size", 10)
	config.Set("boost.factory.mongo.pool.max_idle_time_ms", 30000) // 30 seconds
	
	// Timeout settings
	config.Set("boost.factory.mongo.timeout.connect_ms", 5000)      // 5 seconds
	config.Set("boost.factory.mongo.timeout.socket_ms", 10000)      // 10 seconds
	config.Set("boost.factory.mongo.timeout.server_selection_ms", 5000) // 5 seconds
	
	// Read concern and preference
	config.Set("boost.factory.mongo.read_concern.level", "majority")
	config.Set("boost.factory.mongo.read_preference.mode", "secondaryPreferred")
	
	// Write concern
	config.Set("boost.factory.mongo.write_concern.level", "majority")
	config.Set("boost.factory.mongo.write_concern.j", true)
	config.Set("boost.factory.mongo.write_concern.wtimeout_ms", 5000) // 5 seconds
	
	// Retry settings
	config.Set("boost.factory.mongo.retry.reads", true)
	config.Set("boost.factory.mongo.retry.writes", true)
	
	// Application name for monitoring
	config.Set("boost.factory.mongo.app_name", "MyBoostApplication")
}
