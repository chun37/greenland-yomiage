package tts

import (
	"bytes"
	"log"

	"github.com/chun37/greenland-yomiage/internal/opus"
	"github.com/chun37/greenland-yomiage/internal/wavgenerator"
	"golang.org/x/xerrors"
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
	Speaking   <-chan bool
}

func (u *Usecase) Do(param UsecaseParam) error {
	wav, err := u.deps.WavGenerator.Generate(param.Text)
	if err != nil {
		return xerrors.Errorf("failed to generate wav: %w", err)
	}

	go func() {
		err := opus.Encode(bytes.NewReader(wav), param.OpusChunks, param.Done, param.Speaking)
		if err != nil {
			log.Println("failed to encode audio:", err)
			return
		}
	}()

	return nil
}
