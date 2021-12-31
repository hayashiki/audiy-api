package user_entity

import (
	"bytes"
	"context"
	"fmt"
	"github.com/google/go-cmp/cmp"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"testing"
	"time"

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

func shutdownDatastoreEmulator() error {
	if addr := os.Getenv("DATASTORE_EMULATOR_HOST"); addr != "" {
		resp, err := http.Post("http://"+addr+"/shutdown", "application/json", nil)
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Errorf("invalid response: %s", err.Error())
		}

		fmt.Errorf("datastore emulator: %s", string(body))
	}
	return nil
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

func TestUserCRUD(t *testing.T) {
	ctx := context.Background()
	userRepo := repo{client: _dbClient}
	if err := userRepo.Put(ctx, testUser1); err != nil {
		t.Fatal(err)
	}

	exists, err := userRepo.Exists(ctx, testUser1.ID)
	if err != nil {
		t.Fatal(err)
	}
	if !exists {
		t.Fatal("should not exists user data")
	}
	err = userRepo.Put(ctx, testUser1)
	if err != nil {
		t.Fatal(err)
	}
	exists, err = userRepo.Exists(ctx, testUser1.ID)
	if err != nil {
		t.Fatal(err)
	}
	if !exists {
		t.Fatalf("should exists this audio id: %s", testUser1.ID)
	}
	got, err := userRepo.Get(ctx, testUser1.ID)
	if err != nil {
		t.Fatal(err)
	}
	want := testUser1; if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Audio value is mismatch (-got +want):\n%s", diff)
	}

	err = userRepo.Delete(ctx, testUser1.ID)
	if err != nil {
		t.Fatal(err)
	}
	exists, err = userRepo.Exists(ctx, testUser1.ID)
	if err != nil {
		t.Fatal(err)
	}
	if exists {
		t.Fatalf("should not exists this audio id: %s", testUser1.ID)
	}	

	users, err := userRepo.GetAll(ctx)
	for _, u := range users {
		log.Println(u.ID)
		userRepo.Put(ctx, u)
	}
}

var (
	testUser1 = &model.User{
		ID:        "1",
		Email:     "testuseremail1@example.com",
		Name:      "Test User Name 1",
		PhotoURL:  "http://example.com/profile.png",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
)
