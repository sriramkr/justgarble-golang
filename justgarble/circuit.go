package main

import "math/rand"

type Circuit struct {
	n, m, q     int
	A, B        []int
	isXOR       []bool
	numAndGates int
}

func MakeRandomCircuit(n, m, q int) *Circuit {
	c := new(Circuit)
	c.n = n
	c.m = m
	c.q = q
	c.A = make([]int, c.q)
	c.B = make([]int, c.q)
	c.isXOR = make([]bool, c.q)
	c.numAndGates = 0
	for i := 0; i < q; i++ {
		c.isXOR[i] = rand.Intn(2) == 0
		if !c.isXOR[i] {
			c.numAndGates++
		}
		a := rand.Intn(n + i)
		b := rand.Intn(n + i)
		if a < b {
			c.A[i] = a
			c.B[i] = b
		} else {
			c.A[i] = b
			c.B[i] = a
		}
	}
	return c
}

func Compare(circuit *Circuit, f func([]bool) []bool) bool {
	inputs := make([]bool, circuit.n)
	for i := 0; i < (1 << circuit.n); i++ {
		for j := 0; j < circuit.n; j++ {
			inputs[j] = ((i) & (1 << j)) != 0
		}
		outputs1 := circuit.evalCircuit(inputs)
		outputs2 := f(inputs)
		if len(outputs1) != len(outputs2) {
			return false
		}
		for j := 0; j < circuit.m; j++ {
			if outputs1[j] != outputs2[j] {
				return false
			}
		}
	}
	return true
}

func (circuit *Circuit) evalCircuit(inputs []bool) []bool {
	wireValues := make([]bool, circuit.n+circuit.q)
	for i := 0; i < circuit.n; i++ {
		wireValues[i] = inputs[i]
	}
	for i := 0; i < circuit.q; i++ {
		if circuit.isXOR[i] {
			wireValues[circuit.n+i] = wireValues[circuit.A[i]] != wireValues[circuit.B[i]]
		} else {
			wireValues[circuit.n+i] = wireValues[circuit.A[i]] && wireValues[circuit.B[i]]
		}
	}
	return wireValues[len(wireValues)-circuit.m:]
}
