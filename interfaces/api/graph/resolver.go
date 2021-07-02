package graph

import "github.com/hayashiki/audiy-api/application/usecase"

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	audioUsecase     usecase.AudioUsecase
	audioUserUsecase usecase.AudioUserUsecase
	starUsecase      usecase.StarUsecase
	likeUsecase      usecase.LikeUsecase
	commentUsecase   usecase.CommentUsecase
}

func NewResolver(
	audioUsecase usecase.AudioUsecase,
	audioUserUsecase usecase.AudioUserUsecase,
	starUsecase usecase.StarUsecase,
	likeUsecase usecase.LikeUsecase,
	commentUsecase usecase.CommentUsecase,
) *Resolver {
	return &Resolver{
		audioUsecase:     audioUsecase,
		audioUserUsecase: audioUserUsecase,
		starUsecase:      starUsecase,
		likeUsecase:      likeUsecase,
		commentUsecase:   commentUsecase,
	}
}
