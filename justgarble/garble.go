package main

type GarbledCircuit struct {
	Circuit          Circuit
	Gates            []Wire
	Masterkey        []byte
	R                Wire
	InputWires       []Wire
	AllWires         []Wire
	OutputShouldFlip []bool
}

type Garbler interface {
	Garble(*Circuit) *GarbledCircuit
	Eval(*GarbledCircuit, []bool) []bool
}

type GarblerImpl struct {
}

func (g *GarblerImpl) Garble(circuit *Circuit) *GarbledCircuit {
	return g.GarbleInternal(circuit, AndGateGarbler)
}

func (g *GarblerImpl) GarbleInternal(circuit *Circuit, andGateGarbler func(Wire, Wire, Wire, Wire, int, FixedKeyCipher, []Wire)) *GarbledCircuit {
	gc := new(GarbledCircuit)
	gc.Circuit = *circuit
	gc.Gates = make([]Wire, circuit.q*4)
	var err error
	gc.Masterkey, err = NewKey()
	if err != nil {
		panic(err)
	}
	gc.InputWires = make([]Wire, circuit.n)
	gc.AllWires = make([]Wire, circuit.n+circuit.q)
	gc.OutputShouldFlip = make([]bool, circuit.m)
	gc.R = NewWireWithLsb1()

	f := MakeNewFixedKeyCipher(gc.Masterkey)

	for i := 0; i < circuit.n; i++ {
		wire0 := NewWire()
		gc.AllWires[i] = wire0
		gc.InputWires[i] = wire0
	}

	for i := 0; i < circuit.q; i++ {
		leftWire := gc.AllWires[circuit.A[i]]
		rightWire := gc.AllWires[circuit.B[i]]
		if circuit.isXOR[i] {
			gc.AllWires[i+circuit.n].copyFromXor(&leftWire, &rightWire)
		} else {
			newWire := NewWire()
			gc.AllWires[i+circuit.n] = newWire
			andGateGarbler(leftWire, rightWire, newWire, gc.R, i, f, gc.Gates)
		}
	}

	outputOffset := len(gc.AllWires) - circuit.m
	for i := 0; i < circuit.m; i++ {
		gc.OutputShouldFlip[i] = (gc.AllWires[outputOffset+i].lsb() == 1)
	}
	return gc
}

func AndGateGarbler(leftWire, rightWire, outputWire, R Wire, i int, f FixedKeyCipher, w []Wire) {
	offset := i * 4
	a := leftWire.lsb()
	b := rightWire.lsb()

	leftWire1 := Wire{}
	leftWire1.copyFromXor(&leftWire, &R)
	rightWire1 := Wire{}
	rightWire1.copyFromXor(&rightWire, &R)
	outputWire1 := Wire{}
	outputWire1.copyFromXor(&outputWire, &R)

	//0,0 -> 0
	mw00 := f.CreateMask(leftWire, rightWire, i)
	ow00 := Wire{}
	ow00.copyFromXor(&mw00, &outputWire)
	w[offset+(a<<1+b)].copyFrom(&ow00)

	//0,1 -> 0
	b = 1 - b
	mw01 := f.CreateMask(leftWire, rightWire1, i)
	ow01 := Wire{}
	ow01.copyFromXor(&mw01, &outputWire)
	w[offset+(a<<1+b)].copyFrom(&ow01)

	//1,0 -> 0
	a = 1 - a
	b = 1 - b
	mw10 := f.CreateMask(leftWire1, rightWire, i)
	ow10 := Wire{}
	ow10.copyFromXor(&mw10, &outputWire)
	w[offset+(a<<1+b)].copyFrom(&ow10)

	//1,1 -> 1
	b = 1 - b
	mw11 := f.CreateMask(leftWire1, rightWire1, i)
	ow11 := Wire{}
	ow11.copyFromXor(&mw11, &outputWire1)
	w[offset+(a<<1+b)].copyFrom(&ow11)
}

func (g *GarblerImpl) Eval(gc *GarbledCircuit, inputs []bool) []bool {
	allWires := make([]Wire, gc.Circuit.n+gc.Circuit.q)
	f := MakeNewFixedKeyCipher(gc.Masterkey)

	for i := 0; i < gc.Circuit.n; i++ {
		if inputs[i] {
			allWires[i].copyFromXor(&gc.InputWires[i], &gc.R)
		} else {
			allWires[i].copyFrom(&gc.InputWires[i])
		}
	}
	for i := 0; i < gc.Circuit.q; i++ {
		leftWire := allWires[gc.Circuit.A[i]]
		rightWire := allWires[gc.Circuit.B[i]]
		if gc.Circuit.isXOR[i] {
			allWires[i+gc.Circuit.n].copyFromXor(&leftWire, &rightWire)
		} else {
			leftWireLsb := leftWire.lsb()
			rightWireLsb := rightWire.lsb()
			offset := i * 4
			entry := gc.Gates[offset+(leftWireLsb<<1)+rightWireLsb]
			mw := f.CreateMask(leftWire, rightWire, i)
			ow := Wire{}
			ow.copyFromXor(&mw, &entry)
			allWires[i+gc.Circuit.n] = ow
		}
	}

	outputOffset := len(gc.AllWires) - gc.Circuit.m
	outputs := make([]bool, gc.Circuit.m)
	for i := 0; i < gc.Circuit.m; i++ {
		wire := allWires[outputOffset+i]
		outputs[i] = (wire.lsb() == 1)
		if gc.OutputShouldFlip[i] {
			outputs[i] = !outputs[i]
		}
	}
	return outputs
}
