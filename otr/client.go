package otr

import (
	"crypto/sha256"
	"errors"
	"math/big"

	"github.com/juniorz/smp"
	"github.com/twstrike/otr3"
)

//TODO: replicate the otr3.SMPEvents but with structs
//this way, Ask* events will have the "question"; Progress event will have the percentage and we stop calling event handlers with "zeroed" parameters (question = "", percent = 0)

type clientState int

const (
	notStarted clientState = iota
	inProgress
)

type SecretParams interface {
	SSID() []byte
	OurFingerprint() []byte
	TheirFingerprint() []byte
}

type Client struct {
	smp            *smp.Protocol
	pendingMessage *smp.SMP1
	secParams      SecretParams

	smpEventHandler otr3.SMPEventHandler
	state           clientState
}

func NewClient(conv otr3.Conversation) *Client {
	ops := &otr3.SmpOptions{conv}
	sec := &otr3.SmpSecretParams{conv}

	c := &Client{
		smp:       smp.NewProtocol(ops),
		secParams: sec,
		state:     notStarted,

		//smpEventHandler: handler,
	}

	// go c.watchEvents()

	return c
}

func (c *Client) watchEvents() {
	//TODO: should emit events to the smpHandler
	for e := range c.smp.Events() {
		switch e {
		case smp.InProgress:
			c.state = inProgress
		default:
			c.state = notStarted
		}
	}
}

func (c *Client) Start(question, secret string) (TLV, error) {
	// The spec says:
	// If you wish to restart SMP, send a type 6 TLV (SMP abort) to the other party and then proceed as if smpstate was SMPSTATE_EXPECT1. Otherwise, you may simply continue the current SMP instance.
	// Essentialy is up to the client implementation to decide is a SMP start request should be interpreted as a RESTART or if it should be interpreted as a ENSURE SMP is happening.
	// I have decided the later.
	if c.state == inProgress {
		return nil, nil
	}

	//FIXME should it be strduped?
	c.smp.Question = question

	// we are the initiator
	c.smp.Secret = generateSecret(c.secParams.OurFingerprint(),
		c.secParams.TheirFingerprint(), c.secParams.SSID(), []byte(secret))

	m, err := c.smp.Compare()
	if err != nil {
		return nil, err
	}

	return Encode(m)
}

func (c *Client) Continue(secret string) (TLV, error) {
	if c.smp.Secret != nil {
		return nil, errors.New("SMP already in progress")
	}

	var m *smp.SMP1
	m, c.pendingMessage = c.pendingMessage, nil
	if m == nil {
		return nil, errors.New("can't continue without having received a SMP1")
	}

	// they are the initiator
	c.smp.Secret = generateSecret(c.secParams.TheirFingerprint(),
		c.secParams.OurFingerprint(), c.secParams.SSID(), []byte(secret))

	ret, err := c.smp.Receive(m)
	if err != nil {
		return nil, err
	}

	return Encode(ret)
}

func (c *Client) Abort() (TLV, error) {
	// Per spec, always send the abort
	return Encode(c.smp.Abort())
}

func (c *Client) Receive(tlv TLV) (TLV, error) {
	dec, err := Decode(tlv)
	if err != nil {
		return nil, err
	}

	if m, ok := dec.(*smp.SMP1); ok && c.smp.Secret == nil {
		c.pendingMessage = m
		c.emitEvent(otr3.SMPEventAskForSecret, 0, "")
		return nil, nil
	}

	if m, ok := dec.(*smp.SMP1Q); ok && c.smp.Secret == nil {
		//TODO: fixme
		c.pendingMessage = &(m.SMP1)
		c.emitEvent(otr3.SMPEventAskForAnswer, 0, m.Question())
		return nil, nil
	}

	ret, err := c.smp.Receive(dec)
	if err != nil {
		return nil, err
	}

	if ret == nil {
		return nil, nil
	}

	return Encode(ret)
}

//FIXME: why does the event have to be handled with a percent and a question?
func (c *Client) emitEvent(e otr3.SMPEvent, percent int, question string) {
	if c.smpEventHandler != nil {
		c.smpEventHandler.HandleSMPEvent(e, percent, question)
	}
}

func generateSecret(initiatorFingerprint, recipientFingerprint, ssid, secret []byte) *big.Int {
	h := sha256.New()
	h.Write([]byte{smp.Version})
	h.Write(initiatorFingerprint)
	h.Write(recipientFingerprint)
	h.Write(ssid)
	h.Write(secret)
	return new(big.Int).SetBytes(h.Sum(nil))
}
