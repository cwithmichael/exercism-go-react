package react

import (
	"math/rand"
	"time"
)

type MyReactor struct {
	computeCells []*MyComputeCell
}

func (r *MyReactor) CreateInput(value int) InputCell {
	return &MyInputCell{value: value, notifyWorld: func() { r.NotifyComputeCells() }}
}

func (r *MyReactor) CreateCompute1(cell Cell, cb func(int) int) ComputeCell {
	r.computeCells = append(r.computeCells, &MyComputeCell{
		cells:           []Cell{cell},
		computeFunction: func() int { return cb(cell.Value()) },
		prevComputed:    cb(cell.Value()),
	})
	return r.computeCells[len(r.computeCells)-1]
}

func (r *MyReactor) CreateCompute2(cell Cell, cell2 Cell, cb func(int, int) int) ComputeCell {
	r.computeCells = append(r.computeCells, &MyComputeCell{
		cells:           []Cell{cell, cell2},
		computeFunction: func() int { return cb(cell.Value(), cell2.Value()) },
		prevComputed:    cb(cell.Value(), cell2.Value()),
	})
	return r.computeCells[len(r.computeCells)-1]
}

func (r *MyReactor) NotifyComputeCells() {
	for _, cc := range r.computeCells {
		// If the computed value has changed, then this will trigger the callbacks
		cc.Value()
	}
}

func New() Reactor {
	rand.Seed(time.Now().UnixNano())
	return &MyReactor{}
}

type MyInputCell struct {
	value       int
	notifyWorld func()
}

func (ic *MyInputCell) Value() int {
	return ic.value
}

func (ic *MyInputCell) SetValue(value int) {
	prev := ic.value
	ic.value = value
	// Have to do it this way because once we notify the world
	// the value of the InputCell needs to already be changed
	if prev != value {
		ic.notifyWorld()
	}
}

type MyCancel struct {
	cbId int
	cc   *MyComputeCell
}

func (c *MyCancel) Cancel() {
	delete(c.cc.cbs, c.cbId)
}

type MyComputeCell struct {
	prevComputed    int
	cbs             map[int]func(int)
	cells           []Cell
	computeFunction func() int
	MyCancel        MyCancel
}

func (cc *MyComputeCell) Value() int {
	computedValue := cc.computeFunction()
	// Only call the callbacks if the computed value has changed
	if cc.prevComputed != computedValue {
		for _, cb := range cc.cbs {
			cb(computedValue)
		}
	}
	cc.prevComputed = computedValue
	return computedValue
}

func (cc *MyComputeCell) AddCallback(cb func(value int)) Canceler {
	if cc.cbs == nil {
		cc.cbs = make(map[int]func(int))
	}
	cbId := rand.Int()
	cc.cbs[cbId] = cb
	return &MyCancel{cbId: cbId, cc: cc}
}
