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

package float64

import (
	"math/rand"
	"time"
	"reflect"
	"../pso"
)

// An implement of pso.EvalValue for float64
type EvalValue float64

// Implement pso.EvalValue.CompareTo
func (e EvalValue) CompareTo(evalValue pso.EvalValue) int {

	e2, ok := c.(EvalValue)
	if !ok {
		panic("Cannot compare to another type.")
	}
	
	switch {
	case float64(e) > float64(e2):
		return 1
	case float64(e) < float64(e2):
		return -1
	}
	
	return 0
}

// An implement of pso.Values for []float64
type Float64Array []float64

// implement pso.Values.Add()
func (arry Float64Array) Add(values Values) Values {

	arry2, ok := values.(Float64Array)
	if !ok {
		panic("Cannot add another type.")
	}

	a := []float64(arry)
	a2 := []float64(arry2)
	for i := range a {
		a[i] += a2[i]
	}
	
	return arry
}

// implement pso.Values.Sub()
func (arry Float64Array) Sub(values Values) Values {

	arry2, ok := values.(Float64Array)
	if !ok {
		panic("Cannot substract another type.")
	}

	a := []float64(arry)
	a2 := []float64(arry2)
	for i := range a {
		a[i] -= a2[i]
	}
	
	return arry
}

// implement pso.Values.Mul()
func (arry Float64Array) Mul(values Values) Values {

	arry2, ok := values.(Float64Array)
	if !ok {
		panic("Cannot multiply another type.")
	}

	a := []float64(arry)
	a2 := []float64(arry2)
	for i := range a {
		a[i] *= a2[i]
	}
	
	return arry
}

// implement pso.Values.Div()
func (arry Float64Array) Div(values Values) Values {

	arry2, ok := values.(Float64Array)
	if !ok {
		panic("Cannot divide another type.")
	}

	a := []float64(arry)
	a2 := []float64(arry2)
	for i := range a {
		a[i] /= a2[i]
	}
	
	return arry
}

// implement pso.Values.Cone()
func (arry Float64Array) Clone() Values {

	a := []float64(arry)
	a2 := make([]float64, len(a), cap(a))
	copy(a, a2)
	
	return Float64Array(a2)
}

var rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
// implement pso.Values.Random()
func (arry Float64Array) Random() Values {

	a := []float64(arry)
	a2 := make([]float64, len(a), cap(a))
	for i := range a2 {
		a2[i] = rnd.Float64()
	}

	return Float64Array(a2)
}

// an implement pso.Range for []float64
type Range struct {
	min Float64Array
	max Float64Array
}

// Create a new range.
func NewRange(min, max Float64Array) *Range {

	switch {
	case min == nil:
		panic("min cannot be nil.")
	case max == nil:
		panic("max cannot be nil.")
	case len(min) != len(max):
		panic("length of min and max have to be same.")
	}

	&Range{min, max}
}

// Create a new range with inf
func NewRangeWithInf(length uint) *Range {

	min := make([]float64, length)
	max := make([]float64, length)
	for i := 0; i < length; i++ {
		min[i] = math.Inf(-1)
		max[i] = math.Inf(1)
	}

	return NewRange(Float64Array(min), Float64Array(max))
}

// Get either values are in this range or not.
func (r *Range) In(values Values) bool {

	arry, ok := values.(Float64Array)
	if !ok {
		panic("values have to be float64.Float64Array")
	}

	if len(arry) != len(r.min) {
		panic("length of values have to be same with minx and max.")
	}

	a := []float64(arry)
	min := []float64(r.min)
	max := []float64(r.max)
	for i := range a {
		if a[i] < min[i] || a[i] > max[i] {
			return false
		}
	}

	return true
}

// Get type of min and max
func (r *Range) Type() reflect.Type {
	return reflect.TypeOf(r.min)
}

// an implement pso.TargetFunc
type  TargetFunc func(values Float64Array) EvalValue
func (f TargetFunc) Eval(values pso.Values) pso.EvalValue {
	if values != nil {
		arry, ok := values.(Float64Array)
		if !ok {
			panic("values have to be Float64Array.")
		}
		
		return f(arry)
	} else {
		return f(nil)
	}
}