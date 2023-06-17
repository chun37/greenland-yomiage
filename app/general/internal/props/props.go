package props

import (
	"github.com/chun37/greenland-yomiage/general/internal/config"
	"github.com/chun37/greenland-yomiage/internal/usecase/tts"
)

type HandlerProps struct {
	Config     *config.Config
	TTSUsecase *tts.Usecase
}
