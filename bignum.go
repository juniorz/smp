package smp

import "math/big"

func sub(l, r *big.Int) *big.Int {
	return new(big.Int).Sub(l, r)
}

func subMod(l, r, m *big.Int) *big.Int {
	res := sub(l, r)
	res.Mod(res, m)
	return res
}

func lte(l, r *big.Int) bool {
	return l.Cmp(r) != 1
}

func gte(l, r *big.Int) bool {
	return l.Cmp(r) != -1
}

func eq(l, r *big.Int) bool {
	return l.Cmp(r) == 0
}

func mul(l, r *big.Int) *big.Int {
	return new(big.Int).Mul(l, r)
}

func mulMod(l, r, m *big.Int) *big.Int {
	res := mul(l, r)
	res.Mod(res, m)
	return res
}

// Fast division over a modular field, without using division
func divMod(l, r, m *big.Int) *big.Int {
	return mulMod(l, modInverse(r, m), m)
}

func modInverse(g, x *big.Int) *big.Int {
	return new(big.Int).ModInverse(g, x)
}

func modExp(g, x *big.Int) *big.Int {
	return new(big.Int).Exp(g, x, P)
}
