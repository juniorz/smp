package smp

import (
	"errors"
	"math/big"
)

type smp4State struct {
	y  *big.Int
	r7 *big.Int
}

func (s *smp4State) message(s2 *smp2State, msg3 SMP3) SMP4 {
	var m SMP4

	qaqb := divMod(msg3.qa, s2.qb, P)

	m.rb = modExp(qaqb, s2.b3)
	m.cr = hashMPIsBN(nil, 8, modExp(G1, s.r7), modExp(qaqb, s.r7))
	m.d7 = subMod(s.r7, mul(s2.b3, m.cr), Q)

	return m
}

// SMP4 represents the fourth message in the SMP protocol
type SMP4 struct {
	cr *big.Int
	d7 *big.Int
	rb *big.Int
}

func NewSMP4(mpis ...*big.Int) (*SMP4, error) {
	m := &SMP4{
		rb: new(big.Int),
		cr: new(big.Int),
		d7: new(big.Int),
	}

	if err := assignMPIs(m, mpis); err != nil {
		return nil, err
	}

	return m, nil
}

func (m SMP4) MPIs() []*big.Int {
	return []*big.Int{
		m.rb,
		m.cr,
		m.d7,
	}
}

func (p *Protocol) newSMP4Message(m3 SMP3) (m SMP4, err error) {
	if p.s4, err = p.newSMP4State(); err != nil {
		p.event(Failure)
		return
	}

	m = p.s4.message(p.s2, m3)

	return
}

func (p Protocol) newSMP4State() (s *smp4State, err error) {
	s = &smp4State{
		y: p.Secret,

		r7: new(big.Int),
	}

	err = p.generateRandMPIs([]*big.Int{
		s.r7,
	})

	return
}

func (p Protocol) verifySMP4(msg SMP4) error {
	s3 := p.s3

	if !p.IsGroupElement(msg.rb) {
		return errors.New("Rb is an invalid group element")
	}

	if !verifyZKP4(msg.cr, s3.g3b, msg.d7, s3.qaqb, msg.rb, 8) {
		return errors.New("cR is not a valid zero knowledge proof")
	}

	return nil
}

func (p Protocol) verifySMP4ProtocolSuccess(msg SMP4) error {
	s1 := p.s1
	s3 := p.s3

	rab := modExp(msg.rb, s1.a3)
	if !eq(rab, s3.papb) {
		return errors.New("protocol failed: x != y")
	}

	return nil
}
