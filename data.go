package smp

import (
	"crypto/sha256"
	"hash"
	"math/big"
)

func appendData(l, r []byte) []byte {
	return append(appendWord(l, uint32(len(r))), r...)
}

func appendWord(l []byte, r uint32) []byte {
	return append(l, byte(r>>24), byte(r>>16), byte(r>>8), byte(r))
}

func appendMPI(l []byte, r *big.Int) []byte {
	return appendData(l, r.Bytes())
}

func hashMPIsBN(h hash.Hash, magic byte, mpis ...*big.Int) *big.Int {
	return new(big.Int).SetBytes(hashMPIs(h, magic, mpis...))
}

func hashMPIs(h hash.Hash, magic byte, mpis ...*big.Int) []byte {
	if h != nil {
		h.Reset()
	} else {
		h = sha256.New()
	}

	h.Write([]byte{magic})
	for _, mpi := range mpis {
		h.Write(appendMPI(nil, mpi))
	}
	return h.Sum(nil)
}
