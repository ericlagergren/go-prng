package prng

type Twister interface {
	Seed()
	SeedArray()
	Warmup()
	Int()
	IntN()
}
