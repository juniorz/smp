package smp

import (
	"crypto/rand"
	"errors"
	"io"
	"math/big"
)

const Version = 1

// Event represents SMP events
type Event int

const (
	// Success means the SMP completed with success and the secrets match
	Success Event = iota
	// InProgress means update the auth progress dialog with progress_percent
	InProgress
	// Abort means the SMP protocol has been aborted
	Abort
	// Cheated means the SMP protocol has been cheated
	Cheated
	// Error means the SMP protocol terminated due an error in the protocol, like a verification failure
	Error
	// Failure means the SMP protocol failed due errors extrinsic to the protocol
	Failure
)

var (
	errUnspecifiedSecret = errors.New("missing secret")
	errShortRandomRead   = errors.New("short read from rand source")
)

// Message represents an SMP message
type Message interface {
	MPIs() []*big.Int
	received(*Protocol) (Message, error)
}

// Options represents configuration options for the SMP protocol
type Options interface {
	ParameterLength() int
	IsGroupElement(*big.Int) bool
}

// Protocol represents the SMP protocol
type Protocol struct {
	Options
	Rand     io.Reader
	Question string
	Secret   *big.Int

	eventC chan Event
	smpState
	s1 *smp1State
	s2 *smp2State
	s3 *smp3State
	s4 *smp4State
}

// NewProtocol returns an SMP protocol
func NewProtocol(options Options) *Protocol {
	return &Protocol{
		Options:  options,
		smpState: smpStateExpect1{},
		Rand:     rand.Reader,
		eventC:   make(chan Event, 1),
	}
}

// Receive process the incoming message and potentially returns a message
// addressed to the other peer
func (p *Protocol) Receive(m Message) (Message, error) {
	send, err := m.received(p)
	if err != nil {
		p.event(Failure)
		return nil, err
	}

	return send, nil
}

func (p *Protocol) Abort() (ret Message) {
	//err is always nil
	p.smpState, ret, _ = sendSMPAbortAndRestartStateMachine()
	p.event(Abort)
	return
}

// Compare starts the protocol and generates a message addressed to the other
// peer
func (p *Protocol) Compare() (Message, error) {
	if p.Secret == nil {
		p.event(Failure)
		return nil, errUnspecifiedSecret
	}

	m, err := p.newSMP1Message()
	if err != nil {
		p.event(Failure)
		return nil, err
	}

	p.smpState = smpStateExpect2{}
	p.event(InProgress)

	return m, nil
}

func (p *Protocol) startMessage() (Message, error) {
	if len(p.Question) > 0 {
		return p.newSMP1QMessage(p.Question)
	}

	return p.newSMP1Message()
}

// Events returns the events channel for this Protocol
func (p Protocol) Events() <-chan Event {
	return p.eventC
}

func (p Protocol) event(e Event) {
	go func() { p.eventC <- e }()
}
