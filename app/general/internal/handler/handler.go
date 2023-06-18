package handler

import (
	"github.com/chun37/greenland-yomiage/general/internal/props"
	"github.com/chun37/greenland-yomiage/general/internal/speaker"
)

type Handler struct {
	props    *props.HandlerProps
	messages chan speaker.SpeechMessage
}

func New(props *props.HandlerProps, messages chan speaker.SpeechMessage) *Handler {
	return &Handler{
		props:    props,
		messages: messages,
	}
}
