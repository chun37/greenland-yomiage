package initialize

import (
	"github.com/chun37/greenland-yomiage/internal/usecase/tts"
)

type Usecases struct {
	TTSUsecase *tts.Usecase
}

func NewUsecases(dependencies *ExternalDependencies) Usecases {
	uc := new(Usecases)
	uc.TTSUsecase = tts.New(tts.Dependencies{dependencies.VoiceVox})
	return *uc
}
