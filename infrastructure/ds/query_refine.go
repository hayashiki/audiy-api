package ds

import (
	"cloud.google.com/go/datastore"
	"log"
)

type (
	Query2 interface {
		Filter(string, interface{}) Query2
		Limit(value int) Query2
		Cursor(cursor string) Query2
		KeysOnly() Query2
		Order(orders []string) Query2
	}

	dsQuery struct {
		*datastore.Query
	}
)

// NewQuery returns PersistenceQuery wraps datastore.Query
func NewQuery(kind string) Query2 {
	return &dsQuery{
		Query: datastore.NewQuery(kind),
	}
}

func (q dsQuery) Filter(filter string, value interface{}) Query2 {
	q.Query = q.Query.Filter(filter, value)
	return q
}

func (q dsQuery) Limit(value int) Query2 {
	q.Query = q.Query.Limit(value)
	return q
}

func (q dsQuery) Cursor(cursor string) Query2 {
	dsCursor, err := datastore.DecodeCursor(cursor)
	if err != nil {
		//TODO
		log.Printf("failed to decode %v", err)
	}
	q.Query = q.Query.Start(dsCursor)
	return q
}

// KeysOnly sets the query to return only keys
func (q dsQuery) KeysOnly() Query2 {
	q.Query = q.Query.KeysOnly()
	return q
}

func (q dsQuery) Order(orders []string) Query2 {
	for _, order := range orders {
		q.Query = q.Query.Order(order)
	}
	return q
}
