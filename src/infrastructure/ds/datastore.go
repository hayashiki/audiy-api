package ds

import (
	"cloud.google.com/go/datastore"

	"go.mercari.io/datastore/boom"
	"go.mercari.io/datastore/clouddatastore"

	"context"
	"github.com/hayashiki/audiy-api/src/config"
	"github.com/pkg/errors"
	"log"
)

const defaultLimit = 100

var ErrNoSuchEntity = errors.New("datastore: no such entity")

type client struct {}

type datastoreTransactor struct {
}

func NewDatastoreTransactor() Transactor {
	return &datastoreTransactor{}
}

func (t *datastoreTransactor) RunInTransaction(ctx context.Context, fn func(tx *boom.Transaction) error) error {
	if _, err := FromContext(ctx).RunInTransaction(fn); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func New() Client {
	return &client{}
}

type Client interface {
	GetAll(ctx context.Context, kind string, filters map[string]interface{}, dst interface{}, cursor string, limit int, orderBy string) error
	GetMulti(ctx context.Context, dst interface{}) error
	Get(ctx context.Context, dst interface{}) error
	Put(tx *boom.Transaction, src interface{}) error
	Delete(tx *boom.Transaction, src interface{}) error
}

type Transactor interface {
	RunInTransaction(context.Context, func(tx *boom.Transaction) error,) error
}

func FromContext(ctx context.Context) *boom.Boom {
	cli, err := datastore.NewClient(ctx, config.GetProject())
	if err != nil {
		log.Println("cli", err)
		panic(err)
	}
	ds, err := clouddatastore.FromClient(ctx, cli)
	if err != nil {
		panic(err)
	}
	return boom.FromClient(ctx, ds)
}

func (c *client) GetAll(ctx context.Context, kind string, filters map[string]interface{}, dst interface{}, cursor string, limit int, orderBy string) error {
	b := FromContext(ctx)
	q := b.Client.NewQuery(kind)

	if len(filters) > 0 {
		for k, v := range filters {
			q = q.Filter(k, v)
		}
	}
	if orderBy != "" {
		q = q.Order(orderBy)
	}
	if cursor != "" {
		dsCursor, err := datastore.DecodeCursor(cursor)
		if err != nil {
			//TODO
			log.Printf("failed to decode %v", err)
		}
		q = q.Start(dsCursor)
	}
	if limit != 0 {
		q = q.Limit(limit)
	} else {
		limit = defaultLimit
	}
	if _, err := b.GetAll(q, dst); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (c *client) Get(ctx context.Context, dst interface{}) error {
	b := FromContext(ctx)
	if err := b.Get(dst); err != nil {
		if err == datastore.ErrNoSuchEntity {
			return errors.WithStack(ErrNoSuchEntity)
		}
		return errors.WithStack(err)
	}

	return nil
}

func (c *client) GetMulti(ctx context.Context, dst interface{}) error {
	b := FromContext(ctx)

	if err := b.GetMulti(dst); err != nil {
		multiErr, ok := err.(datastore.MultiError)
		if !ok {
			return errors.WithStack(err)
		}

		for _, e := range multiErr {
			if e == datastore.ErrNoSuchEntity {
				return errors.WithStack(ErrNoSuchEntity)
			}
		}
		return errors.WithStack(err)
	}

	return nil
}

func (c *client) Put(tx *boom.Transaction, src interface{}) error {
	if _, err := tx.Put(src); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (c *client) Delete(tx *boom.Transaction, src interface{}) error {
	if err := tx.Delete(src); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

