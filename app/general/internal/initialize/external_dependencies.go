package initialize

import (
	voicevoxcore "github.com/chun37/greenland-yomiage/internal/voicevox_core"
)

type ExternalDependencies struct {
	VoiceVox *voicevoxcore.VoiceVoxCore
}

func NewExternalDependencies() *ExternalDependencies {
	externalDependencies := new(ExternalDependencies)

	{
		externalDependencies.VoiceVox = voicevoxcore.New()
	}

	return externalDependencies
}
