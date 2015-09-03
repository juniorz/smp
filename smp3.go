package smp

import (
	"errors"
	"math/big"
)

type smp3State struct {
	x              *big.Int
	g3b            *big.Int
	r4, r5, r6, r7 *big.Int
	qaqb, papb     *big.Int
}

func (s *smp3State) message(s1 *smp1State, m2 SMP2) SMP3 {
	var m SMP3

	g2 := modExp(m2.g2b, s1.a2)
	g3 := modExp(m2.g3b, s1.a3)

	m.pa = modExp(g3, s.r4)
	m.qa = mulMod(modExp(G1, s.r4), modExp(g2, s.x), P)

	s.g3b = m2.g3b
	s.qaqb = divMod(m.qa, m2.qb, P)
	s.papb = divMod(m.pa, m2.pb, P)

	m.cp = hashMPIsBN(nil, 6, modExp(g3, s.r5), mulMod(modExp(G1, s.r5), modExp(g2, s.r6), P))
	m.d5 = generateDZKP(s.r5, s.r4, m.cp)
	m.d6 = generateDZKP(s.r6, s.x, m.cp)

	m.ra = modExp(s.qaqb, s1.a3)

	m.cr = hashMPIsBN(nil, 7, modExp(G1, s.r7), modExp(s.qaqb, s.r7))
	m.d7 = subMod(s.r7, mul(s1.a3, m.cr), Q)

	return m
}

// SMP3 represents the third message in the SMP protocol
type SMP3 struct {
	pa, qa     *big.Int
	cp         *big.Int
	d5, d6, d7 *big.Int
	ra         *big.Int
	cr         *big.Int
}

func NewSMP3(mpis ...*big.Int) (*SMP3, error) {
	m := &SMP3{
		pa: new(big.Int),
		qa: new(big.Int),
		cp: new(big.Int),
		d5: new(big.Int),
		d6: new(big.Int),
		ra: new(big.Int),
		cr: new(big.Int),
		d7: new(big.Int),
	}

	if err := assignMPIs(m, mpis); err != nil {
		return nil, err
	}

	return m, nil
}

func (m SMP3) MPIs() []*big.Int {
	return []*big.Int{
		m.pa, m.qa,
		m.cp,
		m.d5, m.d6,
		m.ra,
		m.cr,
		m.d7,
	}
}

func (p *Protocol) newSMP3Message(m2 SMP2) (m SMP3, err error) {
	if p.s3, err = p.newSMP3State(); err != nil {
		p.event(Failure)
		return
	}

	m = p.s3.message(p.s1, m2)

	return
}

func (p Protocol) newSMP3State() (s *smp3State, err error) {
	s = &smp3State{
		x: p.Secret,

		r4: new(big.Int),
		r5: new(big.Int),
		r6: new(big.Int),
		r7: new(big.Int),
	}

	err = p.generateRandMPIs([]*big.Int{
		s.r4, s.r5, s.r6, s.r7,
	})

	return
}

func (p Protocol) verifySMP3(msg SMP3) error {
	if !p.IsGroupElement(msg.pa) {
		return errors.New("Pa is an invalid group element")
	}

	if !p.IsGroupElement(msg.qa) {
		return errors.New("Qa is an invalid group element")
	}

	if !p.IsGroupElement(msg.ra) {
		return errors.New("Ra is an invalid group element")
	}

	if !verifyZKP3(msg.cp, p.s2.g2, p.s2.g3, msg.d5, msg.d6, msg.pa, msg.qa, 6) {
		return errors.New("cP is not a valid zero knowledge proof")
	}

	qaqb := divMod(msg.qa, p.s2.qb, P)

	if !verifyZKP4(msg.cr, p.s2.g3a, msg.d7, qaqb, msg.ra, 7) {
		return errors.New("cR is not a valid zero knowledge proof")
	}

	return nil
}

func (p Protocol) verifySMP3ProtocolSuccess(msg SMP3) error {
	s2 := p.s2

	papb := divMod(msg.pa, s2.pb, P)
	rab := modExp(msg.ra, s2.b3)

	if !eq(rab, papb) {
		return errors.New("protocol failed: x != y")
	}

	return nil
}
