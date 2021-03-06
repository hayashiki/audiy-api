package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/hayashiki/audiy-api/src/domain/model"
	"github.com/hayashiki/audiy-api/src/graph/generated"
)

func (r *monologueElementResolver) Confidence(ctx context.Context, obj *model.MonologueElement) (string, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) CreateTranscript(ctx context.Context, input model.CreateTranscriptInput) (*model.Transcript, error) {
	err := r.transcriptUsecase.Do(ctx, input)
	en := &model.Transcript{}
	return en, err
}

func (r *transcriptResolver) Audio(ctx context.Context, obj *model.Transcript) (*model.Audio, error) {
	panic(fmt.Errorf("not implemented"))
}

// MonologueElement returns generated.MonologueElementResolver implementation.
func (r *Resolver) MonologueElement() generated.MonologueElementResolver {
	return &monologueElementResolver{r}
}

// Transcript returns generated.TranscriptResolver implementation.
func (r *Resolver) Transcript() generated.TranscriptResolver { return &transcriptResolver{r} }

type monologueElementResolver struct{ *Resolver }
type transcriptResolver struct{ *Resolver }
