package graph

import (
	usecase2 "github.com/hayashiki/audiy-api/src/usecase"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	userUsecase    usecase2.UserUsecase
	audioUsecase   usecase2.AudioUsecase
	commentUsecase usecase2.CommentUsecase
	feedUseCase    usecase2.FeedUsecase
}

func NewResolver(
	userUsecase usecase2.UserUsecase,
	audioUsecase usecase2.AudioUsecase,
	commentUsecase usecase2.CommentUsecase,
	feedUseCase usecase2.FeedUsecase,
) *Resolver {
	return &Resolver{
		userUsecase:    userUsecase,
		audioUsecase:   audioUsecase,
		commentUsecase: commentUsecase,
		feedUseCase:    feedUseCase,
	}
}
