package model

import "cloud.google.com/go/datastore"

func NewID(kind string) int64 {
	key := datastore.IncompleteKey(kind, nil)
	return key.ID
}
