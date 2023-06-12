package tts

import (
	"context"
	"fmt"
)

type Usecase struct{}

func NewUsecase() *Usecase {
	return &Usecase{}
}

func (u *Usecase) Do(_ context.Context, messageText string) {
	fmt.Println(messageText)
}
