package initialize

import "github.com/chun37/greenland-yomiage/general/internal/props"

func NewHandlerProps(usecases Usecases) *props.HandlerProps {
	hp := &props.HandlerProps{
		TTSUsecase: usecases.TTSUsecase,
	}
	return hp
}
