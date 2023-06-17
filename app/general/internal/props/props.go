package props

import (
	"github.com/chun37/greenland-yomiage/internal/usecase/tts"
)

type HandlerProps struct {
	TTSUsecase *tts.Usecase
}
