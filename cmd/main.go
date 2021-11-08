package main

import (
	"log"

	"github.com/hayashiki/audiy-api/src/domain/entity"
	"github.com/hayashiki/audiy-api/src/validator"
)

func main() {

	v := entity.User{
		Email: "hjrke",
		Name:  "1",
	}

	if err := validator.Validate(v); err != nil {
		if errString := validator.GetErrorMessages(err); err != nil {
			log.Println(errString)
		}
	}
}
