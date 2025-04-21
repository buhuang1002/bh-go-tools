package bigfloat

import (
	"fmt"
	"math/big"
)

// BigIntMaxSerializedLen is the max length of a byte slice representing a CBOR serialized big.
const BigIntMaxSerializedLen = 128

func NewInt64(x int64) *big.Float {
	return new(big.Float).SetInt64(x)
}

func NewBigInt(x *big.Int) *big.Float {
	return new(big.Float).SetInt(x)
}

func NewFloat(x float64) *big.Float {
	return big.NewFloat(x)
}

func Zero() *big.Float {
	return NewFloat(0)
}

// MustFromString convers dec string into big integer and panics if conversion
// is not sucessful.
func MustFromString(s string) *big.Float {
	v, err := FromString(s)
	if err != nil {
		panic(err)
	}
	return v
}

func FromString(s string) (*big.Float, error) {
	v, ok := big.NewFloat(0).SetString(s)
	if !ok {
		return nil, fmt.Errorf("failed to parse string as a big int")
	}

	return v, nil
}

func FromBigInt(x *big.Int) *big.Float {
	return Zero().SetInt(x)
}

func Product(ints ...*big.Float) *big.Float {
	p := NewFloat(1)
	for _, i := range ints {
		p = Mul(p, i)
	}
	return p
}

func Mul(a, b *big.Float) *big.Float {
	return NewFloat(0).Mul(a, b)
}

func Quo(a, b *big.Float) *big.Float {
	return NewFloat(0).Quo(a, b)
}

func Add(a, b *big.Float) *big.Float {
	return NewFloat(0).Add(a, b)
}

func Sum(ints ...*big.Float) *big.Float {
	sum := Zero()
	for _, i := range ints {
		sum = Add(sum, i)
	}
	return sum
}

func Subtract(num1 *big.Float, ints ...*big.Float) *big.Float {
	sub := num1
	for _, i := range ints {
		sub = Sub(sub, i)
	}
	return sub
}

func Sub(a, b *big.Float) *big.Float {
	return NewFloat(0).Sub(a, b)
}

func Neg(a *big.Float) *big.Float {
	return NewFloat(0).Neg(a)
}

func Max(x, y *big.Float) *big.Float {
	if Eq(x, Zero()) && Eq(x, y) {
		if x.Sign() != 0 {
			return y
		}
		return x
	}

	if Gt(x, y) {
		return x
	}
	return y
}

func Min(x, y *big.Float) *big.Float {
	if Eq(x, Zero()) && Eq(x, y) {
		if x.Sign() != 0 {
			return x
		}
		return y
	}

	if Lt(x, y) {
		return x
	}
	return y
}

func Cmp(a, b *big.Float) int {
	return a.Cmp(b)
}

func Eq(a, b *big.Float) bool {
	return Cmp(a, b) == 0
}

func Ne(a, b *big.Float) bool {
	return Cmp(a, b) != 0
}

func Gt(a, b *big.Float) bool {
	return Cmp(a, b) > 0
}

func Lt(a, b *big.Float) bool {
	return Cmp(a, b) < 0
}

func Ge(a, b *big.Float) bool {
	return Cmp(a, b) >= 0
}

func Le(a, b *big.Float) bool {
	return Cmp(a, b) <= 0
}
