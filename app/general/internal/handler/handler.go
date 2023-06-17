package handler

import "github.com/chun37/greenland-yomiage/general/internal/props"

type Handler struct {
	props *props.HandlerProps
}

func NewHandler(props *props.HandlerProps) *Handler {
	return &Handler{
		props: props,
	}
}
