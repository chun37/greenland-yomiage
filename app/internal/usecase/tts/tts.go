package tts

import (
	"fmt"
)

type Config struct {
	TargetChannelID string
}

type Usecase struct {
	cfg Config
}

func NewUsecase(cfg Config) *Usecase {
	return &Usecase{cfg: cfg}
}

func (u *Usecase) Do(messageText string) {
	fmt.Println(messageText)
}
