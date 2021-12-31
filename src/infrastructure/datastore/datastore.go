package datastore

import (
	"context"
	"github.com/hayashiki/audiy-api/src/config"
	"github.com/hayashiki/audiy-api/src/domain/model"
	mdatastore "go.mercari.io/datastore"
	"go.mercari.io/datastore/boom"

	"log"

	"cloud.google.com/go/datastore"
	"github.com/pkg/errors"
	"google.golang.org/api/iterator"
)

type dsClient struct{
	namespace string
	ancestor *datastore.Key
}

type datastoreDSTransactor struct {
}

func NewDSDatastoreTransactor() Transactor {
	return &datastoreDSTransactor{}
}

func (t *datastoreDSTransactor) RunInTransaction(ctx context.Context, fn func(tx *boom.Transaction) error) error {
	if _, err := FromContext(ctx).RunInTransaction(fn); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func NewDS() DSClient {
	cli := &dsClient{}
	return cli
}

//
type DSClient interface {
	GetAll(ctx context.Context, kind string, dst interface{}) error
	GetMulti(ctx context.Context, keys []*datastore.Key, dst interface{}) error
	Run(ctx context.Context, query *model.Query) ([]*datastore.Key, string, bool, error)
	Exists(ctx context.Context, key *datastore.Key, dst interface{}) (bool, error)
	Get(ctx context.Context, key *datastore.Key, dst interface{}) error
	Put(ctx context.Context, key *datastore.Key, src interface{}) error
	PutTx(tx *datastore.Transaction, key *datastore.Key, src interface{}) error
	DeleteTx(tx *datastore.Transaction, key *datastore.Key) error
	PutMulti(ctx context.Context, keys []*datastore.Key, dst interface{}) error
	Delete(ctx context.Context, key *datastore.Key) error
}

type DSTransactor interface {
	RunInTransaction(context.Context, func(tx *datastore.Transaction) error) error
}

func (c *dsClient) GetAll(ctx context.Context, kind string, dst interface{}) error {
	b, err := NewClient(ctx, config.GetProject())
	if err != nil {
		panic(err)
	}

	q := datastore.NewQuery(kind) //.Namespace()
	if _, err := b.GetAll(ctx, q, dst); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (c *dsClient) Run(ctx context.Context, query *model.Query) ([]*datastore.Key, string, bool, error) {
	var keys []*datastore.Key
	count := 0
	hasMore := false
	b, err := NewClient(ctx, config.GetProject())
	if err != nil {
		panic(err)
	}

	q := datastore.NewQuery(query.Kind).KeysOnly()

	if query.Parent != nil {
		q = q.Ancestor(query.Parent)
	}

	if len(query.Filters) > 0 {
		for k, v := range query.Filters {
			q = q.Filter(k, v)
		}
	}
	if query.OrderBy != "" {
		q = q.Order(query.OrderBy)
	}
	if query.Namespace != "" {
		q = q.Namespace(query.Namespace)
	}
	if query.Cursor != "" {
		dsCursor, err := datastore.DecodeCursor(query.Cursor)
		if err != nil {
			//TODO
			log.Printf("failed to decode %v", err)
		}
		q = q.Start(dsCursor)
	}

	if query.Limit > 0 {
		q = q.Limit(query.Limit+1)
	} else {
		// TODO: add test case
		//limit = defaultLimit
	}

	it := b.Run(ctx, q)

	var nextCursorStr string
	for {
		if key, err := it.Next(nil); err == iterator.Done {
			break
		} else if err != nil {
			log.Println(err)
			return keys, nextCursorStr, hasMore, errors.New("iterator error")
		} else {
			count++
			if query.Limit < count {
				hasMore = true
				break
			}
			keys = append(keys, key)
			if query.Limit == count {
				nextCursor, err := it.Cursor()
				nextCursorStr = nextCursor.String()
				if err != nil {
					return keys, nextCursor.String(), hasMore, err
				}
			}
		}
	}
	it.Cursor()
	return keys, nextCursorStr, hasMore, nil
}

func (c *dsClient) Get(ctx context.Context, key *datastore.Key, dst interface{}) error {
	b, err := NewClient(ctx, config.GetProject())
	if err != nil {
		panic(err)
	}

	if err := b.Get(ctx, key, dst); err != nil {
		if err == datastore.ErrNoSuchEntity {
			log.Print(dst)
			return errors.WithStack(err)
		}
		return errors.WithStack(err)
	}

	return nil
}

func (c *dsClient) GetMulti(ctx context.Context, keys []*datastore.Key, dst interface{}) error {
	b, err := NewClient(ctx, config.GetProject())
	if err != nil {
		return nil
	}

	if err := b.GetMulti(ctx, keys, dst); err != nil {
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

func (c *dsClient) Exists(ctx context.Context, key *datastore.Key, dst interface{}) (bool, error) {
	b, err := NewClient(ctx, config.GetProject())
	if err != nil {
		panic(err)
	}

	if err := b.Get(ctx, key, dst); err != nil {
		if err == mdatastore.ErrNoSuchEntity {
			return false, nil
		}
		return false, errors.WithStack(err)
	}

	return true, nil
}

func (c *dsClient) Put(ctx context.Context, key *datastore.Key, src interface{}) error {
	b, err := NewClient(ctx, config.GetProject())
	if err != nil {
		panic(err)
	}
	if _, err := b.Put(ctx, key, src); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (c *dsClient) PutTx(tx *datastore.Transaction, key *datastore.Key, src interface{}) error {
	if _, err := tx.Put(key, src); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (c *dsClient) DeleteTx(tx *datastore.Transaction, key *datastore.Key) error {
	if err := tx.Delete(key); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (c *dsClient) PutMulti(ctx context.Context, keys []*datastore.Key, dst interface{}) error {
	b, err := NewClient(ctx, config.GetProject())
	if err != nil {
		panic(err)
	}

	if _, err := b.PutMulti(ctx, keys, dst); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (c *dsClient) PutMultiTx(tx *datastore.Transaction, keys []*datastore.Key, src interface{}) error {
	if _, err := tx.PutMulti(keys, src); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (c *dsClient) DeleteMulti(tx *datastore.Transaction, keys []*datastore.Key) error {
	if err := tx.DeleteMulti(keys); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (c *dsClient) Delete(ctx context.Context, key *datastore.Key) error {
	b, err := NewClient(ctx, config.GetProject())
	if err != nil {
		return err
	}
	if err := b.Delete(ctx, key); err != nil {
		return errors.WithStack(err)
	}
	return nil
}
