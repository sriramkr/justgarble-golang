package justgarble

import (
	"testing"
)

func MakeCircuit() *Circuit {
	c := new(Circuit)
	c.n = 2
	c.m = 1
	c.q = 1
	c.A = make([]int, c.q)
	c.B = make([]int, c.q)
	c.isXOR = make([]bool, c.q)
	c.A[0] = 0
	c.B[0] = 1
	c.isXOR[0] = true
	return c
}

func MakeCircuit2() *Circuit {
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

func f1(inputs []bool) []bool {
	return []bool{inputs[0] != inputs[1]}
}

func f2(inputs []bool) []bool {
	v1 := (inputs[0] && inputs[1])
	v2 := (inputs[2] && inputs[3])
	v3 := v1 != v2
	v4 := v1 && v3
	return []bool{v3, v4}
}

func TestCircuits(t *testing.T) {

	c1 := MakeCircuit()
	if !Compare(c1, f1) {
		t.Error("circuit creation failed")
	}

	c2 := MakeCircuit2()
	if !Compare(c2, f2) {
		t.Error("circuit creation failed")
	}

}
