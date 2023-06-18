package speaker

import (
	"log"

	"github.com/chun37/greenland-yomiage/internal/usecase/tts"
	"golang.org/x/xerrors"
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
		if err := s.do(message); err != nil {
			log.Println("failed to speak message: %+v", err)
		}
	}
}

func (s *Speaker) do(message SpeechMessage) error {
	if err := message.VoiceConnection.Speaking(true); err != nil {
		return xerrors.Errorf("Couldn't set speaking: %w", err)
	}

	done := make(chan struct{})
	opusChunks := make(chan []byte, 3)
	defer close(done)
	defer close(opusChunks)
	if err := s.usecase.Do(tts.UsecaseParam{
		Text:       message.Text,
		OpusChunks: opusChunks,
		Done:       done,
	}); err != nil {
		return xerrors.Errorf("failed to exec usecase: %w", err)
	}

	// Send not "speaking" packet over the websocket when we finish
	defer func() {
		err := message.VoiceConnection.Speaking(false)
		if err != nil {
			log.Println("Couldn't stop speaking", err)
		}
	}()

	for {
		select {
		case opus := <-opusChunks:
			message.VoiceConnection.OpusSend <- opus
		case <-done:
			return nil
		}
	}
}
