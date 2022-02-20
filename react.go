package react

import (
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type CellReactor struct {
	computeCells []ReactiveComputeCell
}

func (r *CellReactor) CreateInput(value int) InputCell {
	return &StimulusInputCell{value: value, notifyWorld: func() { r.NotifyComputeCells() }}
}

func (r *CellReactor) CreateCompute1(cell Cell, cb func(int) int) ComputeCell {
	r.computeCells = append(r.computeCells, ReactiveComputeCell{
		cells:           []Cell{cell},
		computeFunction: func() int { return cb(cell.Value()) },
		prevComputed:    cb(cell.Value()),
	})
	return &r.computeCells[len(r.computeCells)-1]
}

func (r *CellReactor) CreateCompute2(cell Cell, cell2 Cell, cb func(int, int) int) ComputeCell {
	r.computeCells = append(r.computeCells, ReactiveComputeCell{
		cells:           []Cell{cell, cell2},
		computeFunction: func() int { return cb(cell.Value(), cell2.Value()) },
		prevComputed:    cb(cell.Value(), cell2.Value()),
	})
	return &r.computeCells[len(r.computeCells)-1]
}

func (r *CellReactor) NotifyComputeCells() {
	for _, cc := range r.computeCells {
		// If the computed value has changed, then this will trigger the callbacks
		cc.Value()
	}
}

func New() Reactor {
	return &CellReactor{}
}

type StimulusInputCell struct {
	value       int
	notifyWorld func()
}

func (ic *StimulusInputCell) Value() int {
	return ic.value
}

func (ic *StimulusInputCell) SetValue(value int) {
	prev := ic.value
	ic.value = value
	// Have to do it this way because once we notify the world
	// the value of the InputCell needs to already be changed
	if prev != value {
		ic.notifyWorld()
	}
}

type Cancel struct {
	cbID int
	cc   *ReactiveComputeCell
}

func (c *Cancel) Cancel() {
	delete(c.cc.cbs, c.cbID)
}

type ReactiveComputeCell struct {
	prevComputed    int
	cbs             map[int]func(int)
	cells           []Cell
	computeFunction func() int
	Cancel          Canceler
}

func (cc *ReactiveComputeCell) Value() int {
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

func (cc *ReactiveComputeCell) AddCallback(cb func(value int)) Canceler {
	if cc.cbs == nil {
		cc.cbs = make(map[int]func(int))
	}
	cbID := rand.Int()
	cc.cbs[cbID] = cb
	return &Cancel{cbID: cbID, cc: cc}
}
