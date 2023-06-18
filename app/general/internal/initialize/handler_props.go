package initialize

import (
	"github.com/chun37/greenland-yomiage/general/internal/config"
	"github.com/chun37/greenland-yomiage/general/internal/props"
)

func NewHandlerProps(cfg config.Config, usecases Usecases) *props.HandlerProps {
	hp := &props.HandlerProps{
		Config:               &cfg,
		DictionaryAddUsecase: usecases.DictAddUsecase,
	}
	return hp
}
