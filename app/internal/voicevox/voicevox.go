package voicevox

import (
	"net/http"
)

type VoiceVox struct {
	HTTPClient *http.Client
}

func New() *VoiceVox {
	return &VoiceVox{
		&http.Client{},
	}
}
