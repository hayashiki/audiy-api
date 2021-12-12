package graph

import (
	"github.com/hayashiki/audiy-api/src/graph/dataloaders"
	"github.com/hayashiki/audiy-api/src/usecase"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	dataLoaders       dataloaders.DataLoaderService
	userUsecase       usecase.UserUsecase
	audioUsecase      usecase.AudioUsecase
	commentUsecase    usecase.CommentUsecase
	feedUseCase       usecase.FeedUsecase
	transcriptUsecase usecase.TranscriptAudioUsecase
	fcmUsecase        usecase.FcmUsecase
}

func NewResolver(
	dataLoaders dataloaders.DataLoaderService,
	userUsecase usecase.UserUsecase,
	audioUsecase usecase.AudioUsecase,
	commentUsecase usecase.CommentUsecase,
	feedUseCase usecase.FeedUsecase,
	transcriptUsecase usecase.TranscriptAudioUsecase,
	fcmUsecase usecase.FcmUsecase,
) *Resolver {
	return &Resolver{
		dataLoaders:       dataLoaders,
		userUsecase:       userUsecase,
		audioUsecase:      audioUsecase,
		commentUsecase:    commentUsecase,
		feedUseCase:       feedUseCase,
		transcriptUsecase: transcriptUsecase,
		fcmUsecase:        fcmUsecase,
	}
}
