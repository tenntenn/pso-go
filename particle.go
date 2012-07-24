/*
pso-go - PSO (Particle Swarm Optimization) library for Go.
https://github.com/tenntenn/pso-go

Copyright (c) 2012, Takuya Ueda.
All rights reserved.

Redistribution and use in source and binary forms, with or without modification,
are permitted provided that the following conditions are met:

* Redistributions of source code must retain the above copyright notice,
  this list of conditions and the following disclaimer.
* Redistributions in binary form must reproduce the above copyright notice,
  this list of conditions and the following disclaimer in the documentation
  and/or other materials provided with the distribution.
* Neither the name of the author nor the names of its contributors may be used
  to endorse or promote products derived from this software
  without specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND
ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED
WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR
ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES
(INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES;
LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON
ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
(INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS
SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
*/

package pso

import (
	"reflect"
)

// Evaluated value by target function.
// It provides compareble function.
type EvalValue interface {
	// Compare with other value.
	// It returns int value likes java.lang.Compareble in Java. 
	CompareTo(v EvalValue) int
}

// A type for particle position vector and velocity vector.
// It provides four arithmetic operations.
type Values interface {
	// Add values and return own values
	Add(values Values) Values
	// Subtrac values and return own values
	Sub(values Values) Values
	// Multiply values and return own values
	Mul(values Values) Values
	// Divide by values and return own values
	Div(values Values) Values
	// Create a clone values
	Clone() Values
	// Return random values
	Random() Values
}

// Range of values.
type Range interface {
	// Either values is in this range or not.	
	In(values Values) bool
	// Type of target values.
	Type() reflect.Type
}

// It is type of particle which find a result of optimization.
type Particle struct {
	// Current position
	position Values
	// Current velocity
	velocity Values
	// range of position
	valuesRange Range
	// evaluated value by target function
	evalValue EvalValue
	// local best of this particle
	best Values
	// a channel which is used for communication with solver
	ch chan<- Values
}

// Create a new particle.
func NewParticle(position, velocity Values, valuesRange Range) *Particle {

	switch {
	case position == nil:
		panic("position cannot be nil.")
	case velocity == nil:
		panic("velocity cannot be nil.")
	case reflect.TypeOf(position) != reflect.TypeOf(velocity):
		panic("position and velocity have to be same type.")
	case reflect.TypeOf(position) != valuesRange.Type():
		panic("type of position and valuesRange.Type() have to be same.")
	}

	ch := chan <-Values(make(chan Values))

	return &Particle{position, velocity, valuesRange, nil, nil, ch}
}

// Get the position of particle on the solution space.
func (p *Particle) Position() Values {
	return p.position
}

// Get the position of particle.
func (p *Particle) Velocity() Values {
	return p.velocity
}

// Get the range of position.
func (p *Particle) Range() Range {
	return p.valuesRange
}

// Get the evaluated value by target function.
func (p *Particle) EvalValue() EvalValue {
	return p.evalValue
}

// Get the channel which is used by communication with the solver.
func (p *Particle) Channel() chan<-Values {
	return p.ch
}

func (p *Particle) Best() Values {
	return p.best
}

// Start tparticle process.
// If ch is closed, process stop.
func (p *Particle) Start(f TargetFunc, param *Param, ch <-chan Values) {

	// Update global best
	globalBest, isRun := <- ch
	go func() {
		for {
			globalBest, isRun = <- ch
		}
	}()
	
	for isRun {
		
		// update position with velocity
		newPosition := p.position.Clone().Add(p.velocity)
		if p.valuesRange.In(newPosition) {
			p.position = newPosition
		}
		
		// update velocity
		r1 := p.velocity.Random()
		r2 := p.velocity.Random()
		// v <- w*v
		p.velocity.Mul(param.W())
		// r1 <- r1*c1*(ownBest - x)
		r1.Mul(param.C1()).Mul(p.best.Clone().Sub(p.position))
		// r2 <- r2*c2*(best - x)
		r2.Mul(param.C2()).Mul(globalBest.Clone().Sub(p.position))
		// v <- w*v + c1*r1*(ownBest - x) + c2*r2*(best - x)
		p.velocity.Add(r1).Add(r2)

		// send best value to solver
		p.evalValue = f.Eval(p.position)
		bestValue := f.Eval(p.best)
		if p.evalValue.CompareTo(bestValue) < 0 {
			p.best = p.position
			p.ch <- p.position
		}
		
	}
}