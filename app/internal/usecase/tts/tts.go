package tts

import (
	"bytes"
	"log"

	"github.com/chun37/greenland-yomiage/internal/opus"
	"github.com/chun37/greenland-yomiage/internal/wavgenerator"
)

type Dependencies struct {
	WavGenerator wavgenerator.Service
}

type Usecase struct {
	deps Dependencies
}

func NewUsecase(deps Dependencies) *Usecase {
	return &Usecase{deps: deps}
}

type UsecaseParam struct {
	Text       string
	OpusChunks chan []byte
	Done       chan struct{}
}

func (u *Usecase) Do(param UsecaseParam) {
	wav, err := u.deps.WavGenerator.Generate(param.Text)
	if err != nil {
		log.Println("failed to generate wav:", err)
		return
	}

	go func() {
		err := opus.Encode(bytes.NewReader(wav), param.OpusChunks, param.Done)
		if err != nil {
			log.Println("failed to encode audio:", err)
			return
		}
	}()
}
