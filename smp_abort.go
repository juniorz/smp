package smp

import "math/big"

// SMPAbort represents the abort message in the SMP protocol
type SMPAbort struct{}

func NewSMPAbort(mpis ...*big.Int) (SMPAbort, error) {
	return SMPAbort{}, nil
}

func (m SMPAbort) MPIs() []*big.Int {
	return []*big.Int{}
}
