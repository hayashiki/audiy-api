package datastore

import (
	"cloud.google.com/go/datastore"
	"google.golang.org/api/iterator"

	"context"
	"github.com/pkg/errors"
	mdatastore "go.mercari.io/datastore"
	"go.mercari.io/datastore/boom"
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
	GetAll(ctx context.Context, kind string, dst interface{}) error
	GetMulti(ctx context.Context, dst interface{}) error
	RunQuery(ctx context.Context, kind string, filters map[string]interface{}, cursor string, limit int, orderBy string) ([]mdatastore.Key, string, bool, error)
	Exists(ctx context.Context, dst interface{}) (bool, error)
	Get(ctx context.Context, dst interface{}) error
	Put(ctx context.Context, src interface{}) error
	PutTx(tx *boom.Transaction, src interface{}) error
	DeleteTx(tx *boom.Transaction, src interface{}) error
}

type Transactor interface {
	RunInTransaction(context.Context, func(tx *boom.Transaction) error,) error
}


func (c *client) GetAll(ctx context.Context, kind string, dst interface{}) error {
	b := FromContext(ctx)
	q := b.NewQuery(kind)
	log.Println(kind)
	if _, err := b.GetAll(q, dst); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (c *client) RunQuery(
	ctx context.Context,
	kind string,
	filters map[string]interface{},
	cursor string,
	limit int,
	orderBy string) ([]mdatastore.Key, string, bool, error) {
	var keys []mdatastore.Key
	count := 0
	hasMore := false

	b := FromContext(ctx)
	q := b.NewQuery(kind).KeysOnly()

	if len(filters) > 0 {
		for k, v := range filters {
			q = q.Filter(k, v)
		}
	}
	if orderBy != "" {
		q = q.Order(orderBy)
	}

	if cursor != "" {
		dsCursor, err := b.Client.DecodeCursor(cursor)
		if err != nil {
			//TODO
			log.Printf("failed to decode %v", err)
		}
		q = q.Start(dsCursor)
	}

	if limit > 0 {
		log.Printf("limit %d", limit+1)
		// nextCursorがあるか把握する
		q = q.Limit(limit+1)
	} else {
		// TODO: add test case
		//limit = defaultLimit
	}

	it := b.Client.Run(ctx, q)

	var nextCursorStr string
	for {
		if key, err := it.Next(nil); err == iterator.Done {
			break
		} else if err != nil {
			log.Println(err)
			return keys, nextCursorStr, hasMore, errors.New("iterator error")
		} else {
			count++
			if limit < count {
				hasMore = true
				break
			}
			keys = append(keys, key)
			if limit == count {
				nextCursor, err := it.Cursor()
				if err != nil {
					return keys, nextCursor.String(), hasMore, err
				}
			}
		}
	}
	return keys, nextCursorStr, hasMore, nil
}

type parent struct {
	kind      string `boom:"kind,Audio"`
	ID        string `boom:"id"`
}

func (c *client) Get(ctx context.Context, dst interface{}) error {
	b := FromContext(ctx)

	if err := b.Get(dst); err != nil {
		if err == mdatastore.ErrNoSuchEntity {
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

func (c *client) Exists(ctx context.Context, dst interface{}) (bool, error) {
	b := FromContext(ctx)

	if err := b.Get(dst); err != nil {
		if err == datastore.ErrNoSuchEntity {
			return false, nil
		}
		return false, errors.WithStack(err)
	}

	return true, nil
}

func (c *client) Put(ctx context.Context, src interface{}) error {
	b := FromContext(ctx)
	if _, err := b.Put(src); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (c *client) PutTx(tx *boom.Transaction, src interface{}) error {
	if _, err := tx.Put(src); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (c *client) DeleteTx(tx *boom.Transaction, src interface{}) error {
	if err := tx.Delete(src); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (c *client) PutMulti(tx *boom.Transaction, src interface{}) error {
	if _, err := tx.PutMulti(src); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (c *client) DeleteMulti(tx *boom.Transaction, src interface{}) error {
	if err := tx.DeleteMulti(src); err != nil {
		return errors.WithStack(err)
	}

	return nil
}
