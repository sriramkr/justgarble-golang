package justgarble

import (
	"crypto/aes"
	"crypto/cipher"
)

type FixedKeyCipher interface {
	CreateMask(A, B Wire, T int) Wire
}

type FixedKeyCipherReal struct {
	c      cipher.Block
	input  [16]byte
	output [16]byte
}

func MakeNewFixedKeyCipher(key []byte) (*FixedKeyCipherReal, error) {
	f := new(FixedKeyCipherReal)
	var err error
	f.c, err = aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	return f, nil
}

func (f *FixedKeyCipherReal) CreateMask(A, B Wire, T int) Wire {
	twoA := Wire{}
	twoA.copyFrom(&A)
	twoA.double()
	fourB := Wire{}
	fourB.double()
	fourB.double()
	twoA.xor(&fourB)
	twoA.xorInt(T)
	twoA.toByteArray(f.input[:])
	f.c.Encrypt(f.output[:], f.input[:])
	w := Wire{}
	w.fromByteArray(f.output[:])
	return w
}
