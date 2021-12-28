package fcm_entity

import (
	"context"
	"github.com/hayashiki/audiy-api/src/infrastructure/datastore"
	"log"

	"go.mercari.io/datastore/boom"
	"testing"
)

func TestPut(t *testing.T) {
	ctx := context.Background()
	repo := datastore.New()
	tran := datastore.NewDatastoreTransactor()
	err := tran.RunInTransaction(ctx, func(tx *boom.Transaction) error {
		entity := NewEntity("888", "222", "333")
		err := repo.PutTx(tx, entity)
		log.Printf("%+v", entity)
		return err
	})
	if err != nil {
		t.Error(err)
	}
}

func TestGetAll(t *testing.T) {
	ctx := context.Background()
	ds := datastore.New()
	fcmRepo := repo{client: ds}

	fcms, nextCursor, hasMore, err := fcmRepo.GetAll(ctx, "", 100, "CreatedAt")
	if err != nil {
		t.Error(err)
	}
	log.Println(fcms[0].ID, nextCursor, hasMore, err)
}



