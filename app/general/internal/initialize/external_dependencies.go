package initialize

import (
	"github.com/chun37/greenland-yomiage/internal/voicevox"
)

type ExternalDependencies struct {
	VoiceVox *voicevox.VoiceVox
}

func NewExternalDependencies() *ExternalDependencies {
	externalDependencies := new(ExternalDependencies)

	{
		externalDependencies.VoiceVox = voicevox.New()
	}

	return externalDependencies
}
