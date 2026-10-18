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

// New connects through boost's mongo factory (boost.factory.mongo.*).
// @Provide(index=0)
func New(ctx context.Context, o *Options) (stock.Repository, error) {
	conn, err := mongofact.NewConn(ctx)
	if err != nil {
		return nil, err
	}
	return &datastore{coll: conn.Database.Collection(o.Collection)}, nil
}

func (d *datastore) Take(ctx context.Context, sku string) error {
	res, err := d.coll.UpdateOne(ctx,
		bson.M{"_id": sku, "units": bson.M{"$gt": 0}},
		bson.M{"$inc": bson.M{"units": -1}})
	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		return errors.Conflictf("%s is sold out", sku)
	}
	return nil
}
