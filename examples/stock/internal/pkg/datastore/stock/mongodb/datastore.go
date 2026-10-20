// Package mongodb keeps the stock in MongoDB, chosen by
// INTERNAL_PKG_DOMAIN_REPOSITORY_STOCK_IMPLEMENTATION=mongodb. Its connection
// is the one builder generates from boost's mongo factory
// (boost.factory.mongo.*), with the health and prometheus plugins.
package mongodb

import (
	"context"

	"github.com/xgodev/boost/examples/stock/internal/pkg/domain/model/entity"
	"github.com/xgodev/boost/examples/stock/internal/pkg/domain/repository/stock"
	mongofact "github.com/xgodev/boost/factory/contrib/go.mongodb.org/mongo-driver/v2"
	"github.com/xgodev/boost/model/errors"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type datastore struct {
	coll *mongo.Collection
}

// document is a stock as MongoDB stores it: the sku is the _id.
type document struct {
	SKU   string `bson:"_id"`
	Units int64  `bson:"units"`
}

// New provides the mongo stock on the collection of Options.
// @Provide(index=0)
// @Inject(index=0, plugins=[health, prometheus])
func New(conn *mongofact.Conn, o *Options) stock.Repository {
	return &datastore{coll: conn.Database.Collection(o.Collection)}
}

func (d *datastore) Create(ctx context.Context, s *entity.Stock) error {
	_, err := d.coll.InsertOne(ctx, document{SKU: s.SKU, Units: s.Units})
	if mongo.IsDuplicateKeyError(err) {
		return errors.AlreadyExistsf("%s is already in stock", s.SKU)
	}
	if err != nil {
		return errors.Wrap(err, errors.Internalf("error creating the stock of %s", s.SKU))
	}
	return nil
}

func (d *datastore) Get(ctx context.Context, sku string) (*entity.Stock, error) {
	var doc document
	err := d.coll.FindOne(ctx, bson.M{"_id": bson.M{"$eq": sku}}).Decode(&doc)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, errors.NewNotFound(nil, sku+" is not in stock")
	}
	if err != nil {
		return nil, errors.Wrap(err, errors.Internalf("error reading the stock of %s", sku))
	}
	return &entity.Stock{SKU: doc.SKU, Units: doc.Units}, nil
}

func (d *datastore) Take(ctx context.Context, sku string) (*entity.Stock, error) {
	var doc document
	err := d.coll.FindOneAndUpdate(ctx,
		bson.M{"_id": bson.M{"$eq": sku}, "units": bson.M{"$gt": 0}},
		bson.M{"$inc": bson.M{"units": -1}},
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	).Decode(&doc)
	if errors.Is(err, mongo.ErrNoDocuments) {
		// nothing taken: either the sku is unknown or it is sold out
		if _, err := d.Get(ctx, sku); err != nil {
			return nil, err
		}
		return nil, errors.Conflictf("%s is sold out", sku)
	}
	if err != nil {
		return nil, errors.Wrap(err, errors.Internalf("error taking a unit of %s", sku))
	}
	return &entity.Stock{SKU: doc.SKU, Units: doc.Units}, nil
}
