package graph

import (
	"context"
	"github.com/99designs/gqlgen/client"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/hayashiki/audiy-api/application/usecase"
	"github.com/hayashiki/audiy-api/infrastructure/ds"
	"github.com/hayashiki/audiy-api/interfaces/api/graph/generated"
	"github.com/hayashiki/audiy-api/interfaces/middleware"
	"os"
	"testing"
)

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
		client.Var("audio","F0240GUKN3A"),
		client.Var("body","something"),
	}

	var resp struct {
		CreateComment struct {
			User struct {
				ID        int64
				Email     string
				//CreatedAt time.Time
				//UpdatedAt time.Time
			}
			ID          int64
			Body        string
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
	actual, want := resp.CreateComment.Body, "something"; if actual != want {
		t.Errorf("want %v: actual: %v", want, actual)
	}
}



