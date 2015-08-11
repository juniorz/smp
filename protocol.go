package smp

import (
	"crypto/rand"
	"errors"
	"io"
	"math/big"
)

type Event int

const (
	Success Event = iota
	Abort
	Cheated
	Error
	Failure
)

var (
	errUnspecifiedSecret = errors.New("missing secret")
	errShortRandomRead   = errors.New("short read from rand source")
)

type Message interface {
	received(*Protocol) (Message, error)
}

type Protocol struct {
	Rand     io.Reader
	Question string
	Secret   *big.Int

	eventC chan Event

	version int

	smpState
	s1 *smp1State
	s2 *smp2State
	s3 *smp3State
	s4 *smp4State
}

func NewProtocol(version int) *Protocol {
	return &Protocol{
		version:  version,
		smpState: smpStateExpect1{},
		Rand:     rand.Reader,
		eventC:   make(chan Event, 1),
	}
}

func (p *Protocol) Receive(m Message) (Message, error) {
	send, err := m.received(p)
	if err != nil {
		p.event(Failure)
		return nil, err
	}

	return send, nil
}

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

	return m, nil
}

func (p Protocol) event(e Event) {
	go func() { p.eventC <- e }()
}

func (p Protocol) Events() <-chan Event {
	return p.eventC
}
