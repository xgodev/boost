// Package oracle keeps the stock in Oracle, chosen by
// INTERNAL_PKG_DOMAIN_REPOSITORY_STOCK_IMPLEMENTATION=oracle. Its *sql.DB is
// the one builder generates from boost's godror factory
// (boost.factory.godror.*), with the health plugin; godror builds with cgo
// and loads the Oracle client libraries at run time.
package oracle

import (
	"context"
	"database/sql"

	"github.com/godror/godror"
	"github.com/xgodev/boost/examples/stock/internal/pkg/domain/model/entity"
	"github.com/xgodev/boost/examples/stock/internal/pkg/domain/repository/stock"
	"github.com/xgodev/boost/model/errors"
)

// uniqueViolation is ORA-00001, a duplicate primary key.
const uniqueViolation = 1

type datastore struct {
	db    *sql.DB
	table string
}

// New provides the oracle stock on the table of Options.
// @Provide(index=0)
// @Inject(index=0, drivers=[godror], plugins=[health])
func New(db *sql.DB, o *Options) stock.Repository {
	return &datastore{db: db, table: o.Table}
}

func (d *datastore) Create(ctx context.Context, s *entity.Stock) error {
	_, err := d.db.ExecContext(ctx, "INSERT INTO "+d.table+" (sku, units) VALUES (:1, :2)", s.SKU, s.Units)
	if oraErr, ok := godror.AsOraErr(err); ok && oraErr.Code() == uniqueViolation {
		return errors.AlreadyExistsf("%s is already in stock", s.SKU)
	}
	if err != nil {
		return errors.Wrap(err, errors.Internalf("error creating the stock of %s", s.SKU))
	}
	return nil
}

func (d *datastore) Get(ctx context.Context, sku string) (*entity.Stock, error) {
	return scanUnits(d.db.QueryRowContext(ctx, "SELECT units FROM "+d.table+" WHERE sku = :1", sku), sku)
}

func (d *datastore) Take(ctx context.Context, sku string) (*entity.Stock, error) {
	tx, err := d.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, errors.Wrap(err, errors.Internalf("error taking a unit of %s", sku))
	}
	defer func() { _ = tx.Rollback() }()

	// the row stays locked until the commit, so two takes never sell one unit twice
	s, err := scanUnits(tx.QueryRowContext(ctx, "SELECT units FROM "+d.table+" WHERE sku = :1 FOR UPDATE", sku), sku)
	if err != nil {
		return nil, err
	}
	if s.Units == 0 {
		return nil, errors.Conflictf("%s is sold out", sku)
	}
	if _, err := tx.ExecContext(ctx, "UPDATE "+d.table+" SET units = units - 1 WHERE sku = :1", sku); err != nil {
		return nil, errors.Wrap(err, errors.Internalf("error taking a unit of %s", sku))
	}
	if err := tx.Commit(); err != nil {
		return nil, errors.Wrap(err, errors.Internalf("error taking a unit of %s", sku))
	}
	s.Units--
	return s, nil
}

// scanUnits reads the units of sku from row; no row is a NotFound.
func scanUnits(row *sql.Row, sku string) (*entity.Stock, error) {
	var units int64
	err := row.Scan(&units)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, errors.NewNotFound(nil, sku+" is not in stock")
	}
	if err != nil {
		return nil, errors.Wrap(err, errors.Internalf("error reading the stock of %s", sku))
	}
	return &entity.Stock{SKU: sku, Units: units}, nil
}
