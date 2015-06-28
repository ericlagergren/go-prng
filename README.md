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

```
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

```
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
