// Copyright 2015, Eric Lagergren.

// Original Copyright for the C version of the code is below.

// /* ***************************************************************************** */
// /* Copyright:      Francois Panneton and Pierre L'Ecuyer, University of Montreal */
// /*                 Makoto Matsumoto, Hiroshima University                        */
// /* Notice:         This code can be used freely for personal, academic,          */
// /*                 or non-commercial purposes. For commercial purposes,          */
// /*                 please contact P. L'Ecuyer at: lecuyer@iro.UMontreal.ca       */
// /* ***************************************************************************** */
package prng

const (
	W  = 32
	R  = 16
	P  = 0
	M1 = 13
	M2 = 9
	M3 = 5

	FACT = 2.32830643653869628906e-10
)

func mat0pos(t, v uint64) uint64 {
	return (v ^ (v >> t))
}

func mat0neg(t, v uint64) uint64 {
	return (v ^ (v << (t)))
}

func mat3neg(t, v uint64) uint64 {
	return (v ^ (v << (t)))
}

func mat4neg(t, b, v uint64) uint64 {
	return (v ^ ((v << (t)) & b))
}

var (
	state_i    = 0
	STATE      [R]uint64
	z0, z1, z2 uint64
)

func InitWELLRNG512a(init []uint64) {
	state_i = 0
	for j := 0; j < R; j++ {
		STATE[j] = init[j]
	}
}

func WELLRNG512a() float64 {
	z0 = STATE[(state_i+15)&0x0000000f]
	z1 = mat0neg(16, STATE[state_i]) ^ mat0neg(15, STATE[(state_i+M1)&0x0000000f])
	z2 = mat0pos(11, STATE[(state_i+M2)&0x0000000f])
	STATE[state_i] = z1 ^ z2
	STATE[(state_i+15)&0x0000000f] = mat0neg(2, z0) ^ mat0neg(18, z1) ^ mat3neg(28, z2) ^ mat4neg(5, 0xda442d24, STATE[state_i])
	state_i = (state_i + 15) & 0x0000000f
	return float64(STATE[state_i]) * FACT
}
