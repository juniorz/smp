package smp

import (
	"errors"
	"math/big"
)

type smp1State struct {
	a2, a3 *big.Int
	r2, r3 *big.Int
}

func (s smp1State) message() SMP1 {
	m := SMP1{}
	m.g2a = modExp(G1, s.a2)
	m.g3a = modExp(G1, s.a3)
	m.c2, m.d2 = generateZKP(s.r2, s.a2, 1)
	m.c3, m.d3 = generateZKP(s.r3, s.a3, 2)

	return m
}

// SMP1 represents the first message in the SMP protocol
type SMP1 struct {
	g2a, g3a *big.Int
	c2, c3   *big.Int
	d2, d3   *big.Int
}

func assignMPIs(m Message, src []*big.Int) error {
	dest := m.MPIs()
	if len(src) != len(dest) {
		return errors.New("not enought MPIs")
	}

	for i, mpi := range src {
		*dest[i] = *mpi
	}

	return nil
}

func NewSMP1(mpis ...*big.Int) (*SMP1, error) {
	m := &SMP1{
		g2a: new(big.Int),
		c2:  new(big.Int),
		d2:  new(big.Int),
		g3a: new(big.Int),
		c3:  new(big.Int),
		d3:  new(big.Int),
	}

	if err := assignMPIs(m, mpis); err != nil {
		return nil, err
	}

	return m, nil
}

func (m SMP1) MPIs() []*big.Int {
	return []*big.Int{
		m.g2a,
		m.c2, m.d2,
		m.g3a,
		m.c3, m.d3,
	}
}

func (p *Protocol) newSMP1Message() (m SMP1, err error) {
	if p.s1, err = p.newSMP1State(); err != nil {
		p.event(Failure)
		return
	}

	m = p.s1.message()
	return
}

func (p Protocol) newSMP1State() (s *smp1State, err error) {
	s = &smp1State{
		a2: new(big.Int),
		a3: new(big.Int),
		r2: new(big.Int),
		r3: new(big.Int),
	}

	err = p.generateRandMPIs([]*big.Int{
		s.a2, s.a3, s.r2, s.r3,
	})

	return
}

func (p Protocol) verifySMP1(msg SMP1) error {
	if !p.IsGroupElement(msg.g2a) {
		return errors.New("g2a is an invalid group element")
	}

	if !p.IsGroupElement(msg.g3a) {
		return errors.New("g3a is an invalid group element")
	}

	if !verifyZKP(msg.d2, msg.g2a, msg.c2, 1) {
		return errors.New("c2 is not a valid zero knowledge proof")
	}

	if !verifyZKP(msg.d3, msg.g3a, msg.c3, 2) {
		return errors.New("c3 is not a valid zero knowledge proof")
	}

	return nil
}
