package smp

import "math/big"

// SMP1Q represents the first message in the SMP protocol, but with a question
type SMP1Q struct {
	SMP1
	question string
}

func NewSMP1Q(question string, mpis ...*big.Int) (*SMP1Q, error) {
	m, err := NewSMP1(mpis...)
	if err != nil {
		return nil, err
	}

	return &SMP1Q{
		//FIXME: should SMP1 be a pointer?
		SMP1:     *m,
		question: question,
	}, nil
}

func (m *SMP1Q) Question() string {
	return m.question
}

func (p *Protocol) newSMP1QMessage(question string) (SMP1Q, error) {
	m, err := p.newSMP1Message()
	if err != nil {
		return SMP1Q{}, err
	}

	//TODO should we strdup question?
	return SMP1Q{
		SMP1:     m,
		question: question,
	}, nil
}
