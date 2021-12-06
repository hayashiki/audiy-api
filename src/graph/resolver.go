package graph

import (
	"github.com/hayashiki/audiy-api/src/usecase"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	userUsecase       usecase.UserUsecase
	audioUsecase      usecase.AudioUsecase
	commentUsecase    usecase.CommentUsecase
	feedUseCase       usecase.FeedUsecase
	transcriptUsecase usecase.TranscriptAudioUsecase
}

func NewResolver(
	userUsecase usecase.UserUsecase,
	audioUsecase usecase.AudioUsecase,
	commentUsecase usecase.CommentUsecase,
	feedUseCase usecase.FeedUsecase,
	transcriptUsecase usecase.TranscriptAudioUsecase,
) *Resolver {
	return &Resolver{
		userUsecase:    userUsecase,
		audioUsecase:   audioUsecase,
		commentUsecase: commentUsecase,
		feedUseCase:    feedUseCase,
		transcriptUsecase: transcriptUsecase,
	}
}
