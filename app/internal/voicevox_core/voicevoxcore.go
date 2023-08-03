package voicevoxcore

import voicevoxcorego "github.com/sh1ma/voicevoxcore.go"

type VoiceVoxCore struct {
	voicevoxcorego.VoicevoxCore
}

func New() *VoiceVoxCore {
	core := voicevoxcorego.New()
	options := core.MakeDefaultInitializeOptions()
	options.UpdateOpenJtalkDictDir("./open_jtalk_dic_utf_8-1.11")
	options.UpdateLoadAllModels(true)
	core.Initialize(options)
	return &VoiceVoxCore{*core}
}

func (r *VoiceVoxCore) Generate(text string) ([]byte, error) {
	return r.Tts(text, 8, r.MakeDefaultTtsOotions())
}

func (r *VoiceVoxCore) Add(word, pronunciation string, accent int) error {
	return nil
}
