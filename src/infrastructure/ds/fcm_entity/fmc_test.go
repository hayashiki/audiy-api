package fcm_entity

import (
	"context"
	"github.com/hayashiki/audiy-api/src/infrastructure/ds"
	"log"

	"go.mercari.io/datastore/boom"
	"testing"
)

func TestFromContext(t *testing.T) {
	ctx := context.Background()
	repo := ds.New()
	tran := ds.NewDatastoreTransactor()
	err := tran.RunInTransaction(ctx, func(tx *boom.Transaction) error {
		entity := NewEntity("111", "222", "333")
		err := repo.Put(tx, entity)
		log.Printf("%+v", entity)
		return err
	})
	if err != nil {
		t.Error(err)
	}
}
