package channel

import (
	"math/big"
	"testing"
	"time"

	"github.com/juniorz/smp"
)

func TestCantSendToUnpipedProtocol(t *testing.T) {
	expectedPanic := "sending to unpiped protocol"
	defer func() {
		if r := recover(); r != expectedPanic {
			t.Errorf("Did not panic")
		}
	}()

	p := NewProtocol(3)
	p.Send() <- smp.SMP1{}
}

func TestMessagesSentAreReceivedByTheOtherEnd(t *testing.T) {
	rec := NewProtocol(3)
	p := NewProtocol(3)

	p.Pipe(rec)
	p.Send() <- smp.SMP1{}

	select {
	case <-rec.receiveC:
		//OK
	default:
		t.Errorf("Failed to receive")
	}
}

func TestComparesIdenticalSecrets(t *testing.T) {
	alice := NewProtocol(3)
	bob := NewProtocol(3)

	alice.Secret = big.NewInt(123456)
	alice.Question = "Whats our secret?"
	alice.Pipe(bob)

	bob.Secret = big.NewInt(123456)
	bob.Pipe(alice)

	select {
	case in := <-alice.Compare():
		switch in {
		case smp.Success:
		default:
			t.Errorf("SMP protocol failed", in)
		}
	case <-time.After(1 * time.Second):
		t.Errorf("SMP protocol failed")
	}

}

func TestFailsToCompareDifferentSecrets(t *testing.T) {
	alice := NewProtocol(3)
	bob := NewProtocol(3)

	alice.Secret = big.NewInt(123456)
	alice.Question = "Whats our secret?"
	alice.Pipe(bob)

	bob.Secret = big.NewInt(1234567)
	bob.Pipe(alice)

	select {
	case in := <-alice.Compare():
		switch in {
		default:
			//so many error cases
		case smp.Success:
			t.Errorf("SMP protocol succeeded", in)
		}
	case <-time.After(1 * time.Second):
		t.Errorf("SMP protocol did not fail")
	}

}
