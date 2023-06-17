package tts

import (
	"fmt"

	"github.com/chun37/greenland-yomiage/internal/wavgenerator"
)

type Config struct {
	TargetChannelID string
}

type Dependencies struct {
	WavGenerator wavgenerator.Service
}

type Usecase struct {
	cfg  Config
	deps Dependencies
}

func NewUsecase(cfg Config, deps Dependencies) *Usecase {
	return &Usecase{cfg: cfg, deps: deps}
}

func (u *Usecase) Do(messageText string) {
	//u.deps.WavGenerator.Generate()
	fmt.Println(messageText)
}
