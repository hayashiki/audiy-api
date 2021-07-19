package graph

import (
	"context"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/99designs/gqlgen/client"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/hayashiki/audiy-api/application/usecase"
	"github.com/hayashiki/audiy-api/domain/entity"
	"github.com/hayashiki/audiy-api/infrastructure/ds"
	"github.com/hayashiki/audiy-api/interfaces/api/graph/auth"
	"github.com/hayashiki/audiy-api/interfaces/api/graph/generated"
	"github.com/hayashiki/audiy-api/interfaces/middleware"
)

func TestMockUsecase(t *testing.T) {
	dsCli, _ := ds.NewClient(context.Background(), os.Getenv("GCP_PROJECT"))
	commentRepo := ds.NewCommentRepository(dsCli)

	// usecase
	audioUsecase := usecase.NewAudioUsecase(nil)

	var testUserID int64 = 111111
	var testAudioID = "DummyAudioID"
	mockPlayUsecase := usecase.MockPlayUsecase{}
	mockPlayUsecase.ExistsFunc = func(userID string, audioID string) (bool, error) {
		if got, want := userID, testUserID; got != want {
			t.Errorf("userID: got %v, want %v", got, want)
		}

		if got, want := audioID, testAudioID; got != want {
			t.Errorf("audioID: got %v, want %v", got, want)
		}
		return true, nil
	}
	commentUsecase := usecase.NewCommentUsecase(commentRepo)

	obj := &entity.Audio{
		ID: testAudioID,
	}
	ctx := context.Background()
	ctx = auth.SetAuth(ctx, &auth.Auth{
		ID: testUserID,
	})

	userUsecase := usecase.MockUserUsecase{}
	userUsecase.SaveFunc = func(input entity.CreateUserInput) error {
		return nil
	}
	userUsecase.GetFunc = func(id string) (*entity.User, error) {
		idInt, _ := strconv.Atoi(id)
		return &entity.User{
			ID:        int64(idInt),
			Email:     "hayashiki@example.com",
			CreatedAt: time.Time{},
			UpdatedAt: time.Time{},
		}, nil
	}

	r := NewResolver(userUsecase, audioUsecase, mockPlayUsecase, commentUsecase)

	played, err := r.Audio().Played(ctx, obj)
	if err != nil {
		return
	}
	t.Log(played)
	if got, want := played, true; got != want {
		t.Errorf("played: got %v, want %v", got, want)
	}
}

func TestAudioCollection(t *testing.T) {

	dsCli, _ := ds.NewClient(context.Background(), os.Getenv("GCP_PROJECT"))
	commentRepo := ds.NewCommentRepository(dsCli)
	userRepo := ds.NewUserRepository(dsCli)

	authenticator := middleware.NewAuthenticator()
	// usecase
	audioUsecase := usecase.NewAudioUsecase(nil)
	playUsecase := usecase.NewPlayUsecase(nil)
	commentUsecase := usecase.NewCommentUsecase(commentRepo)
	userUsecase := usecase.NewUserUsecase(userRepo)

	gqlConfig := generated.Config{Resolvers: NewResolver(userUsecase, audioUsecase, playUsecase, commentUsecase)}
	testGqlServer := handler.NewDefaultServer(generated.NewExecutableSchema(gqlConfig))
	//var resp interface{}
	c := client.New(authenticator.AuthMiddleware(testGqlServer))

	options := []client.Option{
		client.Path("/query"),
		client.AddHeader("Authorization", "Bearer dummy"),
		client.Var("audio", "F0240GUKN3A"),
		client.Var("body", "something"),
	}

	var resp struct {
		CreateComment struct {
			User struct {
				ID    int64
				Email string
				//CreatedAt time.Time
				//UpdatedAt time.Time
			}
			ID   int64
			Body string
		}
	}

	err := c.Post(`
mutation($audio: ID!, $body: String!) {
  createComment(input: {audioID: $audio, body: $body}) {
    id
    body
  	user {
      id
      email
    }
  }
}
`, &resp, options...)
	if err != nil {
		t.Error(err)
	}
	actual, want := resp.CreateComment.Body, "something"
	if actual != want {
		t.Errorf("want %v: actual: %v", want, actual)
	}
}
