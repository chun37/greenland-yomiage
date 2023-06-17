package speaker

import (
	"fmt"

	"github.com/chun37/greenland-yomiage/internal/usecase/tts"
)

type Speaker struct {
	usecase  *tts.Usecase
	messages chan SpeechMessage
}

func NewSpeaker(usecase *tts.Usecase, messages chan SpeechMessage) *Speaker {
	return &Speaker{
		usecase:  usecase,
		messages: messages,
	}
}

func (s *Speaker) Run() {
	for {
		message := <-s.messages
		s.do(message)
	}
}

func (s *Speaker) do(message SpeechMessage) {
	err := message.VoiceConnection.Speaking(true)
	if err != nil {
		fmt.Println("Couldn't set speaking", err)
		return
	}

	done := make(chan struct{})
	opusChunks := make(chan []byte, 3)
	defer close(done)
	defer close(opusChunks)
	s.usecase.Do(tts.UsecaseParam{
		Text:       message.Text,
		OpusChunks: opusChunks,
		Done:       done,
	})

	// Send not "speaking" packet over the websocket when we finish
	defer func() {
		err := message.VoiceConnection.Speaking(false)
		if err != nil {
			fmt.Println("Couldn't stop speaking", err)
		}
	}()

	for {
		select {
		case opus := <-opusChunks:
			message.VoiceConnection.OpusSend <- opus
		case <-done:
			return
		}
	}
}
