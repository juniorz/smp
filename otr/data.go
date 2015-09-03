package otr

import "math/big"

func appendWord(l []byte, r uint32) []byte {
	return append(l, byte(r>>24), byte(r>>16), byte(r>>8), byte(r))
}

func appendShort(l []byte, r uint16) []byte {
	return append(l, byte(r>>8), byte(r))
}

func appendData(l, r []byte) []byte {
	return append(appendWord(l, uint32(len(r))), r...)
}

func appendMPI(l []byte, r *big.Int) []byte {
	return appendData(l, r.Bytes())
}

func appendMPIs(l []byte, r ...*big.Int) []byte {
	for _, mpi := range r {
		l = appendMPI(l, mpi)
	}
	return l
}

func extractShort(d []byte) ([]byte, uint16, bool) {
	if len(d) < 2 {
		return nil, 0, false
	}

	return d[2:], uint16(d[0])<<8 |
		uint16(d[1]), true
}

func extractWord(d []byte) ([]byte, uint32, bool) {
	if len(d) < 4 {
		return nil, 0, false
	}

	return d[4:], uint32(d[0])<<24 |
		uint32(d[1])<<16 |
		uint32(d[2])<<8 |
		uint32(d[3]), true
}

func extractMPI(d []byte) (newPoint []byte, mpi *big.Int, ok bool) {
	d, mpiLen, ok := extractWord(d)
	if !ok || len(d) < int(mpiLen) {
		return nil, nil, false
	}

	mpi = new(big.Int).SetBytes(d[:int(mpiLen)])
	newPoint = d[int(mpiLen):]
	ok = true
	return
}

func extractMPIs(d []byte) ([]byte, []*big.Int, bool) {
	current, mpiCount, ok := extractWord(d)
	if !ok {
		return nil, nil, false
	}
	result := make([]*big.Int, int(mpiCount))
	for i := 0; i < int(mpiCount); i++ {
		current, result[i], ok = extractMPI(current)
		if !ok {
			return nil, nil, false
		}
	}
	return current, result, true
}
