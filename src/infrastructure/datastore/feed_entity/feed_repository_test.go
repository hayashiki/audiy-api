package feed_entity

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/google/go-cmp/cmp"
	"github.com/hayashiki/audiy-api/src/domain/model"
	"github.com/hayashiki/audiy-api/src/infrastructure/datastore"
	"github.com/hayashiki/audiy-api/src/infrastructure/datastore/audio_entity"
	"github.com/hayashiki/audiy-api/src/infrastructure/datastore/user_entity"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"testing"
	"time"
)

var _dbClient datastore.DSClient

func TestMain(m *testing.M) {
	os.Setenv("DATASTORE_HOST", "localhost:8081")
	os.Setenv("DATASTORE_PROJECT_ID", "local")
	os.Setenv("DATASTORE_EMULATOR_HOST", "localhost:8081")
	os.Setenv("DATASTORE_EMULATOR_HOST_PATH", "localhost:8081/audiy")
	_dbClient = datastore.NewDS()
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


func TestFeedCRUD(t *testing.T) {
	resetDatastoreEmulator()
	ctx := context.Background()
	feedRepo := repo{client: _dbClient}

	exists, err := feedRepo.Exists(ctx, testFeed1.UserID, string(testFeed1.ID()))
	if err != nil {
		t.Fatal(err)
	}
	if exists {
		t.Fatalf("should exists this feed id: %s", string(testFeed1.ID()))
	}
	err = feedRepo.Put(ctx, "1", testFeed1)
	if err != nil {
		t.Error(err)
	}

	exists, err = feedRepo.Exists(ctx, testFeed1.UserID, string(testFeed1.ID()))
	if err != nil {
		t.Fatal(err)
	}
	if !exists {
		t.Fatalf("should exists this feed id: %s", string(testFeed1.ID()))
	}

	got, err := feedRepo.Get(ctx, "1-1", "1")
	if err != nil {
		t.Fatal(err)
	}

	want := testFeed1; if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Feed value is mismatch (-got +want):\n%s", diff)
	}
	err = feedRepo.Delete(ctx, testFeed1.UserID, string(testFeed1.ID()))
	if err != nil {
		t.Fatal(err)
	}
	exists, err = feedRepo.Exists(ctx, testFeed1.UserID, string(testFeed1.ID()))
	if err != nil && errors.Is(datastore.ErrNoSuchEntity, err){
		t.Fatal(err)
	}
	if exists {
		t.Fatalf("should not exists this audio id: %s", feedRepo)
	}
}

func TestGetAll(t *testing.T) {
	resetDatastoreEmulator()
	ctx := context.Background()
	feedRepo := repo{client: _dbClient}

	// TODO: DRY, PutMulti
	if err := feedRepo.Put(ctx, testFeed1.UserID, testFeed1); err != nil {
		t.Error(err)
	}
	if err := feedRepo.Put(ctx, testFeed2.UserID, testFeed2); err != nil {
		t.Error(err)
	}
	if err := feedRepo.Put(ctx, testFeed3.UserID, testFeed3); err != nil {
		t.Error(err)
	}
	if err := feedRepo.Put(ctx, testFeed4.UserID, testFeed4); err != nil {
		t.Error(err)
	}
	t.Run("should get next cursor", func(t *testing.T) {
		var count int
		limit := 1
		var cursor string

		for {
			feeds, nextCursor, hasMore, err := feedRepo.GetAll(ctx, testFeed1.UserID, nil,cursor, limit, "PublishedAt")
			count++
			if err != nil {
				t.Fatal(err)
			}
			if len(feeds) != 1 {
				t.Fatal("should get 1 record")
			}
			if !hasMore {
				break
			}
			cursor = nextCursor
		}

		if count != 3 {
			t.Fatalf("failed to get next cursor, total record count %d", count)
		}
	})


	t.Run("PublishedAt sorted", func(t *testing.T) {
		got, _, _, err := feedRepo.GetAll(ctx, testFeed1.UserID, nil, "", 3, "PublishedAt")
		if err != nil {
			t.Fatal(err)
		}
		want := []*model.Feed{
			testFeed1,
			testFeed3,
			testFeed4,
		}
		if diff := cmp.Diff(got, want); diff != "" {
			t.Errorf("Feed value is mismatch (-got +want):\n%s", diff)
		}
	})

	t.Run("-PublishedAt sorted", func(t *testing.T) {
		got, _, _, err := feedRepo.GetAll(ctx, testFeed1.UserID, nil, "", 3, "-PublishedAt")
		if err != nil {
			t.Fatal(err)
		}
		want := []*model.Feed{
			testFeed4,
			testFeed3,
			testFeed1,
		}
		if diff := cmp.Diff(got, want); diff != "" {
			t.Errorf("Feed value is mismatch (-got +want):\n%s", diff)
		}
	})
}

func TestPutMulti(t *testing.T)  {
	ctx := context.Background()
	ds := datastore.New()
	ds2 := datastore.NewDS()
	feedRepo := repo{client: ds2}

	userRepo := user_entity.NewUserRepository(ds)
	audioRepo := audio_entity.NewAudioRepository(ds)

	users, err := userRepo.GetAll(ctx)
	if err != nil {
		t.Error(err)
	}
	log.Println(users)

	audio, err := audioRepo.Get(ctx, "F02QES84U0H")
	if err != nil {
		t.Error(err)
	}

	feeds := make([]*model.Feed, len(users))
	for i, u := range users {
		log.Println(u.Name)
		newFeed := model.NewFeed(audio.ID, u.ID, audio.PublishedAt)
		feeds[i] = newFeed
	}
	err = feedRepo.PutMulti(ctx, feeds)
}

func TestUserCreated(t *testing.T) {
	ctx := context.Background()
	ds := datastore.New()
	ds2 := datastore.NewDS()
	feedRepo := repo{client: ds2}
	userRepo := user_entity.NewUserRepository(ds)
	audioRepo := audio_entity.NewAudioRepository(ds)
	user, err := userRepo.Get(ctx, "103843140833205663533")
	if err != nil {
		t.Error(err)
	}

	audios, _, _, _ := audioRepo.GetAll(ctx, "", 1000, "-PublishedAt")
	feeds := make([]*model.Feed, len(audios))

	for i, a := range audios {
		newFeed := model.NewFeed(a.ID, user.ID, a.PublishedAt)
		feeds[i] = newFeed
		//userIDs[i] = newUser.ID
	}
	err = feedRepo.PutMulti(ctx, feeds)
	if err != nil {
		t.Error(err)
	}
}

var (
	t1, _     = time.Parse(time.RFC3339, "2020-10-15T15:00:00+09:00")
	testFeed1 = &model.Feed{
		AudioID:     "1",
		UserID:      "1",
		Played:      false,
		Liked:       false,
		Stared:      false,
		StartTime:   0,
		PublishedAt: t1,
		CreatedAt:   t1,
		UpdatedAt:   t1,
	}
	testFeed2 = &model.Feed{
		AudioID:     "1",
		UserID:      "2",
		Played:      false,
		Liked:       false,
		Stared:      false,
		StartTime:   0,
		PublishedAt: t1.Add(time.Hour * 10),
		CreatedAt:   t1.Add(time.Hour * 10),
		UpdatedAt:   t1.Add(time.Hour * 10),
	}
	testFeed3 = &model.Feed{
		AudioID:     "2",
		UserID:      "1",
		Played:      false,
		Liked:       false,
		Stared:      false,
		StartTime:   0,
		PublishedAt: t1.Add(time.Hour * 30),
		CreatedAt:   t1.Add(time.Hour * 30),
		UpdatedAt:   t1.Add(time.Hour * 30),
	}
	testFeed4 = &model.Feed{
		AudioID:     "3",
		UserID:      "1",
		Played:      false,
		Liked:       false,
		Stared:      false,
		StartTime:   0,
		PublishedAt: t1.Add(time.Hour * 40),
		CreatedAt:   t1.Add(time.Hour * 40),
		UpdatedAt:   t1.Add(time.Hour * 40),
	}
)
