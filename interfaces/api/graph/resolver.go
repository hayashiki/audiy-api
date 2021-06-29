package graph

import "github.com/hayashiki/audiy-api/application/usecase"

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct{
	audioUsecase usecase.AudioUsecase
	audioUserUsecase usecase.AudioUserUsecase
}

func NewResolver(audioUsecase usecase.AudioUsecase, audioUserUsecase usecase.AudioUserUsecase) *Resolver {
	return &Resolver{
		audioUsecase: audioUsecase,
		audioUserUsecase: audioUserUsecase,
	}
}
