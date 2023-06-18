package props

import (
	"github.com/chun37/greenland-yomiage/general/internal/config"
	"github.com/chun37/greenland-yomiage/internal/usecase/dict"
)

type HandlerProps struct {
	Config *config.Config

	DictionaryAddUsecase *dict.AddUsecase
}
