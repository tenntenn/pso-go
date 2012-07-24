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
	"./float64"
	"fmt"
	"testing"
)

func TestFloat64(t *testing.T) {

	// x + y - 3 = 0
	// 2x + 5y - 9 = 0
	f := func(values float64.Float64Array) float64.EvalValue {
		arry := []float64(values)
		x := arry[0]
		y := arry[1]
		
		return (x + y - 3) + (2 * x + 5 * y - 9)
	}

	// Create particles
	const count = 100
	particles := make([]*Particle, count)
	a := float64.Float64Array([]float64{0.0, 0.0})
	for i := range particles {
		position := a.Random()
		velocity := a.Random()
		valuesRange := float64.NewRangeWithInf(len(a))
		particles[i] = NewParticle(position, velocity, valuesRange)
	}

	// Create a solver
	w := float64.Float64Array([]float64{0.9, 0.9})
	c1 := float64.Float64Array([]float64{0.9, 0.9})
	c2 := float64.Float64Array([]float64{0.9, 0.9})
	param := NewParam(w, c1, c2)
	solver := NewSolver(float64.TargetFunc(f), particles, param)

	// Start solver process
	done := make(chan bool)
	defer func() {
		close(done)
	}()
	go solver.Start(chan <-bool(done))

	for {
		if solver.Best() != nil {
			bestValue := solver.TargetFunc().Eval(solver.Best())
			endValue := float64.EvalValue(0.0001)
			if bestValue.CompareTo(endValue) < 0 {
				done <- true
				fmt.Printf("%v\n", solver.Best())
				break
			}
		}
	}
}
