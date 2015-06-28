/*
   Before using, initialize the state by using init_genrand(seed)
   or init_by_array(init_key, key_length).

   Copyright (C) 2015, Eric Lagergren
   Copyright (C) 1997 - 2002, Makoto Matsumoto and Takuji Nishimura,
   All rights reserved.

   Redistribution and use in source and binary forms, with or without
   modification, are permitted provided that the following conditions
   are met:

     1. Redistributions of source code must retain the above copyright
        notice, this list of conditions and the following disclaimer.

     2. Redistributions in binary form must reproduce the above copyright
        notice, this list of conditions and the following disclaimer in the
        documentation and/or other materials provided with the distribution.

     3. The names of its contributors may not be used to endorse or promote
        products derived from this software without specific prior written
        permission.

   THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
   "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
   LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
   A PARTICULAR PURPOSE ARE DISCLAIMED.  IN NO EVENT SHALL THE COPYRIGHT OWNER OR
   CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL,
   EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO,
   PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR
   PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF
   LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING
   NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS
   SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
*/

package twister32

import "crypto/rand"

const (
	N         = 624
	M         = 397
	MatrixA32 = 0x9908b0df
	UMask     = 0x80000000 // Most significant w-r bits
	LMask     = 0x7fffffff // Least significant r bits
)

func mixbits(u, v uint32) uint32 {
	return (((u) & UMask) | ((v) & LMask))
}

func twist(u, v uint32) uint32 {
	m := mixbits(u, v) >> 1

	if v&1 != 0 {
		return m ^ MatrixA32
	}
	return m ^ 0
}

// MT19937 holds the state of a Mersenne Twister.
type MT19937 struct {
	state       [N]uint32
	left, initf int
	next        [N]uint32
	nptr        int
}

// NewMersennePrime returns a seeded, initialized MT19937, seeded
// using a large from from 'crypto/rand'.
// Will panic if rand.Prime returns an error.
func NewMersennePrime32() *MT19937 {
	m := New32()
	prime, err := rand.Prime(rand.Reader, 64)
	if err != nil {
		panic(err)
	}
	m.Seed(prime.Int64())
	return m
}

// NewMersenne returns a seeded, initialized MT19937.
func NewMersenne32(seed int64) *MT19937 {
	m := New32()
	m.Seed(seed)
	return m
}

// New returns an unseeded, initialized MT19937.
func New32() *MT19937 {
	return &MT19937{
		left: 1,
	}
}

var globalRand = NewMersenne32(1)

// Seed seeds the global MT19937 with the given seed.
func Seed32(seed int64) {
	globalRand.Seed(seed)
}

// Warmup warms up the PRNG by running a 64-bit prime number of iterations.
func Warmup() {
	globalRand.Warmup()
}

// Int returns a random, non-negative int.
func Int() int {
	return globalRand.Int()
}

// IntN generates a random number on [0,0x7fffffff]-interval within
// the given range, n, from the global MT19937.
func IntN(n int) int {
	return globalRand.IntN(n)
}

// Real1 generates a random number on [0,1]-real-interval from the global
// MT19937.
func Real1() float64 {
	return globalRand.Real1()
}

// Real2 generates a random number on [0,1)-real-interval from the global
// MT19937.
func Real2() float64 {
	return globalRand.Real2()
}

// Real3 generates a random number on (0,1)-real-interval from the global
// MT19937.
func Real3() float64 {
	return globalRand.Real3()
}

// Seed seeds an MT19937 with the given seed.
func (m *MT19937) Seed(seed int64) {

	m.state[0] = uint32(seed) & 0xffffffff
	for j := 1; j < N; j++ {
		m.state[j] = (1812433253*(m.state[j-1]^(m.state[j-1]>>30)) + uint32(j))
		// See Knuth TAOCP Vol2. 3rd Ed. P.106 for multiplier.
		// In the previous versions, MSBs of the seed affect
		// only MSBs of the array state[].
		// 2002/01/09 modified by Makoto MatsUMaskoto
		m.state[j] &= 0xffffffff // for >32 bit machines
	}
	m.left = 1
	m.initf = 1
}

// Warmup warms up the PRNG by running a 64-bit prime number of iterations.
func (m *MT19937) Warmup() {
	prime, err := rand.Prime(rand.Reader, 64)
	if err != nil {
		panic(err)
	}
	n := prime.Uint64()

	for i := uint64(0); i < n; i++ {
		m.Int32()
	}
}

// SeedArray seeds an MT19937 with the given array.
func (m *MT19937) SeedArray(initKey [N]uint32) {

	m.Seed(19650218)

	// i := 1
	// j := 0

	var i, j uint32
	i = 1

	k := N
	if N <= len(initKey) {
		k = len(initKey)
	}

	for ; k != 0; k-- {
		m.state[i] = (m.state[i] ^ ((m.state[i-1] ^ (m.state[i-1] >> 30)) * 1664525)) + initKey[j] + j // non linear

		m.state[i] &= 0xffffffff // for WORDSIZE > 32 machines
		i++
		j++
		if i >= N {
			m.state[0] = m.state[N-1]
			i = 1
		}
		if j >= uint32(len(initKey)) {
			j = 0
		}
	}
	for k = N - 1; k != 0; k-- {
		m.state[i] = (m.state[i] ^ ((m.state[i-1] ^ (m.state[i-1] >> 30)) * 1566083941)) - i // non linear

		m.state[i] &= 0xffffffff // for WORDSIZE > 32 machines
		i++
		if i >= N {
			m.state[0] = m.state[N-1]
			i = 1
		}
	}

	m.state[0] = 0x80000000 // MSB is 1; assuring non-zero initial array
	m.left = 1
	m.initf = 1
}

func (m *MT19937) NextState() {
	p := m.state

	// If Seed has not been called, a default seed is used.
	if m.initf == 0 {
		m.Seed(5489)
	}

	m.left = N
	m.next = m.state

	var i, j int

	for i, j = 0, N-M+1; j-1 != 0; i, j = i+1, j-1 {
		p[i] = p[M] ^ twist(p[0], p[1])
	}

	for i, j = 0, M; j-1 != 0; i, j = i+1, j-1 {
		p[i] = p[i-(M-N)] ^ twist(p[0], p[1])
	}

	p[i] = p[i-(M-N)] ^ twist(p[0], m.state[0])
}

// Int32 generates a random number on [0,0xffffffff]-interval
func (m *MT19937) Int32() uint32 {
	var y uint32

	m.left--
	if m.left == 0 {
		m.NextState()
	}

	m.nptr++
	y = uint32(m.next[m.nptr])

	// Tempering
	y ^= (y >> 11)
	y ^= (y << 7) & 0x9d2c5680
	y ^= (y << 15) & 0xefc60000
	y ^= (y >> 18)

	return y
}

// Int31 generates a random number on [0,0x7fffffff]-interval
func (m *MT19937) Int31() int32 {
	var y uint32

	m.left--
	if m.left == 0 {
		m.NextState()
	}

	m.nptr++
	y = uint32(m.next[m.nptr])

	// Tempering
	y ^= (y >> 11)
	y ^= (y << 7) & 0x9d2c5680
	y ^= (y << 15) & 0xefc60000
	y ^= (y >> 18)

	return int32(y >> 1)
}

// Int returns a random, non-negative int.
func (m *MT19937) Int() int {
	u := uint(m.Int32())
	return int(u << 1 >> 1)
}

// IntN generates a random number on [0,0x7fffffff]-interval within
// the given range, n.
func (m *MT19937) IntN(n int) int {
	return int(m.Int32() % uint32(n))
}

// Real1 generates a random number on [0,1]-real-interval from the given
// MT19937.
func (m *MT19937) Real1() float64 {
	var y uint32

	m.left--
	if m.left == 0 {
		m.NextState()
	}

	m.nptr++
	y = uint32(m.next[m.nptr])

	// Tempering
	y ^= (y >> 11)
	y ^= (y << 7) & 0x9d2c5680
	y ^= (y << 15) & 0xefc60000
	y ^= (y >> 18)

	/* divided by 2^32-1 */
	return float64(y) * (1.0 / 4294967295.0)
}

// Real2 generates a random number on [0,1)-real-interval from the given
// MT19937.
func (m *MT19937) Real2() float64 {
	var y uint32

	m.left--
	if m.left == 0 {
		m.NextState()
	}

	m.nptr++
	y = uint32(m.next[m.nptr])

	// Tempering
	y ^= (y >> 11)
	y ^= (y << 7) & 0x9d2c5680
	y ^= (y << 15) & 0xefc60000
	y ^= (y >> 18)

	/* divided by 2^32 */
	return float64(y) * (1.0 / 4294967296.0)
}

// Real3 generates a random number on (0,1)-real-interval from the given
// MT19937.
func (m *MT19937) Real3() float64 {
	var y uint32

	m.left--
	if m.left == 0 {
		m.NextState()
	}

	m.nptr++
	y = uint32(m.next[m.nptr])

	// Tempering
	y ^= (y >> 11)
	y ^= (y << 7) & 0x9d2c5680
	y ^= (y << 15) & 0xefc60000
	y ^= (y >> 18)

	// divided by 2^32
	return (float64(y) + 0.5) * (1.0 / 4294967296.0)
}

func (m *MT19937) Res53() float64 {
	a := m.Int32() >> 5
	b := m.Int32() >> 6

	// It's written this way in C.
	return (float64(a)*67108864.0 + float64(b)) * (1.0 / 9007199254740992.0)
}
