package initialize

import (
	"github.com/chun37/greenland-yomiage/internal/usecase/dict"
	"github.com/chun37/greenland-yomiage/internal/usecase/tts"
)

type Usecases struct {
	TTSUsecase     *tts.Usecase
	DictAddUsecase *dict.AddUsecase
}

func NewUsecases(dependencies *ExternalDependencies) Usecases {
	uc := new(Usecases)
	uc.TTSUsecase = tts.NewUsecase(tts.Dependencies{dependencies.VoiceVox})
	uc.DictAddUsecase = dict.NewAddUsecase(dict.Dependencies{dependencies.VoiceVox})
	return *uc
}
