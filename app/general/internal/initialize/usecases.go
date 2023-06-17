package initialize

import (
	"github.com/chun37/greenland-yomiage/general/internal/config"
	"github.com/chun37/greenland-yomiage/internal/usecase/tts"
)

type Usecases struct {
	TTSUsecase *tts.Usecase
}

func NewUsecases(cfg config.Config) Usecases {
	uc := new(Usecases)
	uc.TTSUsecase = tts.NewUsecase(tts.Config{TargetChannelID: cfg.TargetChannelID})
	return *uc
}
