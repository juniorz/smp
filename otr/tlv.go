package otr

import (
	"bytes"
	"errors"
	"math/big"

	"github.com/juniorz/smp"
)

const (
	tlvTypeSMP1     = uint16(0x02)
	tlvTypeSMP2     = uint16(0x03)
	tlvTypeSMP3     = uint16(0x04)
	tlvTypeSMP4     = uint16(0x05)
	tlvTypeSMPAbort = uint16(0x06)
	tlvTypeSMP1Q    = uint16(0x07)
)

type TLV []byte

func Encode(m smp.Message) (TLV, error) {
	var tlv TLV

	switch v := m.(type) {
	case smp.SMP1:
		tlv = generateTLV(tlvTypeSMP1, tlvValue(m.MPIs()))
	case smp.SMP1Q:
		data := append(
			append([]byte(v.Question()), 0),
			tlvValue(m.MPIs())...,
		)
		tlv = generateTLV(tlvTypeSMP1Q, data)
	case smp.SMP2:
		tlv = generateTLV(tlvTypeSMP2, tlvValue(m.MPIs()))
	case smp.SMP3:
		tlv = generateTLV(tlvTypeSMP3, tlvValue(m.MPIs()))
	case smp.SMP4:
		tlv = generateTLV(tlvTypeSMP4, tlvValue(m.MPIs()))
	case smp.SMPAbort:
		tlv = generateTLV(tlvTypeSMPAbort, tlvValue(m.MPIs()))
	default:
		panic("unreachable")
	}

	return tlv, nil
}

func tlvValue(mpis []*big.Int) []byte {
	//TODO: data - 4 + ParameterLen() * len(mpis)
	data := make([]byte, 0, 1000)
	data = appendWord(data, uint32(len(mpis)))
	data = appendMPIs(data, mpis...)
	return data
}

func generateTLV(tp uint16, value []byte) TLV {
	data := make([]byte, 0, 1000)
	data = appendShort(data, tp)
	data = appendShort(data, uint16(len(value)))
	return append(data, value...)
}

func Decode(m TLV) (smp.Message, error) {
	tBytes, tType, ok := extractShort(m)
	if !ok {
		return nil, errors.New("wrong tlv type")
	}

	var tLen uint16
	tBytes, tLen, ok = extractShort(tBytes)
	if !ok {
		return nil, errors.New("wrong tlv length")
	}

	if len(tBytes) < int(tLen) {
		return nil, errors.New("wrong tlv value")
	}

	return parseTLV(tType, tLen, tBytes)
}

func parseTLV(t uint16, l uint16, v []byte) (smp.Message, error) {
	var question string
	if t == tlvTypeSMP1Q {
		nulPos := bytes.IndexByte(v, 0)
		if nulPos == -1 {
			return nil, errors.New("wrong tlv value")
		}

		//TODO: should we strdup?
		question = string(v[:nulPos])
		v = v[(nulPos + 1):]
	}

	_, mpis, ok := extractMPIs(v[:int(l)])
	if !ok {
		return nil, errors.New("not enough TLVs")
	}

	switch t {
	case tlvTypeSMP1:
		return smp.NewSMP1(mpis...)
	case tlvTypeSMP1Q:
		return smp.NewSMP1Q(question, mpis...)
	case tlvTypeSMP2:
		return smp.NewSMP2(mpis...)
	case tlvTypeSMP3:
		return smp.NewSMP3(mpis...)
	case tlvTypeSMP4:
		return smp.NewSMP4(mpis...)
	case tlvTypeSMPAbort:
		return smp.NewSMPAbort(mpis...)
	default:
		panic("unreachable")
	}
}
