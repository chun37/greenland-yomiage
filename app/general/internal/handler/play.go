package handler

import (
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
	"layeh.com/gopus"
)

func (h *Handler) Play(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}

	// check if the message is "!airhorn"
	if strings.HasPrefix(m.Content, "!play") {

		// Find the channel that the message came from.
		c, err := s.State.Channel(m.ChannelID)
		if err != nil {
			log.Printf("%v", err)
			return
		}

		// Find the guild for that channel.
		g, err := s.State.Guild(c.GuildID)
		if err != nil {
			log.Printf("%v", err)
			return
		}

		// Look for the message sender in that guild's current voice states.
		for _, vs := range g.VoiceStates {
			if vs.UserID == m.Author.ID {
				dgv, err := s.ChannelVoiceJoin(g.ID, vs.ChannelID, false, true)
				if err != nil {
					fmt.Println("failed to join voice channel:", err)
					return
				}
				defer dgv.Disconnect()

				f, err := os.Open("audio.wav")
				if err != nil {
					fmt.Println("failed to open file", err)
					return
				}
				done := make(chan struct{})
				opusChunks := make(OpusChunks, 3)
				defer close(done)
				defer close(opusChunks)
				go func() {
					err := Encode(f, opusChunks, done)
					if err != nil {

					}
				}()
				if err != nil {
					fmt.Println("failed to encode audio", err)
					return
				}
				err = PlayAudioFile(dgv, opusChunks, done)
				if err != nil {
					fmt.Println("failed to play audio files", err)
					return
				}

				return
			}
		}
	}
}

const (
	channels  int = 2                   // 1 for mono, 2 for stereo
	frameRate int = 48000               // audio sampling rate
	frameSize int = 960                 // uint16 size of each audio frame
	maxBytes  int = (frameSize * 2) * 2 // max size of opus data
)

type OpusChunks = chan []byte

func Encode(src io.Reader, chunks OpusChunks, done chan struct{}) error {
	run := exec.Command("ffmpeg", "-i", "pipe:0", "-f", "s16le", "-ar", strconv.Itoa(frameRate), "-ac", strconv.Itoa(channels), "pipe:1")
	run.Stdin = src
	dst, err := run.StdoutPipe()

	err = run.Start()
	if err != nil {
		fmt.Println("RunStart Error", err)
		return err
	}

	defer run.Process.Kill()

	opusEncoder, err := gopus.NewEncoder(frameRate, channels, gopus.Audio)

	for {
		// read data from ffmpeg stdout
		pcm := make([]int16, frameSize*channels)
		err = binary.Read(dst, binary.LittleEndian, &pcm)
		if err == io.EOF || err == io.ErrUnexpectedEOF {
			break
		}
		if err != nil {
			fmt.Println("error reading from ffmpeg stdout", err)
			break
		}
		opus, err := opusEncoder.Encode(pcm, frameSize, maxBytes)
		if err != nil {
			fmt.Println("opus encode error", err)
		}
		chunks <- opus
	}

	done <- struct{}{}

	return nil
}

// PlayAudioFile will play the given filename to the already connected
// Discord voice server/channel.  voice websocket and udp socket
// must already be setup before this will work.
func PlayAudioFile(v *discordgo.VoiceConnection, chunks OpusChunks, encodeClose chan struct{}) error {
	err := v.Speaking(true)
	if err != nil {
		fmt.Println("Couldn't set speaking", err)
	}

	// Send not "speaking" packet over the websocket when we finish
	defer func() {
		err := v.Speaking(false)
		if err != nil {
			fmt.Println("Couldn't stop speaking", err)
		}
	}()

	for {
		select {
		case opus := <-chunks:
			v.OpusSend <- opus
		case <-encodeClose:
			return nil
		}
	}
}
