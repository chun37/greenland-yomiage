package listener

import (
	"time"

	"github.com/bwmarrin/discordgo"
)

type Listener struct {
	soundPacket chan *discordgo.Packet
	quiet       chan struct{}
}

func NewListener(packet chan *discordgo.Packet, quiet chan struct{}) *Listener {
	return &Listener{
		soundPacket: packet,
		quiet:       quiet,
	}
}

func (l Listener) Run() {
	lastSoundAt := time.Now()

	for {
		select {
		case <-l.soundPacket:
			lastSoundAt = time.Now()
		default:
			if time.Since(lastSoundAt).Milliseconds() > 300 {
				l.quiet <- struct{}{}
			}
		}
	}
}
