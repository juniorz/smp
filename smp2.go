package smp

import (
	"errors"
	"math/big"
)

type smp2State struct {
	y                  *big.Int
	b2, b3             *big.Int
	r2, r3, r4, r5, r6 *big.Int
	g2, g3             *big.Int
	g3a                *big.Int
	pb, qb             *big.Int
}

func (s *smp2State) message(s1 SMP1) SMP2 {
	var m SMP2

	m.g2b = modExp(G1, s.b2)
	m.g3b = modExp(G1, s.b3)

	m.c2, m.d2 = generateZKP(s.r2, s.b2, 3)
	m.c3, m.d3 = generateZKP(s.r3, s.b3, 4)

	s.g3a = s1.g3a
	s.g2 = modExp(s1.g2a, s.b2)
	s.g3 = modExp(s1.g3a, s.b3)

	s.pb = modExp(s.g3, s.r4)
	s.qb = mulMod(modExp(G1, s.r4), modExp(s.g2, s.y), P)

	m.pb = s.pb
	m.qb = s.qb

	m.cp = hashMPIsBN(nil, 5,
		modExp(s.g3, s.r5),
		mulMod(modExp(G1, s.r5), modExp(s.g2, s.r6), P))

	m.d5 = subMod(s.r5, mul(s.r4, m.cp), Q)
	m.d6 = subMod(s.r6, mul(s.y, m.cp), Q)

	return m
}

// SMP2 represents the second message in the SMP protocol
type SMP2 struct {
	g2b, g3b *big.Int
	c2, c3   *big.Int
	d2, d3   *big.Int
	pb, qb   *big.Int
	cp       *big.Int
	d5, d6   *big.Int
}

func NewSMP2(mpis ...*big.Int) (*SMP2, error) {
	m := &SMP2{
		g2b: new(big.Int),
		c2:  new(big.Int),
		d2:  new(big.Int),
		g3b: new(big.Int),
		c3:  new(big.Int),
		d3:  new(big.Int),
		pb:  new(big.Int),
		qb:  new(big.Int),
		cp:  new(big.Int),
		d5:  new(big.Int),
		d6:  new(big.Int),
	}

	if err := assignMPIs(m, mpis); err != nil {
		return nil, err
	}

	return m, nil
}

func (m SMP2) MPIs() []*big.Int {
	return []*big.Int{
		m.g2b,
		m.c2, m.d2,
		m.g3b,
		m.c3, m.d3,
		m.pb, m.qb, m.cp,
		m.d5, m.d6,
	}
}

func (p *Protocol) newSMP2Message(m1 SMP1) (m SMP2, err error) {
	if p.s2, err = p.newSMP2State(); err != nil {
		p.event(Failure)
		return
	}

	m = p.s2.message(m1)

	return
}

func (p Protocol) newSMP2State() (s *smp2State, err error) {
	s = &smp2State{
		y: p.Secret,

		b2: new(big.Int),
		b3: new(big.Int),
		r2: new(big.Int),
		r3: new(big.Int),
		r4: new(big.Int),
		r5: new(big.Int),
		r6: new(big.Int),
	}

	err = p.generateRandMPIs([]*big.Int{
		s.b2, s.b3, s.r2, s.r3, s.r4, s.r5, s.r6,
	})

	return
}

func (p Protocol) verifySMP2(msg SMP2) error {
	if !p.IsGroupElement(msg.g2b) {
		return errors.New("g2b is an invalid group element")
	}

	if !p.IsGroupElement(msg.g3b) {
		return errors.New("g3b is an invalid group element")
	}

	if !p.IsGroupElement(msg.pb) {
		return errors.New("Pb is an invalid group element")
	}

	if !p.IsGroupElement(msg.qb) {
		return errors.New("Qb is an invalid group element")
	}

	if !verifyZKP(msg.d2, msg.g2b, msg.c2, 3) {
		return errors.New("c2 is not a valid zero knowledge proof")
	}

	if !verifyZKP(msg.d3, msg.g3b, msg.c3, 4) {
		return errors.New("c3 is not a valid zero knowledge proof")
	}

	g2 := modExp(msg.g2b, p.s1.a2)
	g3 := modExp(msg.g3b, p.s1.a3)

	if !verifyZKP2(g2, g3, msg.d5, msg.d6, msg.pb, msg.qb, msg.cp, 5) {
		return errors.New("cP is not a valid zero knowledge proof")
	}

	return nil
}
