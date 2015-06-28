# Go-PRNG

## What

Go-PRNG is a collection of pseudo-random number generators. Each PRNG is
given an easy-to-use interface, and some libraries implement math.Rand's
`Source` interface.

## Why

I was generating random strings and was accidentally using the same seed,
causing me to generate the same strings. So, I went to use a quick and
dirty XOR solution, and somehow ended up translating these PRNGs into Go.

## How

#### xorshift

```go
import (
	"fmt"
	
	prng "github.com/ericlagerg/go-prng/xorshift"
)

func main() {
	s := new(prng.Shift128Plus)
	s.Seed()
	fmt.Println(s.Next())
}
```

### mersennes twister

```go
import (
	"fmt"
	
	prng "github.com/ericlagerg/go-prng/mersenne_twister_64"
)

func main() {
	m := NewMersenne()
	m.Seed(787094841)
	fmt.Println(m.Int64())
}
```

## Documentation

[godoc.org] [docs]

## License

It depends on the algorithm.

[Apache 2.0] [license].


[docs]:     https://godoc.org/github.com/EricLagerg/go-prng
[license]:  https://github.com/EricLagerg/go-prng/blob/master/apache.txt
