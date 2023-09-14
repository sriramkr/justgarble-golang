package justgarble

import (
	"math/rand"
	"testing"
)

type DummyCipher struct {
}

func (d *DummyCipher) CreateMask(A, B Wire, T int) Wire {
	return Wire{}
}

func TestGarbledGate(t *testing.T) {
	l := NewWire()
	r := NewWire()
	o := NewWire()
	R := NewWireWithLsb1()
	d := new(DummyCipher)
	w := make([]Wire, 8)
	AndGateGarbler(l, r, o, R, 1, d, w)
	a := l.lsb()
	b := r.lsb()
	zero_pos := (a << 1) + b
	one_pos := (a << 1) + (1 - b)
	two_pos := ((1 - a) << 1) + b
	three_pos := (1-a)<<1 + (1 - b)

	if !w[4+zero_pos].isEqual(&o) {
		t.Error("Garble output mismatch")
	}
	if !w[4+one_pos].isEqual(&o) {
		t.Error("Garble output mismatch")
	}
	if !w[4+two_pos].isEqual(&o) {
		t.Error("Garble output mismatch")
	}
	o1 := Wire{}
	o1.copyFromXor(&o, &R)
	if !w[4+three_pos].isEqual(&o1) {
		t.Error("Garble output mismatch")
	}
}

func MakeCircuit3() *Circuit {
	c := new(Circuit)
	c.n = 4
	c.m = 2
	c.q = 4
	c.A = make([]int, c.q)
	c.B = make([]int, c.q)
	c.isXOR = make([]bool, c.q)
	c.A[0] = 0
	c.B[0] = 1
	c.A[1] = 2
	c.B[1] = 3
	c.A[2] = 4
	c.B[2] = 5
	c.A[3] = 4
	c.B[3] = 6
	c.isXOR[0] = false
	c.isXOR[1] = false
	c.isXOR[2] = true
	c.isXOR[3] = false
	return c
}

func TestEvalCircuit(t *testing.T) {
	c := MakeCircuit3()
	g := GarblerImpl{}
	gc, err := g.Garble(c)
	if err != nil {
		t.Error(("unexpected error"))
	}

	inputs := []bool{true, false, true, true}
	outputs, err := g.Eval(gc, inputs)
	if err != nil {
		t.Error(("unexpected error"))
	}
	outputs2 := c.evalCircuit(inputs)
	if len(outputs) != len(outputs2) {
		t.Error("mismatched output length")
	}
	for i := 0; i < len(outputs); i++ {
		if outputs[i] != outputs2[i] {
			t.Error("mismatched output")
		}
	}

}

func randomInputs(inputs []bool) {
	for i := 0; i < len(inputs); i++ {
		inputs[i] = rand.Intn(2) == 0
	}
}

func TestXorGateGarble(t *testing.T) {
	n := 10
	m := 5
	q := 20

	c := MakeRandomCircuit(n, m, q)
	g := GarblerImpl{}
	gc, err := g.Garble(c)
	if err != nil {
		t.Error(("unexpected error"))
	}

	for i := 0; i < c.q; i++ {
		leftWire := gc.AllWires[c.A[i]]
		rightWire := gc.AllWires[c.B[i]]
		if c.isXOR[i] {
			leftWire.xor(&rightWire)
			if !gc.AllWires[i+c.n].isEqual(&leftWire) {
				t.Error("wire mismatch")
			}
		}
	}
}

func TestEvalRandomCircuit(t *testing.T) {
	n := 10
	m := 5
	q := 20

	inputs := make([]bool, n)
	for j := 0; j < 1000; j++ {
		c := MakeRandomCircuit(n, m, q)
		randomInputs(inputs)
		g := GarblerImpl{}
		gc, err := g.Garble(c)
		if err != nil {
			t.Error(("unexpected error"))
		}

		outputs, err := g.Eval(gc, inputs)
		if err != nil {
			t.Error(("unexpected error"))
		}

		outputs2 := c.evalCircuit(inputs)
		for i := 0; i < len(outputs); i++ {
			if outputs[i] != outputs2[i] {
				t.Error("FAILED")
			}
		}
	}

}
