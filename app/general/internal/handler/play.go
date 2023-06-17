package handler

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/bwmarrin/dgvoice"
	"github.com/bwmarrin/discordgo"
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
				pr, pw := io.Pipe()
				encodeProcess, err := Encode(f, pw)
				defer encodeProcess.Kill()
				if err != nil {
					fmt.Println("failed to encode audio", err)
					return
				}
				err = PlayAudioFile(dgv, pr)
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
	channels  int = 2     // 1 for mono, 2 for stereo
	frameRate int = 48000 // audio sampling rate
	frameSize int = 960   // uint16 size of each audio frame
)

func Encode(src io.Reader, dst io.Writer) (*os.Process, error) {
	run := exec.Command("ffmpeg", "-i", "pipe:0", "-f", "s16le", "-ar", strconv.Itoa(frameRate), "-ac", strconv.Itoa(channels), "pipe:1")
	run.Stdin = src
	run.Stdout = dst

	err := run.Start()
	if err != nil {
		fmt.Println("RunStart Error", err)
		return nil, err
	}

	return run.Process, nil
}

// PlayAudioFile will play the given filename to the already connected
// Discord voice server/channel.  voice websocket and udp socket
// must already be setup before this will work.
func PlayAudioFile(v *discordgo.VoiceConnection, reader io.Reader) error {
	buf := bufio.NewReaderSize(reader, 16384)

	//=======
	// Send "speaking" packet over the voice websocket
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

	send := make(chan []int16, 2)
	defer close(send)

	close := make(chan bool)
	go func() {
		dgvoice.SendPCM(v, send)
		close <- true
	}()

	for {
		// read data from ffmpeg stdout
		audiobuf := make([]int16, frameSize*channels)
		err = binary.Read(buf, binary.LittleEndian, &audiobuf)
		if err == io.EOF || err == io.ErrUnexpectedEOF {
			return nil
		}
		if err != nil {
			fmt.Println("error reading from ffmpeg stdout", err)
			return err
		}

		// Send received PCM to the sendPCM channel
		select {
		case send <- audiobuf:
		case <-close:
			return nil
		}
	}
}
