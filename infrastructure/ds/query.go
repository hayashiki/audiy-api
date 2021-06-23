package ds

import (
	"fmt"

	"cloud.google.com/go/datastore"
)

type Query struct {
	Kind    string
	Filters []*Filter
	Offset  int
	Cursor  string
	Limit   int
	Order   []string
}

type Filter struct {
	key   string
	value interface{}
}

func (f *Filter) String() string {
	return fmt.Sprintf("%s %s", f.key, f.value)
}

type EntitySpec interface {
	GetKey() *datastore.Key
	SetID(id *datastore.Key)
}

func NewCursorQuery(kind string, fs []*Filter, limit int, cursor string, sort ...string) *Query {
	return &Query{
		Kind:    kind,
		Filters: fs,
		Cursor:  cursor,
		Limit:   limit,
		Order:   sort,
	}
}
