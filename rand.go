package smp

import (
	"io"
	"math/big"
)

func (p Protocol) parameterLength() int {
	switch p.version {
	case 2:
		return 16
	case 3:
		return 192
	}

	panic("unsupported version")
}

func (p Protocol) generateRandMPIs(mpis []*big.Int) (err error) {
	b := make([]byte, p.parameterLength())

	for i := range mpis {
		var r *big.Int
		r, err = p.randMPI(b)
		if err != nil {
			return
		}

		*mpis[i] = *r
	}

	return
}

func (p Protocol) randMPI(buf []byte) (*big.Int, error) {
	if _, err := io.ReadFull(p.Rand, buf); err != nil {
		return nil, errShortRandomRead
	}

	return new(big.Int).SetBytes(buf), nil
}
