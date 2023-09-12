package main

import (
	"math/rand"
	"testing"
)

func TestNewWirePair(t *testing.T) {
	w := NewWirePair()
	if w[0].lsb() == w[1].lsb() {
		t.Errorf("got %q, wanted %q", w[0], w[1])
	}
}

func TestNewWireWithLsb(t *testing.T) {
	w := NewWireWithLsb1()
	if w.lsb() != 1 {
		t.Error("bad wire")
	}

}

func TestOperations(t *testing.T) {
	w := Wire{1, 2}
	w.double()
	if w.L != 4 || w.H != 2 {
		t.Error("double failed")
	}
	w2 := NewWire()
	w3 := NewWire()
	w2.duplicate(&w3)
	if w2.L != w3.L || w2.H != w3.H {
		t.Error("duplicate failed")
	}
	w4 := NewWire()
	w5 := NewWire()
	w4.duplicate(&w5)
	w6 := NewWire()
	w4.xor(&w6)
	if (w4.L != (w5.L ^ w6.L)) || (w4.H != (w5.H ^ w6.H)) {
		t.Error("xor failed")
	}
	w7 := NewWire()
	w8 := NewWire()
	w8.copyFrom(&w7)
	if !w7.isEqual(&w8) {
		t.Error("copy from failed")
	}
	w9 := NewWire()
	w10 := NewWire()
	w11 := NewWire()
	w11.copyFromXor(&w9, &w10)
	w9.xor(&w10)
	if !w9.isEqual(&w11) {
		t.Error("copy from xor failed")
	}
	w12 := NewWire()
	w13 := NewWire()
	w13.copyFrom(&w12)
	diff := rand.Intn(1000)
	w12.xorInt(diff)
	w13.L ^= uint64(diff)
	if !w13.isEqual(&w12) {
		t.Error("xor int failed")
	}

}

func TestSerialize(t *testing.T) {
	b := make([]byte, 16)
	w := NewWire()
	w.toByteArray(b)
	w2 := NewWire()
	w2.fromByteArray(b)
	if !w.isEqual(&w2) {
		t.Error("serialization failed")
	}
}
