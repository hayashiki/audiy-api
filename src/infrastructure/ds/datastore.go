package ds

import (
	"context"
	"errors"
	"log"

	"cloud.google.com/go/datastore"
	"google.golang.org/api/iterator"
)

var (
	datastoreClient *datastore.Client
)

type DataStore struct {
	Client *datastore.Client
}

func (d *DataStore) Get(ctx context.Context, entity EntitySpec) error {
	err := d.Client.Get(ctx, entity.GetKey(), entity)
	return err
}

func (d *DataStore) GetAll(
	ctx context.Context,
	opts *Query,
	generator func() interface{},
) ([]interface{}, string, error) {
	query := datastore.NewQuery(opts.Kind)
	if opts.Cursor != "" {
		dsCursor, err := datastore.DecodeCursor(opts.Cursor)
		if err != nil {
			//TODO
			log.Printf("failed to decode %v", err)
		}
		query = query.Start(dsCursor)
	}
	if opts.Limit > 0 {
		query = query.Limit(opts.Limit)
	}
	//for _, filter := range opts.Filters {
	//	query = query.Filter(filter.key, filter.value)
	//}
	for _, order := range opts.Order {
		query = query.Order(order)
	}

	log.Printf("query debug %+v", opts.Order)
	it := d.Client.Run(ctx, query)

	entities := make([]interface{}, 0)

	for {
		entity := generator()

		_, err := it.Next(entity)
		if errors.Is(err, iterator.Done) {
			break
		}
		if _, ok := err.(*datastore.ErrFieldMismatch); ok {
			entities = append(entities, entity)
			continue
		}
		if err != nil {
			return entities, "", err
		}
		entities = append(entities, entity)
	}

	cursor, err := it.Cursor()
	if err != nil {
		return entities, "", err
	}

	return entities, cursor.String(), nil
}

func (d *DataStore) GetAll2(
	ctx context.Context,
	q Query2,
	generator func() interface{},
) ([]interface{}, string, error) {
	v, ok := q.(*dsQuery)
	if !ok {
		return nil, "", errors.New("failed to build query")
	}

	it := d.Client.Run(ctx, v.Query)

	entities := make([]interface{}, 0)

	for {
		entity := generator()

		_, err := it.Next(entity)
		if errors.Is(err, iterator.Done) {
			break
		}
		if _, ok := err.(*datastore.ErrFieldMismatch); ok {
			entities = append(entities, entity)
			continue
		}
		if err != nil {
			return entities, "", err
		}
		entities = append(entities, entity)
	}

	cursor, err := it.Cursor()
	if err != nil {
		return entities, "", err
	}

	return entities, cursor.String(), nil
}

func (d *DataStore) Put(ctx context.Context, doc EntitySpec) (*datastore.Key, error) {
	key := doc.GetKey()
	//val := reflect.ValueOf(doc).Elem()
	//now := reflect.ValueOf(time.Now())
	//
	//if key.Incomplete() {
	//	val.FieldByName("CreatedAt").Set(now)
	//}
	//val.FieldByName("UpdatedAt").Set(now)

	key, err := d.Client.Put(ctx, key, doc)
	if err != nil {
		return nil, err
	}
	doc.SetID(key)

	return key, err
}

func (d *DataStore) Delete(ctx context.Context, doc EntitySpec) (*datastore.Key, error) {
	key := doc.GetKey()
	err := d.Client.Delete(ctx, key)
	if err != nil {
		return nil, err
	}
	doc.SetID(key)

	return key, err
}
