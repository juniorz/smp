package channel

import "github.com/juniorz/smp"

type Protocol struct {
	*smp.Protocol

	sendC    chan<- smp.Message
	receiveC chan smp.Message
}

func NewProtocol(version int) *Protocol {
	return &Protocol{
		Protocol: smp.NewProtocol(version),
	}
}

func (p *Protocol) Receive() chan<- smp.Message {
	if p.receiveC == nil {
		p.receiveC = make(chan smp.Message, 1)
		go p.receiveLoop()
	}

	return p.receiveC
}

func (p *Protocol) Send() chan<- smp.Message {
	if p.sendC == nil {
		panic("sending to unpiped protocol")
	}

	return p.sendC
}

//Pipe ourself to the peer. Future invocations of Send() will be received by the peer
func (p *Protocol) Pipe(peer *Protocol) *Protocol {
	p.sendC = peer.Receive()
	return peer
}

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
