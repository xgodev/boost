package main

import (
	"context"
	"github.com/xgodev/boost/wrapper/config"
	"os"
	"time"

	mgo "github.com/xgodev/boost/factory/contrib/go.mongodb.org/mongo-driver/v1"
	vault "github.com/xgodev/boost/factory/contrib/mittwald/vaultgo/v0"
	"github.com/xgodev/boost/factory/contrib/mittwald/vaultgo/v0/managers/contrib/go.mongodb.org/mongo-driver/v1"
	ilog "github.com/xgodev/boost/factory/local/wrapper/log"
	"github.com/xgodev/boost/wrapper/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func init() {
	os.Setenv("BOOST_FACTORY_MONGO_URI", "mongodb://localhost:27004/database")
	os.Setenv("BOOST_FACTORY_VAULT_MANAGERS_MONGO_SECRET__PATH", "v1/database/creds/recency_vault_hml")
	os.Setenv("BOOST_FACTORY_VAULT_ADDR", "http://vault.example.com")
	os.Setenv("BOOST_FACTORY_VAULT_TOKEN", "s.O12jxPL9bbJhS7laeVP2h6fK")
	os.Setenv("BOOST_FACTORY_VAULT_TYPE", "TOKEN")
	os.Setenv("BOOST_FACTORY_LOGRUS_CONSOLE_LEVEL", "TRACE")
}

var conn *mgo.Conn

func main() {

	config.Load()
	ilog.New()

	ctx := context.Background()

	var err error

	conn, err = mgo.NewConn(ctx)
	if err != nil {
		log.Error(err)
	}

	mgoManager := mongo.NewManager(conn)
	vault.ManageAll(ctx, mgoManager)

	for {
		mongoQuery()
		time.Sleep(2 * time.Second)
	}
}

func mongoQuery() {

	if conn == nil {
		return
	}

	col := conn.Database.Collection("teste")
	objID, _ := primitive.ObjectIDFromHex("5f19a4416e6274c01d474089")
	result := col.FindOne(context.Background(), bson.M{"_id": objID})
	if result.Err() != nil {
		log.Error(result.Err())
	} else {
		log.Infof("success")
	}

}
