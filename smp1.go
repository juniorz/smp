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
	g2a, g3a    *big.Int
	c2, c3      *big.Int
	d2, d3      *big.Int
	hasQuestion bool
	question    string
}

func (p *Protocol) newSMP1Message() (m SMP1, err error) {
	if p.s1, err = p.newSMP1State(); err != nil {
		p.event(Failure)
		return
	}

	m = p.s1.message()

	if p.Question != "" {
		m.hasQuestion = true
		m.question = p.Question
	}

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
	if !p.isGroupElement(msg.g2a) {
		return errors.New("g2a is an invalid group element")
	}

	if !p.isGroupElement(msg.g3a) {
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
