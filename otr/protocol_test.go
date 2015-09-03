package otr

import (
	"crypto/rand"
	"testing"

	"github.com/twstrike/otr3"
)

func TestProtocol(t *testing.T) {
	alice := otr3.Conversation{Rand: rand.Reader}
	alice.Policies.AllowV3()
	aliceKey := &otr3.PrivateKey{}
	aliceKey.Generate(rand.Reader)
	alice.SetKeys(aliceKey, nil)

	bob := otr3.Conversation{Rand: rand.Reader}
	bob.Policies.AllowV3()
	bobKey := &otr3.PrivateKey{}
	bobKey.Generate(rand.Reader)
	bob.SetKeys(bobKey, nil)

	var err error
	var aliceMessages []otr3.ValidMessage
	var bobMessages []otr3.ValidMessage

	aliceMessages = append(bobMessages, alice.QueryMessage())

	for len(aliceMessages)+len(bobMessages) > 0 {
		bobMessages = nil
		for _, m := range aliceMessages {
			_, bobMessages, err = bob.Receive(m)
			if err != nil {
				t.Errorf(err.Error())
			}
		}

		aliceMessages = nil
		for _, m := range bobMessages {
			_, aliceMessages, err = alice.Receive(m)
			if err != nil {
				t.Errorf(err.Error())
			}
		}
	}

	if !bob.IsEncrypted() {
		t.Errorf("Bob is not encrypted")
	}

	if !alice.IsEncrypted() {
		t.Errorf("Alice is not encrypted")
	}

	aliceSMP := NewClient(alice)
	bobSMP := NewClient(bob)

	// <- SMP1
	toSend, err := aliceSMP.Start("what is my pet's name?", "scooby")
	if err != nil {
		t.Errorf(err.Error())
	}

	// SMP1 ->
	toSend, err = bobSMP.Receive(toSend)
	if err != nil {
		t.Errorf(err.Error())
	}

	if toSend != nil {
		t.Errorf("Bob shouldn't have sent any message at this point")
	}

	//TODO: should receive an event asking for the secret
	//containing the question

	// <- SMP2
	toSend, err = bobSMP.Continue("scooby")
	if err != nil {
		t.Errorf(err.Error())
	}

	// SMP2 ->
	// <- SMP3
	toSend, err = aliceSMP.Receive(toSend)
	if err != nil {
		t.Errorf(err.Error())
	}

	// SMP3 ->
	// <- SMP4
	toSend, err = bobSMP.Receive(toSend)
	if err != nil {
		t.Errorf(err.Error())
	}

	//Bob should emit a Completed event

	toSend, err = aliceSMP.Receive(toSend)
	if err != nil {
		t.Errorf(err.Error())
	}

	//Alice should emit a Completed event
}
