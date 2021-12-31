package comment_entity

import (
	"bytes"
	"context"
	"fmt"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/hayashiki/audiy-api/src/domain/model"
	"github.com/hayashiki/audiy-api/src/infrastructure/datastore"
)

var _dbClient datastore.Client

func TestMain(m *testing.M) {
	os.Setenv("DATASTORE_HOST", "localhost:8081")
	os.Setenv("DATASTORE_PROJECT_ID", "local")
	os.Setenv("DATASTORE_EMULATOR_HOST", "localhost:8081")
	os.Setenv("DATASTORE_EMULATOR_HOST_PATH", "localhost:8081/audiy")
	_dbClient = datastore.New()
	m.Run()
}

func resetDatastoreEmulator() error {
	if addr := os.Getenv("DATASTORE_EMULATOR_HOST"); addr != "" {

		var buf bytes.Buffer
		resp, err := http.Post("http://"+addr+"/reset", "application/json", &buf)
		if err != nil {
			return fmt.Errorf("unable to reset datastore emulator: %s", err.Error())
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Errorf("invalid response: %s", err.Error())
		}

		fmt.Errorf("datastore emulator: %s", string(body))
	}
	return nil
}

func TestCommentCRUD(t *testing.T) {
	resetDatastoreEmulator()

	ctx := context.Background()
	commentRepo := repo{client: _dbClient}

	err := commentRepo.Put(ctx, testComment1)
	if err != nil {
		t.Error(err)
	}

	got, err := commentRepo.Get(ctx, testComment1.ID)
	if err != nil {
		t.Error(err)
	}

	want := testComment1; if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Comment value is mismatch (-got +want):\n%s", diff)
	}

	if err := commentRepo.Delete(ctx, testComment1.ID); err != nil {
		t.Fatal(err)
	}
	_, err = commentRepo.Get(ctx, testComment1.ID)
	if err != nil && !errors.Is(err, datastore.ErrNoSuchEntity){
		t.Error(err)
	}
}

func TestQuery(t *testing.T)  {
	resetDatastoreEmulator()

	ctx := context.Background()
	commentRepo := repo{client: _dbClient}
	err := commentRepo.Put(ctx, testComment1)
	if err != nil {
		t.Error(err)
	}
	err = commentRepo.Put(ctx, testComment2)
	if err != nil {
		t.Error(err)
	}

	t.Run("audio filtered", func(t *testing.T) {
		comments, _, _, err := commentRepo.GetAllByAudio(ctx, "1", "", 10, "CreatedAt")
		if err != nil {
			t.Error(err)
		}
		if len(comments) != 1 {
			t.Errorf("want: 1, but got record count %d", len(comments))
		}
	})
}

var (
	t1, _        = time.Parse(time.RFC3339, "2020-10-15T15:00:00+09:00")
	testComment1 = &model.Comment{
		ID:        1,
		AudioID:   "1",
		UserID:    "1",
		Body:      "Test Body 1",
		CreatedAt: t1,
		UpdatedAt: t1,
	}
	testComment2 = &model.Comment{
		ID:        2,
		AudioID:   "2",
		UserID:    "1",
		Body:      "Test Body 2",
		CreatedAt: t1,
		UpdatedAt: t1,
	}

)
