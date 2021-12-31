package model

import "cloud.google.com/go/datastore"

const defaultLimit = 20

type Query struct {
	Kind string
	Parent *datastore.Key
	Limit   int
	Cursor string
	Filters map[string]interface{}
	OrderBy string
	Namespace string
}

//func (q *Query) Limit() int {
//	if q.Limit > 0 {
//		return q.Limit
//	} else {
//		return defaultLimit
//	}
//}
