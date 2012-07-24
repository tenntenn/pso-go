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

// Target function
type TargetFunc interface {
	Eval(values Values) EvalValue
	Type() reflect.Type
}

// Params of a solver.
type Param struct {
	w  Values
	c1 Values
	c2 Values
}

// Create a new Param.
func NewParam(w, c1, c2 Values) *Param {

	switch {
	case w == nil:
		panic("w cannot be nil.")
	case c1 == nil:
		panic("c1 cannot be nil.")
	case c2 == nil:
		panic("c2 cannot be nil.")
	case reflect.TypeOf(w) != reflect.TypeOf(c1) ||
			reflect.TypeOf(w) != reflect.TypeOf(c2) :
		panic("all parameter type have to be same.")
	}

	return &Param{w, c1, c2}
}

// Get w of param
func (p *Param) W() Values {
	return p.w
}

// Get c1 of param
func (p *Param) C1() Values {
	return p.c1
}

// Get c2 of param
func (p *Param) C2() Values {
	return p.c2
}

// Solver of PSO
type Solver struct {
	f      TargetFunc
	particles []*Particle
	param *Param
	best   Values
}

// Create a new solver.
func NewSolver(f TargetFunc, particles []*Particle, param *Param) *Solver {

	if particles == nil {
		panic("particles cannot be nil.")
	} else if len(particles) <= 0 {
		panic("The number of particles have to take more than 0.")
	}

	return &Solver{f, particles, param, nil}
}

// Get the target function.
func (s *Solver) TargetFunc() TargetFunc {
	return s.f
}

// Get the particles.
// It is a particles array which is copy of original one.
func (s *Solver) Particles() []*Particle {

	cpy := make([]*Particle, len(s.particles))
	copy(cpy, s.particles)

	return s.particles
}

// Get the parameter of the solver.
func (s *Solver) Param() *Param {
	return s.param
}

// Get best values.
// Before start steps it takes nil.
func (s *Solver) Best() Values {
	return s.best
}

// Start solver process.
// If it catches true value through ch, process stops.
func (s *Solver) Start(ch <-chan bool) {

	bestCh := make(chan Values)
	defer func() {
		close(bestCh)
	}()

	// start particles process
	for _, p := range s.particles {
		if p == nil {
			continue
		}

		go p.Start(s.f, s.param, (<-chan Values)(bestCh))
	}

	// waiting for  done signal
	done := false
	go func() {
		for !done {
			done = <-ch
		}
	}()
	
	var bestValue EvalValue
	for {
		localBest := <-bestCh
		localBestValue := s.TargetFunc().Eval(localBest)
		if bestValue == nil || localBestValue.CompareTo(bestValue) < 0 {
			s.best = localBest
			bestValue = localBestValue
			for _, p := range s.particles {
				if p == nil {
					continue
				}
				
				// send global best to particle
				p.Channel() <- s.best
			}
		}

		if done {
			for _, p := range s.particles {
				if p == nil {
					continue
				}
				
				close(p.Channel())
			}
			break
		}
	}
}