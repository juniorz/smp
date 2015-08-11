package smp

import "math/big"

func verifyZKP(d, gen, c *big.Int, ix byte) bool {
	r := modExp(G1, d)
	s := modExp(gen, c)
	t := hashMPIsBN(nil, ix, mulMod(r, s, P))
	return eq(c, t)
}

func generateZKP(r, a *big.Int, ix byte) (c, d *big.Int) {
	c = hashMPIsBN(nil, ix, modExp(G1, r))
	d = generateDZKP(r, a, c)
	return
}

func generateDZKP(r, a, c *big.Int) *big.Int {
	return subMod(r, mul(a, c), Q)
}

func verifyZKP2(g2, g3, d5, d6, pb, qb, cp *big.Int, ix byte) bool {
	l := mulMod(
		modExp(g3, d5),
		modExp(pb, cp),
		P)
	r := mulMod(mul(modExp(G1, d5),
		modExp(g2, d6)),
		modExp(qb, cp),
		P)
	t := hashMPIsBN(nil, ix, l, r)
	return eq(cp, t)
}

func verifyZKP3(cp, g2, g3, d5, d6, pa, qa *big.Int, ix byte) bool {
	l := mulMod(modExp(g3, d5), modExp(pa, cp), P)
	r := mulMod(mul(modExp(G1, d5), modExp(g2, d6)), modExp(qa, cp), P)
	t := hashMPIsBN(nil, ix, l, r)
	return eq(cp, t)
}

func verifyZKP4(cr, g3a, d7, qaqb, ra *big.Int, ix byte) bool {
	l := mulMod(modExp(G1, d7), modExp(g3a, cr), P)
	r := mulMod(modExp(qaqb, d7), modExp(ra, cr), P)
	t := hashMPIsBN(nil, ix, l, r)
	return eq(cr, t)
}
