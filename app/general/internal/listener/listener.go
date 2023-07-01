package listener

import (
	"github.com/bwmarrin/discordgo"
)

type Listener struct {
	soundPacket chan *discordgo.Packet
	speaking    chan bool
}

func NewListener(packet chan *discordgo.Packet, speaking chan bool) *Listener {
	return &Listener{
		soundPacket: packet,
		speaking:    speaking,
	}
}

func (l Listener) Run() {
	for {
		select {
		case <-l.soundPacket:
			l.speaking <- true
		default:
			l.speaking <- false
		}
	}
}
