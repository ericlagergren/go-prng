/*
	MT19937-64

   Copyright (C) 2015, Eric Lagergren
   Copyright (C) 2004, Makoto Matsumoto and Takuji Nishimura,
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

package twister64

import "crypto/rand"

const (
	NN      = 312
	MM      = 156
	MatrixA = 0xB5026F5AA96619E9
	UM      = 0xFFFFFFFF80000000 // Most significant 33 bits
	LM      = 0x7FFFFFFF         // Least significant 31 bits
)

// MT19937 holds the state of a Mersenne Twister.
type MT19937 struct {
	mt  [NN]uint64
	mti uint64
}

// NewMersennePrime returns a seeded, initialized MT19937, seeded
// using a large from from 'crypto/rand'.
// Will panic if rand.Prime returns an error.
func NewMersennePrime() *MT19937 {
	m := New()
	prime, err := rand.Prime(rand.Reader, 64)
	if err != nil {
		panic(err)
	}
	m.Seed(prime.Int64())
	return m
}

// NewMersenne returns a seeded, initialized MT19937.
func NewMersenne(seed int64) *MT19937 {
	m := New()
	m.Seed(seed)
	return m
}

// New returns an unseeded, initialized MT19937.
func New() *MT19937 {
	return &MT19937{
		mti: NN + 1,
	}
}

var globalRand = NewMersenne(1)

// Seed seeds the global MT19937 with the given seed.
func Seed(seed int64) {
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

// Int64 generates a random number on [0, 2^64-1]-interval from the global
// MT19937.
func Int64() uint64 {
	return globalRand.Int64()
}

// Int63 generates a random number on [0, 2^63-1]-interval from the global
// MT19937.
func Int63() int64 {
	return globalRand.Int63()
}

// IntN generates a random number on [0, 2^64-1]-interval within
// the given range, n, from the global MT19937.
func IntN(n uint64) uint64 {
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
	m.mt[0] = uint64(seed)
	for m.mti = 1; m.mti < NN; m.mti++ {
		m.mt[m.mti] = (6364136223846793005*(m.mt[m.mti-1]^(m.mt[m.mti-1]>>62)) + m.mti)
	}
}

// SeedArray seeds an MT19937 with the given array.
func (m *MT19937) SeedArray(initKey [NN]uint64) {
	var i, j uint64

	m.Seed(19650218)
	i = 1
	j = 0

	k := NN
	if NN <= len(initKey) {
		k = len(initKey)
	}

	for ; k != 0; k-- {
		m.mt[i] = (m.mt[i] ^ ((m.mt[i-1] ^ (m.mt[i-1] >> 62)) * 3935559000370003845)) + initKey[j] + j /* non linear */

		i++
		j++
		if i >= NN {
			m.mt[0] = m.mt[NN-1]
			i = 1
		}
		if j >= uint64(len(initKey)) {
			j = 0
		}
	}

	for k = NN - 1; k != 0; k-- {
		m.mt[i] = (m.mt[i] ^ ((m.mt[i-1] ^ (m.mt[i-1] >> 62)) * 2862933555777941757)) - i /* non linear */
		i++
		if i >= NN {
			m.mt[0] = m.mt[NN-1]
			i = 1
		}
	}

	m.mt[0] = 1 << 63 /* MSB is 1; assuring non-zero initial array */
}

// Warmup warms up the PRNG by running a 64-bit prime number of iterations.
func (m *MT19937) Warmup() {
	prime, err := rand.Prime(rand.Reader, 64)
	if err != nil {
		panic(err)
	}
	n := prime.Uint64()

	for i := uint64(0); i < n; i++ {
		m.Int64()
	}
}

// Int returns a random, non-negative int.
func (m *MT19937) Int() int {
	u := uint(m.Int64())
	return int(u << 1 >> 1)
}

// Int64 generates a random number on [0, 2^64-1]-interval.
func (m *MT19937) Int64() uint64 {

	var (
		i int
		x uint64
	)

	mag01 := [2]uint64{0, MatrixA}

	if m.mti >= NN { /* generate NN words at one time */

		/* if Seed() has not been called, */
		/* a default initial seed is used */
		if m.mti == NN+1 {
			m.Seed(5489)
		}

		for i = 0; i < NN-MM; i++ {
			x = (m.mt[i] & UM) | (m.mt[i+1] & LM)
			m.mt[i] = m.mt[i+MM] ^ (x >> 1) ^ mag01[(int)(x&1)]
		}
		for ; i < NN-1; i++ {
			x = (m.mt[i] & UM) | (m.mt[i+1] & LM)
			m.mt[i] = m.mt[i+(MM-NN)] ^ (x >> 1) ^ mag01[(int)(x&1)]
		}
		x = (m.mt[NN-1] & UM) | (m.mt[0] & LM)
		m.mt[NN-1] = m.mt[MM-1] ^ (x >> 1) ^ mag01[(int)(x&1)]

		m.mti = 0
	}

	x = m.mt[m.mti]
	m.mti++

	x ^= (x >> 29) & 0x5555555555555555
	x ^= (x << 17) & 0x71D67FFFEDA60000
	x ^= (x << 37) & 0xFFF7EEE000000000
	x ^= (x >> 43)

	return x
}

// Int63 generates a random number on [0, 2^63-1]-interval from the given
// MT19937.
func (m *MT19937) Int63() int64 {
	return int64((m.Int64() << 1 >> 1))
}

// IntN generates a random number on [0, 2^64-1]-interval within
// the given range, n.
func (m *MT19937) IntN(n uint64) uint64 {
	return m.Int64() % n
}

// Real1 generates a random number on [0,1]-real-interval from the given
// MT19937.
func (m *MT19937) Real1() float64 {
	return float64(m.Int64()>>11) * (1.0 / 9007199254740991.0)
}

// Real2 generates a random number on [0,1)-real-interval from the given
// MT19937.
func (m *MT19937) Real2() float64 {
	return float64(m.Int64()>>11) * (1.0 / 9007199254740992.0)
}

// Real3 generates a random number on (0,1)-real-interval from the given
// MT19937.
func (m *MT19937) Real3() float64 {
	return (float64(m.Int64()>>12) + 0.5) * (1.0 / 4503599627370496.0)
}
