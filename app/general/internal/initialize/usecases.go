package initialize

import (
	"github.com/chun37/greenland-yomiage/internal/usecase/tts"
	"github.com/chun37/greenland-yomiage/internal/wavgenerator"
)

type Usecases struct {
	TTSUsecase *tts.Usecase
}

func NewUsecases() Usecases {
	uc := new(Usecases)
	uc.TTSUsecase = tts.NewUsecase(tts.Dependencies{wavgenerator.NewVoiceVox()})
	return *uc
}
