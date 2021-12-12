package transcript_entity

import (
	"context"
	"github.com/hayashiki/audiy-api/src/infrastructure/datastore"
	"log"
	"os"
	"testing"
)

func TestGetAll(t *testing.T) {
	log.Println(os.Getenv("GCP_PROJECT"))
	ctx := context.Background()
	ds := datastore.New()
	r := repo{client: ds}
	trans, err := r.GetAll(ctx)
	if err != nil {
		t.Error(err)
	}

	for _, t := range trans {
		log.Println(t.ID)
	}
}
