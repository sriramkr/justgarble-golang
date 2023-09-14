package justgarble

import (
	"encoding/binary"
	"math/rand"
)

func NewKey() ([]byte, error) {
	buf := make([]byte, 16)
	_, err := rand.Read(buf)
	if err != nil {
		return nil, err
	}
	return buf, nil
}

type Wire struct {
	H, L uint64
}

func NewWire() Wire {
	w := Wire{}
	w.H = rand.Uint64()
	w.L = rand.Uint64()
	return w
}

type Wirepair [2]Wire

func NewWirePair() Wirepair {
	wire0 := NewWire()
	wire1 := NewWire()
	if wire0.lsb() == wire1.lsb() {
		wire1.L = ^wire1.L
	}
	return Wirepair{wire0, wire1}
}

func NewWireWithLsb1() Wire {
	w := NewWire()
	if w.lsb() == 0 {
		w.L = ^w.L
	}
	return w
}

func (w *Wire) lsb() int {
	return int(w.L & 1)
}

func (w *Wire) double() {
	w.H = w.H << 1
	w.L = w.L << 1
}

func (w *Wire) duplicate(w2 *Wire) {
	w2.H = w.H
	w2.L = w.L
}

func (w *Wire) copyFrom(w2 *Wire) {
	w.H = w2.H
	w.L = w2.L
}

func (w *Wire) copyFromXor(w2, w3 *Wire) {
	w.H = w2.H ^ w3.H
	w.L = w2.L ^ w3.L
}

func (w *Wire) xor(w2 *Wire) {
	w.H = w.H ^ w2.H
	w.L = w.L ^ w2.L
}

func (w *Wire) xorInt(i int) {
	w.L = w.L ^ uint64(i)
}

func (w *Wire) isEqual(w2 *Wire) bool {
	return (w.L == w2.L) && (w.H == w2.H)
}

func (w *Wire) toByteArray(b []byte) {
	binary.LittleEndian.PutUint64(b, w.L)
	binary.LittleEndian.PutUint64(b[8:], w.H)
}

func (w *Wire) fromByteArray(b []byte) {
	w.L = binary.LittleEndian.Uint64(b)
	w.H = binary.LittleEndian.Uint64(b[8:])

}
