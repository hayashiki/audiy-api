package audio_entity

import (
	"bytes"
	"context"
	"fmt"
	"github.com/google/go-cmp/cmp"
	"github.com/hayashiki/audiy-api/src/domain/model"
	"io/ioutil"
	"net/http"
	"os"
	"testing"
	"time"

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

func TestAudioCRUD(t *testing.T) {
	resetDatastoreEmulator()

	ctx := context.Background()
	audioRepo := repo{client: _dbClient}

	testID1 := "1"
	exists, err := audioRepo.Exists(ctx, testID1)
	if err != nil {
		t.Fatal(err)
	}
	if exists {
		t.Fatalf("should not exists this audio id: %s", testID1)
	}
	err = audioRepo.Put(ctx, testAudio1)
	if err != nil {
		t.Fatal(err)
	}
	exists, err = audioRepo.Exists(ctx, testID1)
	if err != nil {
		t.Fatal(err)
	}
	if !exists {
		t.Fatalf("should exists this audio id: %s", testID1)
	}
	got, err := audioRepo.Get(ctx, testID1)
	if err != nil {
		t.Fatal(err)
	}
	want := testAudio1; if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Audio value is mismatch (-got +want):\n%s", diff)
	}

	err = audioRepo.Delete(ctx, testID1)
	if err != nil {
		t.Fatal(err)
	}
	exists, err = audioRepo.Exists(ctx, testID1)
	if err != nil {
		t.Fatal(err)
	}
	if exists {
		t.Fatalf("should not exists this audio id: %s", testID1)
	}
}

func TestAudiosQuery(t *testing.T) {
	resetDatastoreEmulator()

	ctx := context.Background()
	audioRepo := repo{client: _dbClient}

	audios := []*model.Audio{
		testAudio1,
		testAudio2,
		testAudio3,
	}
	// TODO: PutMulti
	for _, a := range audios {
		err := audioRepo.Put(ctx, a)
		if err != nil {
			t.Fatal(err)
		}
	}
	t.Run("should get next cursor", func(t *testing.T) {
		// orderテスト and Limit test
		var count int
		limit := 1
		var cursor string

		for {
			audios, nextCursor, hasMore, err := audioRepo.GetAll(ctx, cursor, limit, "-PublishedAt")
			count++
			if err != nil {
				t.Fatal(err)
			}
			if len(audios) != 1 {
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

	t.Run(" PublishedAt sorted", func(t *testing.T) {
		got, _, _, _ := audioRepo.GetAll(ctx, "", 3, "PublishedAt")
		want := []*model.Audio{
			testAudio1,
			testAudio2,
			testAudio3,
		}
		if diff := cmp.Diff(got, want); diff != "" {
			t.Errorf("Audio value is mismatch (-got +want):\n%s", diff)
		}
	})

	t.Run(" -PublishedAt sorted", func(t *testing.T) {
		got, _, _, _ := audioRepo.GetAll(ctx, "", 3, "-PublishedAt")
		want := []*model.Audio{
			testAudio3,
			testAudio2,
			testAudio1,
		}
		cmp.Diff(got, want)
		if diff := cmp.Diff(got, want); diff != "" {
			t.Errorf("Audio value is mismatch (-got +want):\n%s", diff)
		}
	})
}

var (
	t1, _ = time.Parse(time.RFC3339, "2020-10-15T15:00:00+09:00")
	testAudio1 = &model.Audio{
		ID:     "1",
		Name:   "Test Audio 1",
		Length: 50.1,
		Mimetype: "mp4",
		PublishedAt: t1,
		CreatedAt:   t1,
		UpdatedAt:   t1,
		Transcribed: false,
		LikeCount: 0,
		PlayCount: 0,
		CommentCount: 0,
	}
	testAudio2 = &model.Audio{
		ID:     "2",
		Name:   "Test Audio 2",
		Length: 50.1,
		Mimetype: "mp4",
		PublishedAt: t1.Add(time.Hour * 10),
		CreatedAt:   t1.Add(time.Hour * 10),
		UpdatedAt:   t1.Add(time.Hour * 10),
		Transcribed: false,
	}
	testAudio3 = &model.Audio{
		ID:     "3",
		Name:   "Test Audio 3",
		Length: 51.1,
		Mimetype: "mp4",
		PublishedAt: t1.Add(time.Hour * 30),
		CreatedAt:   t1.Add(time.Hour * 30),
		UpdatedAt:   t1.Add(time.Hour * 30),
		Transcribed: false,
	}
)
