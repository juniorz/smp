package channel

import "github.com/juniorz/smp"

// Protocol represents a channel-based SMP protocol
type Protocol struct {
	*smp.Protocol

	sendC    chan<- smp.Message
	receiveC chan smp.Message
}

// NewProtocol returns a channel-based SMP protocol
func NewProtocol(version int) *Protocol {
	return &Protocol{
		Protocol: smp.NewProtocol(version),
	}
}

// Receive returns the peer's receive channel.
// It should be used to send messages to this peer.
func (p *Protocol) Receive() chan<- smp.Message {
	if p.receiveC == nil {
		p.receiveC = make(chan smp.Message, 1)
		go p.receiveLoop()
	}

	return p.receiveC
}

// Send returns the peer's send channel.
// It should be used when this peer has a message to send to the other peer.
func (p *Protocol) Send() chan<- smp.Message {
	if p.sendC == nil {
		panic("sending to unpiped protocol")
	}

	return p.sendC
}

// Pipe ourself to the peer. Future invocations of Send() will be received by the peer
func (p *Protocol) Pipe(peer *Protocol) *Protocol {
	p.sendC = peer.Receive()
	return peer
}

// Compare this peer's secret value with the other peer. It starts the protocol.
func (p *Protocol) Compare() <-chan smp.Event {
	m, err := p.Protocol.Compare()
	if err != nil {
		panic(err)
	}

	p.Send() <- m

	return p.Events()
}

func (p *Protocol) receiveLoop() {
	for m := range p.receiveC {
		send, _ := p.Protocol.Receive(m)
		if send == nil {
			continue
		}

		p.Send() <- send
	}
}
