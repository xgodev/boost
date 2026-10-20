// Package mongodb keeps the stock in MongoDB, on boost's mongo factory.
package mongodb

import (
	"context"

	"github.com/acme/shop/internal/pkg/domain/repository/stock"
	mongofact "github.com/xgodev/boost/factory/contrib/go.mongodb.org/mongo-driver/v2"
	"github.com/xgodev/boost/model/errors"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type datastore struct {
	coll *mongo.Collection
}

// New keeps the stock on the conn inject builds from boost.factory.mongo.*,
// which the app closes when it stops.
// @Provide(index=0)
// @Inject(index=0, plugins=[health, prometheus])
func New(conn *mongofact.Conn, o *Options) stock.Repository {
	return &datastore{coll: conn.Database.Collection(o.Collection)}
}

func (d *datastore) Take(ctx context.Context, sku string) error {
	res, err := d.coll.UpdateOne(ctx,
		bson.M{"_id": bson.M{"$eq": sku}, "units": bson.M{"$gt": 0}},
		bson.M{"$inc": bson.M{"units": -1}})
	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		return errors.Conflictf("%s is sold out", sku)
	}
	return nil
}
