package ds

import "cloud.google.com/go/datastore"

// NewKey creates a new datastore key
func (d *DataStore) NewKey(kind string) *datastore.Key {
	return datastore.IncompleteKey(kind, nil)
}

// Key builds a new datastore key based on provided id
func (d *DataStore) Key(kind string, id string) *datastore.Key {
	return datastore.NameKey(kind, id, nil)
}

// DecodeKey decode datastore key based on provided id
func (d *DataStore) DecodeKey(id string) *datastore.Key {
	key, err := datastore.DecodeKey(id)
	if err != nil {
		//	 TODO
	}
	return key
}