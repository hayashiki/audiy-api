package graph

import "github.com/hayashiki/audiy-api/application/usecase"

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	userUsecase    usecase.UserUsecase
	audioUsecase   usecase.AudioUsecase
	playUsecase    usecase.PlayUsecase
	starUsecase    usecase.StarUsecase
	likeUsecase    usecase.LikeUsecase
	commentUsecase usecase.CommentUsecase
}

func NewResolver(
	userUsecase    usecase.UserUsecase,
	audioUsecase usecase.AudioUsecase,
	audioUserUsecase usecase.PlayUsecase,
	//starUsecase usecase.StarUsecase,
	//likeUsecase usecase.LikeUsecase,
	commentUsecase usecase.CommentUsecase,
) *Resolver {
	return &Resolver{
		userUsecase: userUsecase,
		audioUsecase: audioUsecase,
		playUsecase:  audioUserUsecase,
		//starUsecase:      starUsecase,
		//likeUsecase:      likeUsecase,
		commentUsecase:   commentUsecase,
	}
}
