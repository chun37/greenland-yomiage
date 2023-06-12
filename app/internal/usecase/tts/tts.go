package tts

import (
	"fmt"
)

type Usecase struct{}

func NewUsecase() *Usecase {
	return &Usecase{}
}

func (u *Usecase) Do(messageText string) {
	fmt.Println(messageText)
}
